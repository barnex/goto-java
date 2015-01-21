package main

import (
	"bufio"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"

	"golang.org/x/tools/go/types"
)

var (
	flagPrint = flag.Bool("print", false, "Print ast")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	for _, f := range flag.Args() {
		handleFile(f)
	}

}

func checkUserErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

func handleFile(fname string) {

	// read and parse file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	checkUserErr(err)

	var config types.Config
	info := types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	_, err = config.Check(fname, fset, []*ast.File{f}, &info)
	checkUserErr(err)

	// print ast if requested
	if *flagPrint {
		ast.Print(fset, f)
	}

	//// prepare output file
	outFile := fname[:len(fname)-len(path.Ext(fname))]
	out_, errOut := os.Create(outFile + ".java")
	checkUserErr(errOut)
	defer out_.Close()
	out := bufio.NewWriter(out_)
	defer out.Flush()

	// transpile
	w := &writer{out: out, fset: fset, info: info}
	w.putClass(outFile, f)
}
