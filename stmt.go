package gotojava

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Emit code for a statement.
// Statements are emitted without final semicolon,
// it is up the caller to append a semicolon where needed.
// This allows the statement to be put inside, e.g., a for loop.
func (w *Writer) PutStmt(s ast.Stmt) {
	switch s := s.(type) {
	default:
		Error(s, "cannot handle ", reflect.TypeOf(s))
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
func (w *Writer) PutSwitchStmt(s *ast.SwitchStmt) {
	if s.Init != nil {
		Error(s, "switch init not supported")
	}
	if s.Tag == nil {
		Error(s, "switch w/o tag not supported")
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
func (w *Writer) PutBranchStmt(b *ast.BranchStmt) {
	if b.Label != nil {
		Error(b, b.Tok.String(), " with label not allowed")
	}
	switch b.Tok {
	default:
		Error(b, "cannot handle ", b.Tok)
	case token.BREAK, token.CONTINUE:
		w.Put(b.Tok.String())
	case token.FALLTHROUGH:
		// fallthrough does not exist in java, it should never be emitted.
		// Instead, PutSwitchStmt handles it as a special thing.
		// If we do reach this code, it's either a bug or
		// a misplaced fallthrough that slipped through the parser.
		Error(b, b.Tok, "not allowed here")
	}
}

// Emit ++ or -- statement.
// IncDecStmt godoc:
// 	type IncDecStmt struct {
// 	        X      Expr
// 	        TokPos token.Pos   // position of Tok
// 	        Tok    token.Token // INC or DEC
// 	}
func (w *Writer) PutIncDecStmt(s *ast.IncDecStmt) {
	w.Put(LValue(s.X), s.Tok.String())
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
func (w *Writer) PutForStmt(f *ast.ForStmt) {
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
func (w *Writer) PutIfStmt(i *ast.IfStmt) {

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

// Emit a return statement. Multiple return values are wrapped in a Tuple.
// ReturnStmt godoc:
// 	type ReturnStmt struct {
// 	        Return  token.Pos // position of "return" keyword
// 	        Results []Expr    // result expressions; or nil
// 	}
func (w *Writer) PutReturnStmt(r *ast.ReturnStmt) {
	results := r.Results
	names, types := FlattenFields(FuncDeclOf(r).Type.Results) // function declaration belonging to this return

	// dress a naked return with its results
	if len(results) == 0 && len(names) != 0 {
		results = make([]ast.Expr, len(names))
		for i := range results {
			results[i] = names[i]
		}
	}

	switch len(results) {
	case 0:
		w.Put("return")
	case 1:
		w.Put("return ", results[0])
	default:
		w.Put("return new ", JavaTupleType(types))
		w.PutArgs(results, 0)
	}
}

// Emit a braced statement list.
// BlockStmt godoc:
// 	type BlockStmt struct {
// 	        Lbrace token.Pos // position of "{"
// 	        List   []Stmt
// 	        Rbrace token.Pos // position of "}"
// 	}
func (w *Writer) PutBlockStmt(b *ast.BlockStmt) {
	w.Putln("{")
	w.indent++
	w.PutStmtList(b.List)
	w.indent--
	w.Put("}")
}

func (w *Writer) PutStmtList(list []ast.Stmt) {
	prevLine := -1
	for _, n := range list {
		w.PutStmt(n)

		if NeedSemicolon(n) {
			w.Putln(";")
		} else {
			w.Putln()
		}

		// put empty lines if present in original source
		// TODO: writer: global lastLine, handle everywhere
		line := PosOf(n).Line
		if prevLine != -1 {
			emptylines := line - prevLine - 1
			for i := 0; i < emptylines; i++ {
				w.Putln()
			}
		}
		prevLine = line
	}
}

// does this statement need a terminating semicolon if part of a statement list?
func NeedSemicolon(s ast.Stmt) bool {
	switch s := s.(type) {
	default:
		return true
	case *ast.BlockStmt, *ast.ForStmt, *ast.IfStmt, *ast.SwitchStmt:
		return false
	case *ast.DeclStmt:
		return s.Decl.(*ast.GenDecl).Tok != token.TYPE
	}
}

// Emit a declaration in a statement list.
// DeclStmt godoc:
// 	type DeclStmt struct {
// 	        Decl Decl // *GenDecl with CONST, TYPE, or VAR token
// 	}
func (w *Writer) PutDeclStmt(d *ast.DeclStmt) {
	w.PutDecl(NONE, d.Decl)
}

// Emit a (stand-alone) expression in a statement list.
// ExprStmt godoc:
// 	type ExprStmt struct {
// 	        X Expr // expression
// 	}
func (w *Writer) PutExprStmt(n *ast.ExprStmt) {
	w.Put(n.X)
}
