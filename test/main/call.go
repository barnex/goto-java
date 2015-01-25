package main

func f(x int) int {
	return x
}
func main() {
	println(f(x))
	println((f)(x))
	println(((f))(x))
}
