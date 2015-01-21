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

func (w *writer) typeConv(t types.Type) string {
	return t.String()
}
