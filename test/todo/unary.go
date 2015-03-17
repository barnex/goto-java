package main

func main() {
	a := 1
	b := 234
	c := 6
	println(+a)
	println(-a)
	println(!false)
	println(!false && true)
	println(^a)
	println(^a + b)
	println(^a * b)
	println(^a*^b + c)
	println(^a | c)
	println(^a & c)
	println(a & ^c)
	println(a | ^c)
}
