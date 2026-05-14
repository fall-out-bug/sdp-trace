package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

type witnessRecordBuilder func(witnessOptions) (witness.Record, error)

func witnessRecordBuilders() map[string]witnessRecordBuilder {
	// Builders map closed witness kinds to the package function that owns their
	// evidence interpretation.
	return map[string]witnessRecordBuilder{
		witness.KindGitHubActions: buildGitHubActionsWitness,
		witness.KindGitLabCI:      buildEnvelopeWitness,
		witness.KindBuildkite:     buildEnvelopeWitness,
		witness.KindCustomerPKI:   buildCustomerPKIWitness,
	}
}
