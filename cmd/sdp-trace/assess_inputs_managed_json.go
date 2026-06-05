package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"path/filepath"
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

func readManagedJSONInputs(opts *flagSet) (managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var registry managed.Registry
	if err := readJSONFile(opts.stringValue("adapter-registry"), &registry); err != nil {
		// Registry read failures stop before run or witness evidence is combined.
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return readManagedRunAndWitness(opts, registry)
}

func readManagedRunAndWitness(opts *flagSet, registry managed.Registry) (managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var runEvidence managed.RunEvidence
	// Run evidence lives under the run directory contract; callers pass the run
	// root rather than a free-form JSON path.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	var witness managed.Witness
	if err := readJSONFile(opts.stringValue("managed-witness"), &witness); err != nil {
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return registry, runEvidence, witness, nil
}
