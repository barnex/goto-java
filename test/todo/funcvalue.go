package main

func f(x int) int {
	return x
}

func main() {
	x := 1
	println(f(x))
	println((f)(x))
	println((f)(x))
}
