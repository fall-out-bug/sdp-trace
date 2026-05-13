package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func readManagedJSONInputs(opts *flagSet) (managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var registry managed.Registry
	if err := readJSONFile(opts.stringValue("adapter-registry"), &registry); err != nil {
		// Registry read failures stop before run or witness evidence is combined.
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return readManagedRunAndWitness(opts, registry)
}
