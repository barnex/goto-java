package main

// This file provides functionality to rename identifiers.
// Some valid Go identifiers cannot be used in java (e.g. keywords like "static").
// Sometimes we need to rename a variable because of scope rules.

import (
	"fmt"
	"go/ast"
	"log"

	"golang.org/x/tools/go/types"
)

var (
	renamed = map[types.Object]string{} // maps some objects (typ. identifiers) to a new name for java.
	idents  map[string]int              // holds all identifier names and a counter to create a new, non-conflicting name if needed.
	UNUSED  string                      // base name for translating the blank identifier (flag -blank)
)

// Collect the names of all identifiers in the AST and store in idents.
// Used to ensure identifier renaming returns an unused name.
// In principle new names only need to be unique in their scope,
// but we make them globally unique to avoid potential scope subtleties.
func CollectIdents(n ast.Node) {
	idents = make(map[string]int) // init here ensures CollectIdents has been called
	ast.Walk(&identCollector{}, n)
	idents[UNUSED] = idents[UNUSED] // make sure it's in the map for makeNewName(UNUSED) to work.
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

// java keywords and pre-defined literals, cannot be used as java identifier names.
// http://docs.oracle.com/javase/tutorial/java/nutsandbolts/_keywords.html
var javaKeyword = map[string]bool{
	"abstract":     true,
	"assert":       true,
	"boolean":      true,
	"break":        true,
	"byte":         true,
	"case":         true,
	"catch":        true,
	"char":         true,
	"class":        true,
	"const":        true,
	"continue":     true,
	"default":      true,
	"do":           true,
	"double":       true,
	"else":         true,
	"enum":         true,
	"extends":      true,
	"false":        true,
	"final":        true,
	"finally":      true,
	"float":        true,
	"for":          true,
	"goto":         true,
	"if":           true,
	"implements":   true,
	"import":       true,
	"instanceof":   true,
	"int":          true,
	"interface":    true,
	"long":         true,
	"native":       true,
	"new":          true,
	"null":         true,
	"package":      true,
	"private":      true,
	"protected":    true,
	"public":       true,
	"return":       true,
	"short":        true,
	"static":       true,
	"strictfp":     true,
	"super":        true,
	"switch":       true,
	"synchronized": true,
	"this":         true,
	"throw":        true,
	"throws":       true,
	"transient":    true,
	"true":         true,
	"try":          true,
	"void":         true,
	"volatile":     true,
	"while":        true,
}
