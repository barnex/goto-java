package main

type S struct{ v int }
type T struct{ v int }

//type U S

func main() {
	var (
		x1 struct{ v int }
		x2 struct{ v int }
		s1 S
		s2 S
		t1 T
		t2 T
	)

	println(x1.v)
	println(x2.v)
	println(s1.v)
	println(s2.v)
	println(t1.v)
	println(t2.v)

	var sptr *S = new(S)
	var tptr *T = new(T)

	tptr = (*T)(sptr)
	sptr = (*S)(tptr)

	sptr.v = 17
	println(sptr.v)
	println(tptr.v)
}
