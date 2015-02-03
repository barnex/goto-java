package main

import "go/ast"

func (w *writer) PutEmptyInterfaceCast(e ast.Expr) {
	goType := w.TypeOf(e).Underlying().String()

	switch goType {
	default:
		w.error(e, "cannot assign "+goType+" to interface{}")
	case "int":
		w.Put("new go.Int(", e, ")")
	}
}
