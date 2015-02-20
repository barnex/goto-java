package gotojava

// Handle Go built-ins.

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

// IsBuiltinIdent returns true if id refers to a Go built-in identifier.
// The result is scope-sensitive, as built-ins may be shadowed by
// other declarations (e.g. len := 7).
func IsBuiltinIdent(id *ast.Ident) bool {
	obj := ObjectOf(id)
	return (obj.Parent() == types.Universe) && (builtins[id.Name] == true)
}

// IsBuiltinExpr returns true if expression e refers to a built-in identifer. E.g.:
// 	print
// 	(print)
// 	...
func IsBuiltinExpr(e ast.Expr) bool {
	e = StripParens(e)
	// identifier
	if id, ok := e.(*ast.Ident); ok {
		return IsBuiltinIdent(id)
	}
	return false
}

// maps go built-in identifiers to Java
var lit2java = map[string]string{
	"false": "false",
	"nil":   "null", //? need to type!
	"true":  "true",
}

// Emit code for a built-in identifer
func (w *Writer) PutBuiltinIdent(id *ast.Ident) {
	if transl, ok := lit2java[id.Name]; ok {
		w.Put(transl)
	} else {
		Error(id, "built-in identifier not supported: ", id.Name)
	}
}

// Generate code for built-in call, like len(x)
func (w *Writer) PutBuiltinCall(c *ast.CallExpr) {
	name := StripParens(c.Fun).(*ast.Ident).Name
	switch name {
	default:
		Error(c, "cannot handle builtin: ", name)
	case "len":
		w.PutLenExpr(c)
	case "print", "println":
		w.PutBuiltinPrintCall(c)
		//case "byte", "uint8", "int8", "uint16", "int16", "uint32", "int32", "uint", "int", "uint64", "int64":
		//	if len(c.Args) != 1 {
		//		Error(c, "too many arguments to conversion to", name)
		//	}
		//	w.PutTypecast(name, c.Args[0])
	}
}

// Emit code for built-in print, prinln calls.
func (w *Writer) PutBuiltinPrintCall(c *ast.CallExpr) {
	name := StripParens(c.Fun).(*ast.Ident).Name
	switch name {
	default:
		Error(c, "bug: not a print/println call:", name)
	case "print":
		name = "System.out.print"
	case "println":
		name = "System.out.println"
	}
	w.Put(name)
	w.PutArgs(c.Args, c.Ellipsis)
}

// Generate code for len(x)
func (w *Writer) PutLenExpr(n *ast.CallExpr) {
	if len(n.Args) != 1 {
		Error(n, "too many arguments to len")
	}
	argT := JavaTypeOfExpr(n.Args[0])
	switch argT {
	default:
		Error(n, "invalid argument (type ", argT, ") for len")
	case "String":
		w.Put("(")
		w.PutExpr(n.Args[0])
		w.Put(").length()")
	}
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
