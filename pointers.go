package gotojava

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

func (w *Writer) PutAddressOf(x ast.Expr) {
	switch t := TypeOf(x).Underlying().(type) {
	default:
		panic("cannot take address of " + reflect.TypeOf(t).String())
	case *types.Basic:
		w.putAddressOfBasic(x)
	case *types.Struct:
		w.putAddressOfStruct(x)
	}
}

func (w *Writer) putAddressOfBasic(x ast.Expr) {
	switch x := x.(type) {
	default:
		panic("cannot take address of " + reflect.TypeOf(x).String())
	case *ast.Ident:
		if IsLocal(x) {
			w.putAddressOfLocal(x)
		} else {
			w.putAddressOfGlobal(x)
		}
	}
}

// Checks whether the identifier is locally defined.
// E.g., for escape analysis: address of local variable etc.
func IsLocal(id *ast.Ident) bool {
	scope := ObjectOf(id).Parent()
	global := scope.Parent() == types.Universe
	return !global
}

func (w *Writer) putAddressOfLocal(id *ast.Ident) {
	w.putAddressOfStruct(id) // id implemented as struct value
}

func (w *Writer) putAddressOfGlobal(id *ast.Ident) {
	// TODO: indent
	w.Putf(`new go.IntPtr(){
			public int value(){return %s;}
			public void set(int v){%s = v;} // TODO: name collision
		}`, id, id)
}

func (w *Writer) putAddressOfStruct(x ast.Expr) {
	w.Put(x, ".addr()") // TODO: ".addr" -> const
}

func (w *Writer) PutStarExpr(x *ast.StarExpr) {
	w.Put(x.X, ".value()")
}
