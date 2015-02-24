package gotojava

import (
	"go/ast"
	"go/token"
)

// Emit an assignment or a short variable declaration.
// AssignStmt godoc:
// 	type AssignStmt struct {
// 	        Lhs    []Expr
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // assignment token, DEFINE
// 	        Rhs    []Expr
// 	}
func (w *Writer) PutAssignStmt(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		Error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
		// TODO: function with multiple returns
	}

	// java does not have &^=, translate
	if n.Tok == token.AND_NOT_ASSIGN {
		if len(n.Lhs) != 1 || len(n.Rhs) != 1 {
			// should have been caught by type checker.
			Error(n, n.Tok.String(), " requires single argument")
		}
		w.Put(n.Lhs[0], " &= ", " ~", "(", n.Rhs[0], ")") // TODO: implicit conv
		return
	}

	if n.Tok == token.DEFINE {
		w.putDefine(JModifier(NONE), n)
		return
	}

	// multiple assign: put one per line
	for i, lhs := range n.Lhs {
		if i != 0 {
			w.Putln(";")
		}
		// blank identifer: need to put type. E.g.:
		// 	int _4 = f(x);
		if IsBlank(lhs) {
			w.Put(JavaTypeOfExpr(n.Rhs[i]), " ")
			lhs = StripParens(lhs) // border case, go allows "(_) = x"
		}
		w.Put(lhs, " ", n.Tok, " ")
		w.PutImplicitCast(TypeOf(lhs), n.Rhs[i])
	}
}

// Emit a short variable declaration, e.g.:
// 	a := 1
// Special case of PutAssignStmt
func (w *Writer) putDefine(mod JModifier, a *ast.AssignStmt) {
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
