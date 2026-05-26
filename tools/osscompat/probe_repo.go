package main

import (
	"os"
	"path/filepath"
)

// repoRoot returns the repository root by walking up from the current
// working directory until it finds a .git directory or reaches the filesystem
// root. It falls back to "." if the root cannot be determined.
func repoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return findGitRoot(cwd)
}

func findGitRoot(cwd string) string {
	for {
		if hasGitDir(cwd) {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			// Reaching filesystem root means the probe is not inside a checkout.
			return "."
		}
		cwd = parent
	}
}

func hasGitDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}
