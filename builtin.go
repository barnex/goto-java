package main

import "go/ast"

var builtin = map[string]bool{
	"len": true,
}

func isBuiltinExpr(n ast.Expr) bool {
	if ident, ok := n.(*ast.Ident); ok {
		return builtin[ident.Name]
	}
	return false
}

func (w *writer) putBuiltinCall(n *ast.CallExpr) {
	name := n.Fun.(*ast.Ident).Name
	switch name {
	default:
		w.error(n, "cannot handle builtin: ", name)
	case "len":
		w.putLenExpr(n)
	}
}

func (w *writer) putLenExpr(n *ast.CallExpr) {
	if len(n.Args) != 1 {
		w.error(n, "too many arguments to len")
	}
	argT := w.typeConv(w.typeOf(n.Args[0]))
	switch argT {
	default:
		w.error(n, "invalid argument (type ", argT, ")for len")
	}
}
