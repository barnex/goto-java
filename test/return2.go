package main

// test return of multiple values, named and unnamed.

func main() {
	//a, b := div(50, 6)
	//println(a)
	//println(b)

	//a, b = div2(10, 3)
	//println(a)
	//println(b)

	//a, b = div3(20, 3)
	//println(a)
	//println(b)

	//a, b = div4(50, 3)
	//println(a)
	//println(b)

	//a, b = div5(70, 4)
	//println(a)
	//println(b)
}

func div(a, b int) (int, int) {
	return a / b, a % b
}

func div2(a, b int) (quo, rem int) {
	return a / b, a % b
}

func div3(a, b int) (quo int, rem int) {
	return a / b, a % b
}

func div4(a, b int) (quo, rem int) {
	quo = a / b
	rem = a % b
	return quo, rem
}

func div5(a, b int) (quo, rem int) {
	quo = a / b
	rem = a % b
	return
}
