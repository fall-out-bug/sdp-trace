package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func unreadableDoctorContractResult(defaultContract trace.Contract, contractPath string) (trace.Contract, doctorCheck) {
	// A requested contract that cannot load keeps doctor in cannot_verify.
	return defaultContract, doctorCheck{
		ID:        "contract",
		State:     string(trace.VerdictCannotVerify),
		Reason:    "contract cannot be loaded",
		Reference: contractPath,
	}
}
