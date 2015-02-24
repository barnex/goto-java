package main

// Test blank identifier (_)

// Try to induce name collision with names generated for blank identifier
// (flag -blank)
var _ = 3
var _0 = 4
var _1 = 5

var _ = new(Padded)

func main() {
	_ = f(1)
	_ = f(1)
	(_) = f(1)
	(_) = 2
	_, i := f(2), f(3)
	println(i)

	println(f2(1, 2))
}

func f(x int) int {
	println(x)
	return x + 1
}

func f2(_ int, _ int) int {
	var _ = 3
	var _0 = 4
	var _1 = 5
	var _2 = 6
	println(_0)
	println(_1)
	println(_2)
	return 42
}

type Padded struct {
	x int
	_ int
	_ int
}

func (Padded) f1()    {}
func (_ Padded) f2()  {}
func (*Padded) f3()   {}
func (_ *Padded) f4() {}
