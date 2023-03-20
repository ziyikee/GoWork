package linters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"honnef.co/go/tools/analysis/code"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/analysis/report"
)

var ClosureErrAnalyzer = &analysis.Analyzer{
	Name: "ClosureErrorAnalyzer",
	Doc:  " Data race due to loop index variable capture",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		generated.Analyzer,
	},
}

type identVisitor struct {
	objs   map[*ast.Object]struct{}
	gostmt *ast.GoStmt
	pass   *analysis.Pass
}

// Visit：对gostmt中的所有ident进行检测，判断是否为循环变量
func (c *identVisitor) Visit(node ast.Node) ast.Visitor {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return c
	}
	if _, ok := c.objs[ident.Obj]; ok {
		report.Report(c.pass, ident, "Variable ["+ident.Name+"] is freely referenced by different goroutines, causing a data race")
	}
	return c
	//存在变量被引用，则判断是自有引用还是参数引用
	//上面的判断没有必要，因为如果是参数引用的话，函数体内的obj1引用的是函数的形参obj1，而尽管把for循环的变量obj2当做形参obj1，obj1和obj2还是两个obj，即值传递，传递obj副本
}

func run(pass *analysis.Pass) (interface{}, error) {
	fn := func(node ast.Node) {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println("recover error")
		//	}
		//}()
		switch node.(type) {
		case *ast.RangeStmt:
			rangeStmt, _ := node.(*ast.RangeStmt)
			idVisitor := &identVisitor{pass: pass}
			//1.确定循环变量
			idVisitor.objs = findVarsInRange(rangeStmt)
			if len(idVisitor.objs) == 0 {
				break
			}
			stmtList := rangeStmt.Body.List
			//寻找for的block直接调用的gostmt
			for _, stmt := range stmtList {
				gostmt, ok := stmt.(*ast.GoStmt)
				if !ok {
					continue
				}
				idVisitor.gostmt = gostmt
				//只遍历block部分
				ast.Walk(idVisitor, gostmt.Call.Fun)
			}
		case *ast.ForStmt:
			forStmt, _ := node.(*ast.ForStmt)
			idVisitor := &identVisitor{pass: pass}
			//1.确定循环变量
			idVisitor.objs = findVarsInFor(forStmt)
			if len(idVisitor.objs) == 0 {
				break
			}
			stmtList := forStmt.Body.List
			//寻找for的block直接调用的gostmt
			for _, stmt := range stmtList {
				gostmt, ok := stmt.(*ast.GoStmt)
				if !ok {
					continue
				}
				idVisitor.gostmt = gostmt
				//只遍历block部分
				ast.Walk(idVisitor, gostmt.Call.Fun)
			}
		}
	}

	code.Preorder(pass, fn, (*ast.RangeStmt)(nil), (*ast.ForStmt)(nil))
	return nil, nil
}

// 可优化，进一步细化循环中变量的收集规则，目前不全，可能变量在block中更新
// 收集for range循环中的每次更新的变量,for range中迭代赋值语句为AssignStmt，收集其左值
func findVarsInRange(ident *ast.RangeStmt) map[*ast.Object]struct{} {
	keys := make(map[*ast.Object]struct{})
	obj, ok := ident.Key.(*ast.Ident)
	if !ok {
		return keys
	}
	astmt, ok := obj.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return keys
	}
	for _, expr := range astmt.Lhs {
		if id, ok := expr.(*ast.Ident); ok {
			keys[id.Obj] = struct {
			}{}
		}
	}
	return keys
}

// 可优化，进一步细化循环中变量的收集规则，目前不全，可能变量在block中更新
// 常规for循环中，这里主要收集的位置:
// Init部分，对AssignStmt的左值lhs进行收集
// Cond部分，主要考虑二元表达式的左右两边的ident，Cond可能很复杂，目前只考虑简单的
// Post部分，即对变量做修改的部分，主要考虑自增自减类型的表达式的X部分和赋值表达式的左值
func findVarsInFor(ident *ast.ForStmt) map[*ast.Object]struct{} {
	keys := make(map[*ast.Object]struct{})
	//Init部分
	if ident.Init != nil {
		if astmt, ok := ident.Init.(*ast.AssignStmt); ok {
			for _, expr := range astmt.Lhs {
				if id, ok := expr.(*ast.Ident); ok {
					keys[id.Obj] = struct{}{}
				}
			}
		}
	}
	//Cond部分
	if ident.Cond != nil {
		if binStmt, ok := ident.Cond.(*ast.BinaryExpr); ok {
			if left, ok := binStmt.X.(*ast.Ident); ok {
				keys[left.Obj] = struct{}{}
			}
			if right, ok := binStmt.Y.(*ast.Ident); ok {
				keys[right.Obj] = struct{}{}
			}
		}
	}
	//Post部分
	if ident.Post != nil {
		switch ident.Post.(type) {
		case *ast.AssignStmt:
			if astmt, ok := ident.Post.(*ast.AssignStmt); ok {
				for _, expr := range astmt.Lhs {
					if id, ok := expr.(*ast.Ident); ok {
						keys[id.Obj] = struct{}{}
					}
				}
			}
		case *ast.IncDecStmt:
			if idStmt, ok := ident.Post.(*ast.IncDecStmt); ok {
				if x, ok := idStmt.X.(*ast.Ident); ok {
					keys[x.Obj] = struct{}{}
				}
			}
		}
	}

	return keys
}
