package main

func print(a interface{}) {
	println(a)
}

func main() {

	var a interface{}
	a = 1
	//a = byte(1)
	//a = "hello"

	print(a)
}
