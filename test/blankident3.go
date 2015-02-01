package main

// try to induce name collision with "unused"
var unused = 3
var unused_0 = 4
var unused_1 = 5

func f(_ int, _ int) int {
	var unused = 3
	var unused_0 = 4
	var unused_1 = 5
	var unused_2 = 6
	println(unused)
	println(unused_0)
	println(unused_1)
	println(unused_2)
	return 42
}

func main() {
	println(f(1, 2))
}
