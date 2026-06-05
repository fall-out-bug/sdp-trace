package prreview

import (
	"fmt"
	"strings"
)

func Summarize(validation Validation, ledger Ledger) string {
	// Summary is a rendered view over validation and ledger state; it does not
	// recompute coverage or mutate dispositions.
	var b strings.Builder
	writeSummaryHeader(&b, validation)
	writeSummaryPlanes(&b, validation.PlaneResults)
	writeSummaryFindings(&b, ledger.Findings)
	return b.String()
}

// writeSummaryHeader renders the authority boundary explicitly so a review
// summary cannot be mistaken for merge, release, or risk approval.
func writeSummaryHeader(b *strings.Builder, validation Validation) {
	fmt.Fprintf(b, "Review coverage: %s\n", validation.ReviewCoverageState)
	fmt.Fprintf(b, "CI state: %s\n", validation.CIState)
	fmt.Fprintf(b, "Authority scope: %s\n", validation.AuthorityScope)
	fmt.Fprintf(b, "Merge decision: %s\n", validation.MergeDecision)
	fmt.Fprintf(b, "Release decision: %s\n", validation.ReleaseDecision)
	fmt.Fprintf(b, "Risk acceptance: %s\n", validation.RiskAcceptance)
	b.WriteString("This is review-record evidence only; merge, release, and risk decisions remain external.\n")
}
