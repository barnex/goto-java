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
	case *ast.BranchStmt:
		w.PutBranchStmt(s)
	case *ast.DeclStmt:
		w.PutDeclStmt(s)
	case *ast.ExprStmt:
		w.PutExprStmt(s)
	case *ast.ForStmt:
		w.PutForStmt(s)
	case *ast.IfStmt:
		w.PutIfStmt(s)
	case *ast.IncDecStmt:
		w.PutIncDecStmt(s)
	case *ast.ReturnStmt:
		w.PutReturnStmt(s)
	case *ast.SwitchStmt:
		w.PutSwitchStmt(s)
	}
}

func (w *writer) PutSwitchStmt(s *ast.SwitchStmt) {
	if s.Init != nil {
		w.error(s, "switch init not supported")
	}

	w.Putln("switch(", s.Tag, "){")
	w.indent++

	body := s.Body.List
	for _, stmt := range body {
		clause := stmt.(*ast.CaseClause)
		body := clause.Body

		if clause.List == nil {
			w.Putln("default:")
		} else {
			for _, e := range clause.List {
				w.Putln("case ", e, ":")
			}
		}

		w.indent++
		haveFallThrough := false
		for _, s := range body {
			if branch, ok := s.(*ast.BranchStmt); ok {
				if branch.Tok == token.FALLTHROUGH {
					haveFallThrough = true
					continue // do not put "fallthrough" in java
				}
			}
			w.Putln(s, ";")
		}
		w.indent--

		if !haveFallThrough {
			w.Putln("break;")
		}
	}

	w.indent--
	w.Putln("}")
}

// Emit branch statement (breat, continue, goto, fallthrough)
func (w *writer) PutBranchStmt(b *ast.BranchStmt) {
	switch b.Tok {
	default:
		w.error(b, "cannot handle ", b.Tok)
	case token.BREAK, token.CONTINUE:
		w.Put(b.Tok.String())
	}
}

// Emit ++ or -- statement
func (w *writer) PutIncDecStmt(s *ast.IncDecStmt) {
	w.Put(s.X, s.Tok.String())
}

// Emit for statement
func (w *writer) PutForStmt(f *ast.ForStmt) {
	w.Put("for (", nnil(f.Init), "; ", nnil(f.Cond), "; ", nnil(f.Post), ")")
	w.Putln(f.Body)
}

// Emit if statement
func (w *writer) PutIfStmt(i *ast.IfStmt) {

	// put init statement in front
	// guard scope with braces
	//if i.Init != nil {
	//	w.Putln("{")
	//	w.indent++
	//	w.Putln(i.Init)
	//}
	if i.Init != nil {
		w.error(i, "if init statement not supported")
	}

	w.Put("if (", i.Cond, ")", i.Body)

	if i.Else != nil {
		w.Put("else ", i.Else)
	}

	//if i.Init != nil {
	//	w.indent--
	//	w.Putln()
	//	w.Putln("}")
	//}
}

// Emit return statement
func (w *writer) PutReturnStmt(r *ast.ReturnStmt) {
	if len(r.Results) > 1 {
		w.error(r, "cannot handle multiple return values")
	}
	w.Put("return ", r.Results[0])
}

func (w *writer) PutBlockStmt(n *ast.BlockStmt) {
	w.Putln("{")
	w.indent++

	for _, n := range n.List {
		w.PutStmt(n)
		w.Putln(";")
	}

	w.indent--
	w.Put("}")
}

func (w *writer) PutDeclStmt(d *ast.DeclStmt) {
	w.PutDecl(d.Decl)
}

func (w *writer) PutExprStmt(n *ast.ExprStmt) {
	w.Put(n.X)
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
			w.Putln(";")
		}
		if n.Tok == token.DEFINE {
			w.Put(w.javaTypeOf(n.Rhs[i]), " ")
		}
		w.Put(n.Lhs[i], " ", tok, " ", n.Rhs[i])
	}
}
