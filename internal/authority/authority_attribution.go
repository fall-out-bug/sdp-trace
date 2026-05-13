package authority

import (
	"strings"
)

func actorAttributionState(action ObservedAction) string {
	// actorAttributionState keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if strings.TrimSpace(action.ActorID) == "" {

		return AttributionNotAssessed
	}
	switch action.SourceType {
	case "harness_log", "manual_import", "pr_api", "ci_artifact":

		return AttributionVerified
	default:
		return AttributionNotAssessed
	}
}
