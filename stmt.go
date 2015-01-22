package main

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (w *writer) PutStmt(s ast.Stmt) {
	switch s := s.(type) {
	default:
		w.error(s, "cannot handle ", reflect.TypeOf(s))
	case *ast.AssignStmt:
		w.PutAssignStmt(s)
	case *ast.BlockStmt:
		w.PutBlockStmt(s)
	case *ast.DeclStmt:
		w.PutDeclStmt(s)
	case *ast.ExprStmt:
		w.PutExprStmt(s)
	case *ast.IfStmt:
		w.PutIfStmt(s)
	case *ast.ReturnStmt:
		w.PutReturnStmt(s)
	}
}

// Emit if statement
func (w *writer) PutIfStmt(i *ast.IfStmt) {
	if i.Init != nil {
		w.PutStmt(i.Init)
	}

	w.Put("if (", i.Cond, ")", i.Body)

	if i.Else != nil {
		w.Put("else ", i.Else)
	}
}

// Emit return statement
func (w *writer) PutReturnStmt(r *ast.ReturnStmt) {
	if len(r.Results) > 1 {
		w.error(r, "cannot handle multiple return values")
	}
	w.Put("return ", r.Results[0], ";")
}

func (w *writer) PutBlockStmt(n *ast.BlockStmt) {
	w.Putln("{")
	w.indent++

	for _, n := range n.List {
		w.PutStmt(n)
		w.Putln()
	}

	w.indent--
	w.Put("}")
}

func (w *writer) PutDeclStmt(d *ast.DeclStmt) {
	w.PutDecl(d.Decl)
}

func (w *writer) PutExprStmt(n *ast.ExprStmt) {
	w.Put(n.X, ";")
}

func (w *writer) PutAssignStmt(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		w.error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
	}

	// translate := to =
	tok := n.Tok.String()
	if n.Tok == token.DEFINE {
		tok = "="
	}

	// multiple assign: put one per line
	for i := range n.Lhs {
		if i != 0 {
			w.Putln()
		}
		if n.Tok == token.DEFINE {
			w.Put(w.javaTypeOf(n.Rhs[i]), " ")
		}
		w.Put(n.Lhs[i], " ", tok, " ", n.Rhs[i], ";")
	}
}
