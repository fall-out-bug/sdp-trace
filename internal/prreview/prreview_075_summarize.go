package prreview

import (
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
