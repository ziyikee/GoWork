package linters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/analysis/report"
	"honnef.co/go/tools/pattern"
)

var WaitInLoopAnalyzer = &analysis.Analyzer{
	Name: "waitInLoopAnalyzer",
	Doc:  "Calling `$WG.Wait()` inside a loop blocks the call to `$WG.Done()`",
	Run:  waitInLoopAnalyzerRun,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		generated.Analyzer,
	},
}

var (
	//识别sync.WaitGroup变量, 1. wg := sync.WaitGroup{}  2. var wg =var wg = sync.WaitGroup{}
	wgDef = pattern.MustParse(`(Or
		(ValueSpec x (SelectorExpr (Ident "sync") (Ident "WaitGroup")) _)
		(AssignStmt x _ (CompositeLit (SelectorExpr (Ident "sync") (Ident "WaitGroup")) _))
	)`)

	wgDoneStmt = pattern.MustParse(`(Or
		(SelectorExpr x (Ident "Done")) 
		(CallExpr (SelectorExpr x (Ident "Done")) _)
	)`)

	wgWaitStmt = pattern.MustParse(`(CallExpr (SelectorExpr x (Ident "Wait")) _)`)
)

func waitInLoopAnalyzerRun(pass *analysis.Pass) (interface{}, error) {
	fn := func(node ast.Node) {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println("recover error")
		//	}
		//}()
		var stmtListInLoop []ast.Stmt
		switch loopStmt := node.(type) {
		case *ast.ForStmt:
			stmtListInLoop = loopStmt.Body.List
		case *ast.RangeStmt:
			stmtListInLoop = loopStmt.Body.List
		}
		identsInLoop := make(map[*ast.Object]struct{})

		for j := 0; j < len(stmtListInLoop); j++ {
			collectWgIdent(stmtListInLoop[j], identsInLoop)

			targetIdent, ok := isFindDoneInGoStmt(pass, stmtListInLoop[j])
			if !ok {
				continue
			}
			//循环内定义的没问题，跳过
			if _, ok := identsInLoop[targetIdent.Obj]; ok {
				continue
			}
			for k := j + 1; k < len(stmtListInLoop); k++ {
				if ident, ok := isFindWaitAfterGoStmt(pass, stmtListInLoop[k], targetIdent); ok {
					report.Report(pass, ident, "Variable "+ident.Name+" calls `$WG.Wait()` inside a loop may block the call to `$WG.Done()`")
				}
			}

		}
	}
	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, fn)
	return nil, nil
}

// 判断stmt是不是wg变量，如果是，加入到identList中
func collectWgIdent(stmt ast.Stmt, identList map[*ast.Object]struct{}) {
	switch stmt := stmt.(type) {
	case *ast.AssignStmt:
		if value, ok := pattern.Match(wgDef, stmt); ok {
			if idents, ok := value.State["x"].([]ast.Expr); ok {
				for _, x := range idents {
					if ident, ok := x.(*ast.Ident); ok {
						identList[ident.Obj] = struct{}{}
					}
				}
			}
		}
	case *ast.DeclStmt:
		genDecl, ok := stmt.Decl.(*ast.GenDecl)
		if !ok {
			return
		}
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				if value, ok := pattern.Match(wgDef, valueSpec); ok {
					if idents, ok := value.State["x"].([]*ast.Ident); ok {
						for _, ident := range idents {
							identList[ident.Obj] = struct{}{}
						}
					}
				}
			}
		}
	}
}

// 判断stmt是不是GoStmt，且是不是匿名函数，且函数体内有没有执行wg.Done()操作，如果有，则返回wg变量
func isFindDoneInGoStmt(pass *analysis.Pass, stmt ast.Stmt) (*ast.Ident, bool) {
	var ident *ast.Ident
	goStmt, ok := stmt.(*ast.GoStmt)
	if !ok {
		return ident, false
	}
	goFuncList, ok := goStmt.Call.Fun.(*ast.FuncLit)
	if !ok {
		return ident, false
	}
	//遍历goStmt中的stmt
	for _, bodyStmt := range goFuncList.Body.List {
		deferStmt, ok := bodyStmt.(*ast.DeferStmt)
		if ok {
			switch expr := deferStmt.Call.Fun.(type) {
			case *ast.SelectorExpr:
				if m, ok := pattern.Match(wgDoneStmt, expr); ok {
					if ident, ok := m.State["x"].(*ast.Ident); ok && pass.TypesInfo.TypeOf(ident) != nil && pass.TypesInfo.TypeOf(ident).String() == "sync.WaitGroup" {
						return ident, true
					}
				}
			case *ast.FuncLit:
				for _, fnStmt := range expr.Body.List {
					if m, ok := pattern.Match(wgDoneStmt, fnStmt); ok {
						if ident, ok := m.State["x"].(*ast.Ident); ok && pass.TypesInfo.TypeOf(ident) != nil && pass.TypesInfo.TypeOf(ident).String() == "sync.WaitGroup" {
							return ident, true
						}
					}
				}
			}
		} else {
			if m, ok := pattern.Match(wgDoneStmt, bodyStmt); ok {
				if ident, ok := m.State["x"].(*ast.Ident); ok && pass.TypesInfo.TypeOf(ident) != nil && pass.TypesInfo.TypeOf(ident).String() == "sync.WaitGroup" {
					return ident, true
				}
			}
		}
	}
	return ident, false
}

func isFindWaitAfterGoStmt(pass *analysis.Pass, stmt ast.Stmt, targetIdent *ast.Ident) (*ast.Ident, bool) {
	var ident *ast.Ident
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return ident, false
	}
	if m, ok := pattern.Match(wgWaitStmt, exprStmt); ok {
		if ident, ok := m.State["x"].(*ast.Ident); ok && ident.Obj == targetIdent.Obj && pass.TypesInfo.TypeOf(ident) != nil && pass.TypesInfo.TypeOf(ident).String() == "sync.WaitGroup" {
			return ident, true
		}
	}
	return ident, false
}
