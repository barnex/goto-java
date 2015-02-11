package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"

	"golang.org/x/tools/go/types"
)

// TODO: global package, use for class gen unless overridden.
var (
	fset     *token.FileSet          // accessed through PosOf
	info     types.Info              // accessed through TypeOf, ObjectOf, ExactValue
	parent   map[ast.Node]ast.Node   // accessed through ParentOf
	idents   map[string]int          // holds all identifier names and a counter to create a new, non-conflicting name if needed.
	rename   map[types.Object]string // maps some objects (typ. identifiers) to a new name for java.
	typedefs map[types.Object]*TypeDef
)

func HandleFile(fname string) {

	// (1) Parse
	var f *ast.File
	fset, f = ParseFile(fname)
	if *flagPrint {
		ast.Print(fset, f)
	}

	// (2) Determine types
	info = TypeCheck(fname, fset, f)

	// (3) Pre-processing: collect parents and declarations
	parent = CollectParents(f)
	idents = CollectIdents(f)
	rename = RenameReservedIdents(f)
	//typedefs = CollectDefs(f)

	// transpile primary class
	className := fname[:len(fname)-len(path.Ext(fname))]
	w := NewWriter(className + ".java")
	defer w.Close()
	w.PutClass(className, f)

	// generate additional classes from type/method definitions encountered along the way
	GenClasses()
}

func ParseFile(fname string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	checkUserErr(err)
	return fset, f
}

func TypeCheck(fname string, fset *token.FileSet, f *ast.File) types.Info {
	var config types.Config
	info := types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	_, err := config.Check(fname, fset, []*ast.File{f}, &info)
	if !*flagNoTypeCheck {
		checkUserErr(err)
	}
	return info
}

// PosOf returns the position of n using the global fset.
func PosOf(n ast.Node) token.Position {
	return fset.Position(n.Pos())
}

// ObjectOf returns the object denoted by the specified identifier.
func ObjectOf(id *ast.Ident) types.Object {
	obj := objectOf(id)
	if obj == nil {
		Error(id, "undefined:", id.Name)
	}
	return obj
}

func objectOf(id *ast.Ident) types.Object {
	return info.ObjectOf(id)
}

func TypeOf(n ast.Expr) types.Type {
	t := info.TypeOf(n)
	if t == nil {
		Error(n, "cannot infer type")
	}
	return t
}

// ParentOf returns the parent node of n.
// Precondition: CollectParents has been called on the tree containing n.
func ParentOf(n ast.Node) ast.Node {
	if p, ok := parent[n]; ok {
		return p
	} else {
		panic(PosOf(n).String() + ": no parent")
	}
}

// Return the first a ancestor of n that is an ast.FuncDecl.
// Used by return statements to find the function declaration they return from.
func FuncDeclOf(n ast.Node) *ast.FuncDecl {
	for p := ParentOf(n); p != nil; p = ParentOf(p) {
		if f, ok := p.(*ast.FuncDecl); ok {
			return f
		}
	}
	panic(PosOf(n).String() + ": no FuncDecl parent for node")
}

// returun exact value and minimal type for constant expression.
func ExactValue(e ast.Expr) (tv types.TypeAndValue, ok bool) {
	tv, ok = info.Types[e]
	return
}