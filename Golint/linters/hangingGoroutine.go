package linters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"honnef.co/go/tools/analysis/code"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/analysis/report"
	"honnef.co/go/tools/pattern"
)

var HGAnalyzer = &analysis.Analyzer{
	Name: "HangingGouroutineAnalyzer",
	Doc:  "Potential goroutine leak due to unbuffered channel send inside loop or unbuffered channel receive in select block",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		generated.Analyzer,
	},
	Run: hangingGoroutineRun,
}
var channelCheck = pattern.MustParse(`(Or
   (ValueSpec _ (ChanType _ _) _)
   (AssignStmt _ _ (CallExpr (Builtin "make") [(ChanType _ _)]))
)`) // unbuffered channel

// 保存goroutine中有send操作的unbuffered channel
var channelSend []*ast.Ident

type findSendStmtVisitor struct {
	ident []*ast.Ident
	pass  *analysis.Pass
}

func (c *findSendStmtVisitor) Visit(node ast.Node) ast.Visitor {
	sendStmt, ok := node.(*ast.SendStmt)
	if !ok {
		return c
	}
	if ident, ok := sendStmt.Chan.(*ast.Ident); ok && checkUnbuffered(ident) { // 有sendstmt，且该channel创建时为无缓冲状态
		c.ident = append(c.ident, ident)
	}
	return c
}

/**
1: *ast.SendStmt {
  1325  .  .  .  .  .  .  .  .  .  .  .  Chan: *ast.Ident {
  1326  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: 77:3
  1327  .  .  .  .  .  .  .  .  .  .  .  .  Name: "ch"
  1328  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 1246)
  1329  .  .  .  .  .  .  .  .  .  .  .  }
  1330  .  .  .  .  .  .  .  .  .  .  .  Arrow: 77:6
  1331  .  .  .  .  .  .  .  .  .  .  .  Value: *ast.Ident {
  1332  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: 77:9
  1333  .  .  .  .  .  .  .  .  .  .  .  .  Name: "newData"
  1334  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 1302)
  1335  .  .  .  .  .  .  .  .  .  .  .  }
  1336  .  .  .  .  .  .  .  .  .  .  }
*/

func isUnbuffered(ident *ast.Ident) bool {
	//if ident.Obj.Decl.(*ast.ValueSpec).
	Obj := ident.Obj
	if Obj == nil {
		return false
	}
	switch Obj.Decl.(type) {
	case *ast.AssignStmt:
		if _, ok := pattern.Match(channelCheck, ident.Obj.Decl.(*ast.AssignStmt).Rhs[0]); ok {
			return true
		}
	case *ast.ValueSpec:
		if _, ok := pattern.Match(channelCheck, ident.Obj.Decl.(*ast.ValueSpec)); ok {
			return true
		}
	case nil:
	}
	return false
}

func checkUnbuffered(ident *ast.Ident) bool {
	Obj := ident.Obj
	if Obj == nil {
		return false
	}
	switch decl := Obj.Decl.(type) {
	case *ast.ValueSpec:
		for i := 0; i < len(decl.Names); i++ {
			if ident.Name == decl.Names[i].Name {
				if call, ok := decl.Values[i].(*ast.CallExpr); ok {
					if call.Fun.(*ast.Ident).Name == "make" {
						_, ok := call.Args[0].(*ast.ChanType)
						if ok {
							if len(call.Args) == 1 {
								return true
							} else if lit, ok := call.Args[1].(*ast.BasicLit); ok && lit.Value == "0" {
								return true
							}
							return false
						}
					}
				}
			}
		}
	case *ast.AssignStmt:
		for i := 0; i < len(decl.Lhs); i++ {
			if ident.Name == decl.Lhs[i].(*ast.Ident).Name {
				if call, ok := decl.Rhs[i].(*ast.CallExpr); ok {
					if call.Fun.(*ast.Ident).Name == "make" {
						_, ok := call.Args[0].(*ast.ChanType)
						if ok {
							if len(call.Args) == 1 {
								return true
							} else if lit, ok := call.Args[1].(*ast.BasicLit); ok && lit.Value == "0" {
								return true
							}
							return false
						}
					}
				}
			}
		}
	}
	return false
}

func hasChannelReceive(stmt ast.Stmt, targetIdent *ast.Ident) bool {

	if expr, ok := stmt.(*ast.ExprStmt); ok {
		switch e := expr.X.(type) {
		case *ast.UnaryExpr:
			if e.Op.String() == "<-" {
				if ident, ok := e.X.(*ast.Ident); ok && ident.Obj == targetIdent.Obj {
					return true
				}
			}
		}
	}
	if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
		for _, expr := range assignStmt.Rhs {
			switch e := expr.(type) {
			case *ast.UnaryExpr:
				if e.Op.String() == "<-" {
					if ident, ok := e.X.(*ast.Ident); ok && ident.Obj == targetIdent.Obj {
						return true
					}
				}
			}
		}
	}

	return false
}

func hangingGoroutineRun(pass *analysis.Pass) (interface{}, error) {

	fn := func(node ast.Node) {

		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println("recover error")
		//	}
		//}()

		switch typedNode := node.(type) {

		case *ast.GoStmt:
			// 查看当前goroutine是否有send操作
			fcv := &findSendStmtVisitor{pass: pass}
			ast.Walk(fcv, typedNode)
			channelSend = append(channelSend, fcv.ident...)
		case *ast.SelectStmt:
			// select-case中有channelSend的receive操作
			if len(typedNode.Body.List) == 0 {
				return
			}
			for _, i := range typedNode.Body.List {
				commClause, ok := i.(*ast.CommClause)
				if !ok {
					continue
				}
				for _, ident := range channelSend {
					if ok := hasChannelReceive(commClause.Comm, ident); ok {
						report.Report(pass, ident, "Sending message to channel variable "+ident.Name+" in a goroutine may be blocked, resulting in a goroutine leak")
					}
				}
			}
		}
	}

	code.Preorder(pass, fn, (*ast.GoStmt)(nil), (*ast.SelectStmt)(nil))
	return nil, nil
}
