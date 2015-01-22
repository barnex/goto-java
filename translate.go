package main

import "go/ast"

// Return java modifier for Go name. E.g.:
// 	Name -> public
// 	name -> "" (package private)
func ModifierFor(name string) string {
	if ast.IsExported(name) {
		return "public "
	} else {
		return ""
	}
}
