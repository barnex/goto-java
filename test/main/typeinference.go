package main

func f() int {
	return 42
}

func main() {
	a := 1
	b, c := a, 2+3
	f := f()
	println(a)
	println(b)
	println(c)
	println(f)
}
