package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"

	"golang.org/x/tools/go/types"
)

// A writer takes AST nodes and outputs java source
type writer struct {
	out        io.Writer
	indent     int
	needIndent bool
	fset       *token.FileSet
	types.Info
}

func NewWriter(out io.Writer) *writer {
	return &writer{out: out}
}

// outputs a class with given name based on go file.
func (w *writer) putClass(className string, f *ast.File) {
	w.putln("package ", f.Name.Name, ";")
	w.putln()

	w.putln("public final class ", className, " {")
	w.putln()
	w.indent++

	for _, d := range f.Decls {
		w.putDecl(d)
	}

	w.indent--
	w.putln("}")
}

func (w *writer) putln(tokens ...interface{}) {
	w.put(append(tokens, "\n")...)
	w.needIndent = true
}

func (w *writer) put(tokens ...interface{}) {
	w.putIndent()
	for _, t := range tokens {
		fmt.Fprint(w.out, t)
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
