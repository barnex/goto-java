package main

// built-ins

import "go/ast"

// maps built-in Go identifiers to java
var builtinIdentMap = map[string]string{
	"println": "System.out.println",
	"print":   "System.out.print",
}

// Set of Go built-ins
var builtins = map[string]bool{
	"append":  true,
	"cap":     true,
	"close":   true,
	"complex": true,
	"copy":    true,
	"delete":  true,
	"false":   true,
	"imag":    true,
	"iota":    true,
	"len":     true,
	"make":    true,
	"new":     true,
	"nil":     true,
	"panic":   true,
	"print":   true,
	"println": true,
	"real":    true,
	"recover": true,
	"true":    true,

	"bool":       true,
	"byte":       true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"int8":       true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uint8":      true,
	"uintptr":    true,
}

// IsBuiltinIdent returns true if id refers to a Go built-in identifer.
// The resulut is scope-sensitive, as built-ins may be shadowed by
// other declarations (e.g. len := 7).
func (w *writer) IsBuiltinIdent(id *ast.Ident) bool {
	if _, ok := w.info.Defs[id]; ok {
		return false // resolved, so not built-in (though possibly shadowed)
	} else {
		return builtins[id.Name]
	}
}

// IsBuiltinExpr returns true if expression e refers to a built-in identifer. E.g.:
// 	print
// 	(print)
// 	...
func (w *writer) IsBuiltinExpr(e ast.Expr) bool {
	e = StripParens(e)
	// identifier
	if id, ok := e.(*ast.Ident); ok {
		return w.IsBuiltinIdent(id)
	}
	return false
}

// Strip parens from expression, if any. E.g.:
// 	((x)) -> x
func StripParens(e ast.Expr) ast.Expr {
	if par, ok := e.(*ast.ParenExpr); ok {
		return StripParens(par.X)
	} else {
		return e
	}
}

// Emit code for a built-in identifer
func (w *writer) PutBuiltinIdent(id *ast.Ident) {
	panic("todo")
}

// Generate code for built-in call, like len(x)
func (w *writer) PutBuiltinCall(n *ast.CallExpr) {
	name := StripParens(n.Fun).(*ast.Ident).Name
	switch name {
	default:
		w.error(n, "cannot handle builtin: ", name)
	case "len":
		w.PutLenExpr(n)
	case "print", "println":
		//panic("todo")
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
