package main

// test map with pointer type keys
// equals method should not mistakingly compare values instead of addresses

func main(){
	m := make(map[*int]int)
	a, b, c := 1, 1, 2

	m[&a] = 1
	m[&b] = 2
	m[&c] = 3
	
	println(m[&a])
	println(m[&b])
	println(m[&c])
}
