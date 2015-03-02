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
	case *ast.CompositeLit:
		w.PutCompositeLit(e)
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
	typ := TypeOf(n).Underlying().(*types.Basic)
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

// Emit code for a composite literal. Godoc:
// 	type CompositeLit struct {
// 	        Type   Expr      // literal type; or nil
// 	        Lbrace token.Pos // position of "{"
// 	        Elts   []Expr    // list of composite elements; or nil
// 	        Rbrace token.Pos // position of "}"
// 	}
func (w *Writer) PutCompositeLit(lit *ast.CompositeLit) {

	if len(lit.Elts) == 0 {
		w.Put(ZeroValue(JTypeOf(lit.Type)))
		return
	}

	w.Put("new ", javaPointerNameForElem(TypeOf(lit.Type)))
	if _, ok := lit.Elts[0].(*ast.KeyValueExpr); ok {
		w.Put("(")
		tdef := classOf(lit.Type.(*ast.Ident)) // TODO: map java names to tdef?
		names, types := FlattenFields(tdef.typeSpec.Type.(*ast.StructType).Fields)
		values := make(map[string]interface{})
		for i, n := range names {
			values[n.Name] = ZeroValue(types[i])
		}
		for _, e := range lit.Elts {
			e := e.(*ast.KeyValueExpr)
			values[e.Key.(*ast.Ident).Name] = e.Value
		}
		for i, n := range names {
			w.Put(comma(i), values[n.Name])
		}
		w.Put(")")
	} else {
		w.PutArgs(lit.Elts, 0)
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
	panic("no slice expr yet")
	//jType := JavaTypeOf(e.X)
	//switch jType {
	//default:
	//	Error(e, "cannot slice type ", jType)
	//case "String":
	//	w.putStringSlice(e)
	//}
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
	if JTypeOf(e.X).IsValue() && IsPtrMethod(e.Sel) {
		// Pointer method on addressable value:
		// compiler inserts address of receiver.
		// https://golang.org/doc/effective_go.html#pointers_vs_values
		w.PutAddressOf(e.X)
		w.Put(".", e.Sel)
	} else {
		w.Put(e.X, ".", e.Sel)
	}
}

func IsPtrMethod(id *ast.Ident) bool {
	obj := TypeOf(id)
	if sig, ok := obj.(*types.Signature); ok {
		return sig.Recv() != nil && IsPtrType(sig.Recv().Type())
	} else {
		return false
	}
}

func IsPtrType(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
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

	x := RValue(b.X)
	y := RValue(b.Y)

	switch b.Op {
	default:
		w.Put(x, b.Op.String(), y)
	case token.EQL:
		w.PutJEquals(JTypeOf(b.X), x, JTypeOf(b.Y), y)
	case token.LSS, token.GTR, token.LEQ, token.GEQ, token.QUO, token.REM:
		if unsigned {
			w.PutUnsignedOp(b.X, b.Op, b.Y)
		} else {
			w.Put(x, b.Op.String(), y) // default
		}
	case token.SHL, token.SHR, token.AND, token.OR, token.XOR:
		// different precedence in Go and Java, parentisize to be sure
		w.Put("(", x, b.Op.String(), y, ")")
	case token.AND_NOT: //
		// not in java
		w.Put("(", x, "&~", y, ")")
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
func (w *Writer) PutCallExpr(n *ast.CallExpr) {
	if IsBuiltin(n.Fun) {
		w.PutBuiltinCall(n)
	} else {
		w.PutExpr(n.Fun) // TODO: parenthesized = problematic
		w.PutArgs(n.Args, n.Ellipsis)
	}
}

func (w *Writer) PutArgs(args []ast.Expr, ellipsis token.Pos) {
	w.Put("(")
	for i, a := range args {
		w.Put(comma(i), RValue(a)) // TODO: cast
	}
	if ellipsis != 0 {
		w.Put("...")
	}
	w.Put(")")
}

// explicit type cast in input file, e.g.:
// 	a := int(b)
func (w *Writer) PutTypecast(goType types.Type, e ast.Expr) {
	panic("no typecast yet")
}

//func IsType(x ast.Expr)bool{
//tv, err := types.EvalNode(fset, x, pkg,
//func EvalNode(fset *token.FileSet, node ast.Expr, pkg *Package, scope *Scope) (tv TypeAndValue, err error)
//}

// Emit code for Go's "lhs == rhs", with given java types for both sides.
// May emit, e.g.:
// 	lhs == rhs      // basic and pointer types
// 	lhs.equals(rhs) // struct value types
func (w *Writer) PutJEquals(ltyp JType, lhs interface{}, rtyp JType, rhs interface{}) {
	switch {
	default:
		w.Put(lhs, " == ", rhs) // TODO: panic
	case ltyp.IsStructValue() && rtyp.IsStructValue():
		w.Put(lhs, ".equals(", rhs, ")")
	}
}
