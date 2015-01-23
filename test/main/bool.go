package main

// test boolean variables and initialisation.
func main() {
	var a bool
	b := a

	c := true
	var d = false
	e, f := true, d

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)

	println(a || b && c || d)
	println(a || b && c || d == f != a || a == b && c != d == f || f)
}
