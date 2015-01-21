package main

import (
	"go/ast"
	"reflect"
)

func (w *writer) putExpr(n ast.Expr) {
	switch e := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(e))
	case *ast.CallExpr:
		w.putCallExpr(e)
	case *ast.Ident:
		w.putIdent(e)
	case *ast.BasicLit:
		w.putBasicLit(e)
	case *ast.BinaryExpr:
		w.putBinaryExpr(e)
	case *ast.ParenExpr:
		w.putParenExpr(e)
	}
}

var keywordMap = map[string]string{
	"println": "System.out.println",
	"print":   "System.out.print",
}

func (w *writer) putIdent(n *ast.Ident) {
	name := n.Name
	// translate name if keyword
	if trans, ok := keywordMap[name]; ok {
		name = trans
	}
	w.put(name)
}

func (w *writer) putParenExpr(e *ast.ParenExpr) {
	w.put("(")
	w.putExpr(e.X)
	w.put(")")
}

func (w *writer) putBinaryExpr(b *ast.BinaryExpr) {
	// TODO: check unsupported ops
	w.putExpr(b.X)
	w.put(b.Op)
	w.putExpr(b.Y)
}

func (w *writer) putCallExpr(n *ast.CallExpr) {
	w.putExpr(n.Fun)

	w.put("(")
	for i, a := range n.Args {
		if i != 0 {
			w.put(",")
		}
		w.putExpr(a)
	}
	w.put(")")

	if n.Ellipsis != 0 {
		w.error(n, "cannot handle ellipsis")
	}
}

func (w *writer) putBasicLit(n *ast.BasicLit) {
	w.put(n.Value)
	// TODO: translate backquotes, complex etc.
}
