package main

import (
	"os"
	"path/filepath"
)

func findRepoRoot(cwd string) string {
	for cwd != "" {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd
		}
		cwd = parentDir(cwd)
	}
	return "."
}
