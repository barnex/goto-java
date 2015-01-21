package main

import (
	"go/ast"
	"reflect"
)

// Generate code for expression
func (w *writer) PutExpr(n ast.Expr) {
	switch e := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(e))
	case *ast.CallExpr:
		w.PutCallExpr(e)
	case *ast.Ident:
		w.PutIdent(e)
	case *ast.BasicLit:
		w.PutBasicLit(e)
	case *ast.BinaryExpr:
		w.PutBinaryExpr(e)
	case *ast.ParenExpr:
		w.PutParenExpr(e)
	case *ast.SliceExpr:
		w.PutSliceExpr(e)
	}
}

// TODO: overlaps with builtin
var keywordMap = map[string]string{
	"println": "System.out.println",
	"print":   "System.out.print",
}

func (w *writer) PutIdent(n *ast.Ident) {
	name := n.Name
	// translate name if keyword
	if trans, ok := keywordMap[name]; ok {
		name = trans
	}
	w.Put(name)
}

func (w *writer) PutSliceExpr(e *ast.SliceExpr) {
	jType := w.javaTypeOf(e.X)
	switch jType {
	default:
		w.error(e, "cannot slice type ", jType)
	case "String":
		w.putStringSlice(e)
	}
}

// code for slicing a string
func (w *writer) putStringSlice(e *ast.SliceExpr) {
	w.Put(e.X, ".substring(")
	if e.Low == nil {
		w.Put(0)
	} else {
		w.PutExpr(e.Low)
	}
	w.Put(", ")

	if e.High == nil {
		w.Put("(")
		w.PutExpr(e.X)
		w.Put(").length()")
	} else {
		w.PutExpr(e.High)
	}
	w.Put(")")
}

func (w *writer) PutParenExpr(e *ast.ParenExpr) {
	w.Put("(")
	w.PutExpr(e.X)
	w.Put(")")
}

func (w *writer) PutBinaryExpr(b *ast.BinaryExpr) {
	// TODO: check unsupported ops
	w.PutExpr(b.X)
	w.Put(b.Op)
	w.PutExpr(b.Y)
}

func (w *writer) PutCallExpr(n *ast.CallExpr) {
	if IsBuiltinExpr(n.Fun) {
		w.PutBuiltinCall(n)
		return
	}

	w.PutExpr(n.Fun)

	w.Put("(")
	for i, a := range n.Args {
		if i != 0 {
			w.Put(",")
		}
		w.PutExpr(a)
	}
	w.Put(")")

	if n.Ellipsis != 0 {
		w.error(n, "cannot handle ellipsis")
	}
}

func (w *writer) PutBasicLit(n *ast.BasicLit) {
	w.Put(n.Value)
	// TODO: translate backquotes, complex etc.
}
