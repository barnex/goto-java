package main

// Test assignment of untyped int literal to typed variable.
// Java can get very picky here.

func main() {

	var (
		a byte   = 1
		b uint8  = 2
		c int8   = 3
		d int16  = 4
		e uint16 = 5
		f int32  = 6
		g uint32 = 7
		h int    = 8
		i uint   = 9
		j int64  = 10
		k uint64 = 11
	)

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)
	println(g)
	println(h)
	println(i)
	println(j)
	println(k)
}
