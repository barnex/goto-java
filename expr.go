package gotojava

// This file handles expressions.

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/types"
)

// Emit code for an expression.
func (w *Writer) PutExpr(n ast.Expr) {

	if *flagConstFold {
		if tv, ok := ExactValue(n); ok && tv.Value != nil {
			w.Put(tv.Value.String())
			return
		}
	}

	switch e := n.(type) {
	default:
		Error(n, "cannot handle ", reflect.TypeOf(e))
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
	case *ast.SelectorExpr:
		w.PutSelectorExpr(e)
	case *ast.SliceExpr:
		w.PutSliceExpr(e)
	case *ast.StarExpr:
		w.PutStarExpr(e)
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
func (w *Writer) PutBasicLit(n *ast.BasicLit) {
	typ := TypeOf(n).(*types.Basic)
	info := typ.Info()
	switch {
	default:
		panic("cannot handle " + n.Value)
	case info&types.IsUnsigned != 0:
		panic("unsigned")
	case info&types.IsInteger != 0:
		w.Put(n.Value)
	case info&types.IsFloat != 0:
		w.Put(n.Value)
	case info&types.IsString != 0:
		str, err := strconv.Unquote(n.Value)
		checkUserErr(err)
		w.Put(fmt.Sprintf("%q", str)) // TODO: flag for "%q"?
	}
}

// Emit a unary expression, execpt unary "*".
// spec:
// 	unary_op = "+" | "-" | "!" | "^" | "*" | "&" | "<-"
func (w *Writer) PutUnaryExpr(u *ast.UnaryExpr) {
	switch u.Op {
	default:
		Error(u, "unary ", u.Op, " not supported")
	case token.AND:
		w.PutAddressOf(u.X)
	case token.ADD, token.SUB, token.NOT:
		w.Put(u.Op.String(), u.X)
	case token.XOR:
		w.Put("~", u.X)
	}
}

func (w *Writer) PutAddressOf(x ast.Expr) {
	switch t := TypeOf(x).Underlying().(type) {
	default:
		panic("cannot take address of " + reflect.TypeOf(t).String())
	case *types.Struct:
		w.Put(x, ".addr()")
	}
}

func (w *Writer) PutStarExpr(x *ast.StarExpr) {
	w.Put(x.X, ".value()")
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
func (w *Writer) PutSliceExpr(e *ast.SliceExpr) {
	jType := JavaTypeOf(e.X)
	switch jType {
	default:
		Error(e, "cannot slice type ", jType)
	case "String":
		w.putStringSlice(e)
	}
}

// Emit code for slicing a string.
func (w *Writer) putStringSlice(e *ast.SliceExpr) {
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

// Emit code for selector expression: a.b. Godoc:
// 	type SelectorExpr struct {
// 	        X   Expr   // expression
// 	        Sel *Ident // field selector
// 	}
func (w *Writer) PutSelectorExpr(e *ast.SelectorExpr) {
	w.Put(e.X, ".", e.Sel)
}

// Emit code for a parnthesized expression.
// TODO: in many other places parens are inserted,
// do not put parens around parens.
func (w *Writer) PutParenExpr(e *ast.ParenExpr) {
	w.Put("(", e.X, ")")
}

// Emit code for a binary op.
// 	binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
// 	rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
// 	add_op     = "+" | "-" | "|" | "^" .
// 	mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
func (w *Writer) PutBinaryExpr(b *ast.BinaryExpr) {
	if *flagParens {
		w.Put("(")
	}

	unsigned := IsUnsigned(TypeOf(b.X)) || IsUnsigned(TypeOf(b.Y))

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
func (w *Writer) PutCallExpr(n *ast.CallExpr) {
	if n.Ellipsis != 0 {
		Error(n, "cannot handle ellipsis...")
	}
	if IsBuiltin(n.Fun) {
		w.PutBuiltinCall(n)
		return
	}

	w.PutExpr(n.Fun) // TODO: parenthesized = problematic

	//signature := TypeOf(n.Fun).(*types.Signature) // go/types doc says it's always a signature
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

func (w *Writer) PutArgs(args []ast.Expr, ellipsis token.Pos) {
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
