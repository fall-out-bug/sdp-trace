package main

import (
	"go/ast"
)

// cyclomaticComplexity measures Go branch count using the local gate contract.
func cyclomaticComplexity(node ast.Node) int {
	// Cyclomatic complexity starts at one path, then adds each branch point
	// discovered in the function body.
	score := 1
	ast.Inspect(node, func(current ast.Node) bool {
		score += cyclomaticIncrement(current)
		return true
	})
	return score
}
