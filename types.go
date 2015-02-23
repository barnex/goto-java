package gotojava

// Type conversion between Go and Java.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

func JavaTypeOfExpr(typeExpr ast.Expr) string {
	return JavaType(TypeOf(typeExpr))
}

func JavaTypeOfPtr(elemExpr ast.Expr) string {
	return JavaTypeOfExpr(elemExpr) + "Ptr"
}

// Array, Basic, Chan, Signature, Interface, Map, Named, Pointer, Slice, Struct, Tuple
// TODO: pass expr, if *ident: rename?  Or rename types in pre-processing?
func JavaType(t types.Type) string {
	switch t := t.(type) {
	default:
		panic("cannot handle type " + reflect.TypeOf(t).String())
	case *types.Basic:
		return javaBasicType(t)
	case *types.Named:
		return javaNamedType(t)
	case *types.Pointer:
		return javaPointerType(t)
		//case *types.Struct:
		//	return javaStructType(t)
	}
	panic("")
}

func javaBasicType(t *types.Basic) string {
	if transl, ok := builtin2java[t.String()]; ok {
		return transl
	} else {
		panic("cannot handle basic type " + t.String())
	}
}

func javaNamedType(t *types.Named) string {
	obj := t.Obj()
	if r, ok := rename[obj]; ok {
		return r
	} else {
		return obj.Name()
	}
}

func javaPointerType(t *types.Pointer) string {
	switch e := t.Elem().(type) {
	default:
		panic("cannot handle pointer to " + reflect.TypeOf(e).String())
	case *types.Named:
		return javaNamedType(e) + "Ptr" // TODO
	}
}

func javaStructType(t *types.Struct) string {
	panic("")
}

func (w *Writer) PutTypeExpr(typ ast.Expr) {
	w.Put(JavaTypeOfExpr(typ))
}

// explicit type cast in input file, e.g.:
// 	a := int(b)
func (w *Writer) PutTypecast(goType types.Type, e ast.Expr) {
	//Error(e, "TODO: typecast")

	w.PutImplicitCast(goType, e)

	//jType, ok := typeToJava[goType]
	//if !ok {
	//	Error(e, "cannot convert to java:", goType)
	//}
	//w.Put("(", jType, ")(", e, ")")
}

// implicit type cast from untyped to type, e.g.:
// 	f(1)
func (w *Writer) PutImplicitCast(dst types.Type, e ast.Expr) {
	//Error(e, "TODO: typecast")
	//dst = dst.Underlying()
	//src := TypeOf(e).Underlying()
	//log.Println(src, "->", dst)

	//if dst.String() == "interface{}" {
	//	w.PutEmptyInterfaceCast(e)
	//	return
	//}

	w.PutExpr(e)
}

// JavaType returns the java type used to store the given go type. E.g.:
// 	bool   -> boolean
// 	uint32 -> int
//func JavaType(goType types.Type) string {
//	if j, ok := typeToJava[goType.Underlying().String()]; ok {
//		return j
//	} else {
//		panic("cannot convert type to java: " + goType.String())
//	}
//}
