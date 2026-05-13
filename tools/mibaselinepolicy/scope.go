package main

import "strings"

// productionRoots is the portable product-code scope for MI baseline policy.
var productionRoots = []string{"cmd/", "internal/", "tools/"}

// productionGoFile excludes tests and fixture-like paths from ratchet scope.
func productionGoFile(path string) bool {
	// Baseline policy tracks only active product Go code; tests and doc fixtures
	// can change without weakening the maintained production ratchets.
	return strings.HasSuffix(path, ".go") &&
		!strings.HasSuffix(path, "_test.go") &&
		inProductionTree(path)
}

// inProductionTree matches git-reported repository paths against product roots.
func inProductionTree(path string) bool {
	for _, root := range productionRoots {
		// Changed paths come from git and are already slash-separated repo
		// paths, so prefix matching is the intended production-scope contract.
		if strings.HasPrefix(path, root) {
			return true
		}
	}
	return false
}
