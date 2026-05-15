package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkAbsoluteHomePaths returns findings for .md files that contain an
// absolute /home/... path, excluding grandfathered historical entries.
func checkAbsoluteHomePaths(root string, tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if finding := homePathFinding(root, f); finding != "" {
			findings = append(findings, finding)
		}
	}
	return findings
}

// homePathFinding evaluates one tracked file for absolute home paths.
func homePathFinding(root, f string) string {
	if !strings.HasSuffix(f, ".md") {
		return ""
	}
	if homePathAllowlist[f] {
		return ""
	}
	if containsHomePath(filepath.Join(root, f)) {
		return fmt.Sprintf("absolute home path in doc: %s", f)
	}
	return ""
}

// containsHomePath reads a file and reports whether it contains an absolute
// /home/ path segment. I/O errors are treated as negative.
func containsHomePath(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return bytesContainHomePath(data)
}
