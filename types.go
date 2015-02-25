package gotojava

// Type conversion between Go and Java.

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

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
	}
}

func JavaTypeOf(typeExpr ast.Expr) string {
	return JavaType(TypeOf(typeExpr))
}

func JavaTypeOfPtr(elemExpr ast.Expr) string {
	return JavaTypeOf(elemExpr) + "Ptr"
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
	switch e := t.Elem().(type) {
	default:
		panic("cannot handle pointer to " + reflect.TypeOf(e).String())
	case *types.Named, *types.Pointer:
		return JavaType(e) + "Ptr"
	case *types.Basic:
		return Export(javaBasicType(e) + "Ptr")
	}
}

// explicit type cast in input file, e.g.:
// 	a := int(b)
func (w *Writer) PutTypecast(goType types.Type, e ast.Expr) {
	panic("no typecast yet")
}

// Emit code for rhs, possibly converting to make it assignable to lhs.
func (w *Writer) PutRHS(rhs ast.Expr, lhs types.Type, inmethod bool) {
	w.PutExpr(rhs)
}
