package main

const b int = 1
const c = 2
const d, e = 3, 4

const (
	g    int = 5
	h        = 6
	i, j     = 7, 8
)

func main() {

	const a = 12
	var f = a

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
	println(B)
	println(C)
	println(D)
	println(E)
	println(G)
	println(H)
	println(I)
	println(J)
	println(X)
	println(Y)
	println(Z)
	println(U)
}

const B int = 10
const C = 20
const D, E = 30, 40

const (
	G    int = 50
	H        = 60
	I, J     = 70, 80
)

const (
	X = 1
	Y = X + 2
	Z = len("hello")
	U = 2 * X << Y
)
