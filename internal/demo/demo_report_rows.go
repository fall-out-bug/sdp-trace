package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifiedRowsForContract(target, contractPath string) ([]RunRow, trace.Contract, error) {
	// verifiedRowsForContract keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		return nil, trace.Contract{}, err
	}
	rows, err := VerifiedRows(target, contract)
	if err != nil {
		return nil, trace.Contract{}, err
	}
	return rows, contract, nil
}
