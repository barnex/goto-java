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
	case token.LSS, token.GTR, token.LEQ, token.GEQ:
		w.Put("((", x, op.String(), y, ") ^ ((", x, "< 0) != (", y, " < 0)))")
	}
}

// http://www.javamex.com/java_equivalents/unsigned_arithmetic.shtml
