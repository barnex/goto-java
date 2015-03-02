package gotojava

import (
	"go/ast"
	"go/token"
	"reflect"
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
	case token.DEFINE:
		w.putShortDefine(NONE, n)
	case token.ASSIGN:
		w.putMultiAssign(n)
	default:
		assert(len(lhs) == 1 && len(rhs) == 1)
		w.PutAssign(lhs[0], n.Tok, rhs[0])
	}
}

// Emit pure assign statement ('=' token)
func (w *Writer) putMultiAssign(n *ast.AssignStmt) {
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
			w.PutJVarDecl(NONE, JTypeOf(rhs), lhs.(*ast.Ident), rhs, nil)
			//	w.Put(JTypeOf(rhs), " ")
			//	w.PutAssign(JTypeOf(rhs), lhs, JTypeOf(rhs), RValue(rhs))
		} else {
			w.PutAssign(lhs, n.Tok, rhs)
		}
	}
}

func (w *Writer) PutAssign(lhs ast.Expr, op token.Token, rhs ast.Expr) {
	lhs = StripParens(lhs)

	switch lhs := lhs.(type) {
	default:
		panic("unsupported lvalue: " + reflect.TypeOf(lhs).String())
	case *ast.Ident:
		w.putIdentAssign(lhs, op, rhs)
	case *ast.StarExpr:
		//case *ast.SelectorExpr:
	}
}

// TODO: not only for idents
func (w *Writer) putIdentAssign(lhs *ast.Ident, op token.Token, rhs ast.Expr) {
	if JTypeOf(lhs).NeedsSetMethod() {
		rhs := RValue(rhs)
		if meth := opToMeth[op]; meth != "" {
			w.Put(lhs, ".", meth, " (", rhs, ")")
		} else {
			op2 := op.String()[:len(op.String())-1]
			if op == token.AND_NOT_ASSIGN {
				op2 = "&~"
			}
			w.Put(lhs, ".set(", RValue(lhs), " ", op2, " (", rhs, "))")
		}
	} else {
		w.putAssignOp(lhs, op, rhs) // e.g. "lhs = rhs"
	}
}

var opToMeth = map[token.Token]string{
	token.ASSIGN: "set", // =
	token.INC:    "inc", // ++
	token.DEC:    "dec", // --
}

// Emit assign statement with operation, e.g.:
// 	+= -= *= /= ++ ...
// For inc/dec stmt, rhs should be nil.
func (w *Writer) putAssignOp(lhs ast.Expr, tok token.Token, rhs ast.Expr) {
	if tok == token.AND_NOT_ASSIGN {
		w.Put(lhs, " &= ", " ~", "(", RValue(rhs), ")")
	} else {
		w.Put(lhs, tok, RValue(rhs))
	}
}

func (w *Writer) putSelectorAssign(lhs *ast.SelectorExpr, op token.Token, rhs ast.Expr) {
	w.Put(lhs.X, ".")
	w.putIdentAssign(lhs.Sel, op, rhs)
}

func (w *Writer) putStarAssign(lhs *ast.StarExpr, rhs ast.Expr) {
	switch elem := lhs.X.(type) {
	default:
		panic("not supported: assign to *(" + reflect.TypeOf(elem).String() + ")")
	case *ast.Ident:
		w.Put(elem, ".set(", RValue(rhs), ")")
	}
}

// Emit code for rhs, possibly converting to make it assignable to lhs.
//func (w *Writer) PutAutoCast(rhs ast.Expr, lhs JType, inmethod bool) {
//	w.PutExpr(rhs)
//}

func RValue(rhs ast.Expr) interface{} {
	if rhs == nil {
		return ""
	}

	if JTypeOf(rhs).IsEscapedPrimitive() {
		return Transpile(rhs, ".value")
	}

	return rhs
	// TODO: cast
}

// Emit code for Go's "lhs = rhs", with given java types for both sides.
// May emit, e.g.:
// 	a = b    // basic and pointer types
// 	a.set(b) // struct values
func (w *Writer) PutJAssign(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " = ", rhs)
	case ltyp.NeedsSetMethod():
		w.Put(lhs, ".set(", rhs, ")")
	}
}
