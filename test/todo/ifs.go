package main

func main() {
	if 1 != 2 {
		println("all is fine")
	}

	if 1 == 2 {
		println("cosmic rays")
	} else {
		println("good")
	}

	// test init:
	if i := 7; i == 7 {
		println("ok")
	}
	i := 8

	println(i)

}
