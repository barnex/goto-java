package main

// returns msg if i!=0
func ifNonzero(i int, msg string) string {
	if i != 0 {
		return msg
	} else {
		return ""
	}
}

// returns a comma if i!=0
func comma(i int) string {
	return ifNonzero(i, ", ")
}

// concatenate a and b, inserting a space if needed
func cat(a, b string) string {
	if (a != "") && (b != "") {
		return a + " " + b
	} else {
		return a + b
	}
}

func nnil(x interface{}) interface{} {
	if x == nil {
		return ""
	} else {
		return x
	}
}
