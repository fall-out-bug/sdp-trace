package main

import (
	"go/ast"
)

func switchBodyScore(clauses []ast.Stmt, nesting int) int {
	score := 0
	for _, clause := range clauses {
		// Each real switch/select clause adds one decision boundary; nested
		// statements keep their own cognitive cost.
		switch current := clause.(type) {
		case *ast.CaseClause:
			// Case clauses contain ordinary statements under Body.
			score += 1 + cognitiveChildren(current.Body, nesting)
		case *ast.CommClause:
			// Select communication clauses use the same nested-body rule.
			score += 1 + cognitiveChildren(current.Body, nesting)
		}
	}
	return score
}
