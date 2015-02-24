package main

var (
	global    int
	globalptr = &global
)

func inc(i *int) {
	(*i)++
}

func main() {
	println(global)
	inc(&global)
	println(global)
	inc(globalptr)
	println(global)

	var x *int
	println(x == nil)

	y := new(int)
	println(y == nil)

	*y = 7
	println(*y)
	println(x == y)

	x = y
	*x = 8
	println(*y)
	println(x == y)

	z := &x
	**z = 5
	println(**z)

	println(*makePtr())
}

func makePtr() *int {
	i := 42
	return &i
}
