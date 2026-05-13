package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func readCheckpointPolicy(path string) (*checkpoint.TrustedCheckpointPolicy, error) {
	if path == "" {
		// An absent policy leaves signer trust to the checkpoint verifier's
		// default local semantics; it is not treated as a green CI policy.
		return nil, nil
	}
	var loaded checkpoint.TrustedCheckpointPolicy
	if err := readJSONFile(path, &loaded); err != nil {
		return nil, err
	}
	return &loaded, nil
}
