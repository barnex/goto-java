package gotojava

import (
	"fmt"
	"reflect"

	"golang.org/x/tools/go/types"
)

func javaName(orig types.Type) string {
	switch orig := orig.(type) {
	default:
		panic(reflect.TypeOf(orig).String() + ":" + orig.String())
	case *types.Array:
		return javaArrayName(orig)
	case *types.Basic:
		return javaBasicName(orig)
	case *types.Struct:
		return javaStructName(orig)
	case *types.Map:
		return javaMapName(orig)
	case *types.Named:
		return javaNamedName(orig)
	case *types.Pointer:
		return javaPointerName(orig)
	case *types.Interface:
		return javaInterfaceName(orig)
	case *types.Slice:
		return javaSliceName(orig)
	case *types.Signature:
		return javaSignatureName(orig)
	case *types.Tuple:
		return javaTupleName(orig)
	}
}

const (
	ARRAY_PREFIX     = "Array"
	FUNC_PREFIX      = "Func"
	INTERFACE_PREFIX = "Interface"
	MAP_PREFIX       = "Map"
	POINTER_PREFIX   = "Ptr"
	SLICE_PREFIX     = "Slice"
	STRUCT_PREFIX    = "Struct"
)

func javaMapName(t *types.Map) string {
	return MAP_PREFIX + "_" + javaName(t.Key()) + "_" + javaName(t.Elem())
}

func javaSliceName(t *types.Slice) string {
	return fmt.Sprint(SLICE_PREFIX, "_", javaName(t.Elem()))
}

func javaArrayName(t *types.Array) string {
	return fmt.Sprint(ARRAY_PREFIX, "_", t.Len(), "_", javaName(t.Elem()))
}

func javaTupleName(t *types.Tuple) string {
	if t == nil {
		return "void"
	}

	if t.Len() == 1 {
		return javaName(t.At(0).Type())
	}

	name := STRUCT_PREFIX
	for i := 0; i < t.Len(); i++ {
		v := t.At(i)
		name += fmt.Sprint("_", javaName(v.Type()), "_v", i)
	}
	return name
}

func javaSignatureName(t *types.Signature) string {
	return FUNC_PREFIX + javaName(t.Params()) + "_to_" + javaName(t.Results())
}

func javaInterfaceName(t *types.Interface) string {
	if t.NumMethods() == 0 {
		return "Object"
	}
	name := INTERFACE_PREFIX
	for i := 0; i < t.NumMethods(); i++ {
		m := t.Method(i)
		name += "_" + javaName(m.Type()) + "_" + m.Name()
	}

	return name
}

func javaStructName(t *types.Struct) string {
	name := STRUCT_PREFIX
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
	return POINTER_PREFIX + "_" + javaName(e)
}

// Java name for named type.
func javaNamedName(t *types.Named) string {
	//if IsPrimitive(t) {
	//	return javaName(t.Underlying())
	//}

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

func EscapedBasicName(t *types.Basic) string {
	return Export(javaBasicName(t))
}

func javaWrapperName(t types.Type) string {
	if t, ok := t.(*types.Basic); ok {
		return EscapedBasicName(t)
	} else {
		return javaName(t)
	}
}
