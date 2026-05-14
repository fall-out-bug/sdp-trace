package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func loadManagedJSONInputs(opts *flagSet) (managed.Policy, managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var policy managed.Policy
	if err := readJSONFile(opts.stringValue("managed-policy"), &policy); err != nil {
		return managed.Policy{}, managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	// Policy is read first because it defines how registry, run, and witness rows
	// will be interpreted.
	registry, runEvidence, witness, err := readManagedJSONInputs(opts)
	if err != nil {
		return managed.Policy{}, managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return policy, registry, runEvidence, witness, nil
}
