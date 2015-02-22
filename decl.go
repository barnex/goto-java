package gotojava

// This file handles declarations

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/types"
)

// Emit a declaration with optional modifier (like static)
func (w *Writer) PutDecl(mod JModifier, d ast.Decl) {
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
func (w *Writer) PutFuncDecl(n *ast.FuncDecl) {
	if n.Recv == nil {
		w.PutStaticFunc(n)
	} else {
		// ignore method, handled by CollectDefs.
	}
}

// Java return type for a function that returns given types.
// For multiple return types, a Tuple type is returned
func JavaReturnTypeOf(resultTypes []types.Type) string {
	switch len(resultTypes) {
	case 0:
		return "void"
	case 1:
		return JavaTypeOf(resultTypes[0])
	default:
		return JavaTupleType(resultTypes)
	}
}

// Emit code for a top-level function (not method) declaration, e.g.:
// 	func f(a, b int) { ... }
func (w *Writer) PutStaticFunc(f *ast.FuncDecl) {
	w.PutDoc(f.Doc)

	// main is special: need String[] args
	if f.Name.Name == "main" {
		w.PutMainDecl(f)
		return
	}

	// modifier
	mod := STATIC | ModifierFor(f.Name)
	w.Put(mod, " ")

	// return type
	retNames, retTypes := FlattenFields(f.Type.Results)
	w.Put(JavaReturnTypeOf(retTypes), " ", f.Name)

	// arguments
	w.Put("(")

	argNames, argTypes := FlattenFields(f.Type.Params)

	// receiver, if any, is first argument
	if f.Recv != nil {
		assert(len(f.Recv.List) == 1)
		field := f.Recv.List[0]
		name := makeNewName(UNUSED)
		if field.Names != nil {
			assert(len(field.Names) == 1)
			name = JavaNameFor(field.Names[0])
		}

		w.PutTypeExpr(field.Type)

		w.Put(" ", name)
		if len(argNames) != 0 {
			w.Put(", ")
		}
	}

	w.PutParams(argNames, argTypes)
	w.Putln("){")
	w.indent++

	// declare named return values, if any
	for i := range retNames {
		if retNames[i] != nil {
			w.Putln(JavaTypeOf(retTypes[i]), " ", retNames[i], " = ", ZeroValue(retTypes[i]), ";")
		}
	}

	w.PutStmtList(f.Body.List)

	w.indent--
	w.Putln("}")
}

func (w *Writer) PutParams(names []*ast.Ident, types []types.Type) {
	for i := range names {
		w.Put(comma(i), JavaTypeOf(types[i]), " ", names[i])
	}
}

// Emit the main function. Special case in PutStaticFunc.
func (w *Writer) PutMainDecl(n *ast.FuncDecl) {
	w.Put("public static void ", n.Name.Name, "(String[] args)")
	w.PutBlockStmt(n.Body)
	w.Putln()
}

// Emit a generic declaration (import, constant, type or variable)
// with optional modifier
// godoc:
// 	type GenDecl struct {
// 	    Doc    *CommentGroup // associated documentation; or nil
// 	    TokPos token.Pos     // position of Tok
// 	    Tok    token.Token   // IMPORT, CONST, TYPE, VAR
// 	    Lparen token.Pos     // position of '(', if any
// 	    Specs  []Spec
// 	    Rparen token.Pos     // position of ')', if any
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
func (w *Writer) PutGenDecl(mod JModifier, d *ast.GenDecl) {
	switch d.Tok { // IMPORT, CONST, TYPE, VAR
	default:
		Error(d, "cannot handle "+d.Tok.String())
	case token.CONST:
		w.PutValueSpecs(mod|FINAL, d.Specs)
	case token.TYPE:
		// do nothing. already handled by CollectDefs
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
func (w *Writer) PutValueSpecs(mod JModifier, specs []ast.Spec) {
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
func (w *Writer) PutValueSpec(mod JModifier, s *ast.ValueSpec) {
	if s.Type != nil {
		// var with explicit type: everything on one line, e.g.:
		// 	int a = 1, b = 2
		w.PutValueSpecLine(mod, TypeOf(s.Type), s.Names, s.Values, s.Comment)
	} else {
		// var with infered type: specs on separate line, e.g.:
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
			w.PutValueSpecLine(mod, TypeOf(n), s.Names[i:i+1], []ast.Expr{value}, s.Comment)
		}
	}
}

// Put a value spec where all variables have the same, explicit, type, e.g.:
// 	var x, y int = 1, 2
// Translates to java:
// 	int x = 1, y = 2
// Type may be nil to allow short declarations with an existing variable. e.g.:
// 	a := 1
// 	a, b := 2, 3
// becomes:
// 	int a = 1;
// 	a = 2;       // typ = nil
//  int b = 3;
func (w *Writer) PutValueSpecLine(mod JModifier, typ types.Type, names []*ast.Ident, values []ast.Expr, comment *ast.CommentGroup) {

	w.Put(mod)
	if mod != NONE {
		w.Put(" ")
	}

	if typ != nil {
		jType := JavaTypeOf(typ)
		w.Put(jType)
	}

	for i, n := range names {

		w.Put(" ", n, " = ")

		if i < len(values) {
			w.PutExpr(values[i])
		} else {
			w.Put(ZeroValue(TypeOf(n)))
		}

		if i != len(names)-1 {
			w.Put(", ")
		}
	}
}
