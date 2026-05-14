package main

import (
	"fmt"
	"strings"
)

// checkRootArtifactClutter returns findings for tracked files that belong in
// a spec directory rather than at the repository root.
func checkRootArtifactClutter(tracked []string) []string {
	var findings []string
	for _, f := range tracked {
		if rootArtifactClutterReason(f) != "" {
			findings = append(findings, fmt.Sprintf("root artifact clutter: %s", f))
		}
	}
	return findings
}

// rootArtifactClutterReason names the PR-scoped files and directories that
// must not be tracked at the repository root.
func rootArtifactClutterReason(f string) string {
	switch f {
	case "PR_DESCRIPTION.md", "design-note.md":
		return f
	}
	if strings.HasPrefix(f, "reviews/") {
		return "reviews/"
	}
	return ""
}

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
