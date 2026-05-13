package repoobserver

import (
	"fmt"
	"strings"
)

type DiffSummary struct {
	Path    string `json:"path"`
	Action  string `json:"action"`
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
	Summary string `json:"summary"`
	Backup  string `json:"backup,omitempty"`
}

func writeHumanTableDiffSummary(b *strings.Builder, summary []DiffSummary) {
	// Force summaries intentionally expose safe digests/counts, not full file
	// content.
	if len(summary) == 0 {
		return
	}
	b.WriteString("\nForce diff summary\n")
	for _, item := range summary {
		writeHumanTableDiffItem(b, item)
	}
}

func writeHumanTableDiffItem(b *strings.Builder, item DiffSummary) {
	// Diff summaries contain safe action/digest metadata only; file contents are
	// intentionally omitted from the human table.
	fmt.Fprintf(b, "- %s: %s", item.Path, item.Action)
	if item.Before != "" || item.After != "" {
		fmt.Fprintf(b, " [%s -> %s]", item.Before, item.After)
	}
	if item.Backup != "" {
		fmt.Fprintf(b, " (backup: %s)", item.Backup)
	}
	b.WriteString("\n")
}
