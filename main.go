package gotojava

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

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
	fset  = token.NewFileSet() // accessed through PosOf
	files []*ast.File
	info  = types.Info{ // accessed through TypeOf, ObjectOf, ExactValue
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	parent = map[ast.Node]ast.Node{}   // accessed through ParentOf
	idents = map[string]int{}          // holds all identifier names and a counter to create a new, non-conflicting name if needed.
	rename = map[types.Object]string{} // maps some objects (typ. identifiers) to a new name for java. TODO: strings->strings
)

var (
	jFiles JFileSet
)

func Main() {
	log.SetFlags(0)
	flag.Parse()
	UNUSED = *flagBlank

	fnames := flag.Args()
	for _, f := range fnames {
		ParseFile(f)
	}

	if *flagPrint {
		for _, f := range files {
			ast.Print(fset, f)
		}
	}

	for i, f := range files {
		TypeCheck(fnames[i], f)
		CollectParents(f)
		CollectIdents(f)
		RenameReservedIdents(f)
		CollectTypes(f)
		EscapeAnalysis(f)
	}

	//jFiles.Compile()
}

// Parse file, add position info to global fset, add AST to global files.
func ParseFile(fname string) {
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	Check(err)
	files = append(files, f)
}

// Typecheck file and add type information to global info variable.
func TypeCheck(fname string, f *ast.File) {
	var config types.Config
	_, err := config.Check(fname, fset, []*ast.File{f}, &info)
	if !*flagNoTypeCheck {
		Check(err)
	}
}
