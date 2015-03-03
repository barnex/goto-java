package main

//type B int

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

func named_struct() {
	var (
		s1 S
		s2 S = S{}
		s3 S = S{3, nil, S2{}, nil}
		s4 S = S{k: nil, b: new(int), a: 7}
	)

	println(s1.a)
	println(s2.a)
	println(s3.a)
	println(s4.a)

	println((&s1).a)
	println((&s2).a)
	println((&s3).a)
	println((&s4).a)
}

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
	_ = p4
	_ = p5

}

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

}
