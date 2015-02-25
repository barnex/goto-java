package gotojava

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

// Emit an assignment or a short variable declaration. Godoc:
// 	type AssignStmt struct {
// 	        Lhs    []Expr
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // assignment token, DEFINE
// 	        Rhs    []Expr
// 	}
func (w *Writer) PutAssignStmt(n *ast.AssignStmt) {
	lhs, rhs := n.Lhs, n.Rhs
	switch n.Tok {
	case token.ASSIGN:
		w.putAssign(n)
	case token.DEFINE:
		w.putShortDefine(NONE, n)
	default:
		assert(len(lhs) == 1 && len(rhs) == 1)
		w.putAssignOp(lhs[0], n.Tok, rhs[0])
	}
}

// Emit assign statement with operation, e.g.:
// 	+= -= *= /= ...
func (w *Writer) putAssignOp(lhs ast.Expr, tok token.Token, rhs ast.Expr) {
	if tok == token.AND_NOT_ASSIGN {
		w.Put(lhs, " &= ", " ~", "(")
		w.PutRHS(rhs, TypeOf(lhs), false)
		w.Put(")")
	} else {
		w.Put(lhs, tok)
		w.PutRHS(rhs, TypeOf(lhs), false)
	}
}

// Emit pure assign statement ('=' token)
func (w *Writer) putAssign(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		Error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
	}

	// multiple assign: put one per line
	for i, lhs := range n.Lhs {
		rhs := n.Rhs[i]
		if i != 0 {
			w.Putln(";")
		}
		// blank identifer: need to put type. E.g.:
		// 	int _4 = f(x);
		var typeOfLHS types.Type
		if IsBlank(lhs) {
			w.Put(JavaTypeOf(rhs), " ")
			lhs = StripParens(lhs) // border case, go allows "(_) = x"
			typeOfLHS = TypeOf(rhs)
		} else {
			typeOfLHS = TypeOf(lhs)
		}
		w.MakeAssign(JavaType(typeOfLHS), lhs, JavaType(TypeOf(rhs)), rhs)
	}
}

// TODO: other name
// TODO: cast RHS
func (w *Writer) MakeAssign(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " = ", rhs) // TODO: panic
	case ltyp.IsStructValue() && rtyp.IsStructValue():
		w.Put(lhs, ".set(", rhs, ")")
	}
}

// TODO: other name
// TODO: cast RHS
func (w *Writer) MakeEquals(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " == ", rhs) // TODO: panic
	case ltyp.IsStructValue() && rtyp.IsStructValue():
		w.Put(lhs, "._equals(", rhs, ")")
	}
}

// Emit a short variable declaration, e.g.:
// 	a := 1
func (w *Writer) putShortDefine(mod JModifier, a *ast.AssignStmt) {
	if len(a.Lhs) != len(a.Rhs) {
		Error(a, "assignment count mismatch:", len(a.Lhs), "!=", len(a.Rhs))
		// TODO: function with multiple returns
	}
	for i, n := range a.Lhs {
		var value ast.Expr = nil
		if i < len(a.Rhs) {
			value = a.Rhs[i]
		}
		if i != 0 {
			w.Putln(";")
		}
		id := a.Lhs[i].(*ast.Ident) // should be

		typ := TypeOf(n)
		if isShortRedefine(id) {
			typ = nil
		}

		w.PutValueSpecLine(mod, typ, []*ast.Ident{id}, []ast.Expr{value}, nil)
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

// Emit code for rhs, possibly converting to make it assignable to lhs.
func (w *Writer) PutRHS(rhs ast.Expr, lhs types.Type, inmethod bool) {
	w.PutExpr(rhs)
}
