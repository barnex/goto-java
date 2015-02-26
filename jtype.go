package gotojava

// Type conversion from Go to Java.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// JType holds a java type name and go type it represents
type JType struct {
	GoType   types.Type
	JavaName string
}

//func(t JType)String()string{return t.JavaName}

func (t JType) IsStructValue() bool {
	_, ok := t.GoType.Underlying().(*types.Struct)
	return ok
}

func (t JType) IsValue() bool {
	switch t := t.GoType.Underlying().(type) {
	default:
		panic("cannot handle " + reflect.TypeOf(t).String())
	case *types.Struct, *types.Basic:
		return true
	case *types.Pointer:
		return false
	}
}

func JavaType(t types.Type) JType {
	switch t := t.(type) {
	default:
		panic("cannot handle type " + reflect.TypeOf(t).String())
	case *types.Basic:
		return JType{t, javaBasicType(t)}
	case *types.Named:
		return JType{t, javaNamedType(t)}
	case *types.Pointer:
		return JType{t, javaPointerType(t)}
	case *types.Signature:
		Log(nil, "TODO: Signature")
		return JType{t, "**SIGNATURE**"}
	}
}

func JavaTypeOf(typeExpr ast.Expr) JType {
	return JavaType(TypeOf(typeExpr))
}

func JavaPointerName(elemExpr ast.Expr) string {
	return JavaTypeOf(elemExpr).JavaName + "Ptr"
}

// Java return type for a function that returns given types.
// For multiple return types, a Tuple type is returned
func JavaReturnTypeOf(resultTypes []types.Type) JType {
	switch len(resultTypes) {
	case 0:
		return JType{nil, "void"}
	case 1:
		return JavaType(resultTypes[0])
	default:
		panic("multiple returns")
		//return JavaTupleType(resultTypes)
	}
}

func javaBasicType(t *types.Basic) string {
	// remove "untyped "
	name := t.Name()
	if t.Info()&types.IsUntyped != 0 {
		name = name[len("untyped "):]
	}

	if transl, ok := builtin2java[name]; ok {
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

// java type for Go pointer type. E.g.
// 	*int -> IntPtr
func javaPointerType(t *types.Pointer) string {
	return javaPointerNameForElem(t.Elem())
}

func javaPointerNameForElem(e types.Type) string {
	switch e := e.(type) {
	default:
		panic("cannot handle pointer to " + reflect.TypeOf(e).String())
	case *types.Named, *types.Pointer:
		return JavaType(e).JavaName + "Ptr"
	case *types.Basic:
		return Export(javaBasicType(e) + "Ptr")
	}
}
