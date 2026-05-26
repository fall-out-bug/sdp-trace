package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var currentClosureDocs = map[string]bool{
	"docs/closure-decision-ledger.md": true,
	"docs/open-task-breakdown.md":     true,
	"docs/spec-closure-route.md":      true,
	"docs/spec-reality-ledger.md":     true,
}

func checkCurrentDemoRepoDrift(root string, tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if !currentClosureDocs[f] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
		if err != nil {
			findings = append(findings, fmt.Sprintf("read current closure doc: %s: %v", f, err))
			continue
		}
		if strings.Contains(string(data), "sdp-trace-demo-ci-pilot") {
			findings = append(findings, fmt.Sprintf("stale demo repo reference in current closure doc: %s", f))
		}
	}
	return findings
}
