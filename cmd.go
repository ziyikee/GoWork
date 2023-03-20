package main

import (
	"GoLint/linters"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	//singlechecker.Main(linters.PassMutexByValueAnalyzer)
	multichecker.Main(linters.HGAnalyzer, linters.WgAddAnalyzer, linters.WaitInLoopAnalyzer, linters.ClosureErrAnalyzer, linters.PassMutexByValueAnalyzer)
}
