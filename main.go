package gotojava

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
	flagVerbose     = flag.Bool("v", true, "verbose logging")
	UNUSED          string // base name for translating the blank identifier (flag -blank)
)

var (
	fset   *token.FileSet          // accessed through PosOf
	info   types.Info              // accessed through TypeOf, ObjectOf, ExactValue
	parent map[ast.Node]ast.Node   // accessed through ParentOf
	idents map[string]int          // holds all identifier names and a counter to create a new, non-conflicting name if needed.
	rename map[types.Object]string // maps some objects (typ. identifiers) to a new name for java. TODO: strings->strings
)

func Main() {
	log.SetFlags(0)
	flag.Parse()
	UNUSED = *flagBlank

	for _, f := range flag.Args() {
		HandleFile(f)
	}
}

func HandleFile(fname string) {
	// (1) Parse
	var f *ast.File
	fset, f = parseFile(fname)
	if *flagPrint {
		ast.Print(fset, f)
	}

	// (2) Determine types
	info = typeCheck(fname, fset, f)

	// (3) Pre-processing: collect parents and declarations
	parent = CollectParents(f)
	idents = CollectIdents(f)
	rename = RenameReservedIdents(f)
	CollectTypes(f)
	EscapeAnalysis(f)

	// transpile primary class
	className := fname[:len(fname)-len(path.Ext(fname))]
	w := NewWriterFile(className + ".java")
	defer w.Close()
	w.PutClass(className, f)

	// generate additional classes from type/method definitions encountered along the way
	GenClasses()
}

func parseFile(fname string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	Check(err)
	return fset, f
}

func typeCheck(fname string, fset *token.FileSet, f *ast.File) types.Info {
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
		Check(err)
	}
	return info
}
