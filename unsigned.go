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
	typX := w.TypeToJava(w.TypeOf(x))
	typY := w.TypeToJava(w.TypeOf(y))

	if typX != typY {
		w.error(x, "mismatched types", typX, "and", typY)
	}

	operator := map[token.Token]string{
		token.QUO: "quo",
		token.REM: "rem",
		token.LSS: "lss",
		token.GTR: "gtr",
		token.LEQ: "leq",
		token.GEQ: "geq",
	}[op]
	suffix := map[string]string{
		"byte":  "8",
		"short": "32",
		"int":   "32",
		"long":  "64",
	}[typX]
	function := "go.Unsigned." + operator + suffix

	switch op {
	default:
		w.error(x, "unsigned", op.String(), "not supported")
	case token.QUO, token.REM, token.LSS, token.GTR, token.LEQ, token.GEQ:
		w.Put(function+"(", x, ",", y, ")")
	}
}
