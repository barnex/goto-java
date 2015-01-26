package main

// This file handles declarations

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/types"
)

// Emit a declaration with optional modifier (like static)
func (w *writer) PutDecl(mod JModifier, d ast.Decl) {
	switch d := d.(type) {
	default:
		panic("unhandled memeber type: " + reflect.TypeOf(d).String())
	case *ast.FuncDecl:
		w.PutFuncDecl(d)
	case *ast.GenDecl:
		w.PutGenDecl(mod, d)
	}
}

// Emit code for a top-level function/method declaration, e.g.:
// 	func f(a, b int) { ... }
// 	func (x *T) f() { ... }
// ast.FuncDecl godoc:
// 	type FuncDecl struct {
// 	    Doc  *CommentGroup // associated documentation; or nil
// 	    Recv *FieldList    // receiver (methods); or nil (functions)
// 	    Name *Ident        // function/method name
// 	    Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	    Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
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

	mod := STATIC
	if ast.IsExported(f.Name.Name) {
		mod |= PUBLIC
	}

	ret := "void"
	if f.Type.Results != nil && len(f.Type.Results.List) == 1 {
		ret = w.TypeToJava(w.TypeOf(f.Type.Results.List[0].Type)) // todo: multiple names, wtf?
	}
	if f.Type.Results != nil && len(f.Type.Results.List) > 1 {
		w.error(f, "no muliple return values supported")
	}

	w.Put(mod, " ", ret, " ", (f.Name.Name), "(") // TODO: translate funcname

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

// Emit the main function. Special case in PutStaticFunc.
func (w *writer) PutMainDecl(n *ast.FuncDecl) {
	w.Put("public static void ", n.Name.Name, "(String[] args)")
	w.PutBlockStmt(n.Body)
	w.Putln()
}

// Emit code for a method declaration, e.g.:
// 	func (x *T) f() { ... }
func (w *writer) PutMethod(n *ast.FuncDecl) {
	panic("todo: method")
}

//func (w *writer) PutField(f *ast.Field) {
//	w.Put(f.Type, " ")
//	for i, n := range f.Names {
//		w.Put(comma(i), n)
//	}
//}

// Emit a generic declaration (import, constant, type or variable)
// with optional modifier
// godoc:
// 	type GenDecl struct {
// 	    Doc    *CommentGroup // associated documentation; or nil
// 	    TokPos token.Pos     // position of Tok
// 	    Tok    token.Token   // IMPORT, CONST, TYPE, VAR
// 	    Lparen token.Pos     // position of '(', if any
// 	    Specs  []Spec
// 	    Rparen token.Pos // position of ')', if any
// 	}
// A GenDecl node (generic declaration node) represents an import,
// constant, type or variable declaration. A valid Lparen position
// (Lparen.Line > 0) indicates a parenthesized declaration.
//
// Relationship between Tok value and Specs element type:
//
// 	token.IMPORT  *ImportSpec
// 	token.CONST   *ValueSpec
// 	token.TYPE    *TypeSpec
// 	token.VAR     *ValueSpec
//
func (w *writer) PutGenDecl(mod JModifier, d *ast.GenDecl) {
	switch d.Tok { // IMPORT, CONST, TYPE, VAR
	default:
		w.error(d, "cannot handle "+d.Tok.String())
	case token.CONST:
		w.PutValueSpecs(mod|FINAL, d.Specs)
	case token.VAR:
		w.PutValueSpecs(mod, d.Specs)
	}
}

// Emit a list of *ast.ValueSpecs, e.g.:
// 	var(
// 		a int
// 		b, c int
// 	)
// or
// 	const(
// 		a = 1
// 		b, c = 2, 3
// 	)
// with optional modifier (e.g. "static", "final", "static final").
// The concrete type of specs elements must be *ast.ValueSpec.
func (w *writer) PutValueSpecs(mod JModifier, specs []ast.Spec) {
	for i, spec := range specs {
		if i != 0 {
			w.Putln(";")
		}
		w.PutValueSpec(mod, spec.(*ast.ValueSpec)) // doc says it's a valueSpec for Tok == VAR, CONST
	}
}

// Emit a single ValueSpec, e.g.:
// 	var a, b int
// or
// 	const a, b = 0, "hello"
// with optional modifier prefix (e.g. "static", "final", "static final").
//
// ValueSpec godoc:
// 	type ValueSpec struct {
// 	    Doc     *CommentGroup // associated documentation; or nil
// 	    Names   []*Ident      // value names (len(Names) > 0)
// 	    Type    Expr          // value type; or nil
// 	    Values  []Expr        // initial values; or nil
// 	    Comment *CommentGroup // line comments; or nil
// 	}
// A ValueSpec node represents a constant or variable declaration
// (ConstSpec or VarSpec production).
func (w *writer) PutValueSpec(mod JModifier, s *ast.ValueSpec) {
	if s.Type != nil {
		// var with explicit type:
		// Put everything on one line, e.g.:
		// 	int a = 1, b = 2
		w.PutValueSpecLine(mod, w.TypeOf(s.Type), s.Names, s.Values, s.Comment)
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
			w.PutValueSpecLine(mod, w.TypeOf(n), s.Names[i:i+1], []ast.Expr{value}, s.Comment)
		}
	}
}

// Put a value spec where all variables have the same, explicit, type, e.g.:
// 	var x, y int = 1, 2
// Translates to java:
// 	int x = 1, y = 2
func (w *writer) PutValueSpecLine(mod JModifier, typ types.Type, names []*ast.Ident, values []ast.Expr, comment *ast.CommentGroup) {

	w.Put(mod)
	if mod != NONE {
		w.Put(" ")
	}

	// untyped const hack: remove "untyped " from type
	//if mod.Is(FINAL) && strings.HasPrefix(typ, "untyped ") {
	//	typ = typ[len("untyped "):]
	//}

	w.Put(w.TypeToJava(typ))
	for i, n := range names {

		w.Put(" ", w.translate(n), " = ")

		if i < len(values) {
			w.PutExpr(values[i])
		} else {
			w.Put(w.ZeroValue(w.TypeOf(n)))
		}

		if i != len(names)-1 {
			w.Put(", ")
		}
	}
	//w.Put(";")
	//w.PutInlineComment(comment)
}
