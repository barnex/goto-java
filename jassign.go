package gotojava

import "golang.org/x/tools/go/types"

type JAssign struct {
	Lhs []JExpr
	Rhs JExpr
}

func (j *JAssign) jStmt() {}

// 	type Initializer struct {
// 	    Lhs []*Var // var Lhs = Rhs
// 	    Rhs ast.Expr
// 	}
func CompileInitAssign(init *types.Initializer) *JAssign {

	j := new(JAssign)
	for _, l := range init.Lhs {
		j.Lhs = append(j.Lhs, NewJVar(l.Name()))
	}
	j.Rhs = CompileExpr(init.Rhs)
	return j
}
