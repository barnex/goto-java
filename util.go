package gotojava

import (
	"fmt"
	"go/ast"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/types"
)

// FlattenFields turns an ast.FieldList into a list of names and a list of types of the same length. E.g.:
// 	(a, b int) -> names: [a, b], types: [int, int]
func FlattenFields(list *ast.FieldList) (names []*ast.Ident, types []types.Type) {
	if list == nil {
		return
	}
	for _, f := range list.List {
		if f.Names == nil {
			// unnamed field
			names = append(names, nil)
			types = append(types, TypeOf(f.Type))
		} else {
			for _, n := range f.Names {
				names = append(names, n)
				types = append(types, TypeOf(f.Type))
			}
		}
	}
	assert(len(names) == len(types))
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

// Export turns the first character to upper case.
func Export(name string) string {
	rune, width := utf8.DecodeRuneInString(name)
	return fmt.Sprint(unicode.ToUpper(rune), name[width:])
}

// Unexport turns the first character to lower case.
func Unexport(name string) string {
	rune, width := utf8.DecodeRuneInString(name)
	return fmt.Sprint(unicode.ToLower(rune), name[width:])
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
