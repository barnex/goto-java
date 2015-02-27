package gotojava

// Tuples wrap multiple return values in a java class.

import "fmt"

var (
	classGen = make(map[string]bool) // has code for helper class been generated?
)

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
