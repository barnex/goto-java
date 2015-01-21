package main

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

func (w *writer) putTypeOf(n ast.Expr) {
	t := w.info.TypeOf(n)
	if t == nil {
		w.error(n, "cannot infer type")
	}
	w.put(w.typeConv(t))
}

var typeMap = map[string]string{
	"float32": "float",
	"float64": "double",
	"string":  "String",
}

func (w *writer) typeConv(t types.Type) string {
	orig := t.String()
	if conv, ok := typeMap[orig]; ok {
		return conv
	} else {
		return orig
	}
}
