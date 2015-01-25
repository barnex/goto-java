package main

// This file handles expressions.

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Emit code for an expression.
func (w *writer) PutExpr(n ast.Expr) {

	if *flagConstFold {
		if tv, ok := w.exactValue(n); ok && tv.Value != nil {
			w.Put(tv.Value.String())
			return
		}
	}

	switch e := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(e))
	case *ast.BasicLit:
		w.PutBasicLit(e)
	case *ast.BinaryExpr:
		w.PutBinaryExpr(e)
	case *ast.CallExpr:
		w.PutCallExpr(e)
	case *ast.Ident:
		w.PutIdent(e)
	case *ast.ParenExpr:
		w.PutParenExpr(e)
	case *ast.SliceExpr:
		w.PutSliceExpr(e)
	case *ast.UnaryExpr:
		w.PutUnaryExpr(e)
	}
}

// Emit code for a literal of basic type.
// BasicLit godoc:
// 	type BasicLit struct {
// 	        ValuePos token.Pos   // literal position
// 	        Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
// 	        Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
// 	}
func (w *writer) PutBasicLit(n *ast.BasicLit) {
	w.Put(n.Value)
	// TODO: translate backquotes, complex etc.
}

// Emit an identifier, translating built-ins.
// Ident godoc:
// 	type Ident struct {
// 	        NamePos token.Pos // identifier position
// 	        Name    string    // identifier name
// 	        Obj     *Object   // denoted object; or nil
// 	}
func (w *writer) PutIdent(id *ast.Ident) {
	if w.IsBuiltinIdent(id) {
		w.PutBuiltinIdent(id)
	} else {
		if tv, ok := w.info.Defs[id]; ok {
			w.Put(tv.Name())
		} else {
			w.error(id, "undefined: ", id.Name)
		}
	}
}

// Emit a unary expression, execpt unary "*".
// spec:
// 	unary_op = "+" | "-" | "!" | "^" | "*" | "&" | "<-"
func (w *writer) PutUnaryExpr(u *ast.UnaryExpr) {
	switch u.Op {
	default:
		w.error(u, "unary ", u.Op, " not supported")
	case token.ADD, token.SUB, token.NOT: // TODO: xor
		w.Put(u.Op.String(), u.X)
	}
}

// Emit code for a slice expression.
// SliceExpr godoc:
// 	type SliceExpr struct {
// 	        X      Expr      // expression
// 	        Lbrack token.Pos // position of "["
// 	        Low    Expr      // begin of slice range; or nil
// 	        High   Expr      // end of slice range; or nil
// 	        Max    Expr      // maximum capacity of slice; or nil
// 	        Slice3 bool      // true if 3-index slice (2 colons present)
// 	        Rbrack token.Pos // position of "]"
// 	}
func (w *writer) PutSliceExpr(e *ast.SliceExpr) {
	jType := w.javaTypeOf(e.X)
	switch jType {
	default:
		w.error(e, "cannot slice type ", jType)
	case "String":
		w.putStringSlice(e)
	}
}

// Emit code for slicing a string.
func (w *writer) putStringSlice(e *ast.SliceExpr) {
	w.Put(e.X, ".substring(")
	if e.Low == nil {
		w.Put("0")
	} else {
		w.PutExpr(e.Low)
	}
	w.Put(", ")

	if e.High == nil {
		w.Put("(", e.X, ").length()") // need to parenthesize, X may be binary expression.
	} else {
		w.PutExpr(e.High)
	}
	w.Put(")")
}

// Emit code for a parnthesized expression.
// TODO: in many other places parens are inserted,
// do not put parens around parens.
func (w *writer) PutParenExpr(e *ast.ParenExpr) {
	w.Put("(", e.X, ")")
}

// Emit code for a binary op.
// 	binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
// 	rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
// 	add_op     = "+" | "-" | "|" | "^" .
// 	mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
func (w *writer) PutBinaryExpr(b *ast.BinaryExpr) {
	// TODO: check unsupported ops

	if *flagParens {
		w.Put("(")
	}

	switch b.Op {
	default:
		w.Put(b.X, b.Op.String(), b.Y)
	case token.SHL, token.SHR, token.AND, token.OR, token.XOR:
		// different precedence in Go and Java, parentisize to be sure
		w.Put("(", b.X, b.Op.String(), b.Y, ")")
	case token.AND_NOT: //
		// not in java
		w.Put("(", b.X, "&~", b.Y, ")")
	}

	if *flagParens {
		w.Put(")")
	}

}

// Emit code for a call expression.
// CallExpr godoc:
// 	type CallExpr struct {
// 	        Fun      Expr      // function expression
// 	        Lparen   token.Pos // position of "("
// 	        Args     []Expr    // function arguments; or nil
// 	        Ellipsis token.Pos // position of "...", if any
// 	        Rparen   token.Pos // position of ")"
// 	}
// TODO: handle ellipsis.
func (w *writer) PutCallExpr(n *ast.CallExpr) {
	if w.IsBuiltinExpr(n.Fun) {
		w.PutBuiltinCall(n)
		return
	}

	w.PutExpr(n.Fun) // TODO: parenthesized = problematic

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
