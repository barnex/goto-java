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
