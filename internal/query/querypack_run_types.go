package query

import "encoding/json"

type runArtifact struct {
	SchemaVersion   string                     `json:"schema_version,omitempty"`
	RunID           string                     `json:"run_id"`
	RunNonce        string                     `json:"run_nonce,omitempty"`
	SourceBaseline  string                     `json:"source_baseline,omitempty"`
	EventRefs       []eventRef                 `json:"event_refs,omitempty"`
	VerifierStates  map[string]verifierState   `json:"verifier_states,omitempty"`
	RedactionDigest string                     `json:"redaction_policy_digest,omitempty"`
	Raw             map[string]json.RawMessage `json:"-"`
}

type eventRef struct {
	EventType string `json:"event_type"`
}

type verifierState struct {
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}
