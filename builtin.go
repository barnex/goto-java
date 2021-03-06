package gotojava

// Handle Go built-ins and pre-defined types

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

// IsBuiltin returns true if expression e refers to a built-in identifer. E.g.:
// 	print
// 	(print)
// 	bool
// 	...
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

// Emit code for a built-in identifer. E.g.:
// 	bool -> boolean
func (w *Writer) PutBuiltinIdent(id *ast.Ident) {
	if transl, ok := builtin2java[id.Name]; ok {
		w.Put(transl)
	} else {
		panic("cannot handle built-in identifier: " + id.Name)
	}
}

// Generate code for built-in call. E.g.:
// 	len(x)
func (w *Writer) PutBuiltinCall(c *ast.CallExpr) {
	name := StripParens(c.Fun).(*ast.Ident).Name
	switch name {
	default:
		Error(c, "cannot handle builtin: ", name)
	case "len":
		w.PutLenExpr(c)
	case "new":
		w.putNewCall(c)
	case "print", "println":
		w.putPrintCall(c)
	}
}

// Emit code for built-in new(...) call.
func (w *Writer) putNewCall(c *ast.CallExpr) {
	assert(len(c.Args) == 1)
	arg := c.Args[0]

	ptrType := javaPointerNameForElem(TypeOf(arg))
	elemType := JTypeOfExpr(arg).WrapperName()
	addr := NewAddress()
	w.Put("new ", ptrType, "(new ", elemType, "(", addr, ")", ")")

}

func (w *Writer) PutNew(typ interface{}, args ...interface{}) {
	w.Put("new ", typ, "(")
	for i, a := range args {
		w.Put(comma(i))
		w.Put(a)
	}
	w.Put(")")
}

// Emit code for built-in print, prinln calls.
func (w *Writer) putPrintCall(c *ast.CallExpr) {
	assert(c.Ellipsis == 0) // ... not allowed in print/println call
	name := StripParens(c.Fun).(*ast.Ident).Name

	switch name {
	default:
		panic("not a print/println call:" + name)
	case "print":
		name = "System.out.print"
	case "println":
		name = "System.out.println"
	}

	for i, a := range c.Args {
		a := RValue(a)
		if i == len(c.Args)-1 {
			w.Put(name, "(", a, ")")
		} else {
			w.Putln("System.out.print(", a, ");")
			w.Putln("System.out.print(\" \");")
		}
	}
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
}
