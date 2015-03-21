package gotojava

import "go/ast"

type JFileSet struct {
	Decls []ast.Decl
}

func NewJFileSet() *JFileSet {
	return new(JFileSet)
}

//type File struct {
//        Doc        *CommentGroup   // associated documentation; or nil
//        Package    token.Pos       // position of "package" keyword
//        Name       *Ident          // package name
//        Decls      []Decl          // top-level declarations; or nil
//        Scope      *Scope          // package scope (this file only)
//        Imports    []*ImportSpec   // imports in this file
//        Unresolved []*Ident        // unresolved identifiers in this file
//        Comments   []*CommentGroup // list of all comments in the source file
//}
func (j *JFileSet) Add(f *ast.File) {
	j.Decls = append(j.Decls, f.Decls...)
}

func (j *JFileSet) Compile() {

}

func (j *JFileSet) Write(w *Writer) {
	for _, d := range j.Decls {
		w.PutDecl(NONE, d)
	}
}
