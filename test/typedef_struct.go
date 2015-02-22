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

	var s Struct
	var sptr := &s

}
