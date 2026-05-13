package main

import (
	"go/ast"
	"go/token"
)

// schemaVersion returns the discriminator stored in function-MI baselines.
func (baseline functionMIBaseline) schemaVersion() string {
	return baseline.SchemaVersion
}

// schemaVersion returns the discriminator stored in file-MI baselines.
func (baseline fileMIBaseline) schemaVersion() string {
	return baseline.SchemaVersion
}

func cognitiveExpressionStatement(statement ast.Stmt) int {
	// Standalone expressions can only add boolean short-circuit breaks.
	if current, ok := statement.(*ast.ExprStmt); ok {
		// Expression statements do not add structural complexity on their own.
		return booleanBreaks(current.X)
	}
	// Assignments and returns are handled by the result-bearing statement path.
	return cognitiveResultExpressionStatement(statement)
}

func cognitiveResultExpressionStatement(statement ast.Stmt) int {
	// Assignments and returns are the result-bearing statements that can hide
	// boolean breaks outside explicit branch syntax.
	if current, ok := statement.(*ast.AssignStmt); ok {
		// Only right-hand expressions can contain new boolean breaks.
		return exprListBooleanBreaks(current.Rhs)
	}
	if current, ok := statement.(*ast.ReturnStmt); ok {
		// Return expressions are scored the same way as assignment results.
		return exprListBooleanBreaks(current.Results)
	}
	return 0
}

func skipHalsteadToken(tok token.Token) bool {
	return tok == token.COMMENT || tok == token.SEMICOLON
}
