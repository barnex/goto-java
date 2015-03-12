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
		//Error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
		assert(len(n.Rhs) == 1)
		rhs := n.Rhs[0]

		id := "bla" // TODO
		w.Putln(JTypeOfExpr(rhs), " ", id, "=", rhs)
		for i, lhs := range n.Lhs {
			w.PutSemi(i)
			lhs = StripParens(lhs)
			w.Putln(lhs, n.Tok, id, ".v", i)
			// TODO: blank
		}
	} else {

		// multiple assign: put one per line
		for i, lhs := range n.Lhs {
			w.PutSemi(i)

			lhs = StripParens(lhs) // border case, go allows "(x) = y" //?
			rhs := n.Rhs[i]

			// blank identifer: it's actually a define. E.g.: int _4 = f(x);
			if IsBlank(lhs) {
				w.PutJVarDecl(NONE, JTypeOfExpr(rhs), lhs.(*ast.Ident), rhs, nil)
			} else {
				w.PutAssign(lhs, n.Tok, rhs)
			}
		}
	}
}

func (w *Writer) PutAssign(lhs ast.Expr, op token.Token, rhs ast.Expr) {
	lhs = StripParens(lhs)

	if NeedsSetMethod(lhs) {
		w.putAssignMethod(lhs, op, rhs)
	} else {
		w.putAssignOp(lhs, op, rhs)
	}
}

func NeedsSetMethod(lvalue ast.Expr) bool {
	switch lvalue := lvalue.(type) {
	default:
		return true
		//panic("unsupported lvalue: " + reflect.TypeOf(lvalue).String())
	case *ast.Ident:
		return JTypeOfExpr(lvalue).NeedsMethods()
	case *ast.SelectorExpr:
		return NeedsSetMethod(lvalue.Sel)
	}
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

func (w *Writer) putAssignMethod(lhs ast.Expr, tok token.Token, rhs ast.Expr) {
	if meth := opToMeth[tok]; meth != "" {
		w.Put(LValue(lhs), ".", meth, " (", RValue(rhs), ")")
	} else {
		panic(tok)
	}
}

var opToMeth = map[token.Token]string{
	token.ASSIGN:         "set",    // =
	token.INC:            "inc",    // ++
	token.DEC:            "dec",    // --
	token.ADD_ASSIGN:     "add",    // +=
	token.SUB_ASSIGN:     "sub",    // -=
	token.MUL_ASSIGN:     "mul",    // *=
	token.QUO_ASSIGN:     "quo",    // /=
	token.REM_ASSIGN:     "rem",    // %=
	token.AND_ASSIGN:     "and",    // &=
	token.OR_ASSIGN:      "or",     // |=
	token.XOR_ASSIGN:     "xor",    // ^=
	token.SHL_ASSIGN:     "shl",    // <<=
	token.SHR_ASSIGN:     "shr",    // >>=
	token.AND_NOT_ASSIGN: "andnot", // &^=
}

//func (w *Writer) putSelectorAssign(lhs *ast.SelectorExpr, op token.Token, rhs ast.Expr) {
//	w.Put(lhs.X, ".")
//	w.putIdentAssign(lhs.Sel, op, rhs)
//}

//func (w *Writer) putStarAssign(lhs *ast.StarExpr, rhs ast.Expr) {
//	switch elem := lhs.X.(type) {
//	default:
//		panic("not supported: assign to *(" + reflect.TypeOf(elem).String() + ")")
//	case *ast.Ident:
//		w.Put(elem, ".set(", RValue(rhs), ")")
//	}
//}

// Emit code for rhs, possibly converting to make it assignable to lhs.
//func (w *Writer) PutAutoCast(rhs ast.Expr, lhs JType, inmethod bool) {
//	w.PutExpr(rhs)
//}

func RValue(rhs ast.Expr) interface{} {
	if rhs == nil {
		return ""
	}

	if !JTypeOfExpr(rhs).IsPrimitive() {
		return Transpile(rhs, ".value()")
	}

	return rhs
	// TODO: cast
}

func LValue(lhs ast.Expr) interface{} {
	lhs = StripParens(lhs)
	switch lhs := lhs.(type) {
	default:
		panic("cannot make LValue: " + reflect.TypeOf(lhs).String())
	case *ast.Ident:
		return lhs
	case *ast.StarExpr:
		return lhs.X
	case *ast.SelectorExpr:
		return Transpile(lhs.X, ".", lhs.Sel)
	}
}

// Emit code for Go's "lhs = rhs", with given java types for both sides.
// May emit, e.g.:
// 	a = b    // basic and pointer types
// 	a.set(b) // struct values
func (w *Writer) PutJAssign(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " = ", rhs)
	case ltyp.NeedsMethods():
		w.Put(lhs, ".set(", rhs, ")")
	}
}
