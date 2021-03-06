package gotojava

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

func IsUnsigned(t types.Type) bool {
	if b, ok := t.(*types.Basic); ok {
		return b.Info()&types.IsUnsigned != 0
	} else {
		return false
	}
}

func (w *Writer) PutUnsignedOp(x ast.Expr, op token.Token, y ast.Expr) {
	typX := JTypeOfExpr(x)
	typY := JTypeOfExpr(y)

	if typX != typY {
		Error(x, "mismatched types", typX, "and", typY)
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
		"short": "16",
		"int":   "32",
		"long":  "64",
	}[typX.JName()]
	function := "go.Unsigned." + operator + suffix

	switch op {
	default:
		Error(x, "unsigned", op.String(), "not supported")
	case token.QUO, token.REM, token.LSS, token.GTR, token.LEQ, token.GEQ:
		w.Put(function, "(", RValue(x), ", ", RValue(y), ")")
	}
}
