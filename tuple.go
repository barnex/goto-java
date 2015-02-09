package main

// Tuples wrap multiple return values in a java class.

import (
	"fmt"

	"golang.org/x/tools/go/types"
)

var (
	classGen = make(map[string]bool) // has code for helper class been generated?
)

// JavaTupleType returns the java type used to wrap a tuple of go types for multiple return values. E.g.:
// 	return 1, 2 -> return new Tuple_int_int(1, 2)
// Calling this function also ensure code for the tuple has been generated.
func JavaTupleType(types []types.Type) string {
	name := "Tuple"
	for _, t := range types {
		name += "_" + t.String() // not java name as we want to discriminate, e.g., int from uint
	}

	if !classGen[name] {
		GenTupleDef(name, types)
	}
	return name
}

// TODO class{implements, members, ...}.render
func GenTupleDef(name string, types []types.Type) {
	classGen[name] = true

	w := NewWriter(name + ".java")
	defer w.Close()

	w.Putln("public final class ", name, "{\n")
	w.indent++

	for i, t := range types {
		w.Putln("public ", javaTypeOf(t), " ", fmt.Sprint("v", i), ";")
	}

	w.Putln()
	w.Put("public ", name, "(")
	for i, t := range types {
		w.Put(comma(i), javaTypeOf(t), " ", fmt.Sprint("v", i))
	}
	w.Putln("){")
	w.indent++
	for i := range types {
		w.Putln("this.", fmt.Sprint("v", i), " = ", fmt.Sprint("v", i), ";")
	}
	w.indent--
	w.Putln("}")

	w.indent--
	w.Put("}")
}
