package gotojava

// TODO: collect all comments and put in appropriate place

import "go/ast"

func (w *Writer) PutDoc(g *ast.CommentGroup) {
	w.PutComment(g) //TODO: translate to slashstarstar
	w.Putln()
}

func (w *Writer) PutInlineComment(g *ast.CommentGroup) {
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

func (w *Writer) PutComment(g *ast.CommentGroup) {
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
