package gotojava

// Walk the ast and map nodes to their parent

import (
	"fmt"
	. "go/ast"
)

// Walk the ast rooted at n and populate the global parents map,
// mapping each node to its parent.
func CollectParents(n Node) map[Node]Node {
	parents := make(map[Node]Node)
	walk(&parentCollector{parents: parents}, n)
	return parents
}

// Visitor that collects parents
type parentCollector struct {
	parents map[Node]Node
	stack   []Node // stack of parent nodes
}

func (v *parentCollector) push(n Node) { v.stack = append(v.stack, n) }
func (v *parentCollector) pop()        { v.stack = v.stack[:len(v.stack)-1] }

// this ast visitor enters n into its parents map,
// mapping it to its parent obtained from the visitor's stack
func (v *parentCollector) Visit(n Node) *parentCollector {
	if n == nil {
		return v
	}
	var parent Node = nil
	if len(v.stack) > 0 {
		parent = v.stack[len(v.stack)-1]
	}
	v.parents[n] = parent
	return v
}

// Copied from go/ast/walk.go, copyright The Go Authors.
// Modified to call v.push(node)/v.pop() so that v keeps a parent stack.
func walk(v *parentCollector, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	v.push(node)
	defer v.pop()

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		// nothing to do

	case *CommentGroup:
		for _, c := range n.List {
			walk(v, c)
		}

	case *Field:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		walk(v, n.Type)
		if n.Tag != nil {
			walk(v, n.Tag)
		}
		if n.Comment != nil {
			walk(v, n.Comment)
		}

	case *FieldList:
		for _, f := range n.List {
			walk(v, f)
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	case *Ellipsis:
		if n.Elt != nil {
			walk(v, n.Elt)
		}

	case *FuncLit:
		walk(v, n.Type)
		walk(v, n.Body)

	case *CompositeLit:
		if n.Type != nil {
			walk(v, n.Type)
		}
		walkExprList(v, n.Elts)

	case *ParenExpr:
		walk(v, n.X)

	case *SelectorExpr:
		walk(v, n.X)
		walk(v, n.Sel)

	case *IndexExpr:
		walk(v, n.X)
		walk(v, n.Index)

	case *SliceExpr:
		walk(v, n.X)
		if n.Low != nil {
			walk(v, n.Low)
		}
		if n.High != nil {
			walk(v, n.High)
		}
		if n.Max != nil {
			walk(v, n.Max)
		}

	case *TypeAssertExpr:
		walk(v, n.X)
		if n.Type != nil {
			walk(v, n.Type)
		}

	case *CallExpr:
		walk(v, n.Fun)
		walkExprList(v, n.Args)

	case *StarExpr:
		walk(v, n.X)

	case *UnaryExpr:
		walk(v, n.X)

	case *BinaryExpr:
		walk(v, n.X)
		walk(v, n.Y)

	case *KeyValueExpr:
		walk(v, n.Key)
		walk(v, n.Value)

	// Types
	case *ArrayType:
		if n.Len != nil {
			walk(v, n.Len)
		}
		walk(v, n.Elt)

	case *StructType:
		walk(v, n.Fields)

	case *FuncType:
		if n.Params != nil {
			walk(v, n.Params)
		}
		if n.Results != nil {
			walk(v, n.Results)
		}

	case *InterfaceType:
		walk(v, n.Methods)

	case *MapType:
		walk(v, n.Key)
		walk(v, n.Value)

	case *ChanType:
		walk(v, n.Value)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		walk(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *LabeledStmt:
		walk(v, n.Label)
		walk(v, n.Stmt)

	case *ExprStmt:
		walk(v, n.X)

	case *SendStmt:
		walk(v, n.Chan)
		walk(v, n.Value)

	case *IncDecStmt:
		walk(v, n.X)

	case *AssignStmt:
		walkExprList(v, n.Lhs)
		walkExprList(v, n.Rhs)

	case *GoStmt:
		walk(v, n.Call)

	case *DeferStmt:
		walk(v, n.Call)

	case *ReturnStmt:
		walkExprList(v, n.Results)

	case *BranchStmt:
		if n.Label != nil {
			walk(v, n.Label)
		}

	case *BlockStmt:
		walkStmtList(v, n.List)

	case *IfStmt:
		if n.Init != nil {
			walk(v, n.Init)
		}
		walk(v, n.Cond)
		walk(v, n.Body)
		if n.Else != nil {
			walk(v, n.Else)
		}

	case *CaseClause:
		walkExprList(v, n.List)
		walkStmtList(v, n.Body)

	case *SwitchStmt:
		if n.Init != nil {
			walk(v, n.Init)
		}
		if n.Tag != nil {
			walk(v, n.Tag)
		}
		walk(v, n.Body)

	case *TypeSwitchStmt:
		if n.Init != nil {
			walk(v, n.Init)
		}
		walk(v, n.Assign)
		walk(v, n.Body)

	case *CommClause:
		if n.Comm != nil {
			walk(v, n.Comm)
		}
		walkStmtList(v, n.Body)

	case *SelectStmt:
		walk(v, n.Body)

	case *ForStmt:
		if n.Init != nil {
			walk(v, n.Init)
		}
		if n.Cond != nil {
			walk(v, n.Cond)
		}
		if n.Post != nil {
			walk(v, n.Post)
		}
		walk(v, n.Body)

	case *RangeStmt:
		if n.Key != nil {
			walk(v, n.Key)
		}
		if n.Value != nil {
			walk(v, n.Value)
		}
		walk(v, n.X)
		walk(v, n.Body)

	// Declarations
	case *ImportSpec:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		if n.Name != nil {
			walk(v, n.Name)
		}
		walk(v, n.Path)
		if n.Comment != nil {
			walk(v, n.Comment)
		}

	case *ValueSpec:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		if n.Type != nil {
			walk(v, n.Type)
		}
		walkExprList(v, n.Values)
		if n.Comment != nil {
			walk(v, n.Comment)
		}

	case *TypeSpec:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		walk(v, n.Name)
		walk(v, n.Type)
		if n.Comment != nil {
			walk(v, n.Comment)
		}

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		for _, s := range n.Specs {
			walk(v, s)
		}

	case *FuncDecl:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		if n.Recv != nil {
			walk(v, n.Recv)
		}
		walk(v, n.Name)
		walk(v, n.Type)
		if n.Body != nil {
			walk(v, n.Body)
		}

	// Files and packages
	case *File:
		if n.Doc != nil {
			walk(v, n.Doc)
		}
		walk(v, n.Name)
		walkDeclList(v, n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for _, f := range n.Files {
			walk(v, f)
		}

	default:
		fmt.Printf("ast.walk: unexpected node type %T", n)
		panic("ast.walk")
	}

	v.Visit(nil)
}

func walkIdentList(v *parentCollector, list []*Ident) {
	for _, x := range list {
		walk(v, x)
	}
}

func walkExprList(v *parentCollector, list []Expr) {
	for _, x := range list {
		walk(v, x)
	}
}

func walkStmtList(v *parentCollector, list []Stmt) {
	for _, x := range list {
		walk(v, x)
	}
}

func walkDeclList(v *parentCollector, list []Decl) {
	for _, x := range list {
		walk(v, x)
	}
}
