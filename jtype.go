package gotojava

// Type conversion from Go to Java.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// JType represents the java type for a Go identifier.
// Go types do not map to java types one-to-one. The java
// type for a certain Go type may also depend on the context
// (e.g. local variable whose address is taken.)
type JType struct {
	Orig  types.Type
	Ident *ast.Ident
	JName string
}

func JTypeOf(x ast.Expr) JType {

	t := JType{Orig: TypeOf(x)}
	if id, ok := x.(*ast.Ident); ok {
		t.Ident = id
	}

	if t.IsEscapedBasic() {
		panic("")
	} else {

	}
	return t
}

// Java return type for a function that returns given types.
// For multiple return types, a Tuple type is returned
func JavaReturnTypeOf(resultTypes []JType) JType {
	switch len(resultTypes) {
	case 0:
		return JType{Orig: nil, Ident: nil, JName: "void"}
	case 1:
		return resultTypes[0]
	default:
		panic("multiple returns")
		//return JavaTupleType(resultTypes)
	}
}

func (t JType) IsEscapedBasic() bool {
	_, basic := t.Orig.Underlying().(*types.Basic)
	return basic && t.Ident != nil && Escapes(t.Ident)
}

func (t JType) IsStructValue() bool {
	_, ok := t.Orig.Underlying().(*types.Struct)
	return ok
}

func (t JType) IsValue() bool {
	switch t := t.Orig.Underlying().(type) {
	default:
		panic("cannot handle " + reflect.TypeOf(t).String())
	case *types.Struct, *types.Basic:
		return true
	case *types.Pointer:
		return false
	}
}

func javaType(orig types.Type) string {
	switch orig := orig.(type) {
	default:
		panic("cannot handle type " + reflect.TypeOf(orig).String())
	case *types.Basic:
		return javaBasicType(orig)
	case *types.Named:
		return javaNamedType(orig)
	case *types.Pointer:
		return javaPointerType(orig)
	case *types.Signature:
		Log(nil, "TODO: Signature")
		return "**SIGNATURE**"
	}
}

func JavaPointerName(elemExpr ast.Expr) string {
	// TODO: what for pointer to already escaped basic?
	return JTypeOf(elemExpr).JName + "Ptr"
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
		return javaType(e) + "Ptr"
	case *types.Basic:
		return Export(javaBasicType(e) + "Ptr")
	}
}

// JavaTupleType returns the java type used to wrap a tuple of go types for multiple return values. E.g.:
// 	return 1, 2 -> return new Tuple_int_int(1, 2)
// Calling this function also ensure code for the tuple has been generated.
// TODO: JType
func JavaTupleType(types []JType) string {
	name := "Tuple"
	for _, t := range types {
		name += "_" + t.Orig.String() // not java name as we want to discriminate, e.g., int from uint
	}

	if !classGen[name] {
		GenTupleDef(name, types)
	}
	return name
}
