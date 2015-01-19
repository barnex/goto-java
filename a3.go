package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
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
	f, err := parser.ParseFile(fset, fname, nil, 0)
	checkUserErr(err)

	// print ast if requested
	if *flagPrint {
		ast.Print(fset, f)
	}

	//// prepare output file
	outFile := fname[:len(fname)-len(path.Ext(fname))]

	// transpile
	w := &writer{fname: outFile}
	w.parseFile(f)
	w.genCode()
}

type writer struct {
	fname string // output file name for principal class (no extension)
	pkg   string
	out   io.Writer
}

func (w *writer) parseFile(f *ast.File) {
	w.pkg = f.Name.Name
}

const (
	CLASS   = "class"
	FINAL   = "final"
	PACKAGE = "package"
	PUBLIC  = "public"
)

func (w *writer) initOut() {
	if w.out != nil {
		panic("already inited")
	}
	out, err := os.Create(w.fname + ".java")
	checkUserErr(err)
	w.out = out
	// TODO: buffer
}

func (w *writer) genCode() {
	w.initOut()

	className := path.Base(w.fname)
	w.putln(PACKAGE, className)
	w.println()

	w.put(PUBLIC, FINAL, CLASS, className, "{\n")

	w.put("}\n")
}

func (w *writer) put(tokens ...interface{}) {
	for i, t := range tokens {
		if i != 0 {
			fmt.Fprint(w.out, " ")
		}
		fmt.Fprint(w.out, t)
	}
}

func (w *writer) putln(tokens ...interface{}) {
	w.put(tokens...)
	w.put(";\n")
}

func (w *writer) println(x ...interface{}) {
	w.put(x...)
	w.put("\n")
}

func transpile(out io.Writer, f *ast.File) {

}
