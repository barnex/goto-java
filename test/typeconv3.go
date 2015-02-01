package main

// Test for java's "cannot convert by method invocation":
// constant arguments need to be converted before passed to function.

func print_byte(x byte)     { println(x) }
func print_uint8(x uint8)   { println(x) }
func print_int8(x int8)     { println(x) }
func print_uint16(x uint16) { println(x) }
func print_int16(x int16)   { println(x) }
func print_uint(x uint)     { println(x) }
func print_int(x int)       { println(x) }
func print_uint32(x uint32) { println(x) }
func print_int32(x int32)   { println(x) }
func print_uint64(x uint64) { println(x) }
func print_int64(x int64)   { println(x) }

func main() {
	print_byte(1)
	print_uint8(2)
	print_int8(3)
	print_uint16(4)
	print_int16(5)
	print_uint(6)
	print_int(7)
	print_uint32(8)
	print_int32(9)
	print_uint64(12)
	print_int64(11)

	print_byte(1 + 1)
	print_uint8(2 + 1)
	print_int8(3 + 1)
	print_uint16(4 + 1)
	print_int16(5 + 1)
	print_uint(6 + 1)
	print_int(7 + 1)
	print_uint32(8 + 1)
	print_int32(9 + 1)
	print_uint64(12 + 1)
	print_int64(11 + 1)

	print_int8(-3)
	print_int16(-5)
	print_int(-7)
	print_int32(-9)
	print_int64(-11)
}
