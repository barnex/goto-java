package main

func main() {
	a, b := div(50, 6)
	println(a)
	println(b)
}

func div(a, b int) (int, int) {
	return a / b, a % b
}
