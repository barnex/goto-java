package main

func main() {
	sum := 0
	for i := 0; i < 10; i++ {
		if i == 7 {
			continue
		}
		sum += i
	}
	println(sum)

	i := 0
	for {
		i++
		if i == 10 {
			break
		}
	}
	println(i)
}
