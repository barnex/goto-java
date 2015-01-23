package main

var (
	i = 47
	j = 48
)

func main() {
	// test scope
	i := 1
	{
		i := 2
		{
			i := 5
			println(i)
		}
		println(i)
	}
	println(i)

	{
		j := 3
		println(j)
	}

	j := 4
	println(j)

	println(getI())
	println(getJ())
}

func getI() int { return i }
func getJ() int { return j }
