package unmanaged

/*
#cgo CFLAGS: -I${SRCDIR}/../../third-party/tree-sitter/src
#cgo CFLAGS: -I${SRCDIR}/../../third-party/tree-sitter/include
#cgo CFLAGS: -I${SRCDIR}/../../third-party/tree-sitter/externals/utf8proc

#include "tree_sitter/compiler.h"
#include "tree_sitter/parser.h"
#include "tree_sitter/runtime.h"
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

const TREE_SITTER_LANGUAGE_VERSION = 9

type Language uintptr

type Point struct {
	Row    uint32
	Column uint32
}

type InputEdit struct {
	StartByte      uint32
	OldEndByte     uint32
	NewEndByte     uint32
	StartPosition  Point
	OldEndPosition Point
	NewEndPosition Point
}

//------------------------------------------------------------------------------

type Parser struct{ ptr *C.TSParser }

func NewParser() Parser {
	var parser *C.TSParser = C.ts_parser_new()
	return Parser{parser}
}

func (self Parser) SetLanguage(language Language) error {
	c_language := (*C.TSLanguage)(unsafe.Pointer(language))
	version := C.ts_language_version(c_language)
	if version == TREE_SITTER_LANGUAGE_VERSION {
		C.ts_parser_set_language(self.ptr, c_language)
		return nil
	} else {
		return fmt.Errorf("Incompatible language version %d. Expected %d.", version, TREE_SITTER_LANGUAGE_VERSION)
	}
}

func (self Parser) Reset() {
	C.ts_parser_reset(self.ptr)
}

func (self Parser) SetOperationLimit(limit uint64) {
	C.ts_parser_set_operation_limit(self.ptr, C.size_t(limit))
}

func (self Parser) Parse(input []byte, oldTree *Tree) (Tree, bool) {
	c_old_tree := (*C.TSTree)(nil)
	if oldTree != nil {
		c_old_tree = oldTree.ptr
	}

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&input))
	c_tree := C.ts_parser_parse_string(
		self.ptr,
		c_old_tree,
		(*C.char)(unsafe.Pointer(hdr.Data)),
		C.uint(hdr.Len),
	)
	runtime.KeepAlive(&input)

	if c_tree != nil {
		return Tree{c_tree}, true
	} else {
		return Tree{}, false
	}
}

func (self Parser) Delete() {
	C.ts_parser_delete(self.ptr)
}

//------------------------------------------------------------------------------

type Tree struct{ ptr *C.TSTree }

func (self Tree) RootNode() Node {
	c_node := C.ts_tree_root_node(self.ptr)
	return Node{c_node}
}

func (self Tree) Edit(edit InputEdit) {
	c_edit := C.TSInputEdit{
		start_byte:   C.uint(edit.StartByte),
		old_end_byte: C.uint(edit.OldEndByte),
		new_end_byte: C.uint(edit.NewEndByte),
		start_point: C.TSPoint{
			column: C.uint(edit.StartPosition.Column),
			row:    C.uint(edit.StartPosition.Row),
		},
		old_end_point: C.TSPoint{
			column: C.uint(edit.OldEndPosition.Column),
			row:    C.uint(edit.OldEndPosition.Row),
		},
		new_end_point: C.TSPoint{
			column: C.uint(edit.NewEndPosition.Column),
			row:    C.uint(edit.NewEndPosition.Row),
		},
	}

	C.ts_tree_edit(self.ptr, &c_edit)
}

func (self Tree) Walk() TreeCursor {
	return self.RootNode().Walk()
}

func (self Tree) Delete() {
	C.ts_tree_delete(self.ptr)
}

func (self Tree) Copy() Tree {
	return Tree{C.ts_tree_copy(self.ptr)}
}

//------------------------------------------------------------------------------

type Node struct{ val C.TSNode }

func newNode(node C.TSNode) (Node, bool) {
	if node.id == nil {
		return Node{}, false
	} else {
		return Node{node}, true
	}
}

func (self Node) KindId() uint16 {
	return uint16(C.ts_node_symbol(self.val))
}

func (self Node) Kind() string {
	c_str := C.ts_node_type(self.val)
	return C.GoString(c_str)
}

func (self Node) IsNamed() bool {
	return bool(C.ts_node_is_named(self.val))
}

func (self Node) HasChanges() bool {
	return bool(C.ts_node_has_changes(self.val))
}

func (self Node) HasError() bool {
	return bool(C.ts_node_has_error(self.val))
}

func (self Node) StartByte() uint32 {
	return uint32(C.ts_node_start_byte(self.val))
}

func (self Node) EndByte() uint32 {
	return uint32(C.ts_node_end_byte(self.val))
}

func (self Node) StartPosition() Point {
	result := C.ts_node_start_point(self.val)
	return Point{
		Row:    uint32(result.row),
		Column: uint32(result.column),
	}
}

func (self Node) EndPosition() Point {
	result := C.ts_node_end_point(self.val)
	return Point{
		Row:    uint32(result.row),
		Column: uint32(result.column),
	}
}

func (self Node) Child(i uint32) (Node, bool) {
	return newNode(C.ts_node_child(self.val, C.uint32_t(i)))
}

func (self Node) ChildCount() uint32 {
	return uint32(C.ts_node_child_count(self.val))
}

func (self Node) NamedChild(i uint32) (Node, bool) {
	return newNode(C.ts_node_named_child(self.val, C.uint32_t(i)))
}

func (self Node) NamedChildCount() uint32 {
	return uint32(C.ts_node_named_child_count(self.val))
}

func (self Node) Parent() (Node, bool) {
	return newNode(C.ts_node_parent(self.val))
}

func (self Node) NextSibling() (Node, bool) {
	return newNode(C.ts_node_next_sibling(self.val))
}

func (self Node) PrevSibling() (Node, bool) {
	return newNode(C.ts_node_prev_sibling(self.val))
}

func (self Node) NextNamedSibling() (Node, bool) {
	return newNode(C.ts_node_next_named_sibling(self.val))
}

func (self Node) PrevNamedSibling() (Node, bool) {
	return newNode(C.ts_node_prev_named_sibling(self.val))
}

func (self Node) ToSexp() string {
	c_str := C.ts_node_string(self.val)
	defer C.free(unsafe.Pointer(c_str))
	return C.GoString(c_str)
}

func (self Node) Equals(node Node) bool {
	return self.val.id == node.val.id
}

func (self Node) Walk() TreeCursor {
	return TreeCursor{C.ts_tree_cursor_new(self.val)}
}

//------------------------------------------------------------------------------

type TreeCursor struct{ val C.TSTreeCursor }

func (self TreeCursor) Node() Node {
	return Node{C.ts_tree_cursor_current_node(&self.val)}
}

func (self TreeCursor) GoToFirstChild() bool {
	return bool(C.ts_tree_cursor_goto_first_child(&self.val))
}

func (self TreeCursor) GoToParent() bool {
	return bool(C.ts_tree_cursor_goto_parent(&self.val))
}

func (self TreeCursor) GoToNextSibling() bool {
	return bool(C.ts_tree_cursor_goto_next_sibling(&self.val))
}

func (self TreeCursor) GoToFirstChildForIndex(index uint32) int64 {
	return int64(C.ts_tree_cursor_goto_first_child_for_byte(&self.val, C.uint32_t(index)))
}

func (self TreeCursor) Delete() {
	C.ts_tree_cursor_delete(&self.val)
}
