package gotojava

// Handling of type definitions and methods

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// TypeDef represents a type+methods definition.
type TypeDef struct {
	typeSpec   *ast.TypeSpec   // type definition
	valMethods []*ast.FuncDecl // value methods
	ptrMethods []*ast.FuncDecl // pointer methods
}

// Collect all type definitions in the AST rooted at root.
// Save them to global typedefs map.
func CollectDefs(root ast.Node) {
	typedefs = make(map[types.Object]*TypeDef)
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		default:
			return true
		case *ast.TypeSpec:
			CollectTypeSpec(n)
		case *ast.FuncDecl:
			if n.Recv != nil {
				CollectMethodDecl(n)
			}
		}
		return true
	})
}

// CollectTypeSpec sets the type declaration of the corresponding class (in global typedefs variable).
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

// CollectMethodDecl adds a method declaration to the corresponding class's method set (in global typedefs variable).
// Code generation is deferred until all methods are known.
// 	type FuncDecl struct {
// 	        Doc  *CommentGroup // associated documentation; or nil
// 	        Recv *FieldList    // receiver (methods); or nil (functions)
// 	        Name *Ident        // function/method name
// 	        Type *FuncType     // function signature: parameters, results, and position of "func" keyword
// 	        Body *BlockStmt    // function body; or nil (forward declaration)
// 	}
func CollectMethodDecl(s *ast.FuncDecl) {
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

// generate code for all defs in global typedefs variable
func GenClasses() {
	for _, c := range typedefs {
		Log(nil, c.typeSpec.Name)

		switch tdef := c.typeSpec.Type.(type) {
		default:
			panic("cannot handle: " + reflect.TypeOf(tdef).String())
		case *ast.StructType:
			GenStructPointerClass(c)
			GenStructValueClass(c)
		}
	}
}

// Generate java class for Go pointer-to-named-struct type.
func GenStructPointerClass(d *TypeDef) {
	spec := d.typeSpec

	name := JavaTypeOfPtr(spec.Name)
	base := JavaTypeOfExpr(spec.Name)

	w := NewWriter(name + ".java")
	defer w.Close()

	w.PutDoc(spec.Doc)
	w.Putln("public final class ", name, " extends ", base, "{")
	w.Putln()
	w.indent++

	// Methods on pointer
	for _, m := range d.ptrMethods {
		w.PutMethodDecl(m, false)
	}

	w.indent--
	w.Putln("}")
}

// Generate java class for Go named struct type (value semantics).
func GenStructValueClass(d *TypeDef) {
	spec := d.typeSpec
	name := JavaTypeOfExpr(spec.Name)
	w := NewWriter(name + ".java")
	defer w.Close()

	w.Putln("public class ", name, "{")
	w.Putln()
	w.indent++

	// Fields
	fields := spec.Type.(*ast.StructType).Fields
	w.PutStructFields(fields)
	w.Putln()

	// Constructors
	w.Putln("public ", name, "(){}\n")

	names, types := FlattenFields(fields)
	w.Put("public ", name, "(")
	w.PutParams(names, types)
	w.Putln("){")
	w.indent++
	for _, n := range names {
		w.Putln("this.", n, " = ", n, ";")
	}
	w.indent--
	w.Putln("}")

	w.Putln("public ", name, "(", name, " other", "){")
	w.indent++
	for _, n := range names {
		w.Putln("this.", n, " = ", "other.", n, ";")
	}
	w.indent--
	w.Putln("}")

	w.Putln("public ", name, " copy(){")
	w.indent++
	w.Put("return new ", name, "(this);")
	w.indent--
	w.Putln("}")

	// Methods on value
	for _, m := range d.valMethods {
		w.PutMethodDecl(m, true)
	}

	w.indent--
	w.Putln("}")
}

func (w *Writer) PutStructFields(fields *ast.FieldList) {
	names, types := FlattenFields(fields)
	for i, n := range names {
		t := types[i]
		w.Put(ModifierFor(n), " ")
		w.Put(JavaTypeOf(t))
		w.Putln(" ", n, " = ", ZeroValue(t), ";")
		// TODO Docs
	}
}

// 	type StructType struct {
// 	        Struct     token.Pos  // position of "struct" keyword
// 	        Fields     *FieldList // list of field declarations
// 	        Incomplete bool       // true if (source) fields are missing in the Fields list
// 	}
func (w *Writer) PutStructDef(def *TypeDef) {
}

func (w *Writer) PutMethodDecl(f *ast.FuncDecl, copyRecv bool) {

	// (1) Put static implementation with "this" as first receiver
	// TODO: some doc
	w.PutStaticFunc(f)
	w.Putln()

	// (2) Put method, calling static implementation
	w.PutDoc(f.Doc)
	w.Put(ModifierFor(f.Name), " ")

	// return type
	_, retTypes := FlattenFields(f.Type.Results)
	w.Put(JavaReturnTypeOf(retTypes), " ", f.Name)

	// arguments
	w.Put("(")
	argNames, argTypes := FlattenFields(f.Type.Params)
	w.PutParams(argNames, argTypes)
	w.Putln("){")
	w.indent++

	// body calls static implementation with this as first arg
	if len(retTypes) > 0 {
		w.Put("return ")
	}
	w.Put(f.Name, "(", "this")
	if copyRecv {
		w.Put(".copy()")
	}
	for i := range argNames {
		w.Put(", ", argNames[i])
	}
	w.Putln(");")

	w.indent--
	w.Putln("}")

}

//func ClassNameFor(typ ast.Expr) string {
//	switch typ := typ.(type) {
//	default:
//		Error(typ, "cannot handle", reflect.TypeOf(typ))
//		panic("")
//	case *ast.Ident:
//		return typ.Name // TODO: rename
//	}
//}

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
