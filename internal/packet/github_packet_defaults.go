package packet

// Generated residual gaps are derived only for non-pass rows other than the
// residual-gaps control row itself.
func residualGapsForRows(rows []Row) []ResidualGap {
	gaps := []ResidualGap{}
	for _, row := range rows {
		if rowNeedsResidualGap(row) {
			gaps = append(gaps, residualGapForRow(row))
		}
	}
	return gaps
}

func rowNeedsResidualGap(row Row) bool {
	return row.State != StatePass && row.ID != "PC-RESIDUAL-GAPS"
}

func residualGapForRow(row Row) ResidualGap {
	return ResidualGap{RowID: row.ID, State: row.State, Reason: row.Reason, ClosureEvidence: "provide retained evidence for " + row.ID}
}

// Default decision owners are placeholders, not merge or release approval.
func defaultDecisionOwners() []DecisionOwner {
	return []DecisionOwner{
		{Decision: "merge", Owner: "maintainer", State: StateNotAssessed, Reason: "packet is not approval"},
		{Decision: "release", Owner: "release owner", State: StateNotAssessed, Reason: "packet is not release approval"},
		{Decision: "risk_acceptance", Owner: "risk owner", State: StateNotAssessed, Reason: "packet is not risk acceptance"},
		{Decision: "security_review", Owner: "security owner", State: StateNotAssessed, Reason: "packet is not security review"},
	}
}
