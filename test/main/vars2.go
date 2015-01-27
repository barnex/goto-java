package main

// short variable declaration with one already declared.
func main() {
	a := 1
	a, b := 2, 3
	c, a, b := 4, 5, 6

	println(a)
	println(b)
	println(c)
}
