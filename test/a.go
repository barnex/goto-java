package main

type B int

type S struct {
	a int
	b *int
	//c B
	//d *B
	//e struct{ v int }
	//g *struct{ v int }
	//h struct{ int }
	//i *struct{ int }
	j S2
	k *S
}

func (s S) sum() int {
	sum := s.a
	sum += *s.b
	sum += s.j.v
	//sum += s.k.a
	// try to mutate receiver value
	s.a = 88
	s.b = new(int)
	s.j = S2{77}
	s.k = new(S)
	return sum
}

//var (
//	g_a int
//	g_b *int
//	g_c B
//	g_d *B
//	g_e struct{ v int }
//	g_g *struct{ v int }
//	g_h struct{ int }
//	g_i *struct{ int }
//	g_j S2
//	g_k *S
//)

type S2 struct {
	v int
}

func main() {

	unnamed_basic()
	named_basic()

	named_struct()

	unnamed_pointer()

	//		l_b *int
	//		l_c B
	//		l_d *B
	//		l_e struct{ v int }
	//		l_g *struct{ v int }
	//		l_h struct{ int }
	//		l_i *struct{ int }
	//		l_j S2

}

func unnamed_struct() {
	//var (
	//	s1 struct{v int}
	//	s2 struct{v int} = struct{v int}{}
	//	s3 struct{v int} = struct{v int}{3, nil, S2{}, nil}
	//	s4 struct{v int} = struct{v int}{k: nil, b: new(int), a: 7}
	//)
}

func makeSPtr() *S {
	return &S{a: 12}
}

func named_struct_pointer() {
	var s S
	var (
		s1 *S
		s2 *S = &S{}
		s3 *S = &S{3, nil, S2{}, nil}
		s4 *S = &S{k: &S{a: 19}, b: new(int), a: 7}
		s5 *S = makeSPtr()
		s6 *S = &s
	)
	s7 := &S{}
	s8 := &S{3, nil, S2{}, nil}
	s9 := &S{k: nil, b: new(int), a: 7}
	s10 := makeSPtr()

	println(s1.a)
	println(s2.a)
	println(s3.a)
	println(s4.a)
	println(s5.a)
	println(s6.a)
	println(s7.a)
	println(s8.a)
	println(s9.a)
	println(s10.a)

	s1 = s6
	s1.b = s6.b

	println(s1 == s1)
	println(s1 == s2)
	println(s1 == s3)
	println(s1 == s4)
	println(s1 == s5)
	println(s1 == s6)
}

func makeS() S {
	return S{a: 12}
}

func named_struct() {
	var (
		s1 S
		s2 S = S{}
		s3 S = S{3, nil, S2{}, nil}
		s4 S = S{k: &S{a: 12}, b: new(int), a: 7}
		s5 S = makeS()
	)
	s6 := S{}
	s7 := S{3, nil, S2{}, nil}
	s8 := S{k: nil, b: new(int), a: 7}
	s9 := makeS()

	println(s1.a)
	println(s2.a)
	println(s3.a)
	println(s4.a)
	println(s5.a)
	println(s6.a)
	println(s7.a)
	println(s8.a)
	println(s9.a)

	// value method should not mutate
	println(s4.a)
	println(*(s4.b))
	println(s4.k.a)
	println(s4.j.v)

	println(s4.sum())

	println(s4.a)
	println(*(s4.b))
	println(s4.k.a)
	println(s4.j.v)

	s1 = s6
	s1.b = s6.b

	println(s1 == s1)
	println(s1 == s2)
	println(s1 == s3)
	println(s1 == s4)
	println(s1 == s5)
	println(s1 == s6)

	println((&s1).a)
	println((&s2).a)
	println((&s3).a)
	println((&s4).a)
}

func eat_intptr(*int) {}

func unnamed_pointer() {
	var (
		i  int
		p1 *int
		p2 *int = nil
		p3 *int = new(int)
		p4 *int = p3
		p5 *int = &i
	)

	p6 := p1
	p7 := &i

	println(p6 == p1)
	println(p6 == p7)
	println(p6 == &i)

	p1 = nil
	p1 = new(int)
	p1 = p2
	p1 = &i

	*p1 = *p1
	*p1 = 1

	println(*p1)
	*p1++
	println(*p1)
	*p1--
	println(*p1)
	*p1 += 1
	println(*p1)
	*p1 -= 2
	println(*p1)
	*p1 *= 3
	println(*p1)
	*p1 /= 4
	println(*p1)
	*p1 %= 5
	println(*p1)
	*p1 &= 6
	println(*p1)
	*p1 |= 7
	println(*p1)
	*p1 ^= 8
	println(*p1)
	*p1 &^= 9
	println(*p1)

	// use variables
	eat_intptr(p4)
	_ = p5

}

func eat_int(int) {}

func unnamed_basic() {
	var (
		l_a1 int
		l_a2 int = 1
		l_a3 int = l_a2
		e_a1 int = l_a2
		e_a2 int = e_a1
		l_a4 int = e_a1
	)
	println(l_a1)
	println(l_a2)
	println(l_a3)
	println(e_a1)
	println(e_a2)
	println(l_a4)

	s1 := 1
	s2 := l_a2
	s3 := l_a2
	s4 := e_a1
	s5 := e_a1

	println(s1)
	println(s2)
	println(s3)
	println(s4)
	println(s5)

	// make them escape
	_ = &e_a1
	_ = &e_a2

	l_a1 = 2
	e_a1 = 3
	println(l_a1)
	println(e_a1)

	l_a1 = e_a1
	e_a1 = e_a2
	println(l_a1)
	println(e_a1)

	l_a1 = l_a1
	e_a1 = l_a2
	println(l_a1)
	println(e_a1)

	l_a1++
	e_a1++
	println(l_a1)
	println(e_a1)

	l_a1--
	e_a1--
	println(l_a1)
	println(e_a1)

	l_a1 += l_a2 + 2
	e_a1 += l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 -= l_a2 + 2
	e_a1 -= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 *= l_a2 + 2
	e_a1 *= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 /= l_a2 + 2
	e_a1 /= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 += e_a2 + 2
	e_a1 += e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 -= e_a2 + 2
	e_a1 -= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 *= e_a2 + 2
	e_a1 *= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 /= e_a2 + 2
	e_a1 /= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 %= l_a2 + 2
	e_a1 %= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 &= l_a2 + 2
	e_a1 &= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 |= l_a2 + 2
	e_a1 |= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 ^= l_a2 + 2
	e_a1 ^= l_a2 + 2
	println(l_a1)
	println(e_a1)

	//l_a1 <<= e_a2 // TODO: uint
	//e_a1 <<= e_a2

	//l_a1 >>= e_a2
	//e_a1 >>= e_a2

	l_a1 &^= e_a2
	e_a1 &^= e_a2
	println(l_a1)
	println(e_a1)

	println(l_a1 == e_a2)
	println(e_a1 == e_a2)
	println(l_a1 < e_a2)
	println(e_a1 < e_a2)
	println(l_a1 > e_a2)
	println(e_a1 > e_a2)
	println(l_a1 != e_a2)
	println(e_a1 != e_a2)
	println(l_a1 >= e_a2)
	println(e_a1 >= e_a2)
	println(l_a1 <= e_a2)
	println(e_a1 <= e_a2)

	println(l_a1 == l_a2)
	println(e_a1 == l_a2)
	println(l_a1 < l_a2)
	println(e_a1 < l_a2)
	println(l_a1 > l_a2)
	println(e_a1 > l_a2)
	println(l_a1 != l_a2)
	println(e_a1 != l_a2)
	println(l_a1 >= l_a2)
	println(e_a1 >= l_a2)
	println(l_a1 <= l_a2)
	println(e_a1 <= l_a2)

	println(l_a1 == 1+1)
	println(e_a1 == 1+1)
	println(l_a1 < 1+1)
	println(e_a1 < 1+1)
	println(l_a1 > 1+1)
	println(e_a1 > 1+1)
	println(l_a1 != 1+1)
	println(e_a1 != 1+1)
	println(l_a1 >= 1+1)
	println(e_a1 >= 1+1)
	println(l_a1 <= 1+1)
	println(e_a1 <= 1+1)

	println(l_a1)
	println(l_a2)
	println(l_a3)
	println(e_a1)
	println(e_a2)
	println(l_a4)

	println(s1)
	println(s2)
	println(s3)
	println(s4)
	println(s5)

	eat_int(s5)
}

func named_basic() {
	var (
		l_a1 B
		l_a2 B = 1
		l_a3 B = l_a2
		e_a1 B = l_a2
		e_a2 B = e_a1
		l_a4 B = e_a1
	)
	println(l_a1)
	println(l_a2)
	println(l_a3)
	println(e_a1)
	println(e_a2)
	println(l_a4)

	//s1 := 1
	s2 := l_a2
	s3 := l_a2
	s4 := e_a1
	s5 := e_a1

	//println(s1)
	println(s2)
	println(s3)
	println(s4)
	println(s5)

	// make them escape
	_ = &e_a1
	_ = &e_a2

	l_a1 = 2
	e_a1 = 3
	println(l_a1)
	println(e_a1)

	l_a1 = e_a1
	e_a1 = e_a2
	println(l_a1)
	println(e_a1)

	l_a1 = l_a1
	e_a1 = l_a2
	println(l_a1)
	println(e_a1)

	l_a1++
	e_a1++
	println(l_a1)
	println(e_a1)

	l_a1--
	e_a1--
	println(l_a1)
	println(e_a1)

	l_a1 += l_a2 + 2
	e_a1 += l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 -= l_a2 + 2
	e_a1 -= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 *= l_a2 + 2
	e_a1 *= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 /= l_a2 + 2
	e_a1 /= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 += e_a2 + 2
	e_a1 += e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 -= e_a2 + 2
	e_a1 -= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 *= e_a2 + 2
	e_a1 *= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 /= e_a2 + 2
	e_a1 /= e_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 %= l_a2 + 2
	e_a1 %= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 &= l_a2 + 2
	e_a1 &= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 |= l_a2 + 2
	e_a1 |= l_a2 + 2
	println(l_a1)
	println(e_a1)

	l_a1 ^= l_a2 + 2
	e_a1 ^= l_a2 + 2
	println(l_a1)
	println(e_a1)

	//l_a1 <<= e_a2 // TODO: uint
	//e_a1 <<= e_a2

	//l_a1 >>= e_a2
	//e_a1 >>= e_a2

	l_a1 &^= e_a2
	e_a1 &^= e_a2
	println(l_a1)
	println(e_a1)

	println(l_a1 == e_a2)
	println(e_a1 == e_a2)
	println(l_a1 < e_a2)
	println(e_a1 < e_a2)
	println(l_a1 > e_a2)
	println(e_a1 > e_a2)
	println(l_a1 != e_a2)
	println(e_a1 != e_a2)
	println(l_a1 >= e_a2)
	println(e_a1 >= e_a2)
	println(l_a1 <= e_a2)
	println(e_a1 <= e_a2)

	println(l_a1 == l_a2)
	println(e_a1 == l_a2)
	println(l_a1 < l_a2)
	println(e_a1 < l_a2)
	println(l_a1 > l_a2)
	println(e_a1 > l_a2)
	println(l_a1 != l_a2)
	println(e_a1 != l_a2)
	println(l_a1 >= l_a2)
	println(e_a1 >= l_a2)
	println(l_a1 <= l_a2)
	println(e_a1 <= l_a2)

	println(l_a1 == 1+1)
	println(e_a1 == 1+1)
	println(l_a1 < 1+1)
	println(e_a1 < 1+1)
	println(l_a1 > 1+1)
	println(e_a1 > 1+1)
	println(l_a1 != 1+1)
	println(e_a1 != 1+1)
	println(l_a1 >= 1+1)
	println(e_a1 >= 1+1)
	println(l_a1 <= 1+1)
	println(e_a1 <= 1+1)

	println(l_a1)
	println(l_a2)
	println(l_a3)
	println(e_a1)
	println(e_a2)
	println(l_a4)

	//println(s1)
	println(s2)
	println(s3)
	println(s4)
	println(s5)

	//eat_int(s5)
}
