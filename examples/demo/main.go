package main

import (
	"fmt"

	"github.com/cstrahan/go-tree-sitter/cgo"
)

func main() {
	src := []byte("package foo\nfunc main() { }")

	parser := cgo.NewParser()
	parser.SetLanguage(langGo)
	tree, _ := parser.Parse(src, nil)
	node := tree.RootNode()
	sexp := node.ToSexp()
	fmt.Println(sexp)
}
