package main

func main() {
	var (
		a [3]int
		b [3]int = a
		c [3]int = [3]int{1, 2, 3}
		//d [3]int = [...]int{1, 2, 3}
		e [3]int = [3]int{1, 2}
	)

	println(a[0])
	println(b[0])
	println(c[0])
	//	println(d[0])
	println(e[0])
	println(e[2])
}
