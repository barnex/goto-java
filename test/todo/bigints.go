package main

const (
	big1 = 1 << 31
	big2 = 1 << 32
	big3 = 1 << 33
	big4 = 1 >> 31
	big5 = 1 >> 32
	big6 = 1 >> 33
)

func main() {
	println(1 << 31)
	println(1 << 32)
	println(1 << 33)
	println(1 >> 31)
	println(1 >> 32)
	println(1 >> 33)

	println(big1)
	println(big2)
	println(big3)
	println(big4)
	println(big5)
	println(big6)
}
