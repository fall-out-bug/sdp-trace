package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadedDoctorContractResult(contract trace.Contract, contractPath string) (trace.Contract, doctorCheck) {
	// Loaded override contracts are reported with their id and path so doctor
	// output remains replayable from the local command invocation.
	return contract, doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "contract can be loaded",
		Contract:  contract.ContractID,
		Reference: contractPath,
	}
}
