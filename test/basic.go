package main

// test basic type

var glob int

func main() {
	testDeclare()
	testAssign()
	testCompare()
}

func testDeclare() {
	var (
		l1 int
		l2 int = 2
		l3     = 3
	)
	l4 := 4

	println(l1, l2, l3, l4)
}

func testAssign() {
	// global
	glob = 7
	println(glob)

	// local
	var i int
	i = 1
	println(i)

	// escaped
	var e int
	_ = &e
	e = i
	e = 666
	println(e, i)
}

func testCompare() {
	a, A := 2, 2
	b, B := 3, 3

	_ = &A
	_ = &B

	println(a == A, A == a)
	println(a == B, B == a)
	println(b == B, B == b)
	println(b == A, A == b)

	println(a >= A, A >= a)
	println(a >= B, B >= a)
	println(b >= B, B >= b)
	println(b >= A, A >= b)

	println(a <= A, A <= a)
	println(a <= B, B <= a)
	println(b <= B, B <= b)
	println(b <= A, A <= b)

	println(a < A, A < a)
	println(a < B, B < a)
	println(b < B, B < b)
	println(b < A, A < b)

	println(a > A, A > a)
	println(a > B, B > a)
	println(b > B, B > b)
	println(b > A, A > b)

}
