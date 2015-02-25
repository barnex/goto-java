package gotojava

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"reflect"
)

// Convenience function transpiles (presumably tiny) code snippets to java.
// Inteded to feed to Writer.Putf(). E.g.:
// 	w.Putf(`public class %s {`, spec.Name)
// Large code blocks should be transpiled with Writer.Put(...).
func Transpile(tokens ...interface{}) string {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	defer w.Close()
	w.Put(tokens...)
	return buf.String()
}

func (w *Writer) Putf(format string, tokens ...interface{}) {
	compiled := make([]interface{}, 0, len(tokens))
	for i := range tokens {
		compiled = append(compiled, Transpile(tokens[i]))
	}
	fmt.Fprintf(w.out, format, compiled...)
}

// A writer takes AST nodes and outPuts java source
type Writer struct {
	out        io.Writer
	indent     int
	needIndent bool
}

func NewWriterFile(fname string) *Writer {
	f, err := os.Create(fname)
	checkUserErr(err)
	//out := bufio.NewWriter(f)
	return &Writer{out: f}
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{out: out}
}

func (w *Writer) Close() {
	if closer, ok := w.out.(io.Closer); ok {
		err := closer.Close()
		checkUserErr(err)
	}
}

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
