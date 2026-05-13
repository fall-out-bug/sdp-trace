package main

import (
	"go/ast"
)

// exprListBooleanBreaks sums boolean decision breaks across sibling expressions.
func exprListBooleanBreaks(expressions []ast.Expr) int {
	// Each expression is scored independently so caller nesting rules decide how
	// boolean breaks combine with surrounding control flow.
	score := 0
	for _, expression := range expressions {
		score += booleanBreaks(expression)
	}
	return score
}
