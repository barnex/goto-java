package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/types"
)

var typeToJava = map[string]string{
	"bool":    "boolean",
	"float32": "float",
	"float64": "double",
	"int":     "int", //?
	"int16":   "short",
	"int32":   "int",
	"int64":   "long",
	"int8":    "byte",
	"string":  "String", //?
	"uint":    "int",
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

	// remove untyped. TODO: cli switch?
	if strings.HasPrefix(ident, "untyped ") {
		ident = ident[len("untyped "):]
	}

	if j, ok := typeToJava[ident]; ok {
		return j
	}
	panic("cannot convert type to java: " + ident)
}

// returun exact value and minimal type for constant expression.
func (w *writer) exactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = w.info.Types[e]
	return
}
