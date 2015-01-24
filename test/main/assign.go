package main

func main(){
	a := 5
	b := 1
	println(a)
	println(b)
	
	a++
	println(a)
	println(b)

	a--
	println(a)
	println(b)

	a += 2
	b += 3
	println(a)
	println(b)

	a -= 3
	b -= 7
	println(a)
	println(b)

	a *= 4
	b *= -4
	println(a)
	println(b)

	a /= 5
	b /= -6
	println(a)
	println(b)

	a %= 6
	b %= -7
	println(a)
	println(b)

	a ^= 7
	b ^= -8
	println(a)
	println(b)

	a &= 9
	b &= -10
	println(a)
	println(b)

	a |= 10
	b |= -11
	println(a)
	println(b)

	a <<= 11
	b <<= 3
	println(a)
	println(b)

	a >>= 12
	b >>= 4
	println(a)
	println(b)

	// TODO
	//a &^= 13
	//b &^= -14
	//println(a)
	//println(b)

}
