package main

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/types"
)

// maps go primitives to java
var typeToJava = map[string]string{
	"bool":        "boolean",
	"byte":        "byte",
	"float32":     "float",
	"float64":     "double",
	"int":         "int", //?
	"int16":       "short",
	"int32":       "int",
	"int64":       "long",
	"int8":        "byte",
	"interface{}": "Object",
	"string":      "String", //?
	"uint":        "int",    //?
	"uint16":      "short",  //?
	"uint32":      "int",    //?
	"uint64":      "long",   //?
	"uint8":       "byte",   //?
}

// explicit type cast in input file, e.g.:
// 	a := int(b)
func (w *writer) PutTypecast(goType string, e ast.Expr) {
	jType, ok := typeToJava[goType]
	if !ok {
		Error(e, "cannot convert to java:", goType)
	}
	w.Put("(", jType, ")(", e, ")")
}

// implicit type cast from untyped to type, e.g.:
// 	f(1)
func (w *writer) PutImplicitCast(dst types.Type, e ast.Expr) {

	dst = dst.Underlying()
	src := TypeOf(e).Underlying()

	log.Println(src, "->", dst)

	if dst.String() == "interface{}" {
		w.PutEmptyInterfaceCast(e)
		return
	}

	w.PutExpr(e)
}

func TypeOf(n ast.Expr) types.Type {
	t := info.TypeOf(n)
	if t == nil {
		Error(n, "cannot infer type")
	}
	return t
}

func (w *writer) TypeToJava(goType types.Type) string {
	return w.typeToJava(goType.String())

	// remove untyped.
	//if strings.HasPrefix(ident, "untyped ") {
	//	ident = ident[len("untyped "):]
	//}

}

func (w *writer) typeToJava(goType string) string {
	if j, ok := typeToJava[goType]; ok {
		return j
	}
	panic("cannot convert type to java: " + goType)
}

// ObjectOf returns the object denoted by the specified identifier.
func ObjectOf(id *ast.Ident) types.Object {
	obj := info.ObjectOf(id)
	if obj == nil {
		Error(id, "undefined:", id.Name)
	}
	return obj
}

// returun exact value and minimal type for constant expression.
func ExactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = info.Types[e]
	return
}
