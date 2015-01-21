package main

import (
	"go/ast"
)

func (w *writer) TypeOf(n ast.Expr) string {
	t := w.info.TypeOf(n)
	if t == nil {
		w.error(n, "cannot infer type")
	}
	return t.String()
}
