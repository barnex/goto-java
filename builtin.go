package main

import "go/ast"

// list of Go builtin identifiers
var builtin = map[string]bool{
	"len": true,
}

// returns true if n is a Go built-in identifier
func IsBuiltinExpr(n ast.Expr) bool {
	if ident, ok := n.(*ast.Ident); ok {
		return builtin[ident.Name]
	}
	return false
}

// Generate code for built-in call, like len(x)
func (w *writer) PutBuiltinCall(n *ast.CallExpr) {
	name := n.Fun.(*ast.Ident).Name
	switch name {
	default:
		w.error(n, "cannot handle builtin: ", name)
	case "len":
		w.PutLenExpr(n)
	}
}

// Generate code for len(x)
func (w *writer) PutLenExpr(n *ast.CallExpr) {
	if len(n.Args) != 1 {
		w.error(n, "too many arguments to len")
	}
	argT := w.javaTypeOf(n.Args[0])
	switch argT {
	default:
		w.error(n, "invalid argument (type ", argT, ") for len")
	case "String":
		w.Put("(")
		w.PutExpr(n.Args[0])
		w.Put(").length()")
	}
}
