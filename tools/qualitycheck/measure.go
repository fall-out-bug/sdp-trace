package main

import (
	"go/ast"
	"go/token"
)

func measureFunctions(path string, src []byte, fset *token.FileSet, parsed *ast.File) []functionMetric {
	// Function measurement walks top-level declarations only; nested literals do
	// not become independent ratchet subjects.
	metrics := make([]functionMetric, 0, len(parsed.Decls))
	for _, decl := range parsed.Decls {
		// Only declarations with function bodies become function ratchet
		// subjects; file MI omits declaration-only and tiny files separately.
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		metrics = append(metrics, measureFunction(path, src, fset, parsed, fn))
	}
	return metrics
}
