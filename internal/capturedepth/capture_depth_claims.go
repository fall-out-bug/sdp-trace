package capturedepth

import "github.com/fall_out_bug/sdp-trace/internal/adaptercapture"

func unverifiedClaims(conditions []adaptercapture.Condition) []string {
	// Capture-depth is an investigation query, so every non-authoritative state
	// remains visible instead of being collapsed into a pass/fail summary.
	// Retention-limited and unsupported states are included because both block
	// complete replay even when the adapter emitted some evidence.
	out := []string{}
	for _, condition := range conditions {
		switch condition.State {
		case adaptercapture.StateCannotVerify,
			adaptercapture.StateNotAssessed,
			adaptercapture.StateMissingTelemetry,
			adaptercapture.StateNotIntegrated,
			adaptercapture.StateUnsupported,
			adaptercapture.StateRetentionLimited:
			out = append(out, condition.ID)
		}
	}
	return out
}
