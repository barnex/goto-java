package main

import (
	"go/ast"
	"reflect"
)

func (w *writer) putStmt(n ast.Stmt) {
	switch s := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(s))
	case *ast.ExprStmt:
		w.putExprStmt(s)
	case *ast.DeclStmt:
		w.putDeclStmt(s)
	}
}

func (w *writer) putBlockStmt(n *ast.BlockStmt) {
	w.putln(LBRACE)
	w.indent++

	for _, n := range n.List {
		w.putStmt(n)
	}

	w.indent--
	w.putln(RBRACE)
}

func (w *writer) putDeclStmt(d *ast.DeclStmt) {
	w.putDecl(d.Decl)
}

func (w *writer) putExprStmt(n *ast.ExprStmt) {
	w.putExpr(n.X)
	//w.putComment(n.Comment)
	w.putln(";")
}
