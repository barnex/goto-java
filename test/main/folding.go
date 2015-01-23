package main

const (
	A = 1
	B = A + 2
	//C = len("hello")
	D = 2 * A << B
)

func main() {
	println(A)
	println(B)
	//println(C)
	println(D)
}
