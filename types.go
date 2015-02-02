package main

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

// maps go primitives to java
var typeToJava = map[string]string{
	"bool":    "boolean",
	"byte":    "byte",
	"float32": "float",
	"float64": "double",
	"int":     "int", //?
	"int16":   "short",
	"int32":   "int",
	"int64":   "long",
	"int8":    "byte",
	"string":  "String", //?
	"uint":    "int",    //?
	"uint16":  "short",  //?
	"uint32":  "int",    //?
	"uint64":  "long",   //?
	"uint8":   "byte",   //?
}

func (w *writer) PutTypecast(goType string, e ast.Expr) {
	jType, ok := typeToJava[goType]
	if !ok {
		w.error(e, "cannot convert to java:", goType)
	}
	w.Put("(", jType, ")(", e, ")")
}

func (w *writer) PutImplicitCast(e ast.Expr, goType string) {
	if tv, ok := w.exactValue(e); ok && tv.Value != nil {
		w.PutTypecast(goType, e)
	} else {
		w.PutExpr(e)
	}
}

func (w *writer) TypeOf(n ast.Expr) types.Type {
	t := w.info.TypeOf(n)
	if t == nil {
		w.error(n, "cannot infer type")
	}
	return t
}

func (w *writer) TypeToJava(t types.Type) string {
	ident := t.String() // TODO: underlying?

	// remove untyped.
	//if strings.HasPrefix(ident, "untyped ") {
	//	ident = ident[len("untyped "):]
	//}

	if j, ok := typeToJava[ident]; ok {
		return j
	}
	panic("cannot convert type to java: " + ident)
}

// ObjectOf returns the object denoted by the specified identifier.
func (w *writer) ObjectOf(id *ast.Ident) types.Object {
	obj := w.info.ObjectOf(id)
	if obj == nil {
		w.error(id, "undefined:", id.Name)
	}
	return obj
}

// returun exact value and minimal type for constant expression.
func (w *writer) exactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = w.info.Types[e]
	return
}
