package gotojava

import "go/ast"

// Rudimentary escape analysis detects:
// 	address of local variable

var escapes = make(map[*ast.Ident]bool)

func EscapeAnalysis(root ast.Node) {
	ast.Inspect(root, func(n ast.Node) bool {

		switch n := n.(type) {
		default:
			return true
		case *ast.UnaryExpr:
			if id, ok := n.X.(*ast.Ident); ok {
				if IsLocal(id) {
					Log(id, id.Name, "escapes to java heap")
					escapes[id] = true
				}
			}
		}
		return true
	})
}
