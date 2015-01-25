package main

type JModifier uint32

const (
	NONE    JModifier = 0
	PRIVATE JModifier = 1 << iota
	PROTECTED
	PUBLIC
	STATIC
	FINAL
)


