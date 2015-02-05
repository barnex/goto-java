package main

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

func (m JModifier) String() string {
	str := ""
	switch {
	case m.Is(PRIVATE):
		str = cat(str, "private")
	case m.Is(PROTECTED):
		str = cat(str, "protected")
	case m.Is(PUBLIC):
		str = cat(str, "public")
	case m.Is(STATIC):
		str = cat(str, "static")
	case m.Is(FINAL):
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
