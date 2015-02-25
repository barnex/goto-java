package gotojava

// Handle Go built-ins.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// IsBuiltin returns true if expression e refers to a built-in identifer. E.g.:
// 	print
// 	(print)
// The result is scope-sensitive, as built-ins may be shadowed by
// other declarations (e.g. len := 7).
func IsBuiltin(e ast.Expr) bool {
	e = StripParens(e)
	if id, ok := e.(*ast.Ident); ok {
		obj := ObjectOf(id)
		_, ok := obj.(*types.Builtin)
		return ok
	}
	return false
}

// Emit code for a built-in identifer
func (w *Writer) PutBuiltinIdent(id *ast.Ident) {
	if transl, ok := builtin2java[id.Name]; ok {
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
	case "new":
		w.PutNewCall(c)
	case "print", "println":
		w.PutPrintCall(c)
	// IsType(c): cast
	}
}

// Emit code for built-in new(...) call.
func (w *Writer) PutNewCall(c *ast.CallExpr) {
	assert(len(c.Args) == 1)
	arg := c.Args[0]
	switch t := TypeOf(arg).(type) {
	default:
		panic("cannot handle new " + reflect.TypeOf(t).String())
	case *types.Named:
		w.Put("new ", JavaTypeOfPtr(arg), "()")
	case *types.Basic:
		w.Put("new ", JavaTypeOfPtr(arg), "()")
	case *types.Pointer:
		w.Put("new ", JavaTypeOfPtr(arg), "()")
	}
}

// Emit code for built-in print, prinln calls.
func (w *Writer) PutPrintCall(c *ast.CallExpr) {
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
	arg := n.Args[0]

	switch t := TypeOf(arg).(type) {
	default:
		goto ERR
	case *types.Basic:
		if t.Info()&types.IsString != 0 {
			w.Put("(", n.Args[0], ").length()")
		} else {
			goto ERR
		}
	}

	return

ERR:
	Error(n, "invalid argument (type ", TypeOf(arg), ") for len")
}

// maps Go primitives to java
var builtin2java = map[string]string{
	"bool":    "boolean",
	"byte":    "byte",
	"float32": "float",
	"float64": "double",
	"int":     "int",
	"int16":   "short",
	"int32":   "int",
	"int64":   "long",
	"int8":    "byte",
	"string":  "String",
	"uint":    "int",
	"uint16":  "short",
	"uint32":  "int",
	"uint64":  "long",
	"uint8":   "byte",

	"true":  "true",
	"false": "false",
	"nil":   "null",
}
