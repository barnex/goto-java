package main

func main() {
	//var a struct{ v int } = struct{ v int }{}
	var a struct{ v int } = struct{ v int }{7}
	var b struct{ v int } = a
	println(b.v)
}
