package main

type MyInt int
type A MyInt 

func (i MyInt) square() int {
	i = i * i // try to mute
	return int(i)
}

func (i *MyInt) inc() {
	*i++
}

func main() {
	var i MyInt = 4
	println(i)
	println(i.square())
	println(i)
	(&i).inc()
	i.inc()
	println(i)

	j := i
	println(j)
	(&i).inc()
	println(i)

	x := 0
	var a A
	a = 0
	a = A(0)
	a = A(MyInt(0))
	a = A(MyInt(x))
	println(a)
}
