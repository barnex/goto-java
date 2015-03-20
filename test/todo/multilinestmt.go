package main

// Some single-line Go statements are translated to multi-line java statements.
// This can be problematic with for clauses etc.

func main() {
	for i, j := f(); i < 10; println(i, j) {
		i++
	}
}

func f() (int, int) { return 1, 2 }
