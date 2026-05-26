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
