package gotojava

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"reflect"
)

// A writer takes AST nodes and outPuts java source
type Writer struct {
	out        *bufio.Writer
	f          io.WriteCloser
	indent     int
	needIndent bool
}

func NewWriter(fname string) *Writer {
	f, err := os.Create(fname)
	checkUserErr(err)
	out := bufio.NewWriter(f)
	return &Writer{out: out, f: f}
}

func (w *Writer) Close() {
	err := w.out.Flush()
	checkUserErr(err)
	err = w.f.Close()
	checkUserErr(err)
}

//func SPutln(tokens ...interface{})string{
//	return SPut(append(tokens, "\n")...)
//}
//
//func SPut(tokens ...interface{})string{
//	buf := bytes.NewBuffer()
//
//}

func (w *Writer) Putln(tokens ...interface{}) {
	w.Put(append(tokens, "\n")...)
	w.needIndent = true
}

func (w *Writer) Put(tokens ...interface{}) {
	w.putIndent()
	for _, t := range tokens {
		w.put(t)
	}
}

func (w *Writer) put(t interface{}) {
	switch t.(type) {
	case string, JModifier, token.Token:
		fmt.Fprint(w.out, t)
		return
	}

	if t, ok := t.(ast.Node); ok {
		w.PutNode(t)
		return
	}

	if t, ok := t.(JType); ok {
		w.Put(t.JavaName)
		return
	}
	panic("writer: cannot put type " + reflect.TypeOf(t).String())
}

func (w *Writer) PutNode(n ast.Node) {
	if n, ok := n.(ast.Stmt); ok {
		w.PutStmt(n)
		return
	}

	switch n := n.(type) {
	default:
		panic("putnode: need to handle: " + reflect.TypeOf(n).String())
	case ast.Expr:
		w.PutExpr(n)
	}
}

func (w *Writer) putIndent() {
	if w.needIndent == false {
		return
	}
	for i := 0; i < w.indent; i++ {
		fmt.Fprint(w.out, "\t")
	}
	w.needIndent = false
}

// exit with fatal error, print token position of node n and msg.
func Error(n ast.Node, msg ...interface{}) {
	panic(fmt.Sprint(append([]interface{}{PosOf(n), ": "}, msg...)...))
}

// prints parentisized argument list: "(elem[0], elem[1], ...)"
//func parens(elem ...interface{}) string {
//	str := "("
//	for i, e := range elem {
//		if i > 0 {
//			str += ", "
//		}
//		str += fmt.Sprint(e)
//	}
//	str += ")"
//	return str
//}
