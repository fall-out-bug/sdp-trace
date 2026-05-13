package main

import (
	"go/ast"
	"io/fs"
	"path/filepath"
)

func collectGoFile(walkPath string, entry fs.DirEntry, walkErr error, files *[]string) error {
	if walkErr != nil {
		return walkErr
	}
	if entry.IsDir() {
		return walkDirSkip(entry.Name())
	}
	// The walker appends only active production Go files; path ordering comes
	// from filepath.WalkDir and is preserved for deterministic reports.
	if isProductionGo(walkPath) {
		*files = append(*files, walkPath)
	}
	return nil
}

func walkDirSkip(name string) error {
	if shouldSkipDir(name) {
		// Skipped directories contain dependencies, generated output, or local
		// harness state; they are outside the active product ratchet surface.
		return filepath.SkipDir
	}
	return nil
}

func cognitiveSwitchOrSelect(statement ast.Stmt, nesting int) (int, bool) {
	// Selects use the same clause body scoring as expression and type switches.
	if score, ok := cognitiveSwitchStatement(statement, nesting); ok {
		// Switch helpers already know whether the statement matched.
		return score, true
	}
	if current, ok := statement.(*ast.SelectStmt); ok {
		// Select cases have communication clauses but share switch body scoring.
		return switchStatementScore(current.Body.List, nesting), true
	}
	return 0, false
}
