package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"

	"golang.org/x/tools/go/types"
)

var (
	flagBlank       = flag.String("blank", "_", "Java name for the blank (underscore) identifier")
	flagConstFold   = flag.Bool("foldconst", false, "Fold constants")
	flagNoPkg       = flag.Bool("nopkg", false, "Do not output package clause")
	flagNoTypeCheck = flag.Bool("nocheck", false, "Don't do type check")
	flagParens      = flag.Bool("parens", false, "Emit superfluous parens")
	flagPrint       = flag.Bool("print", false, "Print ast")
	flagRenameAll   = flag.Bool("renameall", false, "Rename all variables (debug)")
)

// TODO: global package, use for class gen unless overridden.

var (
	fset    *token.FileSet
	info    types.Info
	parents map[ast.Node]ast.Node // maps every node to his parent node. Populated by CollectParents
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	UNUSED = *flagBlank

	for _, f := range flag.Args() {
		handleFile(f)
	}

}

func checkUserErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

func handleFile(fname string) {

	// read and parse file
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	checkUserErr(err)

	var config types.Config
	info = types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	_, err = config.Check(fname, fset, []*ast.File{f}, &info)
	if !*flagNoTypeCheck {
		checkUserErr(err)
	}

	// print ast if requested
	if *flagPrint {
		ast.Print(fset, f)
	}

	// first pass collects all declarations
	CollectDefs(f)
	idents[UNUSED] = idents[UNUSED] // make sure it's in the map for makeNewName(UNUSED) to work.

	// transpile primary class
	className := fname[:len(fname)-len(path.Ext(fname))]
	w := NewWriter(className + ".java")
	defer w.Close()
	w.PutClass(className, f)

	// generate additional classes from type/method definitions encountered along the way
	GenClasses()
}
