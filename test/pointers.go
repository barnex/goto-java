package main

var (
	global    int
	globalptr = &global
)

//func inc(i *int) {
//	(*i)++
//}

func f()int{return 42}

func main() {

	// Test declarations of escaping basic
	i0 := 4
	var i1 int 
	var i2 int = 1
	var i3, i4 int
	var i5, xx int = 2, 3
	var i7 int = i1
	var i8 = i1
	var i9 = f()
	ia, yy := 5, 6

	var p *int

	p = &i0
	p = &i1
	p = &i2
	p = &i3
	p = &i4
	p = &i5
	p = &i7
	p = &i8
	p = &i9
	p = &ia

	println(*p)

	println(i0)
	println(i1)
	println(i2)
	println(i3)
	println(i4)
	println(i5)
	println(xx)
	println(i7)
	println(i8)
	println(i9)
	println(ia)
	println(yy)


//	i := 0 // escapes!
//	i = 1
//	i++
//	println(i)
//	println(i == 1)

//	x := &i
//	(*x)++
//	println(*x)
//	println(i)
//	*x = 666
//	println(i)
//
//	var j int
//	x = &j
//	(*x)++
//	println(j)
//
//	var k int = 387
//	x = &k
//	(*x)++
//	println(k)

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
