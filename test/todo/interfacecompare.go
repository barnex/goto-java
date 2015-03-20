package main

type A interface{}
type B interface{}

func main() {
	var a A
	var b B
	println(a == b)
	println(a == nil)
	println(nil == a)
	a = 1
	println(a == nil)
	println(a == b)
	println(b == a)
}
