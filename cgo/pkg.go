package cgo

import (
	"runtime"

	"github.com/cstrahan/go-tree-sitter/cgo/unmanaged"
)

const TREE_SITTER_LANGUAGE_VERSION = unmanaged.TREE_SITTER_LANGUAGE_VERSION

type Language = unmanaged.Language

type Point = unmanaged.Point

type InputEdit = unmanaged.InputEdit

//------------------------------------------------------------------------------

type Parser struct{ u *unmanaged.Parser }

func freeParser(p *unmanaged.Parser) {
	p.Delete()
}

func NewParser() Parser {
	uparser := unmanaged.NewParser()
	parser := Parser{&uparser}
	runtime.SetFinalizer(parser.u, freeParser)
	return parser
}

func (self Parser) SetLanguage(language Language) error {
	return self.u.SetLanguage(language)
}

func (self Parser) Reset() {
	self.u.Reset()
}

func (self Parser) SetOperationLimit(limit uint64) {
	self.u.SetOperationLimit(limit)
}

func (self Parser) Parse(input []byte, oldTree *Tree) (Tree, bool) {
	var uoldTree *unmanaged.Tree
	if oldTree != nil {
		uoldTree = oldTree.u
	}

	utree, ok := self.u.Parse(input, uoldTree)
	runtime.KeepAlive(uoldTree)

	if ok {
		return toManagedTree(utree), true
	} else {
		return Tree{}, false
	}
}

//------------------------------------------------------------------------------

type Tree struct{ u *unmanaged.Tree }

func toManagedTree(utree unmanaged.Tree) Tree {
	tree := Tree{&utree}
	runtime.SetFinalizer(tree.u, freeTree)
	return tree
}

func freeTree(t *unmanaged.Tree) {
	t.Delete()
}

func (self Tree) RootNode() Node {
	return Node{self.u.RootNode()}
}

func (self Tree) Edit(edit InputEdit) {
	self.u.Edit(edit)
}

func (self Tree) Walk() TreeCursor {
	return toManagedTreeCursor(self.u.Walk())
}

func (self Tree) Copy() Tree {
	return toManagedTree(self.u.Copy())
}

//------------------------------------------------------------------------------

type Node struct{ u unmanaged.Node }

func toManagedNode(unode unmanaged.Node, ok bool) (Node, bool) {
	return Node{unode}, ok
}

func (self Node) KindId() uint16 {
	return self.u.KindId()
}

func (self Node) Kind() string {
	return self.u.Kind()
}

func (self Node) IsNamed() bool {
	return self.u.IsNamed()
}

func (self Node) HasChanges() bool {
	return self.u.HasChanges()
}

func (self Node) HasError() bool {
	return self.u.HasError()
}

func (self Node) StartByte() uint32 {
	return self.u.StartByte()
}

func (self Node) EndByte() uint32 {
	return self.u.EndByte()
}

func (self Node) StartPosition() Point {
	return self.u.StartPosition()
}

func (self Node) EndPosition() Point {
	return self.u.EndPosition()
}

func (self Node) Child(i uint32) (Node, bool) {
	return toManagedNode(self.u.Child(i))
}

func (self Node) ChildCount() uint32 {
	return self.u.ChildCount()
}

func (self Node) NamedChild(i uint32) (Node, bool) {
	return toManagedNode(self.u.NamedChild(i))
}

func (self Node) NamedChildCount() uint32 {
	return self.u.NamedChildCount()
}

func (self Node) Parent() (Node, bool) {
	return toManagedNode(self.u.Parent())
}

func (self Node) NextSibling() (Node, bool) {
	return toManagedNode(self.u.NextSibling())
}

func (self Node) PrevSibling() (Node, bool) {
	return toManagedNode(self.u.PrevSibling())
}

func (self Node) NextNamedSibling() (Node, bool) {
	return toManagedNode(self.u.NextNamedSibling())
}

func (self Node) PrevNamedSibling() (Node, bool) {
	return toManagedNode(self.u.PrevNamedSibling())
}

func (self Node) ToSexp() string {
	return self.u.ToSexp()
}

func (self Node) Equals(node Node) bool {
	return self.u.Equals(node.u)
}

func (self Node) Walk() TreeCursor {
	return toManagedTreeCursor(self.u.Walk())
}

//------------------------------------------------------------------------------

type TreeCursor struct{ u *unmanaged.TreeCursor }

func freeTreeCursor(t *unmanaged.TreeCursor) {
	t.Delete()
}

func toManagedTreeCursor(ucursor unmanaged.TreeCursor) TreeCursor {
	cursor := TreeCursor{&ucursor}
	runtime.SetFinalizer(cursor.u, freeTreeCursor)
	return cursor
}

func (self TreeCursor) Node() Node {
	return Node{self.u.Node()}
}

func (self TreeCursor) GoToFirstChild() bool {
	return self.u.GoToFirstChild()
}

func (self TreeCursor) GoToParent() bool {
	return self.u.GoToParent()
}

func (self TreeCursor) GoToNextSibling() bool {
	return self.u.GoToNextSibling()
}

func (self TreeCursor) GoToFirstChildForIndex(index uint32) int64 {
	return self.u.GoToFirstChildForIndex(index)
}
