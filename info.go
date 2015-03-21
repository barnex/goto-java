package gotojava

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

// PosOf returns the position of n using the global fset.
func PosOf(n ast.Node) token.Position {
	return fset.Position(n.Pos())
}

// ObjectOf returns the object denoted by the specified identifier.
func ObjectOf(id *ast.Ident) types.Object {
	obj := objectOf(id)
	if obj == nil {
		Error(id, "undefined:", id.Name)
	}
	return obj
}

func objectOf(id *ast.Ident) types.Object {
	return info.ObjectOf(id)
}

func TypeOf(n ast.Expr) types.Type {
	t := info.TypeOf(n)
	if t == nil {
		Error(n, "cannot infer type")
	}
	return t
}

// ParentOf returns the parent node of n.
// Precondition: CollectParents has been called on the tree containing n.
func ParentOf(n ast.Node) ast.Node {
	if p, ok := parent[n]; ok {
		return p
	} else {
		panic(PosOf(n).String() + ": no parent")
	}
}

// Return the first a ancestor of n that is an ast.FuncDecl.
// Used by return statements to find the function declaration they return from.
func FuncDeclOf(n ast.Node) *ast.FuncDecl {
	for p := ParentOf(n); p != nil; p = ParentOf(p) {
		if f, ok := p.(*ast.FuncDecl); ok {
			return f
		}
	}
	panic(PosOf(n).String() + ": no FuncDecl parent for node")
}

// returun exact value and minimal type for constant expression.
func ExactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = info.Types[e]
	return
}
