package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func defaultDoctorContractResult(defaultContract trace.Contract) (trace.Contract, doctorCheck) {
	// The default contract branch is explicitly local evidence; it says nothing
	// about a repository-specific contract file.
	return defaultContract, doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "default contract is available",
		Contract:  defaultContract.ContractID,
		Reference: "local-default-v1",
	}
}
