package trace

// MissingEvidenceRow describes one required-event mismatch.
type MissingEvidenceRow struct {
	ExpectedEvent       string `json:"expected_event"`
	ObservedState       string `json:"observed_state"`
	Reason              string `json:"reason"`
	PolicyReference     string `json:"policy_reference,omitempty"`
	ReplayabilityImpact string `json:"replayability_impact"`
}

// MissingEvidenceTable summarizes observed-vs-required events.
type MissingEvidenceTable struct {
	ContractID string               `json:"contract_id"`
	Rows       []MissingEvidenceRow `json:"rows"`
}

// VerifierResult is the live machine-readable verification envelope.
type VerifierResult struct {
	RunID         string          `json:"run_id"`
	Result        VerifierVerdict `json:"result"`
	TrustScope    TrustScope      `json:"trust_scope"`
	Completeness  Completeness    `json:"completeness"`
	Replayability Replayability   `json:"replayability"`
	Reason        string          `json:"reason"`
	RunDir        string          `json:"run_dir"`
}

// IntegrityAudit records one structural issue for explainability.
type IntegrityAudit struct {
	RunID   string            `json:"run_id"`
	Issue   string            `json:"issue"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details,omitempty"`
}
