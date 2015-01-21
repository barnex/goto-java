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
	case *ast.GenDecl:
		w.putGenDecl(d)
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

func (w *writer) putMainDecl(n *ast.FuncDecl) {
	w.put("public static void ", n.Name.Name, "(String[] args)")
	w.putBlockStmt(n.Body)
	w.putln()
}

func (w *writer) putGenDecl(d *ast.GenDecl) {
	for _, s := range d.Specs {
		w.putSpec(s)
	}
}

func (w *writer) putSpec(s ast.Spec) {
	switch s := s.(type) {
	default:
		w.error(s, "cannot handle ", reflect.TypeOf(s))
	case *ast.ValueSpec:
		w.putValueSpec(s)
	}
}

func (w *writer) putValueSpec(s *ast.ValueSpec) {
	w.putExpr(s.Type)

	for i, n := range s.Names {
		w.put(" ", n.Name, " = ")
		if i < len(s.Values) {
			w.putExpr(s.Values[i])
		} else {
			w.put(n.Obj.Data)
		}

		if i != len(s.Names)-1 {
			w.put(", ")
		}
	}
	w.put(";")
	w.putComment(s.Comment)
	w.putln()
}
