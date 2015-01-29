package main

// This file provides functionality to rename identifiers.

import (
	"fmt"
	"go/ast"
	"log"

	"golang.org/x/tools/go/types"
)

var (
	renamed = map[types.Object]string{}
	idents  map[string]int
)

// Collect the names of all identifiers in the AST and store in idents.
// Used to ensure identifier renaming returns an unused name.
// In principle new names only need to be unique in their scope,
// but we make them globally unique to avoid potential scope subtleties.
func CollectIdents(n ast.Node) {
	idents = make(map[string]int) // init here ensures CollectIdents has been called
	ast.Walk(&identCollector{}, n)
}

type identCollector struct{}

func (f identCollector) Visit(n ast.Node) ast.Visitor {
	if id, ok := n.(*ast.Ident); ok {
		idents[id.Name] = 0
	}
	return f
}

func (w *writer) translate(id *ast.Ident) string {
	obj := w.info.ObjectOf(id)

	if tr, ok := renamed[obj]; ok {
		return tr
	}

	if javaKeyword[obj.Name()] {
		new := makeNewName(obj.Name())
		log.Println("renmaing", obj.Name(), "->", new)
		renamed[obj] = new
		return new
	}

	return obj.Name()
}

func makeNewName(orig string) string {
	idents[orig]++
	return fmt.Sprint(orig, "_", idents[orig])
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
