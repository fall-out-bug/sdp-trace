package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func doctorDefaultContractCheck(defaultContract trace.Contract) doctorCheck {
	// The default contract is a local development fallback, not proof that a
	// specific repository run satisfied its contract.
	return doctorCheck{
		ID:        "default_contract",
		State:     "pass",
		Reason:    "built-in contract is available for local development",
		Contract:  defaultContract.ContractID,
		Reference: defaultContract.Version,
	}
}

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

func unreadableDoctorContractResult(defaultContract trace.Contract, contractPath string) (trace.Contract, doctorCheck) {
	// A requested contract that cannot load keeps doctor in cannot_verify.
	return defaultContract, doctorCheck{
		ID:        "contract",
		State:     string(trace.VerdictCannotVerify),
		Reason:    "contract cannot be loaded",
		Reference: contractPath,
	}
}
