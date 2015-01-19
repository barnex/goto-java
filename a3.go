package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var (
	flagPrint = flag.Bool("print", true, "Print ast")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	for _, f := range flag.Args() {
		handleFile(f)
	}

}

func handleFile(fname string) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, 0)
	checkUserErr(err)

	if *flagPrint {
		ast.Print(fset, f)
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
