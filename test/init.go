package main

var (
	a = b + c
	b = c
	c = 1
	d int
)

const (
	A = B + C
	B = C
	C = 1
)

func init() {
	d = 4
}

func main() {
	println(a, b, c, d)
	println(A, B, C)
}
