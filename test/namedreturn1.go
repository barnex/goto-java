package main

// test named return values

func main() {
	println(inc(1))
}

func inc(a int) (result int) {
	result = a + 1
	return
}
