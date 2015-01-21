package main

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

func (w *writer) typeOf(n ast.Expr) types.Type {
	t := w.info.TypeOf(n)
	if t == nil {
		w.error(n, "cannot infer type")
	}
	return t

}

func (w *writer) putTypeOf(n ast.Expr) {
	w.put(w.typeConv(w.typeOf(n)))
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
