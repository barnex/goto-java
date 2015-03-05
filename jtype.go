package gotojava

// Type conversion from Go to Java.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// JType represents a Java type translated from Go.
// Go types do not map to java types one-to-one. The java
// type for a certain Go type may also depend on the context
// (e.g. local variable whose address is taken.)
type JType struct {
	Orig  types.Type // original Go type translated to Java
	Ident *ast.Ident // identifier having this type, if any (for escape analyis)
}

func JTypeOfExpr(x ast.Expr) JType {
	t := JTypeOfGoType(TypeOf(x))
	if id, ok := x.(*ast.Ident); ok {
		t.Ident = id
	}
	return t
}

func JTypeOfGoType(t types.Type) JType {
	return JType{
		Orig:  t,
		Ident: nil,
	}
}

func (t JType) JName() string {
	if t.Orig == nil && t.Ident == nil {
		return "void"
	}
	if t.IsEscapedPrimitive() {
		return EscapedBasicName(t)
	}
	if ptr, ok := t.Orig.(*types.Pointer); ok {
		return javaName(ptr.Elem().Underlying())
	}
	return javaName(t.Orig.Underlying())
}

func (t JType) ClassName() string {

}

func javaName(orig types.Type) string {

	switch orig := orig.(type) {
	default:
		panic("cannot handle type " + reflect.TypeOf(orig).String() + ":" + orig.String())
	case *types.Basic:
		return javaBasicName(orig)
	case *types.Struct:
		return javaStructName(orig)
	case *types.Named:
		return javaNamedName(orig)
	case *types.Pointer:
		return javaPointerName(orig)
		//case *types.Signature: // TODO
		//	return "**SIGNATURE**"
	}
}

func javaStructName(t *types.Struct) string {
	name := "Struct" // TODO: go.
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		name += "_" + javaName(f.Type()) + "_" + f.Name()
	}
	return name
}

// java type for Go pointer type. E.g.
// 	*int -> IntPtr
func javaPointerName(t *types.Pointer) string {
	return javaPointerNameForElem(t.Elem())
}

func javaPointerNameForElem(e types.Type) string {
	if IsPointer(e) {
		panic("pointer to pointer")
	}
	if IsPrimitive(e) {
		return Export(javaName(e)) // TODO: go., ...
	} else {
		return javaName(e)
	}
}

// Java name for named type.
func javaNamedName(t *types.Named) string {
	if IsPrimitive(t) {
		return javaName(t.Underlying())
	}

	obj := t.Obj()
	if r, ok := rename[obj]; ok {
		return r
	}

	return obj.Name()
}

// Java name for basic type.
func javaBasicName(t *types.Basic) string {
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

/*
struct{v int} -> class Struct0{ int v }
struct{v int} -> final Struct0
*struct{v int} -> Struct0
type S struct{v int} -> class S extends Struct0 implements ... { set(Struct0); get()Struct0; val methods only for iface}
type S struct{v int} -> class SPtr extends Struct0 implements ... { ptr methods only for iface}
S -> final Struct0
*S -> Struct0
s.meth() -> S.static(this, ...)
interface{...} x = s -> IntefaceXXX x = new S(Struct0)
interface{...} x = &s -> IntefaceXXX x = new SPtr(Struct0)
*/

func EscapedBasicName(t JType) string {
	return Export(javaBasicName(t.Orig.Underlying().(*types.Basic))) // TODO: go., ...
}

// Java return type for a function that returns given types.
// For multiple return types, a Tuple type is returned
func JavaReturnTypeOf(resultTypes []JType) JType {
	switch len(resultTypes) {
	case 0:
		return JType{Orig: nil, Ident: nil}
	case 1:
		return resultTypes[0]
	default:
		panic("multiple returns")
		//return JavaTupleType(resultTypes)
	}
}

// Returns whether t represents a Go basic type that was moved to the heap. E.g.:
//  i := 0
// 	x := &i
// 	IsEscapedBasic(JTypeOf(i)) // true
// In such case the java primitive (e.g. int) is replaced by a wrapper (e.g. go.Int)
func (t JType) IsEscapedPrimitive() bool {
	return IsPrimitive(t.Orig) && t.Ident != nil && Escapes(t.Ident)
}

func (t JType) IsNamedPrimitive() bool {
	return IsPrimitive(t.Orig) && t.IsNamed()
}

func IsPrimitive(t types.Type) bool {
	_, basic := t.Underlying().(*types.Basic)
	return basic
}

func (t JType) IsNamed() bool {
	_, named := t.Orig.(*types.Named)
	return named
}

func (t JType) IsPointerToPrimitive() bool {
	if ptr, ok := t.Orig.(*types.Pointer); ok {
		_, ok := ptr.Elem().Underlying().(*types.Basic)
		return ok
	} else {
		return false
	}
}

func IsPointer(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
}

// Returns whether t represents a Go struct type (value semantics).
// This affects assignment and equality:
// 	a = b    ->    a.set(b)
// 	a == b   ->    a.equals(b)
func (t JType) IsStructValue() bool {
	_, ok := t.Orig.Underlying().(*types.Struct)
	return ok
}

// Returns whether java variables of this type should be declared final.
// E.g.: an escaped basic type, so we can access it from an inner class (e.g. closure).
func (t JType) NeedsFinal() bool {
	return t.IsEscapedPrimitive() || t.IsStructValue()
}

func (t JType) NeedsSetMethod() bool {
	return t.NeedsFinal()
}

func (t JType) NeedsEqualsMethod() bool {
	return t.NeedsFinal()
}

// Returns whether t is underlied by a Go value type (as opposed to pointer type).
// In that case, pointer methods on t need insertion of an address-of.
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
