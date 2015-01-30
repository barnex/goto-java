package main

// out-of-order constants initialization
const (
	D = 2 * A << B
	B = A + 2
	A = 1
)

func main() {
	println(A)
	println(B)
	println(D)
}
