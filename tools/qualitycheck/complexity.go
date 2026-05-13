package main

import (
	"go/ast"
)

// cognitiveComplexity scores only executable function bodies.
func cognitiveComplexity(fn *ast.FuncDecl) int {
	if fn.Body == nil {
		// Declarations without bodies, such as extern-like stubs, have no local
		// control-flow burden to score.
		return 0
	}
	return cognitiveChildren(fn.Body.List, 0)
}

// cognitiveChildren accumulates sibling statement cost at a shared nesting depth.
func cognitiveChildren(statements []ast.Stmt, nesting int) int {
	// Sibling statements share the same nesting level; child helpers add the
	// extra nesting cost only when control flow introduces it.
	score := 0
	for _, statement := range statements {
		score += cognitiveStatement(statement, nesting)
	}
	return score
}
