package packet

func residualGapsForRows(rows []Row) []ResidualGap {
	// residualGapsForRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	gaps := []ResidualGap{}
	for _, row := range rows {
		if row.State == StatePass || row.ID == "PC-RESIDUAL-GAPS" {
			continue
		}

		gaps = append(gaps, ResidualGap{RowID: row.ID, State: row.State, Reason: row.Reason, ClosureEvidence: "provide retained evidence for " + row.ID})
	}
	return gaps
}
