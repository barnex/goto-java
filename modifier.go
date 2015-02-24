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

func ModifierFor(ident *ast.Ident) JModifier {
	if ident.IsExported() {
		return PUBLIC
	} else {
		return PROTECTED
	}
}

func (m JModifier) String() string {
	str := ""
	if m.Is(PRIVATE) {
		str = cat(str, "private")
	}
	if m.Is(PROTECTED) {
		str = cat(str, "protected")
	}
	if m.Is(PUBLIC) {
		str = cat(str, "public")
	}
	if m.Is(STATIC) {
		str = cat(str, "static")
	}
	if m.Is(FINAL) {
		str = cat(str, "final")
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

// Returns true if m contains b. E.g.:
// 	m.Is(PUBLIC) // true if m is public
func (m JModifier) Is(b JModifier) bool {
	return m&b != 0
}
