package main

// Handling of type definitions and methods

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

var (
	classes = make(map[types.Object]*TypeDef)
)

type TypeDef struct {
	typeSpec   *ast.TypeSpec
	valMethods []*ast.FuncDecl
	ptrMethods []*ast.FuncDecl
}

// RecordTypeSpec sets the type declaration of the corresponding class (in global classes variable).
// Code generation is deferred until all methods are known.
// 	type TypeSpec struct {
// 	        Doc     *CommentGroup // associated documentation; or nil
// 	        Name    *Ident        // type name
// 	        Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
// 	        Comment *CommentGroup // line comments; or nil
// 	}
func RecordTypeSpec(s *ast.TypeSpec) {
	cls := classOf(s.Name)
	assert(cls.typeSpec == nil)
	cls.typeSpec = s
}

// RecordMethodDecl adds a method declaration to the corresponding class's method set (in global classes variable).
// Code generation is deferred until all methods are known.
// 	type FuncDecl struct {
// 	        Doc  *CommentGroup // associated documentation; or nil
// 	        Recv *FieldList    // receiver (methods); or nil (functions)
// 	        Name *Ident        // function/method name
// 	        Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	        Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
func RecordMethodDecl(s *ast.FuncDecl) {
	rl := s.Recv.List
	assert(len(rl) == 1)
	assert(len(rl[0].Names) == 1)
	//recvName := rl[0].Names[0]
	recvTyp := rl[0].Type

	// method on value, e.g., func(T)M(){}
	if id, ok := recvTyp.(*ast.Ident); ok {
		classDef := classOf(id)
		classDef.valMethods = append(classDef.valMethods, s)
		return
	}

	// method on pointer, e.g., func(*T)M(){}
	if star, ok := recvTyp.(*ast.StarExpr); ok {
		id := star.X.(*ast.Ident)
		classDef := classOf(id)
		classDef.ptrMethods = append(classDef.ptrMethods, s)
		return
	}

	Error(s, "invalid receiver: "+reflect.TypeOf(recvTyp).String())
}

// generate code for all defs in global classes variable
func GenClasses() {
	for _, c := range classes {
		name := ClassNameFor(c.typeSpec.Name)
		w := NewWriter(name + ".java")
		w.PutTypeDef(name, c)
		w.Close()
	}
}

func (w *writer) PutTypeDef(name string, c *TypeDef) {
	w.Putln("public final class ", name, "{")
	w.indent++

	switch typ := c.typeSpec.Type.(type) {
	default:
		Error(typ, "cannot handle", reflect.TypeOf(typ))
	case *ast.StructType:
		w.PutStructDef(typ)
	}

	w.indent--
	w.Putln("}")
}

// 	type StructType struct {
// 	        Struct     token.Pos  // position of "struct" keyword
// 	        Fields     *FieldList // list of field declarations
// 	        Incomplete bool       // true if (source) fields are missing in the Fields list
// 	}
func (w *writer) PutStructDef(c *ast.StructType) {
	for _, f := range c.Fields.List {
		w.Put(JavaType(TypeOf(f.Type)), " ")
		for i, n := range f.Names {
			w.Put(comma(i), n)
		}
		w.Putln(";")
	}
}

func ClassNameFor(typ ast.Expr) string {
	switch typ := typ.(type) {
	default:
		Error(typ, "cannot handle", reflect.TypeOf(typ))
		panic("")
	case *ast.Ident:
		return typ.Name // TODO: rename
	}
}

func classOf(typeId *ast.Ident) *TypeDef {
	cls := ObjectOf(typeId)
	if def, ok := classes[cls]; ok {
		return def
	} else {
		def := new(TypeDef)
		classes[cls] = def
		return def
	}
}
