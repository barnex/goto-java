package gotojava

// Generate zero values (initial values) for each type.

import (
	"reflect"

	"go/ast"

	"golang.org/x/tools/go/types"
)

func InitValue(rhs ast.Expr, typ JType) interface{} {
	if typ.IsEscapedPrimitive() {
		return "new " + typ.JName + "(" + Transpile(RValue(rhs)) + ")"
	} else {
		return RValue(rhs)
	}
}

// ZeroValue returns the zero value for a new variable of java type jType.
// E.g.:
// 	var x int  ->  int x = 0;
func ZeroValue(t JType) interface{} {
	if t.IsEscapedPrimitive() {
		return "new " + t.JName + "()"
	}
	switch typ := t.Orig.(type) {
	default:
		panic("cannot make zero value for " + reflect.TypeOf(typ).String() + ":" + t.JName)
	case *types.Basic:
		return basicZeroValue(typ)
	case *types.Named:
		return namedZeroValue(typ)
	case *types.Pointer:
		return "null"
	case *types.Struct:
		return "new " + t.JName + "()"
	}
}

func namedZeroValue(t *types.Named) string {
	switch u := t.Underlying().(type) {
	default:
		panic("cannot make zero value for named " + reflect.TypeOf(u).String())
	case *types.Struct:
		// always load Struct with StructPtr (extends Struct),
		// so we can take address of value by typecasting.
		return "new " + javaPointerNameForElem(t) + "()"
	case *types.Basic:
		return basicZeroValue(u)
	}
}

func basicZeroValue(t *types.Basic) string {
	info := t.Info()
	switch {
	default:
		panic("cannot make zero value for basic type " + t.String())
	case info&types.IsBoolean != 0:
		return "false"
	case info&types.IsFloat != 0:
		return "0.0"
	case info&types.IsInteger != 0:
		return "0"
	case info&types.IsString != 0:
		return `""`
	}
}
