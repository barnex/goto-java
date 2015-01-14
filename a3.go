package main

import (
	"flag"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	var fset token.FileSet
	pkgs, err := parser.ParseDir(&fset, path, nil, 0)
	checkUserErr(err)

	if len(pkgs) == 0 {
		fatal("no packages found in ", path)
	}
	if len(pkgs) > 1 {
		fatal(len(pkgs), " packages found in ", path)
	}
}

func checkUserErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(msg ...interface{}) {
	log.Fatal(msg...)
}
