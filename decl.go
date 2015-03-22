package gotojava

// This file handles declarations

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

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
	w.Put(mod | GlobalModifierFor(f.Name.String(), JTypeOfExpr(f.Name)))

	// return type
	returnType := TypeOf(f.Name).(*types.Signature).Results() //TODO: putSignature?
	w.Put(JTypeOfGoType(returnType), " ", f.Name)

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
		w.Put(JTypeOfExpr(recv.Type), " ", name)
		if len(argNames) != 0 {
			w.Put(", ")
		}
	}

	w.PutParams(argNames, argTypes)
	w.Putln("){")
	w.indent++

	// declare named return values, if any
	retNames, retTypes := FlattenFields(f.Type.Results)
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
		w.PutJVarDecl(mod, JTypeOfExpr(id), id, value, s.Comment)
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
		rhs := a.Rhs[i]

		id := lhs.(*ast.Ident)
		if isShortRedefine(id) {
			w.PutJAssign(JTypeOfExpr(id), id, JTypeOfExpr(rhs), RValue(rhs))
		} else {
			w.PutJVarDecl(mod, JTypeOfExpr(id), id, rhs, nil)
		}
	}
}

// Put a value java variable declaration:
// 	modifier type ident = value;
// value may be nil.
func (w *Writer) PutJVarDecl(mod JModifier, jType JType, id interface{}, value ast.Expr, comment *ast.CommentGroup) {
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
