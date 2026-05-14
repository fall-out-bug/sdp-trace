package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
