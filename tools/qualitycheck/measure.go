package main

import (
	"go/ast"
	"go/token"
)

func measureFunctions(path string, src []byte, fset *token.FileSet, parsed *ast.File) []functionMetric {
	// Function measurement walks top-level declarations only; nested literals do
	// not become independent ratchet subjects.
	var metrics []functionMetric
	for _, decl := range parsed.Decls {
		// Only declarations with function bodies become function ratchet
		// subjects; package declarations and other top-level nodes are file MI.
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			// Non-function declarations still contribute to file-level MI.
			continue
		}
		metrics = append(metrics, measureFunction(path, src, fset, parsed, fn))
	}
	return metrics
}
