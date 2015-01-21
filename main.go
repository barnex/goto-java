package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
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
	var info types.Info
	pkg, err := config.Check(fname, fset, []*ast.File{f}, &info)
	checkUserErr(err)
	fmt.Println(pkg)

	// print ast if requested
	if *flagPrint {
		ast.Print(fset, f)
	}

	//// prepare output file
	outFile := fname[:len(fname)-len(path.Ext(fname))]

	// transpile
	w := &writer{fname: outFile, fset: fset}
	w.parseFile(f)
	w.genCode()
}
