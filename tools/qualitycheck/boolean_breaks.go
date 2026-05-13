package main

import (
	"go/ast"
	"go/token"
)

// booleanBreaks scans one expression tree for && and || cognitive breaks.
func booleanBreaks(expression ast.Expr) int {
	// Cognitive complexity treats chained && and || operators as decision
	// breaks, so scan only binary expressions and keep operands transparent.
	score := 0
	ast.Inspect(expression, func(node ast.Node) bool {
		binary, ok := node.(*ast.BinaryExpr)
		if ok && (binary.Op == token.LAND || binary.Op == token.LOR) {
			// Each logical operator adds one short-circuit decision point.
			score++
		}
		// Continue through operands so nested boolean expressions are counted.
		return true
	})
	return score
}
