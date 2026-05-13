package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadManagedInput(opts *flagSet) (managed.Input, error) {
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		return managed.Input{}, err
	}
	// Contract required events are copied so managed evaluation cannot mutate the
	// loaded trace contract.
	return loadManagedInputWithContract(opts, contract)
}
