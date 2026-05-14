package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func buildGitHubActionsWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteGitHubActions(opts.out, opts.target, opts.reportDir, witness.EnvironmentFromOS())
}
