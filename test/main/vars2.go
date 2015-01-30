package main

var b = 42 // try to confuse scope

// short variable declaration with one already declared.
func main() {

	{
		// try to confuse scope
		a, b, c := 999, 999, 999
		println(a)
		println(b)
		println(c)
	}

	a := 1
	a, b := 2, 3
	c, a, b := 4, 5, 6

	println(a)
	println(b)
	println(c)

}

var c = 43 // try to confuse position info
