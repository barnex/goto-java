package gotojava

import "go/ast"

type JExpr interface {
	jExpr() // dummy
}

func CompileExpr(e ast.Expr) JExpr {
	panic(0)
}
