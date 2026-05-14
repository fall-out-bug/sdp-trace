package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

func allowedWitnessKind(kind string) bool {
	switch kind {
	case witness.KindGitHubActions, witness.KindGitLabCI, witness.KindBuildkite, witness.KindCustomerPKI:
		// Each allowed kind has an explicit builder and schema contract.
		return true
	default:
		return false
	}
}
