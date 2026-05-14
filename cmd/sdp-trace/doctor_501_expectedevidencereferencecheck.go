package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func expectedEvidenceReferenceCheck(contract trace.Contract) doctorCheck {
	if len(contract.RequiredEvents) == 0 {
		return expectedEvidenceNoRequiredEventsCheck(contract)
	}
	// Required event and evidence references are checked separately so drift is
	// reported with concrete missing keys.
	missing := expectedEvidenceReferenceGaps(contract)
	if len(missing) > 0 {
		return expectedEvidenceUnsupportedReferenceCheck(contract, missing)
	}
	// A pass here means only that the local event vocabulary can represent the
	// contract's required evidence references.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    "pass",
		Reason:   "contract required events and evidence references are supported by the current local event model",
		Contract: contract.ContractID,
	}
}
