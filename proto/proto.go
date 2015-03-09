//
//
package main

var (
	gi int
)

func main() {

	// basic

	var i int
	var i2 int = 2

	i = 7
	i2 = i

	println(i)
	println(i2)

	// escaped
	e := 3

	i = e
	e = i2
	e++

	println(e)

	// pointer->basic
	eptr := &e
	println(*eptr)

	giptr := &gi

	giptr2 := &gi

	*giptr = 3
	println(*giptr2)
	println(giptr == giptr2)

	gi := &gi
	println(giptr == gi)

	// struct
	var xs struct{ v int }
	var xs2 = struct{ v int }{7}
	xs = xs2

	println(xs.v)
	println(xs2.v)
	println(xs == xs2)

	// pointer->struct
	var xsptr *struct{ v int } = new(struct{ v int })
	var xsptr2 = &xs2

	println(xsptr2 == &xs2)
	println(xsptr == &xs2)

	println(*xsptr == *(&xs2))
	*xsptr = xs2
	println(*xsptr == xs2)

	// named->basic

	var mi MyInt

	println(mi)
	println(mi.Square())

	// named->struct

	var s S
	println(s.v)
	s = S(xs)
	println(s.v)
	xs.v = 99
	xs = struct{ v int }(s)
	println(s.Square())

	var t T = T(s)
	println(t.v)

	// pointer->named->struct
	sptr := &s
	s.Inc()
	println(sptr == &s)
	println(s.Square())

	sptr = (*S)(&xs)
	println(sptr == (*S)(&xs))
	*sptr = S{33}
	println(xs.v)

	// interface
	var any interface{}

	any = i
	println(any == 3)
	_, ok := any.(int)
	println(ok)
	_, ok = any.(*int)
	println(ok)

	any = e
	println(any == 3)
	any = eptr
	println(any == 3)
	//...

	// func

	var ncalls int
	var f func(int) int
	f = func(x int) int {
		ncalls++
		return x * x
	}
	a := f
	println(a(3))
	println(ncalls)

}

type MyInt int

func (i MyInt) Square() int { return int(i * i) }

type S struct{ v int }

func (s S) Square() int { return s.v * s.v }
func (s *S) Inc()       { s.v++ }

type T S
