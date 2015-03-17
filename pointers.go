package gotojava

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/types"
)

func (w *Writer) PutAddressOf(x ast.Expr) {

	elemT := TypeOf(x).Underlying()

	// Sole exception:
	// global primitive (not wrapped like escaping local)
	if _, ok := elemT.(*types.Basic); ok {
		if id, ok := x.(*ast.Ident); ok && !isLocal(id) {
			w.PutAddressOfGlobalBasic(id)
			return
		}
	}

	switch x := x.(type) {
	default:
		class := javaPointerNameForElem(elemT)
		w.PutNew(class, x)
	}

}

func (w *Writer) PutAddressOfGlobalBasic(x *ast.Ident) {
	elemT := TypeOf(x).Underlying()
	class := javaPointerNameForElem(elemT)
	prim := JTypeOfGoType(elemT).JName()
	lValueT := "LValue_" + prim

	w.Putln("new ", class, "(new ", lValueT, "(){")
	w.indent++ //?
	w.Putln("public ", prim, " value(){ return ", x, ";}")
	w.Putln("public void set(", prim, " v){ ", x, " = v ;}")
	w.Putln("public ", prim, " addr(){ return ", FakeAddressFor(x), ";}")
	w.indent--
	w.Put("})")
}

var (
	fakeAddr            = map[types.Object]int{}
	fakeAddrCounter int = 1
)

func FakeAddressFor(x *ast.Ident) string {
	obj := ObjectOf(x)
	if addr, ok := fakeAddr[obj]; ok {
		return fmt.Sprint(addr)
	} else {
		fakeAddr[obj] = fakeAddrCounter
		fakeAddrCounter++
		return fmt.Sprint(fakeAddr[obj])
	}
}

//func (w *Writer) putAddressOf(id *ast.Ident) {
//}
//
////func (w *Writer) putAddressOfStruct(x ast.Expr) {
////	w.Put("new ", javaPointerNameForElem(TypeOf(x)), "(", x, ")")
////}

func (w *Writer) PutStarExpr(x *ast.StarExpr) {
	w.Put(x.X, ".value()")
}
