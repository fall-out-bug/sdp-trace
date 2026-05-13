package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// analyzeFile is intentionally the handoff from filesystem bytes to verifier
// metrics. It parses comments because comment ownership is part of both file
// and function MI, while all scoring still uses the same immutable source read.
// Keeping parse, file metrics, and function metrics in this small boundary
// avoids mixing discovery policy with trust measurements.
func analyzeFile(path string) (fileMetric, []functionMetric, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return fileMetric{}, nil, err
	}
	fset := token.NewFileSet()
	// Comments must be parsed because MI uses comment-line counts at both file
	// and function granularity.
	parsed, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return fileMetric{}, nil, err
	}
	return measureFile(path, src, parsed), measureFunctions(path, src, fset, parsed), nil
}

func cognitiveSwitchStatement(statement ast.Stmt, nesting int) (int, bool) {
	// Expression and type switches share the same clause scoring; only AST
	// dispatch differs.
	if current, ok := statement.(*ast.SwitchStmt); ok {
		// Expression switches carry their cases under Body.List, matching the
		// common switch scoring helper.
		return switchStatementScore(current.Body.List, nesting), true
	}
	if current, ok := statement.(*ast.TypeSwitchStmt); ok {
		// Type switches use the same body shape and therefore the same scoring
		// boundary as expression switches.
		return switchStatementScore(current.Body.List, nesting), true
	}
	return 0, false
}

func halsteadTokenKey(tok token.Token, literal string) string {
	if literal != "" {
		return literal
	}
	// Operators and keywords do not always carry a literal from the scanner;
	// token.String keeps their Halstead vocabulary stable across files.
	return tok.String()
}
