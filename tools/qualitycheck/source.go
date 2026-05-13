package main

import (
	"go/ast"
	"go/token"
)

// sourceForNode returns the exact original source slice used by Halstead and
// line-count measurement for a parsed AST node.
func sourceForNode(src []byte, fset *token.FileSet, node ast.Node) string {
	file := fset.File(node.Pos())
	if file == nil {
		// A missing file-set entry means the parser cannot bind node positions
		// back to source bytes, so the caller receives an empty measurement.
		return ""
	}
	// AST token positions are converted through the file set before slicing so
	// function MI measures the original bytes for exactly that declaration.
	start := file.Offset(node.Pos())
	end := file.Offset(node.End())
	if !validSourceRange(start, end, len(src)) {
		return ""
	}
	return string(src[start:end])
}
