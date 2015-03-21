package gotojava

import (
	"go/ast"
	"log"
	"path"
	"runtime"
)

// Check error. Fatal exit if err != nil.
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Log message msg, with position info from n.
func Log(n ast.Node, msg ...interface{}) {
	if *flagVerbose {
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fname := path.Ext(runtime.FuncForPC(pc).Name()) // strip package prefix
			fname = fname[1:]                               // strip "."
			msg = append([]interface{}{fname + ": "}, msg...)
		}
		if n != nil {
			msg = append([]interface{}{PosOf(n).String() + ": "}, msg...)
		}
		log.Println(msg...)
	}
}

func assert(test bool) {
	if !test {
		panic("assertion failed")
	}
}
