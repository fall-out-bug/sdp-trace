package posture

import (
	"errors"
)

func validateRefusalRow(row RefusalRow) error {
	return malformedRowError(malformedRefusalRow(row), "malformed posture export refusal_row")
}

func malformedRowError(malformed bool, message string) error {
	// malformedRowError keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if malformed {
		return errors.New(message)
	}
	return nil
}

func malformedRefusalRow(row RefusalRow) bool {
	return missingRefusalIdentity(row) || malformedRefusalState(row) || malformedOptionalRefusalWindow(row)
}

func missingRefusalIdentity(row RefusalRow) bool {
	return row.ID == "" || !safeLabel(row.InputID)
}

func malformedRefusalState(row RefusalRow) bool {
	return !validRefusalReason(row.RefusalReason) || !validInputTrustState(row.InputTrustState)
}

func malformedOptionalRefusalWindow(row RefusalRow) bool {
	return row.TimeWindow != "" && !safeLabel(row.TimeWindow)
}
