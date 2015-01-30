package main

func main() {
	_ = f(1)
	_, i := f(2), f(3)
	println(i)
}

func f(x int) int {
	println(x)
	return x + 1
}
