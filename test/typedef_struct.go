package main

var s Struct

// A test struct.
type Struct struct {
	a, b int
	s    string
}

func (s Struct) ValMethod1() {}

func (s Struct) ValMethod2(a, b int) int {
	return s.a + b
}

func (Struct) valmethod() {}

func (s *Struct) PtrMethod1() {}

func (s *Struct) PtrMethod2(a, b int) {
	s.a = a
	s.b = b
}

func (*Struct) ptrmethod() {}

func main() {
	var s Struct

	s.ValMethod1()
	println(s.ValMethod2(1, 2))
	s.valmethod()

	s.PtrMethod1()
	s.PtrMethod2(1, 2)
	println(s.a)
	s.ptrmethod()
}
