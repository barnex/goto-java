package gotojava

import (
	"go/ast"
)

// JModifier holds a java modifier. E.g.:
// 	PUBLIC | STATIC | FINAL
type JModifier uint32

const (
	NONE    JModifier = 0
	PRIVATE JModifier = 1 << iota
	PROTECTED
	PUBLIC
	STATIC
	FINAL
)

// Returns a suited java modifier (public, ...) for the definition of ident. E.g.:
// 	x -> NONE   (package private)
// 	X -> PUBLIC (exported)
// Modifier is also final for struct types (implemented as final java references)
func ModifierFor(ident *ast.Ident) JModifier {
	mod := NONE
	if ident.IsExported() {
		mod |= PUBLIC
	}
	if JavaType(TypeOf(ident)).IsStructValue() {
		mod |= FINAL
	}
	return mod
}

// String representation of modifier, followed by space unless empty. E.g.:
// 	"private static"
// 	""                // package private
func (m JModifier) String() string {
	str := ""
	if m&PRIVATE != 0 {
		str = cat(str, "private")
	}
	if m&PROTECTED != 0 {
		str = cat(str, "protected")
	}
	if m&PUBLIC != 0 {
		str = cat(str, "public")
	}
	if m&STATIC != 0 {
		str = cat(str, "static")
	}
	if m&FINAL != 0 {
		str = cat(str, "final")
	}
	if str != "" {
		str += " "
	}
	return str
}

// concatenate a and b, inserting a space if needed
func cat(a, b string) string {
	if (a != "") && (b != "") {
		return a + " " + b
	} else {
		return a + b
	}
}
