package gotojava

import (
	"go/ast"
	"go/token"
)

// Emit code for a binary op.
// 	binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
// 	rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
// 	add_op     = "+" | "-" | "|" | "^" .
// 	mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
func (w *Writer) PutBinaryExpr(b *ast.BinaryExpr) {
	if *flagParens {
		w.Put("(")
	}

	unsigned := IsUnsigned(TypeOf(b.X)) || IsUnsigned(TypeOf(b.Y))

	x := RValue(b.X)
	y := RValue(b.Y)

	switch b.Op {
	default:
		w.Put(x, b.Op.String(), y)
	case token.EQL:
		w.PutJEquals(JTypeOfExpr(b.X), x, JTypeOfExpr(b.Y), y)
	case token.LSS, token.GTR, token.LEQ, token.GEQ, token.QUO, token.REM:
		if unsigned {
			w.PutUnsignedOp(b.X, b.Op, b.Y)
		} else {
			w.Put(x, b.Op.String(), y) // default
		}
	case token.SHL, token.SHR, token.AND, token.OR, token.XOR:
		// different precedence in Go and Java, parentisize to be sure
		w.Put("(", x, b.Op.String(), y, ")")
	case token.AND_NOT: //
		// not in java
		w.Put("(", x, "&~", y, ")")
	}

	if *flagParens {
		w.Put(")")
	}
}

// Emit code for Go's "lhs == rhs", with given java types for both sides.
// May emit, e.g.:
// 	lhs == rhs      // basic and pointer types
// 	lhs.equals(rhs) // struct value types
func (w *Writer) PutJEquals(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " == ", rhs) // TODO: panic
	case ltyp.NeedsEqualsMethod():
		w.Put(lhs, ".equals(", rhs, ")")
	}
}
