package main

import "go/ast"

func (w *writer) putDoc(g *ast.CommentGroup) {
	w.putComment(g) //TODO: translate to slashstarstar
	w.putln()
}

func (w *writer) putInlineComment(g *ast.CommentGroup) {
	if g == nil {
		return
	}
	w.put("\t")
	for i, c := range g.List {
		w.put(c.Text)
		if i != len(g.List)-1 {
			w.putln()
		}
	}
}

func (w *writer) putComment(g *ast.CommentGroup) {
	if g == nil {
		return
	}
	for i, c := range g.List {
		w.put(c.Text)
		if i != len(g.List)-1 {
			w.putln()
		}
	}
}
