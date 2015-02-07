package main

// Type conversion between Go and Java.

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/types"
)

// maps Go primitives to java
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

// JavaType returns the java type used to store the given go type. E.g.:
// 	bool   -> boolean
// 	uint32 -> int
func JavaType(goType types.Type) string {
	if j, ok := typeToJava[goType.Underlying().String()]; ok {
		return j
	} else {
		panic("cannot convert type to java: " + goType.String())
	}
}


// JavaTupleType returns the java type used to wrap a tuple of go types for multiple return values. E.g.:
// 	return 1, 2 -> return new Tuple_int_int(1, 2)
func JavaTupleType(types []types.Type) string {
	name := "Tuple"
	for _, t := range types {
		name += "_" + t.String() // not java name as we want to discriminate, e.g., int from uint
	}
	return name
}
