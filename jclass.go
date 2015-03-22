package gotojava

import "go/ast"

type JClass struct {
	Init JBlockStmt
}

func NewJClass() *JClass {
	return &JClass{}
}

// Add a top-level function declaration, e.g.:
// 	func f(a, b int) { ... }
// Godoc:
// 	type FuncDecl struct {
// 	    Doc  *CommentGroup // associated documentation; or nil
// 	    Recv *FieldList    // receiver (methods); or nil (functions)
// 	    Name *Ident        // function/method name
// 	    Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	    Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
func (j *JClass) AddStaticFunc(d *ast.FuncDecl) {
	assert(d.Recv == nil)
	panic(0)
}

// Add a method declaration, e.g.:
// 	func (r *Recv) f(a, b int) { ... }
// Godoc:
// 	type FuncDecl struct {
// 	    Doc  *CommentGroup // associated documentation; or nil
// 	    Recv *FieldList    // receiver (methods); or nil (functions)
// 	    Name *Ident        // function/method name
// 	    Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	    Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
func (j *JClass) AddMethod(d *ast.FuncDecl) {
	panic(0)
}

// Add a generic declaration (import, constant, type or variable)
// godoc:
// 	type GenDecl struct {
// 	    Doc    *CommentGroup // associated documentation; or nil
// 	    TokPos token.Pos     // position of Tok
// 	    Tok    token.Token   // IMPORT, CONST, TYPE, VAR
// 	    Lparen token.Pos     // position of '(', if any
// 	    Specs  []Spec
// 	    Rparen token.Pos     // position of ')', if any
// 	}
// A GenDecl node (generic declaration node) represents an import,
// constant, type or variable declaration. A valid Lparen position
// (Lparen.Line > 0) indicates a parenthesized declaration.
// Relationship between Tok value and Specs element type:
// 	token.IMPORT  *ImportSpec
// 	token.CONST   *ValueSpec
// 	token.TYPE    *TypeSpec
// 	token.VAR     *ValueSpec
//
func (j *JClass) AddGenDecl(d *ast.GenDecl) {
	switch d.Tok {
	default:
		panic(d.Tok.String())
		//case token.TYPE:
		//	// do nothing. already handled by CollectDefs
		//case token.CONST:
		//	j.addValueSpecs(FINAL, d.Specs)
		//case token.VAR:
		//	for _, spec := range specs {
		//		j.addValueSpec(NONE, spec.(*ast.ValueSpec))
		//	}
	}
}

// Add a list of *ast.ValueSpecs, e.g.:
// 	var(
// 		a int
// 		b, c int
// 	)
// or
// 	const(
// 		a = 1
// 		b, c = 2, 3
// 	)
// with optional modifier (e.g. "static", "final", "static final").
// The concrete type of specs elements must be *ast.ValueSpec.
func (j *JClass) addValueSpecs(mod JModifier, specs []ast.Spec) {
}

func (j *JClass) AddValueSpecs() {

}

func (j *JClass) Write(w *Writer) {
	panic(0)
}
