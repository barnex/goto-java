package gotojava

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

// Rudimentary escape analysis detects:
// 	address of local variable

var escapes = make(map[types.Object]bool)

func EscapeAnalysis(root ast.Node) {
	ast.Inspect(root, func(n ast.Node) bool {

		switch n := n.(type) {
		default:
			return true
		case *ast.UnaryExpr:
			// check address of local identifier
			if id, ok := n.X.(*ast.Ident); ok && n.Op == token.AND {
				if isLocal(id) {
					Log(id, id.Name, "escapes to java heap")
					escapes[ObjectOf(id)] = true
				}
			}
		}
		return true
	})
}

func Escapes(id *ast.Ident) bool {
	if id == nil {
		return false
	}
	return escapes[ObjectOf(id)]
}

// Checks whether the identifier is locally defined.
// E.g., for escape analysis: address of local variable etc.
func isLocal(id *ast.Ident) bool {
	scope := ObjectOf(id).Parent()
	global := scope.Parent() == types.Universe
	return !global
}
