package main

import "os"

func repoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return findRepoRoot(cwd)
}
