package main

// A test struct.
var s Struct

type Struct struct {
	a, b int
	s    string
}

func (s Struct) Method1() {}

func (s Struct) Method2(a, b int) int {
	return s.a + b
}

func (Struct) method3() {}

func (s *Struct) PtrMethod1() {}

func (s *Struct) PtrMethod2(a, b int) {
	s.a = a
	s.b = b
}

func (*Struct) ptrMethod3() {}

func main() {

}
