package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func doctorContractCheck(contractPath string, defaultContract trace.Contract) (trace.Contract, doctorCheck) {
	if contractPath == "" {
		// No override means doctor can only report local default-contract
		// availability.
		return defaultDoctorContractResult(defaultContract)
	}
	// Loading the requested contract is the only source-bound proof for an
	// override contract path.
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		// Load failures lower only this control point to cannot_verify.
		return unreadableDoctorContractResult(defaultContract, contractPath)
	}
	// A loaded override becomes the contract returned to later doctor checks.
	return loadedDoctorContractResult(contract, contractPath)
}
