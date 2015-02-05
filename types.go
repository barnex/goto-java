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

func JavaType(goType types.Type) string {
	if j, ok := typeToJava[goType.Underlying().String()]; ok {
		return j
	} else {
		panic("cannot convert type to java: " + goType.String())
	}
}
