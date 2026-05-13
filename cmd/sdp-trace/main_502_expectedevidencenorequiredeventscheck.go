package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func expectedEvidenceNoRequiredEventsCheck(contract trace.Contract) doctorCheck {
	// A contract with no required events cannot prove evidence coverage.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    string(trace.VerdictCannotVerify),
		Reason:   "contract has no required_events",
		Contract: contract.ContractID,
	}
}
