package main

import (
	"os"
	"path/filepath"
)

func analyzePaths(paths []string) (qualityReport, error) {
	var files []string
	for _, path := range paths {
		// Discovery expands each CLI argument before analysis so package-level
		// metrics are computed from a single deterministic file list.
		found, err := goFilesInPath(path)
		if err != nil {
			return qualityReport{}, err
		}
		files = append(files, found...)
	}
	return analyzeFiles(files)
}

func goFilesInPath(path string) ([]string, error) {
	// Discovery accepts both explicit files and directories because local checks
	// often target a single changed file while CI passes package roots.
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	// Direct file inputs still pass through the production filter; tests,
	// generated files, and unsupported suffixes stay outside gate scope.
	if !info.IsDir() {
		if isProductionGo(path) {
			return []string{path}, nil
		}
		return nil, nil
	}
	return walkGoFiles(path)
}

func walkGoFiles(path string) ([]string, error) {
	var files []string
	// WalkDir preserves deterministic path order while collectGoFile owns the
	// production-file and skipped-directory boundary.
	err := filepath.WalkDir(path, func(walkPath string, entry os.DirEntry, walkErr error) error {
		// The callback is intentionally tiny so traversal errors and scope
		// filtering stay in the reusable collector.
		return collectGoFile(walkPath, entry, walkErr, &files)
	})
	return files, err
}
