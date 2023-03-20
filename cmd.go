package main

import (
	"GoLint/linters"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(linters.ClosureErrAnalyzer)
	//multichecker.Main(linters.HGAnalyzer, linters.WgAddAnalyzer, linters.WaitInLoopAnalyzer, linters.ClosureErrAnalyzer, linters.PassMutexByValueAnalyzer)
}
