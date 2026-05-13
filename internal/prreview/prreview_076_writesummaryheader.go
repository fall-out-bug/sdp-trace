package prreview

import (
	"fmt"

	"strings"
)

func writeSummaryHeader(b *strings.Builder, validation Validation) {
	// writeSummaryHeader keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	fmt.Fprintf(b, "Review coverage: %s\n", validation.ReviewCoverageState)
	fmt.Fprintf(b, "CI state: %s\n", validation.CIState)
	fmt.Fprintf(b, "Authority scope: %s\n", validation.AuthorityScope)
	fmt.Fprintf(b, "Merge decision: %s\n", validation.MergeDecision)
	fmt.Fprintf(b, "Release decision: %s\n", validation.ReleaseDecision)
	fmt.Fprintf(b, "Risk acceptance: %s\n", validation.RiskAcceptance)
	b.WriteString("This is review-record evidence only; merge, release, and risk decisions remain external.\n")
}
