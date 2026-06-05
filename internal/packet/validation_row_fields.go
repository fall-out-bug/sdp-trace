package packet

import "strings"

// Row validation is deliberately ordered: structural row fields are checked
// before evidence refs so callers get stable accumulated diagnostics.
func (v *bundleValidator) validateRow(row Row) {
	v.validateRowRequiredFields(row)
	v.validateRowReason(row)
	v.validatePassRowEvidence(row)
	v.validateRowEvidenceRefs(row)
}

// Required field checks stay together because state, summary, and owner form
// the minimal row identity shown in packet reports.
func (v *bundleValidator) validateRowRequiredFields(row Row) {
	v.validateRowState(row)
	v.validateRowSummary(row)
	v.validateRowOwner(row)
}

func (v *bundleValidator) validateRowState(row Row) {
	if !states[row.State] {
		v.add("%s has unknown state %q", row.ID, row.State)
	}
}

func (v *bundleValidator) validateRowSummary(row Row) {
	if strings.TrimSpace(row.Summary) == "" {
		v.add("%s requires summary", row.ID)
	}
}

func (v *bundleValidator) validateRowOwner(row Row) {
	if strings.TrimSpace(row.Owner) == "" {
		v.add("%s requires owner", row.ID)
	}
}

// Non-pass rows need a human-readable reason so partial or failed gates cannot
// be mistaken for a complete pass in downstream packet rendering.
func (v *bundleValidator) validateRowReason(row Row) {
	if missingReasonStates[row.State] && strings.TrimSpace(row.Reason) == "" {
		v.add("%s state %s requires reason", row.ID, row.State)
	}
}
