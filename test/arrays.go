package main

func main() {
	var (
		a [3]int
		b [3]int = a
		c [3]int = [3]int{1, 2, 3}
	)

	println(a[0])
	println(b[0])
	println(c[0])

}
