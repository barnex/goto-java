package main

type MyInt int

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
	println(i)

	j := i
	println(j)
	(&i).inc()
	println(i)
}
