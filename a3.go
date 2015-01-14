package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"log"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	w := newWorld()

	pkgs, err := parser.ParseDir(&w.fset, path, nil, 0)
	checkUserErr(err)

	if len(pkgs) == 0 {
		fatal("no packages found in ", path)
	}
	if len(pkgs) > 1 {
		fatal(len(pkgs), " packages found in ", path)
	}

	for _, v := range pkgs { // pick the only package
		w.compilePkg(v)
	}
}

type world struct {
	fset token.FileSet
}

func newWorld() *world {
	return &world{}
}

func (w *world) compilePkg(pkg *ast.Package) {
	if len(pkg.Imports) > 0 {
		fatal("imports not supported")
	}

	for _, f := range pkg.Files {
		w.compileFile(f)
	}
}

func (w *world) compileFile(f *ast.File) {

}

func checkUserErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(msg ...interface{}) {
	log.Fatal(msg...)
}
