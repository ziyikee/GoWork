# GoLinter

## Go AST

1. `*ast.Ident`是Go语言AST（Abstract Syntax Tree，抽象语法树）包中的一个结构体类型，表示**代码中的标识符（identifier），如变量名、函数名、类型名**等。

   这个结构体类型有如下几个重要的字段：

   - `Name`: 标识符的名称，是一个字符串类型的字面值。
   - `Obj`: 标识符的对象（Object），是一个指向一个对象的指针，这个对象可能是一个变量、常量、函数、类型等。
   - `Pos`: 标识符在源代码中的位置，是一个token.Pos类型的值。

   在Go语言的AST中，`*ast.Ident`类型经常出现在各种表达式中，例如变量声明、函数调用、运算符等等。通过使用这个结构体类型，我们可以方便地获取标识符的名称、位置等信息，从而进行语法分析和代码检查等操作。

2. `ast.Expr` 是一个表示 Go 语言中**所有表达式的接口类型**。它有**多个实现类型**，包括但不限于 `ast.BasicLit`（表示基本字面量），`ast.Ident`（表示标识符），`ast.SelectorExpr`（表示选择器表达式）等等。在 AST 中，许多节点都可能包含一个 `ast.Expr` 类型的表达式字段，因此该接口类型的实现类别是非常丰富的。

   在代码中，当我们需要处理一个表达式节点时，我们可以使用 `ast.Expr` 类型来表示它，这样就可以接收到**各种可能的表达式类型**了。例如，当我们处理一个函数调用时，其参数列表中的每个参数都是一个表达式，因此我们可以使用 `ast.Expr` 来表示它们。

3. `SelectorExpr`是Go语言AST中的一种节点类型，表示选择表达式，即对一个选择器的访问，例如`x.y`、`pkg.Type.Method`等。

   `SelectorExpr`有以下属性：

   - `X`：表示被选择的表达式，类型为`ast.Expr`，通常是一个标识符或一个字段访问表达式。
   - `Sel`：表示被选择的标识符，类型为`*ast.Ident`，表示被选择的标识符的名称。

   例如，对于表达式`foo.bar`，`foo`是`X`，`bar`是`Sel`。

[Golang的抽象语法树(AST) Step By Step - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/380421057)
https://toutiao.io/posts/fzt61t/preview

https://juejin.cn/post/7130188153792495630

WaitGroup：https://zhuanlan.zhihu.com/p/527193688

https://go101.org/article/concurrent-common-mistakes.html

Go AST结构：https://zhuanlan.zhihu.com/p/380421057，https://zhuanlan.zhihu.com/p/28516587

Go 断言：http://c.biancheng.net/view/4281.html

https://www.infoq.cn/article/n5k2unejil6cjtqghdvh

[https://draveness.me/golang/docs/part1-prerequisite/ch02-compile/golang-lexer-and-parser/#%E8%AF%8D%E6%B3%95%E5%88%86%E6%9E%90](https://draveness.me/golang/docs/part1-prerequisite/ch02-compile/golang-lexer-and-parser/#词法分析)

Goroutine leak and **fix** ：https://juejin.cn/post/7033711399041761311 

Data Race：https://zhuanlan.zhihu.com/p/532060939

Golangcli-lint https://blog.csdn.net/wohu1104/article/details/113751501

## ClosureError

Data Race：https://zhuanlan.zhihu.com/p/532060939

对于循环中每次被更新的变量，在for循环体中如果存在gostmt，且该变量被gostmt引用，则会出现goroutine的延迟绑定，不同goroutine可能拿到的变量是相同的变量值

![image-20230314150006328](C:\Users\yefengyuan\AppData\Roaming\Typora\typora-user-images\image-20230314150006328.png)

检测规则：

1. 检测范围在for循环体中
2. 确定for循环每次更新的变量(对于在block内更新的变量暂不考虑)
3. 确定for循环中是否存在gostmt
4. 确定更新的变量是否在gostmt中被引用，且不是通过参数传参的形式进行引用的

## pass_mutex_by_value 

https://github.com/trailofbits/semgrep-rules/blob/main/go/sync-mutex-value-copied.yaml

https://go101.org/article/concurrent-common-mistakes.html

https://stackoverflow.com/questions/49808622/sync-mutex-and-sync-mutex-which-is-better

主要是锁在作为参数时，不能以值的方式进行传递，因为通常情况下锁是用来做并发控制，不同的线程是共享同一个锁，若通过值传递，会得到锁的副本，就不是同一把锁。

检测思路：

不要主动去检测参数是不是锁，是值还是引用，而是先去找锁使用的地方，再反过来去看是通过什么方式传递的。

传递锁的方式有两种，一种通过函数形参传递，另一个中通过函数接收器传递，锁可能是接收器或者接收器的属性，此时需要判断接收器是值形式还是引用形式。
