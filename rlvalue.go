package gotojava

import (
	"go/ast"
	"reflect"
)

func RValue(rhs ast.Expr) interface{} {
	if rhs == nil {
		return ""
	}

	if JTypeOfExpr(rhs).IsEscaped() {
		return Transpile(rhs, ".value")
	}

	return rhs
}

func LValue(lhs ast.Expr) interface{} {
	lhs = StripParens(lhs)
	switch lhs := lhs.(type) {
	default:
		panic("cannot make LValue: " + reflect.TypeOf(lhs).String())
	case *ast.Ident:
		return lhs
		//case *ast.StarExpr:
		//	return lhs.X
		//case *ast.SelectorExpr:
		//	return Transpile(lhs.X, ".", lhs.Sel)
	}
}
