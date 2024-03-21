package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

const src = `
package main

var answer = 40 + 2

func main() {
	println("Hello, GoCN")
	println(answer)
}
`

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
	ast.Print(nil, f)
}
