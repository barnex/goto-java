package main

var zeroValues = map[string]string{
	"String":  `""`,
	"int":     "0",
	"float":   "0.0f",
	"double":  "0.0",
	"boolean": "false",
}

func ZeroValue(jType string) string {
	if v, ok := zeroValues[jType]; ok {
		return v
	} else {
		return "null"
	}
}
