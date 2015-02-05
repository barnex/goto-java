package main

// This file provides functionality to rename identifiers.

import (
	"fmt"
	"go/ast"
	"log"

	"golang.org/x/tools/go/types"
)

var (
	renamed = map[types.Object]string{} // maps some objects (typ. identifiers) to a new name for java.
	idents  map[string]int              // holds all identifier names and a counter to create a new, non-conflicting name if needed.
	UNUSED  string
)

// Collect the names of all identifiers in the AST and store in idents.
// Used to ensure identifier renaming returns an unused name.
// In principle new names only need to be unique in their scope,
// but we make them globally unique to avoid potential scope subtleties.
func CollectIdents(n ast.Node) {
	idents = make(map[string]int) // init here ensures CollectIdents has been called
	ast.Walk(&identCollector{}, n)
	idents[UNUSED] = idents[UNUSED] // make sure it's in the map for makeNewName("unused") to work.
}

// used by CollectIdents to put all identifier names in idents.
type identCollector struct{}

func (f identCollector) Visit(n ast.Node) ast.Visitor {
	if id, ok := n.(*ast.Ident); ok {
		idents[id.Name] = 0
	}
	return f
}

// Translate an identifier to its java name.
// Usually returns the identifier's name unchanged,
// unless it has been renamed for some reason or when
// the identifier name is a protected java keyword.
func JavaName(id *ast.Ident) string {
	obj := ObjectOf(id)

	if obj == nil {
		Error(id, "undefined:", id.Name)
	}

	// object has been renamed
	if tr, ok := renamed[obj]; ok {
		return tr
	}

	// Name is keyword: rename it and return new name.
	// DEBUG: flag -renameall renames all variables for stress testing.
	if javaKeyword[obj.Name()] || *flagRenameAll {
		new := makeNewName(obj.Name())
		log.Println("renmaing", obj.Name(), "->", new)
		renamed[obj] = new
		return new
	}

	// nothing special: return original name
	// TODO: do we need to filter unicode names?
	return obj.Name()
}

func (w *writer) rename(id *ast.Ident) {
	renamed[ObjectOf(id)] = makeNewName(id.Name)
}

// Construct a new (java) name for a (go) identifier with original name orig.
// The new name is globally unique and can be used in any scope.
func makeNewName(orig string) string {
	new := orig
	for {
		if _, ok := idents[new]; ok {
			idents[orig]++
			new = fmt.Sprint(orig, idents[orig]) // append number
		} else {
			break
		}
	}
	return new
}

var lit2java = map[string]string{
	"false": "false",
	"nil":   "null", //? need to type!
	"true":  "true",
}

var javaKeyword = map[string]bool{
	"abstract":     true,
	"continue":     true,
	"for":          true,
	"new":          true,
	"switch":       true,
	"assert":       true,
	"default":      true,
	"goto":         true,
	"package":      true,
	"synchronized": true,
	"boolean":      true,
	"do":           true,
	"if":           true,
	"private":      true,
	"this":         true,
	"break":        true,
	"double":       true,
	"implements":   true,
	"protected":    true,
	"throw":        true,
	"byte":         true,
	"else":         true,
	"import":       true,
	"public":       true,
	"throws":       true,
	"case":         true,
	"enum":         true,
	"instanceof":   true,
	"return":       true,
	"transient":    true,
	"catch":        true,
	"extends":      true,
	"int":          true,
	"short":        true,
	"try":          true,
	"char":         true,
	"final":        true,
	"interface":    true,
	"static":       true,
	"void":         true,
	"class":        true,
	"finally":      true,
	"long":         true,
	"strictfp":     true,
	"volatile":     true,
	"const":        true,
	"float":        true,
	"native":       true,
	"super":        true,
	"while":        true,
	"true":         true,
	"false":        true,
	"null":         true,
}
