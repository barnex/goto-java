package main

var (
	global    int
	globalptr = &global
)

//func inc(i *int) {
//	(*i)++
//}

func main() {
	i := 0
	x := &i
	(*x)++
	println(i)
	i = 666
	println(i)

	var j int
	x = &j
	(*x)++
	println(j)

	var k int = 387
	x = &k	
	(*x)++
	println(k)

	println(global)
//	inc(&global)
//	println(global)
//	inc(globalptr)
//	println(global)
//
//	var x *int
//	println(x == nil)
//
//	y := new(int)
//	println(y == nil)
//
//	*y = 7
//	println(*y)
//	println(x == y)
//
//	x = y
//	*x = 8
//	println(*y)
//	println(x == y)
//
//	z := &x
//	**z = 5
//	println(**z)
//
//	println(*makePtr())
	//testShortRedefine()
}

//func makePtr() *int {
//	i := 42
//	return &i
//}

//func testShortRedefine(){
//	a := ptrTo(1)
//	a, b := ptrTo(2), ptrTo(3)
//	c, a, b := ptrTo(4), ptrTo(5), ptrTo(6)
//
//	println(*a)
//	println(*b)
//	println(*c)
//}
//
//func ptrTo(x int)*int{
//	a := x
//	return &a
//}
