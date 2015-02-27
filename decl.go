package gotojava

// This file handles declarations

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Emit a declaration with optional modifier (like static)
func (w *Writer) PutDecl(mod JModifier, d ast.Decl) {
	switch d := d.(type) {
	default:
		panic("unhandled memeber type: " + reflect.TypeOf(d).String())
	case *ast.FuncDecl:
		w.PutFuncDecl(mod, d)
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
func (w *Writer) PutFuncDecl(mod JModifier, n *ast.FuncDecl) {
	if n.Recv == nil {
		// main is special: need String[] args
		if n.Name.Name == "main" {
			w.PutMainDecl(n)
		} else {
			w.PutFunc(mod, n)
		}
	} else {
		// ignore method, handled by CollectDefs.
	}
}

// Emit the main function. Special case in PutFuncDecl.
func (w *Writer) PutMainDecl(n *ast.FuncDecl) {
	w.Put("public static void ", n.Name.Name, "(String[] args)")
	w.PutBlockStmt(n.Body)
	w.Putln()
}

// Emit code for a top-level function (not method) declaration, e.g.:
// 	func f(a, b int) { ... }
func (w *Writer) PutFunc(mod JModifier, f *ast.FuncDecl) {
	w.PutDoc(f.Doc)

	// modifier
	w.Put(mod | GlobalModifierFor(f.Name))

	// return type
	retNames, retTypes := FlattenFields(f.Type.Results)
	w.Put(JavaReturnTypeOf(retTypes), " ", f.Name)

	// arguments
	argNames, argTypes := FlattenFields(f.Type.Params)
	w.Put("(")

	// receiver, if any, is first argument
	if f.Recv != nil {
		assert(len(f.Recv.List) == 1)
		recv := f.Recv.List[0]
		name := ""
		if recv.Names != nil {
			assert(len(recv.Names) == 1)
			name = JavaName(recv.Names[0])
		} else {
			name = makeNewName(UNUSED)
		}
		w.Put(JTypeOf(recv.Type), " ", name)
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
			w.Putln(retTypes[i], " ", retNames[i], ";")
		}
	}

	// rest of body
	w.PutStmtList(f.Body.List)

	w.indent--
	w.Putln("}")
}

// Emit a parameter list, e.g.:
// 	(type1 name1, type2 name2)
// Types  and names typically obtained by FlattenFields().
func (w *Writer) PutParams(names []*ast.Ident, types []JType) {
	for i := range names {
		w.Put(comma(i), types[i], " ", names[i])
	}
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
	switch d.Tok {
	default:
		Error(d, "cannot handle "+d.Tok.String())
	case token.TYPE:
		// do nothing. already handled by CollectDefs
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
func (w *Writer) PutValueSpecs(mod JModifier, specs []ast.Spec) {
	for i, spec := range specs {
		w.PutSemi(i)
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
	// One per line
	// TODO: what if spec.Type != JTypeOf(id) e.g. escaped basic
	for i, id := range s.Names {
		w.PutSemi(i)
		var value ast.Expr = nil
		if i < len(s.Values) {
			value = s.Values[i]
		}
		w.PutJVarDecl(mod, JTypeOf(id), id, value, s.Comment)
	}
}

// Emit a short variable declaration, e.g.:
// 	a := 1
func (w *Writer) putShortDefine(mod JModifier, a *ast.AssignStmt) {
	if len(a.Lhs) != len(a.Rhs) {
		Error(a, "assignment count mismatch:", len(a.Lhs), "!=", len(a.Rhs))
		// TODO: function with multiple returns
	}
	for i, lhs := range a.Lhs {
		w.PutSemi(i)

		var rhs ast.Expr = nil
		if i < len(a.Rhs) {
			rhs = a.Rhs[i]
		}

		id := lhs.(*ast.Ident)
		if isShortRedefine(id) {
			w.PutJAssign(JTypeOf(id), id, JTypeOf(rhs), RValue(rhs))
		} else {
			w.PutJVarDecl(mod, JTypeOf(id), id, rhs, nil)
		}
	}
}

// Put a value java variable declaration:
// 	modifier type ident = value;
// value may be nil.
func (w *Writer) PutJVarDecl(mod JModifier, jType JType, id *ast.Ident, value ast.Expr, comment *ast.CommentGroup) {
	if jType.NeedsFinal() {
		mod |= FINAL
	}
	w.Put(mod, jType, " ", id, " = ")
	if value != nil {
		w.Put(InitValue(value, jType))
	} else {
		w.Put(ZeroValue(jType))
	}
}

// Is the identifier already defined its scope?
// Detects redeclaration in a short variable declaration, e.g.:
// 	a := 1
// 	a, b := 2, 3  // IsShortRedefine(a): true
// See: https://golang.org/doc/effective_go.html#redeclaration
func isShortRedefine(id *ast.Ident) bool {
	if IsBlank(id) {
		return false // blank identifier _ is never redefined
	}
	obj := ObjectOf(id)
	pos := id.Pos()
	scope := obj.Parent()
	names := scope.Names()
	// TODO: names is sorted, could binary search
	for _, n := range names {
		if n == id.Name && pos > scope.Lookup(n).Pos() {
			return true
		}
	}
	return false
}
