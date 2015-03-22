package gotojava

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/types"
)

// JFileSet roots the java AST.
type JFileSet struct {
	classes   map[string]*JClass // All java classes, indexed by name
	MainClass *JClass            // The main java class, where Go top-level declarations go
}

// NewJFileSet constructs an emtpy java file set with given main class.
// Go top-level declarations go into the main class as static members,
// typedefs and methods go into their own class.
func NewJFileSet(mainClass string) *JFileSet {
	j := &JFileSet{classes: make(map[string]*JClass)}
	j.MainClass = j.Class(mainClass)
	return j
}

// Transpile a Go file and add to the java AST. Godoc:
// 	type File struct {
// 	        Doc        *CommentGroup   // associated documentation; or nil
// 	        Package    token.Pos       // position of "package" keyword
// 	        Name       *Ident          // package name
// 	        Decls      []Decl          // top-level declarations; or nil
// 	        Scope      *Scope          // package scope (this file only)
// 	        Imports    []*ImportSpec   // imports in this file
// 	        Unresolved []*Ident        // unresolved identifiers in this file
// 	        Comments   []*CommentGroup // list of all comments in the source file
// 	}
func (j *JFileSet) Add(f *ast.File) {
	// add all top-level declarations to the corresponding classes
	for _, d := range f.Decls {
		switch d := d.(type) {
		default:
			panic(reflect.TypeOf(d).String())
		case *ast.FuncDecl:
			if d.Recv == nil {
				j.ClassForPkg(f.Name.Name).AddStaticFunc(d)
			} else {
				assert(len(d.Recv.List) == 1)
				j.ClassForType(TypeOf(d.Recv.List[0].Type)).AddMethod(d)
			}
		case *ast.GenDecl:
			j.ClassForPkg(f.Name.Name).AddGenDecl(d)
		}
	}
}

func (j *JFileSet) ClassForPkg(pkgName string) *JClass {
	className := pkgName // TODO
	return j.Class(className)
}

func (j *JFileSet) ClassForType(t types.Type) *JClass {
	return j.Class(javaName(t))
}

// Class returns the JClass with given class name.
// A new JClass is first created and added to the JFileSet if needed.
func (j *JFileSet) Class(name string) *JClass {
	if c, ok := j.classes[name]; ok {
		return c
	} else {
		j.classes[name] = NewJClass()
		return j.classes[name]
	}
}

func (j *JFileSet) Write() {
	for n, c := range j.classes {
		w := NewWriterFile(n + ".java")
		c.Write(w)
		w.Close()
	}
}
