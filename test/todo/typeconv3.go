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

func return_byte() byte     { return 10 }
func return_uint8() uint8   { return 20 }
func return_int8() int8     { return 30 }
func return_uint16() uint16 { return 40 }
func return_int16() int16   { return 50 }
func return_uint() uint     { return 60 }
func return_int() int       { return 70 }
func return_uint32() uint32 { return 80 }
func return_int32() int32   { return 90 }
func return_uint64() uint64 { return 100 }
func return_int64() int64   { return 110 }

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
	print_int64(11)
	print_uint64(12)

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

	print_byte(return_byte())
	print_uint8(return_uint8())
	print_int8(return_int8())
	print_uint16(return_uint16())
	print_int16(return_int16())
	print_uint(return_uint())
	print_int(return_int())
	print_uint32(return_uint32())
	print_int32(return_int32())
	print_uint64(return_uint64())
	print_int64(return_int64())
}
