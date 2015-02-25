package gotojava

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

func (w *Writer) PutAddressOf(x ast.Expr) {
	switch t := TypeOf(x).Underlying().(type) {
	default:
		panic("cannot take address of " + reflect.TypeOf(t).String())
	case *types.Struct:
		w.Put(x, ".addr()")
	}
}

func (w *Writer) PutStarExpr(x *ast.StarExpr) {
	w.Put(x.X, ".value()")
}
