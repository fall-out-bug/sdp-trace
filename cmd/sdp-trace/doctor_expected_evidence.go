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

func expectedEvidenceNoRequiredEventsCheck(contract trace.Contract) doctorCheck {
	// A contract with no required events cannot prove evidence coverage.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    string(trace.VerdictCannotVerify),
		Reason:   "contract has no required_events",
		Contract: contract.ContractID,
	}
}

func expectedEvidenceUnsupportedReferenceCheck(contract trace.Contract, missing []string) doctorCheck {
	// Unsupported references are reported as explicit gaps, not hidden in a
	// generic contract failure.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    string(trace.VerdictCannotVerify),
		Reason:   "contract references unsupported event types",
		Contract: contract.ContractID,
		Missing:  missing,
	}
}
