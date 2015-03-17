package main

// test basic type

var glob int

func main() {
	testDeclare()
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
