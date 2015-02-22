package main

type Int struct {
	val int
}

// A test struct.
type Struct struct {
	x, y Int
	v    int
}

func (s Struct) ValMethod() int {
	s.v *= 2
	s.x.val += 10
	return s.x.val + s.y.val + s.v
}

func (s *Struct) PtrMethod(v int) {
	s.v = v
}

func main() {

	var s *Struct
	println(s == nil)

	s = new(Struct)
	println(s == nil)

	println(s.v)
	println(s.x.val)
	println(s.y.val)

	s.PtrMethod(42)
	println(s.v)
	println(s.x.val)
	println(s.y.val)

	var v Struct
	println(v.v)
	println(v.x.val)
	println(v.y.val)

	v.v = 7
	v.x.val = 19
	println(v.v)
	println(v.x.val)
	println(v.y.val)

	println(v.ValMethod()) // must not modify v
	println(v.v)
	println(v.x.val)
	println(v.y.val)
}
