package demo

import (
	"fmt"
	"strings"
)

func buildTimeline(rows []RunRow) string {
	// Timeline output is derived from sanitized row fields and escapes table
	// delimiters so command text cannot corrupt Markdown structure.
	var builder strings.Builder
	builder.WriteString("# SDP Trace Timeline\n\n")

	builder.WriteString("| Run | Kind | Result | Trust Scope | Command | Exit |\n")
	builder.WriteString("|-----|------|--------|-------------|---------|------|\n")
	for _, row := range rows {
		builder.WriteString(timelineRow(row))
	}
	return builder.String()
}

func timelineRow(row RunRow) string {
	// timelineRow keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
		escapeMD(row.Name),
		escapeMD(row.Kind),
		escapeMD(string(row.Result)),
		escapeMD(string(row.TrustScope)),
		escapeMD(row.Command),
		escapeMD(timelineExit(row)),
	)
}

func timelineExit(row RunRow) string {
	// timelineExit keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if row.ExitCode == nil {
		return ""
	}

	return fmt.Sprintf("%d", *row.ExitCode)
}

func escapeMD(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}
