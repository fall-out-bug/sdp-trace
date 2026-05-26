package main

import (
	"fmt"
	"strings"
)

// checkForbiddenTrackedPaths returns findings for local or generated paths
// that must never be committed.
func checkForbiddenTrackedPaths(tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if forbiddenPathReason(f) != "" {
			findings = append(findings, fmt.Sprintf("tracked forbidden path: %s", f))
		}
	}
	return findings
}

// forbiddenPathReason returns a non-empty label for paths that are always
// forbidden in the tracked tree.
func forbiddenPathReason(f string) string {
	switch {
	case strings.HasPrefix(f, ".worktrees/"):
		return ".worktrees"
	case strings.HasPrefix(f, ".codex-subagents/runs/"):
		return ".codex-subagents/runs"
	case strings.HasPrefix(f, ".sdp-trace-"):
		return ".sdp-trace-*"
	}
	return ""
}
