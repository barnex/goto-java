package gotojava

// Handling of comments
// TODO: collect all (non-doc) comments and put in appropriate place
// TODO: Emit all doc comments.

import "go/ast"

// Emit a doc comment. ast godoc:
// 	type CommentGroup struct {
// 	        List []*Comment // len(List) > 0
// 	}
// 	type Comment struct {
// 	        Slash token.Pos // position of "/" starting the comment
// 	        Text  string    // comment text (excluding '\n' for //-style comments)
// 	}
func (w *Writer) PutDoc(cg *ast.CommentGroup) {
	if cg != nil {
		w.Putln("/** ", cg.Text(), "*/") // TODO: wrap long lines
	}
}
