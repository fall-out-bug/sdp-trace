package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

type protectedGateInputs struct {
	signed  checkpoint.SignedCheckpoint
	policy  checkpoint.TrustedCheckpointPolicy
	witness demo.WitnessSummary
}
