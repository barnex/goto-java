package main

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Emit code for a statement.
// Statements are emitted without final semicolon,
// it is up the caller to append a semicolon where needed.
// This allows the statement to be put inside, e.g., a for loop.
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

// Emit a switch statement.
// SwitchStmt godoc:
// 	type SwitchStmt struct {
// 	        Switch token.Pos  // position of "switch" keyword
// 	        Init   Stmt       // initialization statement; or nil
// 	        Tag    Expr       // tag expression; or nil
// 	        Body   *BlockStmt // CaseClauses only
// 	}
func (w *writer) PutSwitchStmt(s *ast.SwitchStmt) {
	if s.Init != nil {
		w.error(s, "switch init not supported")
	}
	if s.Tag == nil {
		w.error(s, "switch w/o tag not supported")
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
// BranchStmt godoc:
// 	type BranchStmt struct {
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
// 	        Label  *Ident      // label name; or nil
// 	}
func (w *writer) PutBranchStmt(b *ast.BranchStmt) {
	if b.Label != nil {
		w.error(b, b.Tok.String(), " with label not allowed")
	}
	switch b.Tok {
	default:
		w.error(b, "cannot handle ", b.Tok)
	case token.BREAK, token.CONTINUE:
		w.Put(b.Tok.String())
	case token.FALLTHROUGH:
		// fallthrough does not exist in java, it should never be emitted.
		// Instead, PutSwitchStmt handles it as a special thing.
		// If we do reach this code, it's either a bug or
		// a misplaced fallthrough that slipped through the parser.
		w.error(b, b.Tok, "not allowed here")
	}
}

// Emit ++ or -- statement.
// IncDecStmt godoc:
// 	type IncDecStmt struct {
// 	        X      Expr
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // INC or DEC
// 	}
func (w *writer) PutIncDecStmt(s *ast.IncDecStmt) {
	w.Put(s.X, s.Tok.String())
}

// Emit a for statement.
// FotStmt godoc:
// 	type ForStmt struct {
// 	        For  token.Pos // position of "for" keyword
// 	        Init Stmt      // initialization statement; or nil
// 	        Cond Expr      // condition; or nil
// 	        Post Stmt      // post iteration statement; or nil
// 	        Body *BlockStmt
// 	}
func (w *writer) PutForStmt(f *ast.ForStmt) {
	w.Put("for (", nnil(f.Init), "; ", nnil(f.Cond), "; ", nnil(f.Post), ")")
	w.Putln(f.Body)
}

// Emit if statement.
// IfStmt godoc:
// 	type IfStmt struct {
// 	        If   token.Pos // position of "if" keyword
// 	        Init Stmt      // initialization statement; or nil
// 	        Cond Expr      // condition
// 	        Body *BlockStmt
// 	        Else Stmt // else branch; or nil
// 	}
func (w *writer) PutIfStmt(i *ast.IfStmt) {

	// put init statement in front
	// guard scope with braces
	// TODO: there is still a potential shadowing problem
	if i.Init != nil {
		w.Putln("{")
		w.indent++
		w.Putln(i.Init, ";")
	}

	w.Put("if (", i.Cond, ")", i.Body)

	if i.Else != nil {
		w.Put("else ", i.Else)
	}

	if i.Init != nil {
		w.indent--
		w.Putln()
		w.Putln("}")
	}
}

// Emit a return statement.
// ReturnStmt godoc:
// 	type ReturnStmt struct {
// 	        Return  token.Pos // position of "return" keyword
// 	        Results []Expr    // result expressions; or nil
// 	}
func (w *writer) PutReturnStmt(r *ast.ReturnStmt) {
	if len(r.Results) > 1 {
		w.error(r, "cannot handle multiple return values")
	}
	w.Put("return ", r.Results[0])
}

// Emit a braced statement list.
// BlockStmt godoc:
// 	type BlockStmt struct {
// 	        Lbrace token.Pos // position of "{"
// 	        List   []Stmt
// 	        Rbrace token.Pos // position of "}"
// 	}
func (w *writer) PutBlockStmt(n *ast.BlockStmt) {
	w.Putln("{")
	w.indent++

	for _, n := range n.List {
		w.PutStmt(n)
		if needSemicolon(n) {
			w.Putln(";")
		} else {
			w.Putln()
		}
	}

	w.indent--
	w.Put("}")
}

// does this statement need a terminating semicolon if part of a statement list?
func needSemicolon(s ast.Stmt) bool {
	switch s.(type) {
	default:
		return true
	case *ast.BlockStmt, *ast.ForStmt, *ast.IfStmt, *ast.SwitchStmt:
		return false
	}
}

// Emit a declaration in a statement list.
// DeclStmt godoc:
// 	type DeclStmt struct {
// 	        Decl Decl // *GenDecl with CONST, TYPE, or VAR token
// 	}
func (w *writer) PutDeclStmt(d *ast.DeclStmt) {
	modifier := NONE
	w.PutDecl(modifier, d.Decl)
}

// Emit a (stand-alone) expression in a statement list.
// ExprStmt godoc:
// 	type ExprStmt struct {
// 	        X Expr // expression
// 	}
func (w *writer) PutExprStmt(n *ast.ExprStmt) {
	w.Put(n.X)
}

// Emit an assignment or a short variable declaration.
// AssignStmt godoc:
// 	type AssignStmt struct {
// 	        Lhs    []Expr
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // assignment token, DEFINE
// 	        Rhs    []Expr
// 	}
func (w *writer) PutAssignStmt(n *ast.AssignStmt) {
	if len(n.Lhs) != len(n.Rhs) {
		w.error(n, "assignment count mismatch:", len(n.Lhs), "!=", len(n.Rhs))
		// TODO: function with multiple returns
	}

	// java does not have &^=, translate
	if n.Tok == token.AND_NOT_ASSIGN {
		if len(n.Lhs) != 1 || len(n.Rhs) != 1 {
			// should have been caught by type checker.
			w.error(n, n.Tok.String(), " requires single argument")
		}
		w.Put(n.Lhs[0], " &= ", " ~", "(", n.Rhs[0], ")")
		return
	}

	if n.Tok == token.DEFINE {
		w.PutDefine(JModifier(NONE), n)
		return
	}

	// multiple assign: put one per line
	for i := range n.Lhs {
		if i != 0 {
			w.Putln(";")
		}
		w.Put(n.Lhs[i], " ", n.Tok, " ", n.Rhs[i])
	}
}

func (w *writer) PutDefine(mod JModifier, a *ast.AssignStmt) {
	for i, n := range a.Lhs {
		var value ast.Expr = nil
		if i < len(a.Rhs) {
			value = a.Rhs[i]
		}
		if i != 0 {
			w.Putln(";")
		}
		id := a.Lhs[i].(*ast.Ident) // should be
		w.PutValueSpecLine(mod, w.javaTypeOf(n), []*ast.Ident{id}, []ast.Expr{value}, nil)
	}
}
