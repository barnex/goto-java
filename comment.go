package main

import "go/ast"

func (w *writer) putDoc(g *ast.CommentGroup) {
	w.putComment(g) //TODO: translate to slashstarstar
}

func (w *writer) putComment(g *ast.CommentGroup) {
	if g == nil {
		return
	}
	for _, c := range g.List {
		w.putln(c.Text)
	}
}
