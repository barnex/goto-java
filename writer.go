package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"

	"golang.org/x/tools/go/types"
)

// A writer takes AST nodes and outPuts java source
type writer struct {
	out        io.Writer
	indent     int
	needIndent bool
	fset       *token.FileSet
	info       types.Info
}

func NewWriter(out io.Writer) *writer {
	return &writer{out: out}
}

// Outputs a class with given name based on go file.
func (w *writer) PutClass(className string, f *ast.File) {
	w.Putln("package ", f.Name.Name, ";")
	w.Putln()

	w.Putln("public final class ", className, " {")
	w.Putln()
	w.indent++

	for _, d := range f.Decls {
		w.PutDecl(STATIC, d)

		switch d.(type) {
		default: // no semi
		case *ast.GenDecl:
			w.Putln(";")
		}
	}

	w.indent--
	w.Putln("}")
}

func (w *writer) Putln(tokens ...interface{}) {
	w.Put(append(tokens, "\n")...)
	w.needIndent = true
}

func (w *writer) Put(tokens ...interface{}) {
	w.putIndent()
	for _, t := range tokens {
		w.put(t)
	}
}

func (w *writer) put(t interface{}) {
	switch t.(type) {
	case string, JModifier, token.Token:
		fmt.Fprint(w.out, t)
		return
	}

	if t, ok := t.(ast.Node); ok {
		w.PutNode(t)
		return
	}
	panic("writer: cannot put type " + reflect.TypeOf(t).String())
}

func (w *writer) PutNode(n ast.Node) {
	if n, ok := n.(ast.Stmt); ok {
		w.PutStmt(n)
		return
	}

	switch n := n.(type) {
	default:
		panic("putnode: need to handle: " + reflect.TypeOf(n).String())
	case ast.Expr:
		w.PutExpr(n)
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
	panic(fmt.Sprint(append([]interface{}{w.pos(n), ": "}, msg...)...))
}

// return position of node using this writer's fileset
func (w *writer) pos(n ast.Node) token.Position {
	return w.fset.Position(n.Pos())
}

// prints parentisized argument list: "(elem[0], elem[1], ...)"
//func parens(elem ...interface{}) string {
//	str := "("
//	for i, e := range elem {
//		if i > 0 {
//			str += ", "
//		}
//		str += fmt.Sprint(e)
//	}
//	str += ")"
//	return str
//}
