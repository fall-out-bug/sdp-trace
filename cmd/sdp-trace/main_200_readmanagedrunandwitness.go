package main

import (
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

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
