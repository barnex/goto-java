package main

// test empty return
func f() {
	return
}

func main() {
	f()

	// test return of one value, named and unnamed
	println(div(10, 3))
	println(div2(20, 4))
	println(div3(30, 5))
	println(div4(50, 6))
}

func div(a, b int) int {
	return a / b
}

func div2(a, b int) (result int) {
	result = a / b
	return
}

func div3(a, b int) (result int) {
	result = a / b
	return result
}

func div4(a, b int) (result int) {
	return a / b
}
