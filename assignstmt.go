package gotojava

import (
	"go/ast"
	"go/token"
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
		w.PutAutoCast(rhs, JTypeOf(lhs), false)
		w.Put(")")
	} else {
		w.Put(lhs, tok)
		w.PutAutoCast(rhs, JTypeOf(lhs), false)
	}
}

// Emit pure assign statement ('=' token)
func (w *Writer) putAssign(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		Error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
	}

	// multiple assign: put one per line
	for i, lhs := range n.Lhs {
		lhs = StripParens(lhs) // border case, go allows "(x) = y" //?
		rhs := n.Rhs[i]
		if i != 0 {
			w.Putln(";")
		}
		// blank identifer: need to put type. E.g.:
		// 	int _4 = f(x);
		if IsBlank(lhs) {
			w.Put(JTypeOf(rhs), " ")
			w.PutJAssign(JTypeOf(rhs), lhs, JTypeOf(rhs), RValue(rhs))
		} else {
			w.PutJAssign(JTypeOf(lhs), LValue(lhs), JTypeOf(rhs), RValue(rhs))
		}
	}
}

// Emit code for rhs, possibly converting to make it assignable to lhs.
func (w *Writer) PutAutoCast(rhs ast.Expr, lhs JType, inmethod bool) {
	w.PutExpr(rhs)
}

// Emit code for Go's "lhs = rhs", with given java types for both sides.
// May emit, e.g.:
// 	a = b    // basic and pointer types
// 	a.set(b) // struct values
func (w *Writer) PutJAssign(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " = ", rhs)
	case ltyp.IsStructValue():
		w.Put(lhs, ".set(", rhs, ")")
	}
}
