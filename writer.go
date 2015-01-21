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

const (
	CLASS   = "class"
	EOL     = ";"
	FINAL   = "final"
	LBRACE  = "{"
	PACKAGE = "package"
	PUBLIC  = "public"
	RBRACE  = "}"
	STATIC  = "static"
	VOID    = "void"
	LPAREN  = "("
	RPAREN  = ")"
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

// parsing

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

// code gen

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

	w.putln(PACKAGE, w.pkg, EOL)
	w.putln()

	className := path.Base(w.fname)
	w.putln(PUBLIC, FINAL, CLASS, className, LBRACE)
	w.putln()
	w.indent++

	w.genMembers()

	w.indent--
	w.putln(RBRACE)
}

func (w *writer) genMembers() {
	for _, n := range w.classMembers {
		w.putDecl(n)
	}
}

func (w *writer) putMainDecl(n *ast.FuncDecl) {
	w.put(PUBLIC, STATIC, VOID, n.Name.Name, parens("String[] args"))
	w.putBlockStmt(n.Body)
	w.putln()
}

func (w *writer) putBlockStmt(n *ast.BlockStmt) {
	w.putln(LBRACE)
	w.indent++

	for _, n := range n.List {
		w.putStmt(n)
	}

	w.indent--
	w.putln(RBRACE)
}

func (w *writer) putStmt(n ast.Stmt) {
	switch s := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(s))
	case *ast.ExprStmt:
		w.putExprStmt(s)
	case *ast.DeclStmt:
		w.putDeclStmt(s)
	}
}

func (w *writer) putDeclStmt(d *ast.DeclStmt) {
	switch d := d.Decl.(type) {
	default:
		w.error(d, "cannot handle ", reflect.TypeOf(d))
	case *ast.GenDecl:
		w.putGenDecl(d)
	}
}

func (w *writer) putGenDecl(d *ast.GenDecl) {
	for _, s := range d.Specs {
		w.putSpec(s)
	}
}

func (w *writer) putSpec(s ast.Spec) {
	switch s := s.(type) {
	default:
		w.error(s, "cannot handle ", reflect.TypeOf(s))
	case *ast.ValueSpec:
		w.putValueSpec(s)
	}
}

func (w *writer) putValueSpec(s *ast.ValueSpec) {
	w.putExpr(s.Type)

	for i, n := range s.Names {
		w.put(n.Name, "=")
		if i < len(s.Values) {
			w.putExpr(s.Values[i])
		} else {
			w.put(n.Obj.Data)
		}

		if i != len(s.Names)-1 {
			w.put(", ")
		}
	}
	w.put(";")
	w.putComment(s.Comment)
	w.putln()
}

func (w *writer) putExprStmt(n *ast.ExprStmt) {
	w.putExpr(n.X)
	//w.putComment(n.Comment)
	w.putln(";")
}

func (w *writer) putBasicLit(n *ast.BasicLit) {
	w.put(n.Value)
	// TODO: translate backquotes, complex etc.
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
