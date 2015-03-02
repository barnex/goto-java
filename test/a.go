package main

//type B int
//
//type S struct {
//	a int
//	b *int
//	c B
//	d *B
//	e struct{ v int }
//	g *struct{ v int }
//	h struct{ int }
//	i *struct{ int }
//	j S2
//	k *S
//}
//
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
//
//type S2 struct {
//	v int
//}

func main() {

	unnamed_basic()

	//		l_b *int
	//		l_c B
	//		l_d *B
	//		l_e struct{ v int }
	//		l_g *struct{ v int }
	//		l_h struct{ int }
	//		l_i *struct{ int }
	//		l_j S2

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


	s1 := 1
	s2 := l_a2
	s3 := l_a2
	s4 := e_a1
	s5 := e_a1

	_ = &e_a1
	_ = &e_a2

	l_a1 = 2
	e_a1 = 3

	l_a1 = e_a1
	e_a1 = e_a2

	l_a1 = l_a1
	e_a1 = l_a2

	l_a1++
	e_a1++

	l_a1--
	e_a1--

	l_a1 += l_a2 + 2
	e_a1 += l_a2 + 2
               
	l_a1 -= l_a2 + 2
	e_a1 -= l_a2 + 2
              
	l_a1 *= l_a2 + 2
	e_a1 *= l_a2 + 2
             
	l_a1 /= l_a2 + 2
	e_a1 /= l_a2 + 2
            
	l_a1 += e_a2 + 2
	e_a1 += e_a2 + 2
           
	l_a1 -= e_a2 + 2
	e_a1 -= e_a2 + 2
          
	l_a1 *= e_a2 + 2
	e_a1 *= e_a2 + 2
         
	l_a1 /= e_a2 + 2
	e_a1 /= e_a2 + 2
        
	l_a1 %= l_a2 + 2
	e_a1 %= l_a2 + 2
       
	l_a1 &= l_a2 + 2
	e_a1 &= l_a2 + 2
      
	l_a1 |= l_a2 + 2
	e_a1 |= l_a2 + 2
     
	l_a1 ^= l_a2 + 2
	e_a1 ^= l_a2 + 2

	//l_a1 <<= e_a2 // TODO: uint
	//e_a1 <<= e_a2

	//l_a1 >>= e_a2
	//e_a1 >>= e_a2

	l_a1 &^= e_a2
	e_a1 &^= e_a2

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
