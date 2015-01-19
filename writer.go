package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path"
	"reflect"
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

	fname        string // output file name for principal class (no extension)
	pkg          string
	classMembers []ast.Decl
}

// parsing

func (w *writer) parseFile(f *ast.File) {
	w.pkg = f.Name.Name

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
		switch m := n.(type) {
		default:
			panic("unhandled memeber type: " + reflect.TypeOf(m).String())
		case *ast.FuncDecl:
			w.putFuncDecl(m)
		}
	}
}

func (w *writer) putFuncDecl(n *ast.FuncDecl) {
	if n.Name.Name == "main" {
		w.putMainDecl(n)
		return
	}

	panic("todo")
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
		w.putln(";")
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
	for _, n := range s.Names {
		w.putExpr(s.Type)
		w.putln(n.Name, "=", n.Obj.Data, ";")
	}
}

func (w *writer) putExprStmt(n *ast.ExprStmt) {
	w.putExpr(n.X)
}

func (w *writer) putExpr(n ast.Expr) {
	switch e := n.(type) {
	default:
		w.error(n, "cannot handle ", reflect.TypeOf(e))
	case *ast.CallExpr:
		w.putCallExpr(e)
	case *ast.Ident:
		w.putIdent(e)
	case *ast.BasicLit:
		w.putBasicLit(e)
	case *ast.BinaryExpr:
		w.putBinaryExpr(e)
	case *ast.ParenExpr:
		w.putParenExpr(e)
	}
}

func (w *writer) putParenExpr(e *ast.ParenExpr) {
	w.put(LPAREN)
	w.putExpr(e.X)
	w.put(RPAREN)
}

func (w *writer) putBinaryExpr(b *ast.BinaryExpr) {
	// TODO: check unsupported ops
	w.putExpr(b.X)
	w.put(b.Op)
	w.putExpr(b.Y)
}

func (w *writer) putCallExpr(n *ast.CallExpr) {
	w.putExpr(n.Fun)

	w.put("(")
	for i, a := range n.Args {
		if i != 0 {
			w.put(",")
		}
		w.putExpr(a)
	}
	w.put(")")

	if n.Ellipsis != 0 {
		w.error(n, "cannot handle ellipsis")
	}
}

var keywordMap = map[string]string{
	"println": "System.out.println",
	"print":   "System.out.print",
}

func (w *writer) putIdent(n *ast.Ident) {
	name := n.Name
	// translate name if keyword
	if trans, ok := keywordMap[name]; ok {
		name = trans
	}
	w.put(name)
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
	fatal(append([]interface{}{w.pos(n), ": "}, msg...)...)
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
