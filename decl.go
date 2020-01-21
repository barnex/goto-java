package main

import (
	"go/ast"
	"reflect"
)

func (w *writer) PutDecl(d ast.Decl) {
	switch d := d.(type) {
	default:
		panic("unhandled memeber type: " + reflect.TypeOf(d).String())
	case *ast.FuncDecl:
		w.PutFuncDecl(d)
	case *ast.GenDecl:
		w.PutGenDecl(d)
	}
}

func (w *writer) PutFuncDecl(n *ast.FuncDecl) {
	w.PutDoc(n.Doc)
	if n.Name.Name == "main" {
		w.PutMainDecl(n)
		return
	}

	panic("todo")
}

func (w *writer) PutMainDecl(n *ast.FuncDecl) {
	w.Put("public static void ", n.Name.Name, "(String[] args)")
	w.PutBlockStmt(n.Body)
	w.Putln()
}

func (w *writer) PutGenDecl(d *ast.GenDecl) {
	for _, s := range d.Specs {
		w.PutSpec(s)
	}
}

func (w *writer) PutSpec(s ast.Spec) {
	switch s := s.(type) {
	default:
		w.error(s, "cannot handle ", reflect.TypeOf(s))
	case *ast.ValueSpec:
		w.PutValueSpec(s)
	}
}

func (w *writer) PutValueSpec(s *ast.ValueSpec) {
	// var with explicit type:
	// Put everything on one line, e.g.:
	// 	int a = 1, b = 2
	if s.Type != nil {
		w.Put(w.javaTypeOf(s.Type))
		for i, n := range s.Names {
			w.Put(" ", n.Name, " = ")
			if i < len(s.Values) {
				w.PutExpr(s.Values[i])
			} else {
				w.Put(n.Obj.Data) // TODO
			}

			if i != len(s.Names)-1 {
				w.Put(", ")
			}
		}
		w.Put(";")
		w.PutInlineComment(s.Comment)
		w.Putln()
	} else {
		// var with infered type:
		// Put specs on separate line, e.g.:
		// 	int a = 1;
		// 	String b = "";
		for i, n := range s.Names {
			w.Put(w.javaTypeOf(n))

			w.Put(" ", n.Name, " = ")
			if i < len(s.Values) {
				w.PutExpr(s.Values[i])
			} else {
				w.Put(n.Obj.Data)
			}
			w.Put(";")
			if i == 0 {
				w.PutInlineComment(s.Comment)
			}
			w.Putln()
		}
	}
}
