package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

func (w *writer) IsUnsigned(t types.Type) bool {
	if b, ok := t.(*types.Basic); ok {
		return b.Info()&types.IsUnsigned != 0
	} else {
		return false
	}
}

func (w *writer) PutUnsignedOp(x ast.Expr, op token.Token, y ast.Expr) {
	switch op {
	default:
		w.error(x, "unsigned", op.String(), "not supported")
	case token.LSS:
		w.Put("(", "Integer.compareUnsigned(", x, ",", y, ")<0", ")")
	case token.GTR:
		w.Put("(", "Integer.compareUnsigned(", x, ",", y, ")>0", ")")
	case token.LEQ:
		w.Put("(", "Integer.compareUnsigned(", x, ",", y, ")<=0", ")")
	case token.GEQ:
		w.Put("(", "Integer.compareUnsigned(", x, ",", y, ")>=0", ")")
	}
}
