package repoobserver

import "strings"

func HumanTable(status Status) string {
	// Human output keeps install and proof state separate instead of compressing
	// them into a single health score.
	var b strings.Builder
	writeHumanTableHeader(&b, status)
	writeHumanTableSurfaces(&b, status.Surfaces)
	writeHumanTableDiffSummary(&b, status.ForceDiffSummary)
	b.WriteString("\nNote: core.hooksPath is local checkout configuration and is not committed into the repository.\n")
	return b.String()
}
