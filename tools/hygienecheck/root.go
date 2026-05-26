package main

import (
	"path/filepath"
	"runtime"
)

// repoRoot resolves the repository root from the source file location so that
// the tool checks the same tree regardless of the working directory. It falls
// back to the current directory when caller resolution fails.
func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
