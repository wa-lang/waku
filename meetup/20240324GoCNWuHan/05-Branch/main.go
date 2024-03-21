package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/ssa"
)

const src = `
package main

func add(i int, j int) int{
	return i + j
}

func fib(i0, i1, n int) (ret int) {
	print(i1, " ")
	if n <= 2 {
		ret = i0 + i1
	} else {
		ret = fib(i1, i0+i1, n-1)
	}
	return
}

func main() {
	var i int
	if add(3, 5) < 9{
		i = 13
	} else{
		i = 42
	}
	println(fib(0, 1, i))
}
`

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{Importer: nil}
	pkg, err := conf.Check("test.go", fset, []*ast.File{f}, info)
	if err != nil {
		log.Fatal(err)
	}

	var ssaProg = ssa.NewProgram(fset, ssa.SanityCheckFunctions)
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true)

	ssaPkg.Build()
	ssaPkg.WriteTo(os.Stdout)
	ssaPkg.Func("main").WriteTo(os.Stdout)
	ssaPkg.Func("fib").WriteTo(os.Stdout)

	p := NewEngine(ssaPkg)
	p.initGlobals()

	p.runFunc(ssaPkg.Func("main"), nil)
}
