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

	modifier := "private"
	if ast.IsExported(n.Name.Name) {
		modifier = "public"
	}
	w.Put(modifier, " static ")

	ret := "void"
	if len(n.Type.Results.List) == 1 {
		ret := w.javaTypeOf(n.Type.Results.List[0])
	}
	if len(n.Type.Results.List) > 1 {
		w.error(n, "no muliple return values supported")
	}

	w.Put(ret, "(")

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
	if s.Type != nil {
		// var with explicit type:
		// Put everything on one line, e.g.:
		// 	int a = 1, b = 2
		w.putSpecOneType(w.javaTypeOf(s.Type), s.Names, s.Values, s.Comment)
	} else {
		// var with infered type:
		// Put specs on separate line, e.g.:
		// 	int a = 1;
		// 	String b = "";
		for i, n := range s.Names {
			var value ast.Expr = nil
			if i < len(s.Values) {
				value = s.Values[i]
			}
			w.putSpecOneType(w.javaTypeOf(n), s.Names[i:i+1], []ast.Expr{value}, s.Comment)
		}
	}
}

// Put a value spec where all variables have the same, explicit, type, e.g.:
// 	var x, y int = 1, 2
// Translates to java:
// 	int x = 1, y = 2
func (w *writer) putSpecOneType(typ string, names []*ast.Ident, values []ast.Expr, comment *ast.CommentGroup) {
	w.Put(typ)
	for i, n := range names {
		w.Put(" ", n.Name, " = ")
		if i < len(values) {
			w.PutExpr(values[i])
		} else {
			w.Put(ZeroValue(w.javaTypeOf(n)))
		}

		if i != len(names)-1 {
			w.Put(", ")
		}
	}
	w.Put(";")
	w.PutInlineComment(comment)
	w.Putln()
}
