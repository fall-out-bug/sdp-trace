package main

import "go/ast"

func cognitiveIfScore(statement *ast.IfStmt, nesting int) int {
	// Conditions own the branch decision and any boolean-expression breaks
	// inside it; body statements are scored after the branch nesting increases.
	score := 1 + nesting + booleanBreaks(statement.Cond)
	score += cognitiveChildren(statement.Body.List, nesting+1)
	switch elseBranch := statement.Else.(type) {
	case *ast.IfStmt:
		// Else-if continues the same decision chain, so it does not gain an
		// extra nesting level beyond its own condition.
		score += cognitiveIfScore(elseBranch, nesting)
	case *ast.BlockStmt:
		// A terminal else block is nested under the original branch decision.
		score += cognitiveChildren(elseBranch.List, nesting+1)
	}
	return score
}
