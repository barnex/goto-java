package main

func main() {


	// TODO: with consts, vars int, int64, int32, ...

	println(42)
	println(1 + 1)
	println(1 + 2*3)
	println(1*2 + 3)
	println(1*2 + 3 - 4/5*6)
	println((1 + 1) * 2)
	println(1 << 2)
	println(11 % 3)
	println(1<<8 + 2<<16)
	println(1 << 2 >> 3)
	println(1<<2>>3 | 4)
	println(1<<2>>3 | 4&5)
	println(1<<2>>3 | 4&5 | 6)
	println(1<<2>>3 | 4&5 | 6&7)
	println(1<<2>>3 | 4&5 | 6&7<<1)
	println(1<<2>>3 | 4&5 | 6&7<<1 ^ 1 | 3 | 4&5 ^ 6&7)
	println(1*2/3%4 + 5 - 6>>1<<2&3 | 4)
	println(1 | 2 ^ 3&4<<5>>6 - 1 + 2*3/4)
	println(3 &^ 4)

	// TODO: sizeof int
	//println(1 << 31)
	//println(1 << 32)
	//println(1 << 33)
	//println(1 >> 31)
	//println(1 >> 32)
	//println(1 >> 33)
}
