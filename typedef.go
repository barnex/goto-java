package gotojava

// Collect type definitions and methods

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

var typedefs map[types.Object]*TypeDef // collects named type defintions and methods
var unnamed map[string]types.Type      // collects uses of all unnamed types

// TypeDef represents type+methods definitions.
// Turned into one or more java classes by classgen.go.
type TypeDef struct {
	typeSpec   *ast.TypeSpec   // type definition
	valMethods []*ast.FuncDecl // value methods
	ptrMethods []*ast.FuncDecl // pointer methods
}

// Collect all type definitions in the AST rooted at root.
// Save them to global typedefs/structs maps
func CollectDefs(root ast.Node) {
	typedefs = make(map[types.Object]*TypeDef)
	unnamed = make(map[string]types.Type)
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		default:
			return true
		case *ast.TypeSpec: // named types
			collectTypeSpec(n)
		case *ast.FuncDecl: // methods for named types
			if n.Recv != nil {
				collectMethodDecl(n)
			}
		case *ast.StructType: // unnamed types
			// TODO: others
			if unnamed[TypeOf(n).String()] == nil {
				Log(n, "discovered:", TypeOf(n).String())
			}
			unnamed[TypeOf(n).String()] = TypeOf(n)
		}
		return true
	})
}

// fetches typedef for typeId from global typedefs
func classOf(typeId *ast.Ident) *TypeDef {
	cls := ObjectOf(typeId)
	if def, ok := typedefs[cls]; ok {
		return def
	} else {
		def := new(TypeDef)
		typedefs[cls] = def
		return def
	}
}

// CollectTypeSpec sets the type declaration of the corresponding class (in global typedefs variable).
// Code generation is deferred until all methods are known.
// 	type TypeSpec struct {
// 	        Doc     *CommentGroup // associated documentation; or nil
// 	        Name    *Ident        // type name
// 	        Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
// 	        Comment *CommentGroup // line comments; or nil
// 	}
func collectTypeSpec(s *ast.TypeSpec) {
	Log(s, "discovered:", s.Name)
	cls := classOf(s.Name)
	assert(cls.typeSpec == nil)
	cls.typeSpec = s
}

// CollectMethodDecl adds a method declaration to the corresponding class's method set (in global typedefs variable).
// Code generation is deferred until all methods are known.
// 	type FuncDecl struct {
// 	        Doc  *CommentGroup // associated documentation; or nil
// 	        Recv *FieldList    // receiver (methods); or nil (functions)
// 	        Name *Ident        // function/method name
// 	        Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	        Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
func collectMethodDecl(s *ast.FuncDecl) {
	Log(s, s.Name)
	rl := s.Recv.List
	assert(len(rl) == 1)
	recvTyp := rl[0].Type

	// TODO: switch

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
