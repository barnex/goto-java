package main

func main() {
	var a struct{ v int } = struct{ v int }{7}
	var b struct{ v int } = a
	var c struct{ v int } = struct{ v int }{}
	println(b.v)
	println(c.v)

	c = b
	println(c.v)
	println(c.v)
}
