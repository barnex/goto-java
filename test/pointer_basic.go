package main

// test pointer to basic

var (
	glob, glob2  *int
	globi, globj int
)

func main() {
	testDeclare()
	testAssign()
	testEquals()
}

func testDeclare() {
	var i int
	var (
		l1 *int
		l2 *int = l1
		l3      = new(int)
	)
	l4 := new(int)
	l5 := &i

	println(l1 == nil, l2 == nil, *l3, *l4, *l5)
}

func testAssign() {
	var i *int
	i = new(int)
	println(*i)

	j := 2
	glob = &j
	println(*glob)
}

func testEquals() {
	var a, A, b, B *int
	a = new(int)
	A = a
	b = new(int)
	B = b

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

	i := 1
	b = &i
	B = &i

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

	b = &globi
	B = &globi

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

	a = &globj
	A = &globj

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

	b = glob
	B = glob

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

	a = glob2
	A = glob2

	println(a == a, A == A)
	println(b == b, B == B)

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == A, A == b)
	println(b == B, b == b)

}
