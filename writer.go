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
	EOL     = ";"
	FINAL   = "final"
	LBRACE  = "{"
	PACKAGE = "package"
	PUBLIC  = "public"
	RBRACE  = "}"
	STATIC  = "static"
	VOID    = "void"
)

type writer struct {
	fset         *token.FileSet
	fname        string // output file name for principal class (no extension)
	pkg          string
	out          io.Writer
	classMembers []ast.Decl
	indent       int
}

// parsing

func (w *writer) parseFile(f *ast.File) {
	w.pkg = f.Name.Name

	for _, decl := range f.Decls {
		switch n := decl.(type) {
		default:
			w.error(f, "cannot handle ", reflect.TypeOf(decl))
		case *ast.FuncDecl:
			w.classMembers = append(w.classMembers, n)
		}
	}
}

// code gen

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

	w.putln(PACKAGE, w.pkg, EOL)

	className := path.Base(w.fname)
	w.open(PUBLIC, FINAL, CLASS, className)
	w.putln()

	w.genMembers()

	w.close()
}

func (w *writer) genMembers() {
	for _, n := range w.classMembers {
		switch m := n.(type) {
		default:
			panic("unhandled memeber type: " + reflect.TypeOf(m).String())
		case *ast.FuncDecl:
			w.genFuncDecl(m)
		}
	}
}

func (w *writer) genFuncDecl(n *ast.FuncDecl) {
	if n.Name.Name == "main" {
		w.genMainDecl(n)
		return
	}

	panic("todo")
}

func (w *writer) genMainDecl(n *ast.FuncDecl) {
	w.open(PUBLIC, STATIC, FINAL, VOID, n.Name.Name)
	w.close()
	w.putln()
}

func (w *writer) putln(tokens ...interface{}) {
	w.putIndent()
	w.put(append(tokens, "\n")...)
}

func (w *writer) open(tokens ...interface{}) {
	w.putln(append(tokens, LBRACE)...)
	w.indent++
}

func (w *writer) close() {
	w.indent--
	w.putln(RBRACE)
}

func (w *writer) put(tokens ...interface{}) {
	for i, t := range tokens {
		if i != 0 {
			fmt.Fprint(w.out, " ")
		}
		fmt.Fprint(w.out, t)
	}
}

func (w *writer) putIndent() {
	for i := 0; i < w.indent; i++ {
		fmt.Fprint(w.out, "\t")
	}
}

// exit with fatal error, print token position of node n and msg.
func (w *writer) error(n ast.Node, msg ...interface{}) {
	fatal(append([]interface{}{w.pos(n), ": "}, msg...)...)
}

// return position of node using this writer's fileset
func (w *writer) pos(n ast.Node) token.Position {
	return w.fset.Position(n.Pos())
}
