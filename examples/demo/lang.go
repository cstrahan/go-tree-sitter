package main

/*
#cgo CFLAGS: -I${SRCDIR}

#include "tree_sitter/parser.h"

extern const TSLanguage *tree_sitter_go();
*/
import "C"

import (
	"unsafe"

	"github.com/cstrahan/go-tree-sitter/cgo"
)

var langGo = (*cgo.Language)(unsafe.Pointer(C.tree_sitter_go()))
