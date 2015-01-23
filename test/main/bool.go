package main

// test boolean variables and initialisation.
func main() {
	var a bool
	b := a

	c := true
	var d = false
	e, f := true, d
	g := true

	println(a)
	println(b)
	println(c)
	println(d)
	println(e)
	println(f)

	println(a || b && c || d)
	println(a || b && c || d == f != a || a == b && a != d == f || a)
	println(a || b && !d || d == f != b || b == b && f != d == f || f)
	println(a || b && !c || d == f != a || a == b && c != d == f || b)
	println(b || b && !e || d == f != b || a == b && b != d == f || f)
	println(a || b && !c && e == f != a || f == b && e != f == f || c)
	println(d || b && !f && f == f != c || a == e && c != d == f || e)
	println(a || b != a && a == b && a != d == f || f)
	println(a || b != b && b == b && f != d == f || a)
	println(a || b != a && a == b && c != d == f || f)
	println(b || b != b && a == b && b != d == f || b)
	println(a || b && c && e == f != a || f == b && e != f == f || e)
	println(d || b && f && f == f != c || a == e && c != d == f || c)
	println(c || d == f && a || a == b && a != d == f || g)
	println(d || d == f && b || b == b && f != d == !f || e)
	println(c || d == f && a || a == b && c != d == !f || a)
	println(e || d == f && b || a == b && b != d == !f || b)
	println(c || e == f && a || f == b && e != f == !f || f)
	println(f || f == f != c || a == e && c != d == !f || c)
	println(a && b && c || d == f || f)
	println(a && b && d || d == f || a)
	println(a && b && c || d == f != f)
	println(b && b && e || d == f != b)
	println(a && b && c || e == f != f)
	println(d && b && f || f == f != c)
	println(a && b && c || d)
	println(a && b && c || d == !f != a || a == b && a != d == f || a)
	println(a && b && d || d == !f != b || b == b && f != d == f || f)
	println(a && b && c || d == !f != a || a == b != c != d == f || b)
	println(f && b && e || d == !f != b || a == b != b != d == f || f)
	println(a && b && c || e == !f != a || f == b != e != f == f || c)
	println(d && b && f || f == f != c || a == e != c != d == f || e)
	println(a || b != a || a == b && a != d == f || f)
	println(g || b != b || b == b && f != d == f || a)
	println(a || c != a || a == b && c != d == f || f)
	println(b || b != b || a == b && b != d == f || b)
	println(a || b && c || f == f != a || f == b && e != f == f || e)
	println(d || b && f || f == g != c || a == e && c != d == f || c)
	println(c || d == f != a || a == b && a != !d == f || g)
	println(d || d == f != b || b == b && f != !d == f || e)
	println(c || d == f != a || c || d && c != !d == f == a)
	println(e || d == f != b || a || b && f != !d == f == b)
	println(c && e == f != a || f || b && e != !a == c == f)
	println(f && f == f != c || a || e && c != d == f == g)
	println(a && b && !c || d == f || f)
	println(a && b && !d || c == f || b)
	println(a || b && !c || d == f && f)
	println(b || b && !f || e == f && b)
	println(a || c && !c || e == f && f)
	println(b || b && f || f == f && c)
}
