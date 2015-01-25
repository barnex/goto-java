package main

// test variable renaming
func main() {
	true_ := 1 // try to mix up with rename true -> true_
	println(true)
	true := false
	println(true)
	println(true_)
}
