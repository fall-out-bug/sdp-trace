package trace

// Contract controls required events for milestone verification.
type Contract struct {
	ContractID         string                `json:"contract_id"`
	Version            string                `json:"version"`
	RequiredEvents     []string              `json:"required_events"`
	RequiredEvidence   []EvidenceRequirement `json:"required_evidence,omitempty"`
	RequiredRuns       []RequiredRun         `json:"required_runs,omitempty"`
	LockRequiredBefore string                `json:"lock_required_before,omitempty"`
}

// RequiredRun names a contract-declared run that should be observed for an advisory gate.
type RequiredRun struct {
	ID               string   `json:"id"`
	WrapperName      string   `json:"wrapper_name"`
	RequiredEvidence []string `json:"required_evidence,omitempty"`
	Profile          string   `json:"profile,omitempty"`
}

// EvidenceRequirement names a contract-declared observation that can be
// matched against event payload fields without product-specific classifiers.
type EvidenceRequirement struct {
	ID            string `json:"id"`
	EventType     string `json:"event_type"`
	PayloadField  string `json:"payload_field"`
	PayloadEquals string `json:"payload_equals"`
}

// DefaultContract is the minimal local contract for first-milestone local recorder output.
var DefaultContract = Contract{
	ContractID: "local-default-v1",
	Version:    SchemaVersion,
	RequiredEvents: []string{
		string(EventRecorderAttached),
		string(EventRunStarted),
		string(EventCommandStarted),
		string(EventCommandFinished),
		string(EventRunClosed),
	},
	LockRequiredBefore: "run_started",
}
