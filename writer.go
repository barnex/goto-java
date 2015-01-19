package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"
	"reflect"
)

const (
	CLASS   = "class"
	FINAL   = "final"
	PACKAGE = "package"
	PUBLIC  = "public"
)

type writer struct {
	fset      *token.FileSet
	fname     string // output file name for principal class (no extension)
	pkg       string
	out       io.Writer
	classDecl []*ast.Decl
}

func (w *writer) parseFile(f *ast.File) {
	w.pkg = f.Name.Name

	for _, decl := range f.Decls {
		switch decl.(type) {
		default:
			w.error(f, "cannot handle ", reflect.TypeOf(decl))
		}
	}
}

func (w *writer) error(n ast.Node, msg ...interface{}) {
	fatal(append([]interface{}{w.pos(n), ": "}, msg...)...)
}

func (w *writer) pos(n ast.Node) token.Position {
	return w.fset.Position(n.Pos())
}

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

	w.putln(PACKAGE, w.pkg)
	w.println()

	className := path.Base(w.fname)
	w.put(PUBLIC, FINAL, CLASS, className, "{\n")

	w.genClassDecls()

	w.put("}\n")
}

func (w *writer) genClassDecls() {

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
