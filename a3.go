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

	outFile := fname[:len(fname)-len(path.Ext(fname))] + ".java"
	out_, err := os.Create(outFile)
	checkUserErr(err)
	defer out_.Close()
	out := bufio.NewWriter(out_)
	defer out.Flush()

}

func checkUserErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(msg ...interface{}) {
	log.Fatal(msg...)
}
