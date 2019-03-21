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

func freeParser(p *Parser) {
	p.u.Delete()
}

func NewParser() *Parser {
	uparser := unmanaged.NewParser()
	parser := &Parser{&uparser}
	runtime.SetFinalizer(parser, freeParser)
	return parser
}

func (self *Parser) SetLanguage(language *Language) error {
	return self.u.SetLanguage(language)
}

func (self *Parser) Reset() {
	self.u.Reset()
}

func (self *Parser) SetOperationLimit(limit uint64) {
	self.u.SetOperationLimit(limit)
}

func (self *Parser) Parse(input []byte, oldTree *Tree) (*Tree, bool) {
	var uoldTree unmanaged.Tree
	if oldTree != nil {
		uoldTree = oldTree.u
	}

	utree, ok := self.u.Parse(input, uoldTree)
	runtime.KeepAlive(uoldTree)

	if ok {
		return toManagedTree(utree, self), true
	} else {
		return &Tree{}, false
	}
}

//------------------------------------------------------------------------------

type Tree struct {
	u unmanaged.Tree
	p interface{}
}

func toManagedTree(utree unmanaged.Tree, parent interface{}) *Tree {
	tree := &Tree{utree, parent}
	runtime.SetFinalizer(tree, freeTree)
	return tree
}

func freeTree(t *Tree) {
	t.u.Delete()
}

func (self *Tree) RootNode() Node {
	return Node{self.u.RootNode(), self}
}

func (self *Tree) Edit(edit InputEdit) {
	self.u.Edit(edit)
}

func (self *Tree) Walk() TreeCursor {
	return toManagedTreeCursor(self.u.Walk(), self)
}

func (self *Tree) Copy() *Tree {
	return toManagedTree(self.u.Copy(), self.p)
}

//------------------------------------------------------------------------------

type Node struct {
	u unmanaged.Node
	p interface{}
}

func toManagedNode(unode unmanaged.Node, parent interface{}) Node {
	return Node{unode, parent}
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
	n, ok := self.u.Child(i)
	return toManagedNode(n, self), ok
}

func (self Node) ChildCount() uint32 {
	return self.u.ChildCount()
}

func (self Node) NamedChild(i uint32) (Node, bool) {
	n, ok := self.u.NamedChild(i)
	return toManagedNode(n, self), ok
}

func (self Node) NamedChildCount() uint32 {
	return self.u.NamedChildCount()
}

func (self Node) Parent() (Node, bool) {
	n, ok := (self.u.Parent())
	return toManagedNode(n, self), ok
}

func (self Node) NextSibling() (Node, bool) {
	n, ok := (self.u.NextSibling())
	return toManagedNode(n, self), ok
}

func (self Node) PrevSibling() (Node, bool) {
	n, ok := (self.u.PrevSibling())
	return toManagedNode(n, self), ok
}

func (self Node) NextNamedSibling() (Node, bool) {
	n, ok := (self.u.NextNamedSibling())
	return toManagedNode(n, self), ok
}

func (self Node) PrevNamedSibling() (Node, bool) {
	n, ok := (self.u.PrevNamedSibling())
	return toManagedNode(n, self), ok
}

func (self Node) ToSexp() string {
	return self.u.ToSexp()
}

func (self Node) Equals(node Node) bool {
	return self.u.Equals(node.u)
}

func (self Node) Walk() TreeCursor {
	return toManagedTreeCursor(self.u.Walk(), self)
}

//------------------------------------------------------------------------------

type TreeCursor struct {
	u *unmanaged.TreeCursor
	p interface{}
}

func freeTreeCursor(t *unmanaged.TreeCursor) {
	t.Delete()
}

func toManagedTreeCursor(ucursor unmanaged.TreeCursor, parent interface{}) TreeCursor {
	cursor := TreeCursor{&ucursor, parent}
	runtime.SetFinalizer(cursor.u, freeTreeCursor)
	return cursor
}

func (self TreeCursor) Node() Node {
	return toManagedNode(self.u.Node(), self)
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
