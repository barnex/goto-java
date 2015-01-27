package main

func main() {
	var a, b, c, d uint16 = (1 << 15) - 1, (1 << 15), (1 << 15) + 1, 0

	println(a / a)
	println(a / b)
	println(a / c)

	println(b / a)
	println(b / b)
	println(b / c)

	println(c / a)
	println(c / b)
	println(c / c)

	println(d / a)
	println(d / b)
	println(d / c)

	println(a > a)
	println(a > b)
	println(a > c)
	println(a > d)
	println(a > 0)

	println(b > a)
	println(b > b)
	println(b > c)
	println(b > d)
	println(b > 0)

	println(c > a)
	println(c > b)
	println(c > c)
	println(c > d)
	println(c > 0)

	println(d > a)
	println(d > b)
	println(d > c)
	println(d > d)
	println(d > 0)

	println(0 > a)
	println(0 > b)
	println(0 > c)
	println(0 > d)
	println(0 > 0)

	println(a < a)
	println(a < b)
	println(a < c)
	println(a < d)
	println(a < 0)

	println(b < a)
	println(b < b)
	println(b < c)
	println(b < d)
	println(b < 0)

	println(c < a)
	println(c < b)
	println(c < c)
	println(c < d)
	println(c < 0)

	println(d < a)
	println(d < b)
	println(d < c)
	println(d < d)
	println(d < 0)

	println(0 < a)
	println(0 < b)
	println(0 < c)
	println(0 < d)
	println(0 < 0)

	println(a >= a)
	println(a >= b)
	println(a >= c)
	println(a >= d)
	println(a >= 0)

	println(b >= a)
	println(b >= b)
	println(b >= c)
	println(b >= d)
	println(b >= 0)

	println(c >= a)
	println(c >= b)
	println(c >= c)
	println(c >= d)
	println(c >= 0)

	println(d >= a)
	println(d >= b)
	println(d >= c)
	println(d >= d)
	println(d >= 0)

	println(0 >= a)
	println(0 >= b)
	println(0 >= c)
	println(0 >= d)
	println(0 >= 0)

	println(a <= a)
	println(a <= b)
	println(a <= c)
	println(a <= d)
	println(a <= 0)

	println(b <= a)
	println(b <= b)
	println(b <= c)
	println(b <= d)
	println(b <= 0)

	println(c <= a)
	println(c <= b)
	println(c <= c)
	println(c <= d)
	println(c <= 0)

	println(d <= a)
	println(d <= b)
	println(d <= c)
	println(d <= d)
	println(d <= 0)

	println(0 <= a)
	println(0 <= b)
	println(0 <= c)
	println(0 <= d)
	println(0 <= 0)
}
