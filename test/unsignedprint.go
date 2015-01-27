package main

func main() {
	x := -3

	a := byte(x)
	b := uint8(x)
	c := uint16(x)
	d := uint32(x)
	e := uint64(x)

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)

	e = 1<<63 + 1234
	e = 1<<64 - 1
	println(e)
}
