package main

func main() {

	for i := 0; i < 5; i++ {
		switch i {
		default:
			println("default")
		case 1:
			println("one")
		case 2:
			println("two")
			fallthrough
		case 3,4:
			println("three or four")
		}
	}

	for i := 0; i < 5; i++ {
		switch i {
		case 1:
			println("one")
		case 2:
			println("two")
		case 0,3:
			println("and three")
			fallthrough
		default:
			println("default")
		}
	}

}
