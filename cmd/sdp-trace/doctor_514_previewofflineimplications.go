package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func previewOfflineImplications() []previewOfflineImplication {
	// Offline implications tell the user which evidence must be re-collected in
	// CI or external systems before trust can be upgraded.
	return []previewOfflineImplication{
		{
			Requirement: "ci_witnessed",
			State:       string(trace.ObservationStateOfflineDev),
			Reason:      "rerun in CI with OIDC before using CI witness evidence",
		},
		{
			Requirement: "external_witnessed",
			State:       string(trace.ObservationStateNotIntegrated),
			Reason:      "external witness profile is not implemented in Block 13B",
		},
	}
}
