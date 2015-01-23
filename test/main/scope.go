package main

func main() {
	// test scope

	i := 1
	{
		i := 2
		println(i)
	}
	println(i)

	{
		j := 3
		println(j)
	}

	j := 4
	println(j)
}
