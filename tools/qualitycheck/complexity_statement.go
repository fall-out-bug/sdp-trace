package main

import "go/ast"

func cognitiveStatement(statement ast.Stmt, nesting int) int {
	// Structural forms are checked before leaf expressions so nesting-sensitive
	// score is assigned by the branch that owns it.
	if score, ok := cognitiveIfOrLoop(statement, nesting); ok {
		return score
	}
	if score, ok := cognitiveSwitchOrSelect(statement, nesting); ok {
		return score
	}
	return cognitiveSimpleStatement(statement, nesting)
}

func cognitiveIfOrLoop(statement ast.Stmt, nesting int) (int, bool) {
	// If statements own else-if handling; loops share the body scoring helper.
	if current, ok := statement.(*ast.IfStmt); ok {
		// The if helper handles condition breaks, body nesting, and else chains.
		return cognitiveIfScore(current, nesting), true
	}
	// Loop syntax is checked separately so if handling stays readable.
	return cognitiveLoopStatement(statement, nesting)
}

func cognitiveLoopStatement(statement ast.Stmt, nesting int) (int, bool) {
	// For and range syntax each creates one structural break before children.
	if current, ok := statement.(*ast.ForStmt); ok {
		// Plain for loops and range loops share the same body scoring formula.
		return loopStatementScore(current.Body.List, nesting), true
	}
	if current, ok := statement.(*ast.RangeStmt); ok {
		// Range loops differ syntactically only in the AST node wrapper.
		return loopStatementScore(current.Body.List, nesting), true
	}
	return 0, false
}
