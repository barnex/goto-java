package main

// test variable declarations and type inference

func main() {

	var a int           // no init value
	b := 1              // :=, type inferred from value
	var c, d int = 2, 3 // var, two init values
	var e = 4           // var, type inferred from value
	var f, g = 5, "hi"  // var, two different types
	var h string        // var, string, no init value

	A := 1
	B, C := A, 2+3
	F := F()
	println(A)
	println(B)
	println(C)
	println(F)

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

	testShortRedefine()
}

var b = 42 // try to confuse scope

func F() int { return 56 }

// short variable re-declaration
func testShortRedefine() {
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
