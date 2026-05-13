package trace

import "encoding/json"

// Event captures one append-only record in a run.
type Event struct {
	SchemaVersion    string           `json:"schema_version"`
	RunID            string           `json:"run_id"`
	EventID          string           `json:"event_id"`
	Sequence         int              `json:"sequence"`
	EventType        EventType        `json:"event_type"`
	Timestamp        string           `json:"timestamp"`
	PrevEventHash    string           `json:"prev_event_hash"`
	HashAlgorithm    string           `json:"hash_algorithm"`
	Canonicalization Canonicalization `json:"canonicalization"`
	PayloadDigest    string           `json:"payload_digest"`
	EventPayload     map[string]any   `json:"event_payload"`
	Payload          json.RawMessage  `json:"payload"`
	EventHash        string           `json:"event_hash"`
	ObservedBy       string           `json:"observed_by,omitempty"`
	ClockMonotonic   int64            `json:"clock_monotonic,omitempty"`
}

type Canonicalization struct {
	Algorithm string `json:"algorithm"`
	Version   string `json:"version"`
}
