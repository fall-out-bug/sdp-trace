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
