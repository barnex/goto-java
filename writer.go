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
	out        io.Writer
	indent     int
	needIndent bool
	needSpace  bool

	fset *token.FileSet

	fname        string // output file name for principal class (no extension)
	pkg          string
	classMembers []ast.Decl
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
	w.putln()

	className := path.Base(w.fname)
	w.putln(PUBLIC, FINAL, CLASS, className, LBRACE)
	w.putln()
	w.indent++

	w.genMembers()

	w.indent--
	w.putln(RBRACE)
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
	w.put(PUBLIC, STATIC, VOID, n.Name.Name, parens("String[] args"))
	w.genBlockStmt(n.Body)
	w.putln()
}

func (w *writer) genBlockStmt(n *ast.BlockStmt) {
	w.putln(LBRACE)
	w.indent++

	w.indent--
	w.putln(RBRACE)
}

// fmt utils

func (w *writer) putln(tokens ...interface{}) {
	w.put(append(tokens, "\n")...)
	w.needIndent = true
	w.needSpace = false
}

func (w *writer) put(tokens ...interface{}) {
	w.putIndent()
	for _, t := range tokens {
		if t == ";" || t == "\n" {
			w.needSpace = false
		}
		if w.needSpace {
			fmt.Fprint(w.out, " ")
		}
		fmt.Fprint(w.out, t)
		w.needSpace = true
	}
}

func (w *writer) putIndent() {
	if w.needIndent == false {
		return
	}
	for i := 0; i < w.indent; i++ {
		fmt.Fprint(w.out, "\t")
	}
	w.needIndent = false
}

// exit with fatal error, print token position of node n and msg.
func (w *writer) error(n ast.Node, msg ...interface{}) {
	fatal(append([]interface{}{w.pos(n), ": "}, msg...)...)
}

// utils

// return position of node using this writer's fileset
func (w *writer) pos(n ast.Node) token.Position {
	return w.fset.Position(n.Pos())
}

// prints parentisized argument list: "(elem[0], elem[1], ...)"
func parens(elem ...interface{}) string {
	str := "("
	for i, e := range elem {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprint(e)
	}
	str += ")"
	return str
}
