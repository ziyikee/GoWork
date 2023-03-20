package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var srcCode = `
package test

func example() {
	var ch1 = make(chan string)
	ch2 := make(chan string)
	var ch3 = make(chan string, 1)
	ch4 := make(chan string, 1)

	go func() {
		str := "aaa"
		ch1 <- str
		ch2 <- str
		ch3 <- str
		ch4 <- str
	}()
	a := 1
	select {
	case <-ch1:
		a++
	case b := <-ch2:
		a = len(b)
	}
}
`

/*
*打印AST树
 */
func Test(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "D:\\Projects\\GoLint\\doc\\PassMutexByValue.go", nil, 0)
	if err != nil {
		fmt.Printf("err = %s", err)
	}
	ast.Print(fset, f)
}
