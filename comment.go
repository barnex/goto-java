package main

import "go/ast"

func (w *writer) PutDoc(g *ast.CommentGroup) {
	w.PutComment(g) //TODO: translate to slashstarstar
	w.Putln()
}

func (w *writer) PutInlineComment(g *ast.CommentGroup) {
	if g == nil {
		return
	}
	w.Put("\t")
	for i, c := range g.List {
		w.Put(c.Text)
		if i != len(g.List)-1 {
			w.Putln()
		}
	}
}

func (w *writer) PutComment(g *ast.CommentGroup) {
	if g == nil {
		return
	}
	for i, c := range g.List {
		w.Put(c.Text)
		if i != len(g.List)-1 {
			w.Putln()
		}
	}
}
