package harnessobs

type Limits struct {
	MaxLineBytes int `json:"max_line_bytes,omitempty"`
	MaxEvents    int `json:"max_events,omitempty"`
}

type Rule struct {
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
}

type Profile struct {
	SchemaVersion         string          `json:"schema_version"`
	ProfileID             string          `json:"profile_id"`
	HarnessFamily         string          `json:"harness_family"`
	EventSchemaVersion    string          `json:"event_schema_version"`
	RequiredEventFamilies []string        `json:"required_event_families"`
	OptionalEventFamilies []string        `json:"optional_event_families,omitempty"`
	RawRetentionPolicy    string          `json:"raw_retention_policy"`
	UnsupportedFields     []string        `json:"unsupported_fields,omitempty"`
	DegradationRules      map[string]Rule `json:"degradation_rules"`
	Limits                Limits          `json:"limits,omitempty"`
}
