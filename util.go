package main

import (
	"go/ast"

	"golang.org/x/tools/go/types"
)

// FlattenFields turns an ast.FieldList into a list of names and a list of types of the same length. E.g.:
// 	(a, b int) -> names: [a, b], types: [int, int]
func FlattenFields(list *ast.FieldList) (names []*ast.Ident, types []ast.Expr) {
	if list == nil {
		return
	}
	for _, f := range list.List {
		if f.Names == nil {
			// unnamed field
			names = append(names, nil)
			types = append(types, f.Type)
		} else {
			for _, n := range f.Names {
				names = append(names, n)
				types = append(types, f.Type)
			}
		}
	}
	assert(len(names) == len(types))
	return
}

// ObjectOf returns the object denoted by the specified identifier.
func ObjectOf(id *ast.Ident) types.Object {
	obj := info.ObjectOf(id)
	if obj == nil {
		Error(id, "undefined:", id.Name)
	}
	return obj
}

// returun exact value and minimal type for constant expression.
func ExactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = info.Types[e]
	return
}

// Strip parens from expression, if any. E.g.:
// 	((x)) -> x
func StripParens(e ast.Expr) ast.Expr {
	if par, ok := e.(*ast.ParenExpr); ok {
		return StripParens(par.X)
	} else {
		return e
	}
}

// Is e the blank identifier?
// Also handles the corner case of parenthesized blank ident (_)
func IsBlank(e ast.Expr) bool {
	e = StripParens(e)
	if id, ok := e.(*ast.Ident); ok {
		return id.Name == "_"
	}
	return false
}

func TypeOf(n ast.Expr) types.Type {
	t := info.TypeOf(n)
	if t == nil {
		Error(n, "cannot infer type")
	}
	return t
}

// Returns a comma if i!=0.
// Used for comma-separating values from a loop.
func comma(i int) string {
	if i != 0 {
		return ","
	} else {
		return ""
	}
}

func nnil(x interface{}) interface{} {
	if x == nil {
		return ""
	} else {
		return x
	}
}

func assert(test bool) {
	if !test {
		panic("assertion failed")
	}
}
