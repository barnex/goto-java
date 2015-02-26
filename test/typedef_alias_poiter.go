package main

type S struct{ v int }
type A *S
type B A

// Hooray, no methods allowed on pointer type


func main() {
	var a A
	a = &S{}
	a = &S{0}
	a = &S{v: 0}
	a = A((*S)(&struct{ v int }{0}))

	var b B
	b = B(a)
	b.v++
	
	println(a.v)
}

