package gotojava

// Handling of type definitions and methods

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

type TypeDef struct {
	typeSpec   *ast.TypeSpec
	valMethods []*ast.FuncDecl
	ptrMethods []*ast.FuncDecl
}

func CollectDefs(root ast.Node) {
	typedefs = make(map[types.Object]*TypeDef)
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		default:
			return true
		case *ast.TypeSpec:
			CollectTypeSpec(n)
		}
		return true
	})
}

// RecordTypeSpec sets the type declaration of the corresponding class (in global classes variable).
// Code generation is deferred until all methods are known.
// 	type TypeSpec struct {
// 	        Doc     *CommentGroup // associated documentation; or nil
// 	        Name    *Ident        // type name
// 	        Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
// 	        Comment *CommentGroup // line comments; or nil
// 	}
func CollectTypeSpec(s *ast.TypeSpec) {
	Log(s, s.Name)
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
	for _, c := range typedefs {
		Log(nil, c.typeSpec.Name)

		GenClass(c)
	}
}

func GenClass(c *TypeDef) {
	name := ClassNameFor(c.typeSpec.Name)
	w := NewWriter(name + ".java")
	defer w.Close()

	w.PutDoc(c.typeSpec.Doc)
	w.Putln("public final class ", name, "{")
	w.Putln()
	w.indent++

	typ := c.typeSpec.Type
	switch typ.(type) {
	default:
		Error(typ, "cannot handle", reflect.TypeOf(typ))
	case *ast.StructType:
		w.PutStructDef(c)
	}

	w.indent--
	w.Putln("}")
}

// 	type StructType struct {
// 	        Struct     token.Pos  // position of "struct" keyword
// 	        Fields     *FieldList // list of field declarations
// 	        Incomplete bool       // true if (source) fields are missing in the Fields list
// 	}
func (w *Writer) PutStructDef(def *TypeDef) {
	// Fields
	spec := def.typeSpec.Type.(*ast.StructType)
	names, types := FlattenFields(spec.Fields)
	for i, n := range names {
		t := types[i]
		w.Put(ModifierFor(n), " ")
		w.Put(t)
		w.Put(" ", n, ";")
		// TODO Docs
	}
	w.Putln()

	// Methods on value
	for _, m := range def.valMethods {
		w.PutMethodDecl(m)
	}
	// Methods on pointer
	for _, m := range def.ptrMethods {
		w.PutMethodDecl(m)
	}
}

func (w *Writer) PutMethodDecl(f *ast.FuncDecl) {

	// actual implementation with "this" as first receiver
	w.PutStaticFunc(f)
	w.Putln()

	w.PutDoc(f.Doc)
	w.Put(ModifierFor(f.Name), " ")

	// return type
	_, retTypes := FlattenFields(f.Type.Results)
	w.Put(JavaReturnTypeOf(retTypes), " ", f.Name)

	// arguments
	w.Put("(")
	argNames, argTypes := FlattenFields(f.Type.Params)
	for i := range argNames {
		w.Put(comma(i), javaTypeOf(argTypes[i]), " ", argNames[i])
	}
	w.Put(")")

	w.Putln("{")
	w.indent++

	// body calls static implementation with this as first arg
	if len(retTypes) > 0 {
		w.Put("return ")
	}
	w.Put(f.Name, "(", "this")
	for i := range argNames {
		w.Put(", ", argNames[i])
	}
	w.Putln(");")

	w.indent--
	w.Putln("}")

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
	if def, ok := typedefs[cls]; ok {
		return def
	} else {
		def := new(TypeDef)
		typedefs[cls] = def
		return def
	}
}
