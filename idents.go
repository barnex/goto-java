package gotojava

// This file provides functionality to rename identifiers.
// Some valid Go identifiers cannot be used in java (e.g. keywords like "static").
// Sometimes we need to rename a variable because of scope rules.

import (
	"fmt"
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// Emit an identifier, honoring the global rename map. Godoc:
// 	type Ident struct {
// 	        NamePos token.Pos // identifier position
// 	        Name    string    // identifier name
// 	        Obj     *Object   // denoted object; or nil
// 	}
func (w *Writer) PutIdent(id *ast.Ident) {
	if IsBlank(id) {
		w.Put(makeNewName(UNUSED))
		return
	}

	if new, ok := rename[ObjectOf(id)]; ok {
		w.Put(new)
		return
	}

	switch id := ObjectOf(id).(type) {
	default:
		panic("cannot handle " + reflect.TypeOf(id).String())
	case *types.Const:
		w.Put(id.Name())
	case *types.Func:
		w.Put(id.Name())
	case *types.Nil:
		w.Put("null")
	case *types.Var:
		w.Put(id.Name())
	}
}

// Is e the blank identifier?
// Also handles the corner case of parenthesized blank ident (_)
func IsBlank(e ast.Expr) bool {
	e = StripParens(e)
	if id, ok := e.(*ast.Ident); ok {
		return id == nil || id.Name == "_"
	}
	return false
}

// Collect the names of all identifiers in the AST and maps them to a counter set to 0.
// Later this counter can be incremented to create an new, unused, identifier name.
// In principle new names only need to be unique in their scope,
// but we make them globally unique to avoid potential scope subtleties.
func CollectIdents(n ast.Node) map[string]int {
	idents := make(map[string]int) // init here ensures CollectIdents has been called
	ast.Inspect(n, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			idents[id.Name] = 0
		}
		return true
	})
	idents[UNUSED] = idents[UNUSED] // make sure it's in the map for makeNewName(UNUSED) to work.
	return idents
}

func RenameReservedIdents(n ast.Node) map[types.Object]string {
	rename := make(map[types.Object]string)

	ast.Inspect(n, func(n ast.Node) bool {

		if id, ok := n.(*ast.Ident); ok {
			obj := objectOf(id)
			if obj == nil {
				return true
			}

			// Name is keyword: rename it and return new name.
			// DEBUG: flag -renameall renames all variables for stress testing.
			if canRename(obj) && (javaKeyword[obj.Name()] || *flagRenameAll) {
				new := makeNewName(obj.Name())
				Log(n, obj.Name(), "->", new)
				rename[obj] = new
			}
		}
		return true
	})

	return rename
}

func canRename(obj types.Object) bool {
	isBuiltin := (obj.Parent() == types.Universe)
	return !isBuiltin
}

// JavaNameFor returns the java name for identifier,
// possibly renamed to avoid java keywords.
func JavaName(id *ast.Ident) string {
	obj := ObjectOf(id)
	if renamed, ok := rename[obj]; ok {
		return renamed
	} else {
		return obj.Name()
	}
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
