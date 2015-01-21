package main

import (
	"go/ast"
	"reflect"
)

func (w *writer) putDecl(d ast.Decl) {
	switch d := d.(type) {
	default:
		panic("unhandled memeber type: " + reflect.TypeOf(d).String())
	case *ast.FuncDecl:
		w.putFuncDecl(d)
	}
}

func (w *writer) putFuncDecl(n *ast.FuncDecl) {
	w.putDoc(n.Doc)
	if n.Name.Name == "main" {
		w.putMainDecl(n)
		return
	}

	panic("todo")
}
