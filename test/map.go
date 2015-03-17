package main

func main() {
	var (
		a map[int]int = map[int]int{1: 2, 3: 4}
		b map[int]int
		c map[int]int = b
		d map[int]int = a
	)

	println(a[1])
	println(b == nil)
	println(c == nil)
	println(d[1])

}
