package gotojava

// Generate zero values (initial values) for each type.

import (
	"fmt"
	"reflect"

	"go/ast"

	"golang.org/x/tools/go/types"
)

func InitValue(rhs ast.Expr, typ JType) interface{} {
	if typ.IsPrimitive() {
		return RValue(rhs)
	} else {
		v := "new " + typ.JName() + "("

		if typ.NeedsAddress() {
			v += fmt.Sprint(FakeAddressFor(typ.Ident), ", ")
		}

		v += Transpile(RValue(rhs)) + ")"
		return v
	}
}

// ZeroValue returns the zero value for a new variable of java type jType.
// E.g.:
// 	var x int  ->  int x = 0;
func ZeroValue(t JType) interface{} {
	if t.IsPrimitive() {
		return basicZeroValue(t.Orig.Underlying().(*types.Basic))
	} else {
		v := "new " + t.JName() + "("
		if t.NeedsAddress() {
			v += fmt.Sprint(FakeAddressFor(t.Ident))
		}
		v += ")"
		return v
	}
}

func namedZeroValue(t *types.Named) string {
	switch u := t.Underlying().(type) {
	default:
		panic("cannot make zero value for named " + reflect.TypeOf(u).String())
	case *types.Struct:
		// always load Struct with StructPtr (extends Struct),
		// so we can take address of value by typecasting.
		return "new " + JTypeOfGoType(t).JName() + "()"
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
