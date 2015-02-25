package gotojava

// Generate zero values (initial values) for each type.

import (
	"reflect"

	"golang.org/x/tools/go/types"
)

// ZeroValue returns the zero value for a new variable of java type jType.
// E.g.:
// 	var x int  ->  int x = 0;
func ZeroValue(typ types.Type) string {

	switch typ := typ.(type) {
	default:
		panic("cannot make zero value for " + reflect.TypeOf(typ).String())
	case *types.Basic:
		return basicZeroValue(typ)
	case *types.Named:
		return namedZeroValue(typ)
	case *types.Pointer:
		return "null"
	}
}

func namedZeroValue(t *types.Named) string {
	switch u := t.Underlying().(type) {
	default:
		panic("cannot make zero value for named " + reflect.TypeOf(u).String())
	case *types.Struct:
		return "new " + JavaType(t) + "Ptr" + "()"
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
