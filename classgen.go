package gotojava

// Generate java classes based on type/method definitions.

import (
	"go/ast"
	"reflect"
)

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
	w.Putf("public final class %v extends %v {\n", name, base)
	w.indent++

	fields := spec.Type.(*ast.StructType).Fields
	w.PutConstructors(name, fields) // TODO: superconstructor

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
	w.PutConstructors(name, fields)

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
	fieldNames, _ := FlattenFields(fields)
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

func (w *Writer) PutConstructors(name string, fields *ast.FieldList) {
	// (1) no arguments
	w.Putln("\tpublic ", name, "(){}\n")

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
}

func (w *Writer) PutStructFields(fields *ast.FieldList) {
	names, types := FlattenFields(fields)
	for i, n := range names {
		t := types[i]
		w.Put(GlobalModifierFor(n), t, " ", n)
		if t.NeedsFinal() {
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
	w.Put(GlobalModifierFor(f.Name))

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
