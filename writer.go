package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"
	"reflect"

	"golang.org/x/tools/go/types"
)

type writer struct {
	out        io.Writer
	indent     int
	needIndent bool
	needSpace  bool

	fset *token.FileSet
	types.Info

	fname        string // output file name for principal class (no extension)
	pkg          string
	classMembers []ast.Decl
}

func (w *writer) parseFile(f *ast.File) {
	w.pkg = f.Name.Name

	//w.putDocComment(f.Doc)

	for _, decl := range f.Decls {
		switch n := decl.(type) {
		default:
			w.error(f, "cannot handle ", reflect.TypeOf(decl))
		case *ast.FuncDecl:
			w.classMembers = append(w.classMembers, n)
		}
	}
}

func (w *writer) initOut() {
	if w.out != nil {
		panic("already inited")
	}
	out, err := os.Create(w.fname + ".java")
	checkUserErr(err)
	w.out = out
	// TODO: buffer
}

func (w *writer) genCode() {
	w.initOut()

	w.putln("package", w.pkg, ";")
	w.putln()

	className := path.Base(w.fname)
	w.putln("public final class", className, "{")
	w.putln()
	w.indent++

	w.genMembers()

	w.indent--
	w.putln("}")
}

func (w *writer) genMembers() {
	for _, n := range w.classMembers {
		w.putDecl(n)
	}
}

func (w *writer) putMainDecl(n *ast.FuncDecl) {
	w.put("public static void", n.Name.Name, "(String[] args)")
	w.putBlockStmt(n.Body)
	w.putln()
}

// fmt utils

func (w *writer) putln(tokens ...interface{}) {
	w.put(append(tokens, "\n")...)
	w.needIndent = true
	w.needSpace = false
}

var noSpaceAround = map[string]bool{
	";":  true,
	"\n": true,
	"(":  true,
	")":  true,
	",":  true,
	", ": true,
}

func (w *writer) put(tokens ...interface{}) {
	w.putIndent()
	for _, t := range tokens {
		t := fmt.Sprint(t)
		if noSpaceAround[t] {
			w.needSpace = false
		}
		if w.needSpace {
			fmt.Fprint(w.out, " ")
		}
		fmt.Fprint(w.out, t)
		w.needSpace = !noSpaceAround[t]
	}
}

func (w *writer) putIndent() {
	if w.needIndent == false {
		return
	}
	for i := 0; i < w.indent; i++ {
		fmt.Fprint(w.out, "\t")
	}
	w.needIndent = false
}

// exit with fatal error, print token position of node n and msg.
func (w *writer) error(n ast.Node, msg ...interface{}) {
	panic(fmt.Sprint(append([]interface{}{w.pos(n), ": "}, msg...)...))
}

// utils

// return position of node using this writer's fileset
func (w *writer) pos(n ast.Node) token.Position {
	return w.fset.Position(n.Pos())
}

// prints parentisized argument list: "(elem[0], elem[1], ...)"
func parens(elem ...interface{}) string {
	str := "("
	for i, e := range elem {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprint(e)
	}
	str += ")"
	return str
}
