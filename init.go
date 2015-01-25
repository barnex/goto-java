package main

// ZeroValue returns the zero value for a new variable of java type jType.
// E.g.:
// 	var x int  ->  int x = 0;
func ZeroValue(jType string) string {
	if v, ok := zeroValues[jType]; ok {
		return v
	} else {
		return "null"
	}
}

var zeroValues = map[string]string{
	"String":  `""`,
	"int":     "0",
	"float":   "0.0f",
	"double":  "0.0",
	"boolean": "false",
}
