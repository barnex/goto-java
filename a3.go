package main

import (
	"flag"
	"fmt"
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
	fmt.Println(pkgs)
}

func checkUserErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
