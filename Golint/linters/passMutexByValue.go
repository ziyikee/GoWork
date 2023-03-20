package linters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"honnef.co/go/tools/analysis/code"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/analysis/report"
)

var PassMutexByValueAnalyzer = &analysis.Analyzer{
	Name: "PassMutexByValueAnalyzer",
	Doc:  "Pass or refer to a Mutex or a receiver containing a Mutex as a value type",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		generated.Analyzer,
	},
	Run: passMutexByValueRun,
}

/*
*1. find: selectorExpr: X.Lock() X.Unlock
 */
func passMutexByValueRun(pass *analysis.Pass) (interface{}, error) {
	fn := func(node ast.Node) {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println("recover error")
		//	}
		//}()
		selExpr, ok := node.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if selExpr.Sel != nil && (selExpr.Sel.Name == "Lock" || selExpr.Sel.Name == "Unlock") {

			if ident, ok := selExpr.X.(*ast.Ident); ok {
				if filed, ok := ident.Obj.Decl.(*ast.Field); ok { //排除在函数体内声明的mutex
					if id, ok := filed.Type.(*ast.StarExpr); !ok && id == nil {
						//fmt.Println("variable: " + ident.Obj.Name + " op: " + selExpr.Sel.Name)
						//report针对同一个ident自动去重了，不用担心lock和unlock会报告两次
						report.Report(pass, ident.Obj, "Mutex variable: "+ident.Obj.Name+" may be passed by value.")
					}
				}
			}

			if ident, ok := selExpr.X.(*ast.SelectorExpr); ok {
				if id, ok := ident.X.(*ast.Ident); ok {
					if obj := id.Obj; obj != nil {
						if filed, ok := obj.Decl.(*ast.Field); ok { //排除在函数体内声明的mutex
							if expr, ok := filed.Type.(*ast.StarExpr); !ok && expr == nil {
								//fmt.Println("variable: " + ident.Obj.Name + " op: " + selExpr.Sel.Name)
								//report针对同一个ident自动去重了，不用担心lock和unlock会报告两次
								report.Report(pass, id.Obj, "Mutex variable: "+id.Obj.Name+" may be passed by value.")
							}
						}
					}
				}
			}
		}
	}

	code.Preorder(pass, fn, (*ast.SelectorExpr)(nil))
	return nil, nil
}
