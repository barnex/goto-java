package gotojava

import "go/ast"

func (w *Writer) PutEmptyInterfaceCast(e ast.Expr) {
	goType := TypeOf(e).Underlying().String()

	switch goType {
	default:
		Error(e, "cannot assign "+goType+" to interface{}")
	case "int":
		w.Put("new go.Int(", e, ")")
	}
}
