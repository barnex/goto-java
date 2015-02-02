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
		w.PutResolvedIdent(e)
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
	goType := w.TypeOf(n).String()
	//log.Println(goType, n.Value)

	switch goType {
	default:
		w.Put("(", w.typeToJava(goType), ")", n.Value)
	case "int", "int32", "untyped int":
		w.Put(n.Value)
	case "int64":
		w.Put(n.Value, "L")
	}
}

// Emit an identifier
// Ident godoc:
// 	type Ident struct {
// 	        NamePos token.Pos // identifier position
// 	        Name    string    // identifier name
// 	        Obj     *Object   // denoted object; or nil
// 	}
func (w *writer) PutResolvedIdent(id *ast.Ident) {
	if IsBlank(id) {
		w.Put(makeNewName(UNUSED))
		return
	}
	if w.IsBuiltinIdent(id) {
		w.PutBuiltinIdent(id)
	} else {
		w.Put(w.translate(id))
	}
}

// Emit a unary expression, execpt unary "*".
// spec:
// 	unary_op = "+" | "-" | "!" | "^" | "*" | "&" | "<-"
func (w *writer) PutUnaryExpr(u *ast.UnaryExpr) {
	switch u.Op {
	default:
		w.error(u, "unary ", u.Op, " not supported")
	case token.ADD, token.SUB, token.NOT:
		w.Put(u.Op.String(), u.X)
	case token.XOR:
		w.Put("~", u.X)
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
	jType := w.TypeToJava(w.TypeOf(e.X))
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

	unsigned := w.IsUnsigned(w.TypeOf(b.X)) || w.IsUnsigned(w.TypeOf(b.Y))

	switch b.Op {
	default:
		w.Put(b.X, b.Op.String(), b.Y)
	case token.LSS, token.GTR, token.LEQ, token.GEQ, token.QUO, token.REM:
		if unsigned {
			w.PutUnsignedOp(b.X, b.Op, b.Y)
		} else {
			w.Put(b.X, b.Op.String(), b.Y) // default
		}
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
	if n.Ellipsis != 0 {
		w.error(n, "cannot handle ellipsis...")
	}
	if w.IsBuiltinExpr(n.Fun) {
		w.PutBuiltinCall(n)
		return
	}

	w.PutExpr(n.Fun) // TODO: parenthesized = problematic

	//signature := w.TypeOf(n.Fun).(*types.Signature) // go/types doc says it's always a signature
	//params := signature.Params()

	w.Put("(")
	for i, a := range n.Args {
		if i != 0 {
			w.Put(",")
		}
		w.Put(a)
		//w.PutImplicitCast(a, params.At(i).Type().Underlying().String())
	}
	w.Put(")")
}

func (w *writer) PutArgs(args []ast.Expr, ellipsis token.Pos) {
	w.Put("(")
	for i, a := range args {
		if i != 0 {
			w.Put(",")
		}
		w.PutExpr(a)
	}
	if ellipsis != 0 {
		w.Put("...")
	}
	w.Put(")")
}
