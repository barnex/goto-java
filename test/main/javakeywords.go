package main

var (
	static_2  = 3 // try to trigger collision with renamed variable
	private_1 = 4 // try to trigger collision with renamed variable
	new       = 4
	static    = 7 // try to trigger global/local collision
	final     = 8 // try to trigger global/local collision
)

func private(a int, b int) int {
	return a + b
}

func main() {

	println(static)
	static := 1 // use java keyword as identifier
	println(static)

	static_1 := 2 // try to trigger collision with renamed variable
	println(static_1)

	public := 2
	println(public)

	println(private(5, 7))

	private_2 := 6
	private := 3
	println(private)
	println(private_2)

	final := 4
	println(final)

	class := 5
	println(class)

	println(new)
	new := 6
	println(new)

}
