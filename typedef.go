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

	name := javaPointerNameForElem(TypeOf(spec.Name))
	base := JTypeOf(spec.Name)

	w := NewWriterFile(name + ".java")
	defer w.Close()

	w.Putf("/** %v extends %v with pointer methods. */\n", name, base)
	w.Putf(`public final class %v extends %v {`, name, base)
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
	name := JTypeOf(spec.Name).JName
	ptrname := javaPointerNameForElem(TypeOf(spec.Name))
	w := NewWriterFile(name + ".java")
	defer w.Close()

	w.PutDoc(spec.Doc)
	w.Putln("public class ", name, "{")
	w.Putln()
	w.indent++

	// Fields
	fields := spec.Type.(*ast.StructType).Fields
	w.PutStructFields(fields)
	w.Putln()

	// Constructors:
	// (1) no arguments
	w.Putln("public ", name, "(){}\n")

	// (2) fields as individual values
	fieldNames, fieldTypes := FlattenFields(fields)
	if len(fieldNames) > 0 {
		w.Put("public ", name, "(")
		w.PutParams(fieldNames, fieldTypes)
		w.Putln("){")
		w.indent++
		for i, n := range fieldNames {
			t := fieldTypes[i]
			w.PutJAssign(t, Transpile("this.", n), t, n)
			w.Putln(";")
		}
		w.indent--
		w.Putln("}")
	}

	// (3) copy constructor
	w.Putf(`
	public %s(%s other){
		this.set(other);
	}
`, name, name)

	// Methods on value
	for _, m := range d.valMethods {
		w.PutMethodDecl(m, true)
	}

	// TODO: override some for PtrType, panic if they should not be called.

	// copy method
	w.Putf(`
	public %s copy(){
		return new %s(this);
	}
`, name, name)

	// addr method
	w.Putf(`
	public %s addr(){
		return (%s)this;
	}
`, ptrname, ptrname)

	// set method
	w.Putf(`
	public void set(%v other){
`, name)
	w.indent++
	for _, n := range fieldNames {
		w.PutJAssign(JTypeOf(n), Transpile("this.", n), JTypeOf(n), Transpile("other.", n))
		w.Putln(";")
	}
	w.indent--
	w.Putln("}")

	// equals method
	w.Putf(`
	/** @Override
		Deep equality test of all fields. */
	public boolean equals(Object o){
		if (o instanceof %v){	
			%v other = (%v)o;
`, name, name, name)
	w.indent += 2
	if len(fieldNames) == 0 {
		w.Put("return true") // struct{}{} == struct{}{}
	} else {
		w.Put("return ")
		for i, n := range fieldNames {
			if i > 0 {
				w.Putln(" &&")
				w.Put("\t")
			}
			w.PutJEquals(JTypeOf(n), Transpile("this.", n), JTypeOf(n), Transpile("other.", n))
		}
	}
	w.Putln(";")
	w.indent--
	w.Putln(`} else {
			return false;
		}`)
	w.indent--
	w.Putln("}")

	// todo hashCode

	w.indent--
	w.Putln("}")
}

func (w *Writer) PutStructFields(fields *ast.FieldList) {
	names, types := FlattenFields(fields)
	for i, n := range names {
		t := types[i]
		w.Put(ModifierFor(n), t, " ", n)
		if ModifierFor(n)&FINAL != 0 {
			w.Put(" = ", ZeroValue(t))
		}
		w.Putln(";")
		// TODO Docs
	}
}

func (w *Writer) PutMethodDecl(f *ast.FuncDecl, copyRecv bool) {

	// (1) Put static implementation with "this" as first receiver
	w.Putf("\t/** Implementation for method %v, with receiver as first argument. */\n", f.Name)
	w.PutFunc(STATIC, f)
	w.Putln()

	// (2) Put method, calling static implementation
	w.PutDoc(f.Doc)
	w.Put(ModifierFor(f.Name))

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
