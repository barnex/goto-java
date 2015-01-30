package main

// try to confuse handling of built-ins by parenthesizing.
func main() {
	println(1)
	(println)(2)
	((println))(3)
}
