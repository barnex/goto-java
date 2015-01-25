package main

func len() {

}

// Test shadowing of built-ins
func main() {

	//println(true)
	//true := false
	//println(true)

	println := 1
	print(println)

	len := 2
	print(len)
}
