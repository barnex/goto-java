package main

func main() {
	f1()
	f2()
}

func f1() {
	const true = false
	println(true)
}

func f2() {
	const false = true
	println(false)
}
