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
	// TODO: emit final for struct (value) type?
	if ident.IsExported() {
		return PUBLIC
	} else {
		return PROTECTED
	}
}

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
