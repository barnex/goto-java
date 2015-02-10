package main

import (
	"flag"
	"log"
)

var (
	flagBlank       = flag.String("blank", "_", "Java name for the blank (underscore) identifier")
	flagConstFold   = flag.Bool("foldconst", false, "Fold constants")
	flagNoPkg       = flag.Bool("nopkg", false, "Do not output package clause")
	flagNoTypeCheck = flag.Bool("nocheck", false, "Don't do type check")
	flagParens      = flag.Bool("parens", false, "Emit superfluous parens")
	flagPrint       = flag.Bool("print", false, "Print ast")
	flagRenameAll   = flag.Bool("renameall", false, "Rename all variables (debug)")
	UNUSED          string // base name for translating the blank identifier (flag -blank)
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	UNUSED = *flagBlank

	for _, f := range flag.Args() {
		HandleFile(f)
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
