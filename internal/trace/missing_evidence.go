package trace

// GenerateMissingEvidenceTable emits expected/observed rows for a contract.
func GenerateMissingEvidenceTable(contract Contract, observed map[string]bool) MissingEvidenceTable {
	// Missing rows are evidence gaps, not failed gates; callers decide how the
	// partial replayability impact affects a verdict.
	// Observed events are keyed by event type because the default contract only
	// requires event-class presence.
	rows := make([]MissingEvidenceRow, 0, len(contract.RequiredEvents))
	for _, eventType := range contract.RequiredEvents {
		if observed[eventType] {
			// Observed required events do not produce rows; the table is a gap
			// report rather than a full checklist.
			continue
		}
		rows = append(rows, missingEvidenceRow(eventType))
	}
	return MissingEvidenceTable{
		// The contract ID binds the missing-evidence rows to the loaded spec.
		ContractID: contract.ContractID,
		Rows:       rows,
	}
}

func missingEvidenceRow(eventType string) MissingEvidenceRow {
	// Required-by-contract rows stay portable and avoid harness-specific
	// remediation text.
	return MissingEvidenceRow{
		ExpectedEvent:       eventType,
		ObservedState:       string(EvidenceStateMissing),
		Reason:              "required_by_contract",
		ReplayabilityImpact: string(ReplayabilityPartial),
	}
}
