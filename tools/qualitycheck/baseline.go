package main

import (
	"go/ast"
	"go/token"
)

const (
	// Function and file baselines intentionally use distinct schema ids because
	// their row keys have different meanings.
	functionMIBaselineSchema = "sdp-trace-function-mi-baseline/v1"
	fileMIBaselineSchema     = "sdp-trace-file-mi-baseline/v1"
)

var branchStatementScores = map[token.Token]int{
	token.GOTO:     1,
	token.BREAK:    1,
	token.CONTINUE: 1,
}

// functionMIBaseline is the persisted ratchet for existing function-MI debt.
type functionMIBaseline struct {
	SchemaVersion string                     `json:"schema_version"`
	Threshold     float64                    `json:"threshold"`
	Functions     []functionMIBaselineRecord `json:"functions"`
}

func cognitiveSimpleStatement(statement ast.Stmt, nesting int) int {
	// Leaf statements either recurse through transparent blocks or count
	// boolean breaks in result-bearing expressions.
	if score, ok := cognitiveBranchOrBlock(statement, nesting); ok {
		// Branch and block statements already produced their complete score.
		return score
	}
	// Remaining simple statements can only affect cognitive score through
	// boolean expression breaks.
	return cognitiveExpressionStatement(statement)
}

func cognitiveBranchOrBlock(statement ast.Stmt, nesting int) (int, bool) {
	// Branch statements are leaf jumps; blocks keep the current nesting level.
	if current, ok := statement.(*ast.BranchStmt); ok {
		// Explicit jumps count only when branchScore recognizes the token.
		return branchScore(current), true
	}
	if current, ok := statement.(*ast.BlockStmt); ok {
		// Blocks are transparent containers; children retain the caller nesting.
		return cognitiveChildren(current.List, nesting), true
	}
	return 0, false
}

func branchScore(statement *ast.BranchStmt) int {
	// Fallthrough is intentionally excluded: it belongs to switch clause flow,
	// not to an explicit jump out of nested control.
	return branchStatementScores[statement.Tok]
}
