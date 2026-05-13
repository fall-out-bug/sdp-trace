package adaptercapture

func topLevel(conditions []Condition) string {
	// Top-level state reports the strongest non-pass condition without averaging or
	// hiding adapter capture gaps.
	highest := StatePass
	for _, condition := range conditions {
		if condition.State == StateFail {

			return StateFail
		}
		switch condition.State {
		case StateCannotVerify, StateNotAssessed, StateMissingTelemetry, StateNotIntegrated, StateUnsupported, StateRetentionLimited:

			highest = StateCannotVerify
		}
	}
	return highest
}
