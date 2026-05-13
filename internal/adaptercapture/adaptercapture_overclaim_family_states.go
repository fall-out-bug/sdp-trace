package adaptercapture

var insufficientEventFamilyStates = map[string]bool{
	StateMissingTelemetry: true,
	StateUnsupported:      true,
	StateNotIntegrated:    true,
	StateNotAssessed:      true,
	StateCannotVerify:     true,
	StateRetentionLimited: true,
}
