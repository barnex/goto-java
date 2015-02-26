package gotojava

// Tuples wrap multiple return values in a java class.

import "fmt"

var (
	classGen = make(map[string]bool) // has code for helper class been generated?
)

// JavaTupleType returns the java type used to wrap a tuple of go types for multiple return values. E.g.:
// 	return 1, 2 -> return new Tuple_int_int(1, 2)
// Calling this function also ensure code for the tuple has been generated.
// TODO: JType
func JavaTupleType(types []JType) string {
	name := "Tuple"
	for _, t := range types {
		name += "_" + t.GoType.String() // not java name as we want to discriminate, e.g., int from uint
	}

	if !classGen[name] {
		GenTupleDef(name, types)
	}
	return name
}

// TODO merge with GenStructSomething
func GenTupleDef(name string, types []JType) {
	classGen[name] = true

	w := NewWriterFile(name + ".java")
	defer w.Close()

	w.Putln("public final class ", name, "{\n")
	w.indent++

	for i, t := range types {
		w.Putln("public ", t, " ", fmt.Sprint("v", i), ";")
	}

	w.Putln()
	w.Put("public ", name, "(")
	for i, t := range types {
		w.Put(comma(i), t, " ", fmt.Sprint("v", i))
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
