package main

import "go/ast"

func (w *writer) PutEmptyInterfaceCast(e ast.Expr) {
	goType := TypeOf(e).Underlying().String()

	switch goType {
	default:
		Error(e, "cannot assign "+goType+" to interface{}")
	case "int":
		w.Put("new go.Int(", e, ")")
	}
}
