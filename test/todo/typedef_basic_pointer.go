package main

type I int
type B *I

func main() {
	var x int // escapes

	var i I // escapes
	i = 0
	i = I(0)
	i = I(x)

	var b B
	b = B(&i)
	(*b)++

	b = B((*I)(&x))
	(*b) += 2

	println(*b)
	println(i)
	println(x)
}
