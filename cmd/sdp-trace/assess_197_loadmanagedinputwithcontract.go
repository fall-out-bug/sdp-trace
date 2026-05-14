package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadManagedInputWithContract(opts *flagSet, contract trace.Contract) (managed.Input, error) {
	// Managed assessment combines a trace contract with managed-harness artifacts;
	// each JSON input is decoded before the final package input is assembled.
	policy, registry, runEvidence, witness, err := loadManagedJSONInputs(opts)
	if err != nil {
		return managed.Input{}, err
	}
	// Required events are copied into the managed package shape so downstream
	// evaluation cannot mutate the trace contract.
	return managed.Input{
		Contract: managed.Contract{RequiredEventTypes: append([]string(nil), contract.RequiredEvents...)},
		Policy:   policy,
		Registry: registry,
		Run:      runEvidence,
		Witness:  witness,
	}, nil
}
