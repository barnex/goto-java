package main

func main() {
	var (
		a []int
		b []int = a
		c []int = []int{1, 2, 3}
	)

	println(a[0])
	println(b[0])
	println(c[0])

}
