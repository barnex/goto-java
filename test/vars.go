package main

// test variable declarations, type inference and type conversion
func main() {

	var a int           // no init value
	b := 1              // :=, type inferred from value
	var c, d int = 2, 3 // var, two init values
	var e = 4           // var, type inferred from value
	var f, g = 5, "hi"  // var, two different types
	var h string        // var, string, no init value

	var (
		i    int
		j    = 7
		k, l = 8, "hi"
		m    string
	)

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)
	println(g)
	println(h)
	println(i)
	println(j)
	println(k)
	println(l)
	println(m)
}