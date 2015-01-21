package main

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

func (w *writer) goTypeOf(n ast.Expr) types.Type {
	t := w.info.TypeOf(n)
	if t == nil {
		w.error(n, "cannot infer type")
	}
	return t
}

func (w *writer) javaTypeOf(n ast.Expr) string {
	return w.toJavaType(w.goTypeOf(n))
}

var typeMap = map[string]string{
	"float32": "float",
	"float64": "double",
	"string":  "String",
	"bool":    "boolean",
}

func (w *writer) toJavaType(t types.Type) string {
	orig := t.String()
	if conv, ok := typeMap[orig]; ok {
		return conv
	} else {
		return orig
	}
}
