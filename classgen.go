package gotojava

// Generate java classes based on type/method definitions.

import (
	"go/ast"
	"log"
	"reflect"

	"golang.org/x/tools/go/types"
)

// generate code for all defs in global typedefs variable
func GenClasses() {
	for _, t := range alltypes {
		switch t := t.(type) {
		default:
			panic(reflect.TypeOf(t).String())
		case *types.Basic:
			// nothing to do, handwritten.
		case *types.Pointer:
			genPointer(t)
		case *types.Named:
			genNamed(t)
		case *types.Signature:
			genSignature(t)
		case *types.Struct:
			genStruct(t)
		case *types.Tuple:
			genTuple(t)
		}
	}
}

func genNamed(t *types.Named) {
	w := OpenClass(javaName(t), "")
	defer w.CloseClass()
}

func genPointer(t *types.Pointer) {
	w := OpenClass(javaName(t), "")
	defer w.CloseClass()
}

func genSignature(t *types.Signature) {
	w := OpenClass(javaName(t), "")
	defer w.CloseClass()
}

func genTuple(t *types.Tuple) {
	w := OpenClass(javaName(t), "")
	defer w.CloseClass()
}

func genStruct(t *types.Struct) {
	name := javaName(t)
	w := OpenClass(name, "")
	w.PutEmptyConstructor(name)
	w.PutCompositeConstructor(name)
	defer w.CloseClass()
}

// Open name.java and already write class signature.
// Then, Writer methods can put class body.
// Finally, CloseClass() should be called.
func OpenClass(name, extends string, implements ...string) *Writer {
	log.Println("Generate class ", name)
	w := NewWriterFile(name + ".java")

	mod := PUBLIC
	w.Put(mod, "class ", name)
	if extends != "" {
		w.Put(" extends ", extends)
	}
	if len(implements) != 0 {
		w.Put(" implements")
		for _, x := range implements {
			w.Put(" ", x)
		}
	}
	w.Putln("{\n")
	w.indent++
	return w
}

func (w *Writer) CloseClass() {
	w.indent--
	w.Putln("}")
	w.Close()
}

func (w *Writer) PutEmptyConstructor(name string) {
	w.Putln()
	w.Putln("public ", name, "(){}\n")
}

func (w *Writer) PutCompositeConstructor() {
	// all fields
	if len(c.FieldNames) > 0 {
		w.Putln()
		w.Put("public ", c.Name, "(")

		//w.PutParams(c.FieldNames, c.FieldTypes)
		for i := range c.FieldNames {
			w.Put(comma(i), c.FieldTypes[i], " ", c.FieldNames[i])
		}

		w.Putln("){")
		w.indent++
		for i, n := range c.FieldNames {
			t := c.FieldTypes[i]
			w.PutJAssign(t, Transpile("this.", n), t, n)
			w.Putln(";")
		}
		w.indent--
		w.Putln("}")
	}
}

//func flattenFields(t fieldser) (names []string, types []JType) {
//	names = make([]string, 0, t.NumFields())
//	types = make([]JType, 0, t.NumFields())
//	for i := 0; i < t.NumFields(); i++ {
//		names = append(names, t.Field(i).Name())                // TODO: rename
//		types = append(types, JTypeOfGoType(t.Field(i).Type())) // TODO: rename
//	}
//	return names, types
//}

//// Generate java class for Go named struct type (value semantics).
//func GenStructValueClass(d *TypeDef) {
//	spec := d.typeSpec
//	names, types := FlattenFields(spec.Type.(*ast.StructType).Fields)
//
//	def := &ClassDef{
//		Name:       JTypeOfExpr(spec.Name).JName(),
//		Doc:        spec.Doc,
//		Extends:    JTypeOfGoType(TypeOf(spec.Type).Underlying()).JName(),
//		FieldNames: names,
//		FieldTypes: types,
//		Methods:    append(d.valMethods, d.ptrMethods...),
//	}
//	def.Gen()
//}
//
//// Generate java class for Go pointer-to-named-struct type.
//func GenStructPointerClass(d *TypeDef) {
//	spec := d.typeSpec
//	//fields := spec.Type.(*ast.StructType).Fields
//
//	def := &ClassDef{
//		Name:      javaPointerNameForElem(TypeOf(spec.Name)),
//		ValueName: JTypeOfExpr(spec.Name).JName(),
//		Methods:   d.ptrMethods,
//	}
//
//	def.Gen(
//		(*ClassDef).genFields,
//		(*ClassDef).genPtrConstructor,
//		(*ClassDef).genMethods)
//}
//
////	// copy method
////	w.Putf(`
////	public %s copy(){
////		return new %s(this);
////	}
////`, name, name)
////
////	// equals method
////	w.Putf(`
////	/** @Override
////		Deep equality test of all fields. */
////	public boolean equals(Object o){
////		if (o instanceof %v){
////			%v other = (%v)o;
////`, name, name, name)
////	w.indent += 2
////	if len(fieldNames) == 0 {
////		w.Put("return true") // struct{}{} == struct{}{}
////	} else {
////		w.Put("return ")
////		for i, n := range fieldNames {
////			if i > 0 {
////				w.Putln(" &&")
////				w.Put("\t")
////			}
////			w.PutJEquals(JTypeOfExpr(n), Transpile("this.", n), JTypeOfExpr(n), Transpile("other.", n))
////		}
////	}
////	w.Putln(";")
////	w.indent--
////	w.Putln(`} else {
////			return false;
////		}`)
////	w.indent--
////	w.Putln("}")
////
////	// todo hashCode
////
////	w.indent--
////	w.Putln("}")

//func (c *ClassDef) Gen(f ...func(*ClassDef, *Writer)) {
//
//	w := NewWriterFile(c.Name + ".java")
//	defer w.Close()
//
//	c.genSignature(w)
//
//	for _, f := range f {
//		f(c, w)
//	}
//
//	w.indent--
//	w.Putln("}")
//}

//func (c *ClassDef) genFields(w *Writer) {
//	w.Putln()
//	for i, n := range c.FieldNames {
//		t := c.FieldTypes[i]
//		mod := NONE // TODO: GlobalModifierFor(n)
//		w.Put(mod, t, " ", n)
//		if t.NeedsFinal() {
//			w.Put(" = ", ZeroValue(t))
//		}
//		w.Putln(";")
//		// TODO Docs
//	}
//}
//
//func (c *ClassDef) genMethods(w *Writer) {
//	for _, m := range c.Methods {
//		w.PutMethodDecl(m, false)
//	}
//}
//
//func (c *ClassDef) genEmptyConstructor(w *Writer) {
//	w.Putln()
//	w.Putln("public ", c.Name, "(){}\n")
//}
//
//func (c *ClassDef) genCopyConstructor(w *Writer) {
//	w.Putln()
//	w.Putln("public ", c.Name, "(", c.Name, " other){")
//	w.indent++
//	for i, n := range c.FieldNames {
//		t := c.FieldTypes[i]
//		w.PutJAssign(t, Transpile("this.", n), t, Transpile("other.", n))
//		w.Putln(";")
//	}
//	w.indent--
//	w.Putln("}")
//}
//
//
//func (c *ClassDef) genPtrConstructor(w *Writer) {
//	w.Putln()
//	w.Putln("public ", c.Name, "(", c.ValueName, " other){")
//	w.indent++
//	w.indent--
//	w.Putln("}")
//}
//
//func (c *ClassDef) genSetMethod(w *Writer) {
//	w.Putln()
//	w.Put("public void set(", c.Name, "  other){")
//	w.indent++
//	for i, n := range c.FieldNames {
//		t := c.FieldTypes[i]
//		w.PutJAssign(t, Transpile("this.", n), t, Transpile("other.", n))
//		w.Putln(";")
//	}
//	w.indent--
//	w.Putln("}")
//}

//func genBasicClass(d *TypeDef) {
//	valueDef := &ClassDef{
//		Name: JTypeOfExpr(d.typeSpec.Name).JName(),
//	}
//	valueDef.Gen()
//
//	ptrDef := &ClassDef{
//		Name: javaPointerNameForElem(TypeOf(d.typeSpec.Name)),
//	}
//	ptrDef.Gen()
//}
//
//func (w *Writer) PutConstructors(name string, fields *ast.FieldList) {
//	// (1) no arguments
//	w.Putln("\tpublic ", name, "(){}\n")
//
//	// (2) fields as individual values
//	fieldNames, fieldTypes := FlattenFields(fields)
//	if len(fieldNames) > 0 {
//		w.Put("public ", name, "(")
//		w.PutParams(fieldNames, fieldTypes)
//		w.Putln("){")
//		w.indent++
//		for i, n := range fieldNames {
//			t := fieldTypes[i]
//			w.PutJAssign(t, Transpile("this.", n), t, n)
//			w.Putln(";")
//		}
//		w.indent--
//		w.Putln("}")
//	}
//}
//
//func (w *Writer) PutStructFields(fields *ast.FieldList) {
//	names, types := FlattenFields(fields)
//	for i, n := range names {
//		t := types[i]
//		w.Put(GlobalModifierFor(n), t, " ", n)
//		if t.NeedsFinal() {
//			w.Put(" = ", ZeroValue(t))
//		}
//		w.Putln(";")
//		// TODO Docs
//	}
//}
//
func (w *Writer) PutMethodDecl(f *ast.FuncDecl, copyRecv bool) {

	// (1) Put static implementation with "this" as first receiver
	w.Putf("\t/** Implementation for method %v, with receiver as first argument. */\n", f.Name)
	w.PutFunc(STATIC, f)
	w.Putln()

	// (2) Put method, calling static implementation
	w.PutDoc(f.Doc)
	w.Put(GlobalModifierFor(f.Name))

	// return type
	//_, retTypes := FlattenFields(f.Type.Results)
	returnType := TypeOf(f.Type).(*types.Signature).Results() //TODO: putSignature?
	w.Put(returnType, " ", f.Name)

	// arguments
	w.Put("(")
	argNames, argTypes := FlattenFields(f.Type.Params)
	w.PutParams(argNames, argTypes)
	w.Putln("){")
	w.indent++

	// body calls static implementation with this as first arg
	if returnType.Len() > 0 {
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
