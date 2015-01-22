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

// Emit code for a top-level function/method declaration, e.g.:
// 	func f(a, b int) { ... }
// 	func (x *T) f() { ... }
func (w *writer) PutFuncDecl(n *ast.FuncDecl) {

	if n.Recv == nil {
		w.PutStaticFunc(n)
	} else {
		w.PutMethod(n)
	}
}

// Emit code for a top-level function (not method) declaration, e.g.:
// 	func f(a, b int) { ... }
func (w *writer) PutStaticFunc(f *ast.FuncDecl) {
	w.PutDoc(f.Doc)

	// main is special: need String[] args
	if f.Name.Name == "main" {
		w.PutMainDecl(f)
		return
	}

	w.Put(ModifierFor(f.Name.Name), "static ")

	ret := "void"
	if len(f.Type.Results.List) == 1 {
		ret = w.javaTypeOf(f.Type.Results.List[0].Type) // todo: multiple names, wtf?
	}
	if len(f.Type.Results.List) > 1 {
		w.error(f, "no muliple return values supported")
	}

	w.Put(ret, " ", (f.Name.Name), "(") // TODO: translate funcname
	for i, a := range f.Type.Params.List {
		if len(a.Names) != 1 {
			w.error(f, "cannot handle multiple field names")
		}
		name := a.Names[0] // TODO: more/none?
		w.Put(comma(i), a.Type, " ", name)
	}
	w.Put(")")
	w.Putln(f.Body)
}

//func (w *writer) PutField(f *ast.Field) {
//	w.Put(f.Type, " ")
//	for i, n := range f.Names {
//		w.Put(comma(i), n)
//	}
//}

// Emit code for a method declaration, e.g.:
// 	func (x *T) f() { ... }
func (w *writer) PutMethod(n *ast.FuncDecl) {
	panic("todo: method")
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
			if i != 0 {
				w.Putln(";")
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
	//w.Put(";")
	//w.PutInlineComment(comment)
}
