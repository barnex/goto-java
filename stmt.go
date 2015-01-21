package main

import (
	"go/ast"
	"go/token"
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
	case *ast.AssignStmt:
		w.putAssignStmt(s)
	}
}

func (w *writer) putBlockStmt(n *ast.BlockStmt) {
	w.putln("{")
	w.indent++

	for _, n := range n.List {
		w.putStmt(n)
	}

	w.indent--
	w.putln("}")
}

func (w *writer) putDeclStmt(d *ast.DeclStmt) {
	w.putDecl(d.Decl)
}

func (w *writer) putExprStmt(n *ast.ExprStmt) {
	w.putExpr(n.X)
	w.putln(";")
	//w.putComment(n.Comment)
}

func (w *writer) putAssignStmt(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		w.error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
	}

	tok := n.Tok.String()
	if n.Tok == token.DEFINE {
		tok = "="
	}

	for i := range n.Lhs {
		typ := ""
		if n.Tok == token.DEFINE {
			typ = w.TypeOf(n.Rhs[i]) + " "
		}
		w.put(typ)
		w.putExpr(n.Lhs[i])
		w.put(tok)
		w.putExpr(n.Rhs[i])
		w.putln(";")
	}
}
