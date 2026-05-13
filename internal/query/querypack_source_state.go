package query

var sourceStateRows = map[string]string{
	"pass":                   RowStatePresent,
	"":                       RowStateCannotVerify,
	"fail":                   RowStateIssueObserved,
	RowStateCannotVerify:     RowStateCannotVerify,
	RowStateNotAssessed:      RowStateNotAssessed,
	RowStateMissingTelemetry: RowStateMissingTelemetry,
	RowStateUnsupported:      RowStateUnsupported,
	RowStateNotIntegrated:    RowStateNotIntegrated,
	RowStateRetentionLimited: RowStateRetentionLimited,
}

func mapSourceState(state string) string {
	// mapSourceState keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if mapped, ok := sourceStateRows[state]; ok {
		return mapped
	}
	return RowStateCannotVerify
}
