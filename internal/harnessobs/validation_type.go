package harnessobs

type Dimension struct {
	Family     string `json:"family"`
	Required   bool   `json:"required"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	EventCount int    `json:"event_count"`
}

type Validation struct {
	SchemaVersion      string      `json:"schema_version"`
	ProfileID          string      `json:"profile_id"`
	HarnessFamily      string      `json:"harness_family"`
	EventSchemaVersion string      `json:"event_schema_version"`
	ValidationState    string      `json:"validation_state"`
	ReasonCode         string      `json:"reason_code"`
	Dimensions         []Dimension `json:"dimensions"`
	EventCount         int         `json:"event_count"`
	ValidationDigest   string      `json:"validation_digest,omitempty"`
	NonAuthority       string      `json:"non_authority"`
}
