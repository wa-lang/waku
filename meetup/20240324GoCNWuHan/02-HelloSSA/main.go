package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
)

const src = `
package main

func main() {
	println("Hello, GoCN!")
	println("The answer is:", 42)
}
`

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "test.go", src, parser.AllErrors)

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{Importer: nil}
	pkg, _ := conf.Check("test.go", fset, []*ast.File{f}, info)

	var ssaProg = ssa.NewProgram(fset, ssa.SanityCheckFunctions)
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true)

	ssaPkg.Build()
	ssaPkg.Func("main").WriteTo(os.Stdout)

	runFunc(ssaPkg.Func("main"))
}
