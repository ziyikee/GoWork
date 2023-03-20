package doc

import "sync"

type Holder struct {
	sync.Mutex
	mutex2 sync.Mutex
	count  int
}

/*
*

	Recv: *ast.FieldList {
	   107  .  .  .  .  Opening: D:\Projects\GoLint\doc\PassMutexByValue.go:11:6
	   108  .  .  .  .  List: []*ast.Field (len = 1) {
	   109  .  .  .  .  .  0: *ast.Field {
	   110  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
	   111  .  .  .  .  .  .  .  0: *ast.Ident {
	   112  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:11:7
	   113  .  .  .  .  .  .  .  .  Name: "h"
	   114  .  .  .  .  .  .  .  .  Obj: *ast.Object {
	   115  .  .  .  .  .  .  .  .  .  Kind: var
	   116  .  .  .  .  .  .  .  .  .  Name: "h"
	   117  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 109)
	   118  .  .  .  .  .  .  .  .  }
	   119  .  .  .  .  .  .  .  }
	   120  .  .  .  .  .  .  }
	   121  .  .  .  .  .  .  Type: *ast.Ident {
	   122  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:11:9
	   123  .  .  .  .  .  .  .  Name: "Holder"
	   124  .  .  .  .  .  .  .  Obj: *(obj @ 32)
	   125  .  .  .  .  .  .  }
	   126  .  .  .  .  .  }
	   127  .  .  .  .  }
*/
func (h Holder) func1() {
	h.Lock()
	h.count++
	h.Unlock()
}

/*
*

	Recv: *ast.FieldList {
		   201  .  .  .  .  Opening: D:\Projects\GoLint\doc\PassMutexByValue.go:17:6
		   202  .  .  .  .  List: []*ast.Field (len = 1) {
		   203  .  .  .  .  .  0: *ast.Field {
		   204  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
		   205  .  .  .  .  .  .  .  0: *ast.Ident {
		   206  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:17:7
		   207  .  .  .  .  .  .  .  .  Name: "h"
		   208  .  .  .  .  .  .  .  .  Obj: *ast.Object {
		   209  .  .  .  .  .  .  .  .  .  Kind: var
		   210  .  .  .  .  .  .  .  .  .  Name: "h"
		   211  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 203)
		   212  .  .  .  .  .  .  .  .  }
		   213  .  .  .  .  .  .  .  }
		   214  .  .  .  .  .  .  }
		   215  .  .  .  .  .  .  Type: *ast.StarExpr {
		   216  .  .  .  .  .  .  .  Star: D:\Projects\GoLint\doc\PassMutexByValue.go:17:9
		   217  .  .  .  .  .  .  .  X: *ast.Ident {
		   218  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:17:10
		   219  .  .  .  .  .  .  .  .  Name: "Holder"
		   220  .  .  .  .  .  .  .  .  Obj: *(obj @ 32)
		   221  .  .  .  .  .  .  .  }
		   222  .  .  .  .  .  .  }
		   223  .  .  .  .  .  }
		   224  .  .  .  .  }
		   225  .  .  .  .  Closing: D:\Projects\GoLint\doc\PassMutexByValue.go:17:16
		   226  .  .  .  }
*/
func (h *Holder) func2() {
	h.Lock()
	h.count++
	h.Unlock()
}

func (h Holder) func3() {
	h.mutex2.Lock()
	h.count++
	h.mutex2.Unlock()
}

func (h *Holder) func4() {
	h.mutex2.Lock()
	h.count++
	h.mutex2.Unlock()
}

/*
*

	0: *ast.Field {
	   527  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
	   528  .  .  .  .  .  .  .  .  0: *ast.Ident {
	   529  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:35:12
	   530  .  .  .  .  .  .  .  .  .  Name: "mutex5"
	   531  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
	   532  .  .  .  .  .  .  .  .  .  .  Kind: var
	   533  .  .  .  .  .  .  .  .  .  .  Name: "mutex5"
	   534  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 526)
	   535  .  .  .  .  .  .  .  .  .  }
	   536  .  .  .  .  .  .  .  .  }
	   537  .  .  .  .  .  .  .  }
	   538  .  .  .  .  .  .  .  Type: *ast.SelectorExpr {
	   539  .  .  .  .  .  .  .  .  X: *ast.Ident {
	   540  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:35:19
	   541  .  .  .  .  .  .  .  .  .  Name: "sync"
	   542  .  .  .  .  .  .  .  .  }
	   543  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
	   544  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:35:24
	   545  .  .  .  .  .  .  .  .  .  Name: "Mutex"
	   546  .  .  .  .  .  .  .  .  }
	   547  .  .  .  .  .  .  .  }
	   548  .  .  .  .  .  .  }
*/
func func5(mutex5 sync.Mutex) {
	i := 0
	mutex5.Lock()
	i++
	mutex5.Unlock()
}

/*
*

	0: *ast.Field {
	   643  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
	   644  .  .  .  .  .  .  .  .  0: *ast.Ident {
	   645  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:42:12
	   646  .  .  .  .  .  .  .  .  .  Name: "mutex6"
	   647  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
	   648  .  .  .  .  .  .  .  .  .  .  Kind: var
	   649  .  .  .  .  .  .  .  .  .  .  Name: "mutex6"
	   650  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 642)
	   651  .  .  .  .  .  .  .  .  .  }
	   652  .  .  .  .  .  .  .  .  }
	   653  .  .  .  .  .  .  .  }
	   654  .  .  .  .  .  .  .  Type: *ast.StarExpr {
	   655  .  .  .  .  .  .  .  .  Star: D:\Projects\GoLint\doc\PassMutexByValue.go:42:19
	   656  .  .  .  .  .  .  .  .  X: *ast.SelectorExpr {
	   657  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
	   658  .  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:42:20
	   659  .  .  .  .  .  .  .  .  .  .  Name: "sync"
	   660  .  .  .  .  .  .  .  .  .  }
	   661  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
	   662  .  .  .  .  .  .  .  .  .  .  NamePos: D:\Projects\GoLint\doc\PassMutexByValue.go:42:25
	   663  .  .  .  .  .  .  .  .  .  .  Name: "Mutex"
	   664  .  .  .  .  .  .  .  .  .  }
	   665  .  .  .  .  .  .  .  .  }
	   666  .  .  .  .  .  .  .  }
*/
func func6(mutex6 *sync.Mutex) {
	i := 0
	mutex6.Lock()
	i++
	mutex6.Unlock()
}

func func7() {
	mutex := sync.Mutex{}
	i := 0
	mutex.Lock()
	i++
	mutex.Unlock()
}

func func8() {
	var mutex sync.Mutex
	i := 0
	mutex.Lock()
	i++
	mutex.Unlock()
}
