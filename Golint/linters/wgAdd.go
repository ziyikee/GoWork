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

/*
*检测思路：只判断go func(){}()，即匿名函数中对sync.WaitGroup.Add()的直接调用
 */
var WgAddAnalyzer = &analysis.Analyzer{
	Name: "wgAddAnalyzer",
	Doc:  "Check if directly calling WaitGroup.add() in anonymous goroutine",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		generated.Analyzer,
	},
	Run: wgAddAnalyzerRun,
}

var wgAddStmt = pattern.MustParse(`(SelectorExpr x (Ident "Add") )`)

type findWaitAddVisitor struct {
	ident []*ast.Ident
	pass  *analysis.Pass
}

func (c *findWaitAddVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		if m, ok := pattern.Match(wgAddStmt, n); ok {
			if ident, ok := m.State["x"].(*ast.Ident); ok && c.pass.TypesInfo.TypeOf(ident) != nil && c.pass.TypesInfo.TypeOf(ident).String() == "sync.WaitGroup" {
				c.ident = append(c.ident, ident)
			}
		}
		return c
	case *ast.GoStmt:
		c.ident = nil
	}
	return c
}

func wgAddAnalyzerRun(pass *analysis.Pass) (interface{}, error) {

	fn := func(node ast.Node) {
		//defer func() {
		//	if err := recover(); err != nil {
		//		log.Println("recover error")
		//	}
		//}()
		gostmt, ok := node.(*ast.GoStmt).Call.Fun.(*ast.FuncLit)
		if !ok {
			return
		}
		fcv := &findWaitAddVisitor{pass: pass}
		ast.Walk(fcv, gostmt)
		for _, ident := range fcv.ident {
			report.Report(pass, ident, "Variable "+ident.Name+" calls `$WG.Add()` in anonymous goroutine")
		}

	}

	code.Preorder(pass, fn, (*ast.GoStmt)(nil))
	return nil, nil
}
