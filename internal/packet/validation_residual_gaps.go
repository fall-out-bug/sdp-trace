package packet

import "strings"

// Residual gaps must point at portable required rows and carry a reason; blank
// gap records would otherwise hide unresolved trust work.
func (v *bundleValidator) validateResidualGap(gap ResidualGap) {
	if !requiredRow(gap.RowID) {
		v.add("residual gap has unknown row id %q", gap.RowID)
	}
	if strings.TrimSpace(gap.Reason) == "" {
		v.add("residual gap for %s requires reason", gap.RowID)
	}
}

// Coverage is checked from rows to gaps so every non-pass row has an explicit
// residual explanation unless it is the residual-gaps summary row itself.
func (v *bundleValidator) validateResidualCoverage(rows map[string]Row) {
	for _, row := range rows {
		v.validateResidualCoverageForRow(row)
	}
}

// Pass rows and the residual-gaps rollup row do not need their own residual gap
// record; every other non-pass row must be explicitly explained.
func (v *bundleValidator) validateResidualCoverageForRow(row Row) {
	if residualCoverageExempt(row) {
		return
	}
	if !gapForRow(v.bundle.Packet.ResidualGaps, row.ID) {
		v.add("%s non-pass row requires residual gap explanation", row.ID)
	}
}

// Keep this exemption narrow: adding more exempt rows weakens packet closure.
func residualCoverageExempt(row Row) bool {
	return row.ID == "PC-RESIDUAL-GAPS" || row.State == StatePass
}
