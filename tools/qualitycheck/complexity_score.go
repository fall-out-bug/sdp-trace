package main

import "go/ast"

// loopStatementScore applies the common loop rule used by both for and range
// statements in the cognitive-complexity metric.
func loopStatementScore(body []ast.Stmt, nesting int) int {
	// Loops add one structural point plus nesting; child statements are scored
	// at the next nesting level.
	return 1 + nesting + cognitiveChildren(body, nesting+1)
}

// switchStatementScore applies the common switch/select rule while keeping case
// bodies responsible for their nested decisions.
func switchStatementScore(clauses []ast.Stmt, nesting int) int {
	// Switch/select statements own the outer decision; cases are scored inside
	// the body so empty and populated clauses share one rule.
	return 1 + nesting + switchBodyScore(clauses, nesting+1)
}
