package harnessobs

type UnavailableField struct {
	Field      string `json:"field"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
}

type Event struct {
	EventID            string             `json:"event_id"`
	EventSchemaVersion string             `json:"event_schema_version"`
	EventFamily        string             `json:"event_family"`
	EventType          string             `json:"event_type"`
	ObservedAt         string             `json:"observed_at,omitempty"`
	SourceRef          string             `json:"source_ref"`
	SourceDigest       string             `json:"source_digest"`
	TaskRef            string             `json:"task_ref,omitempty"`
	OperationRef       string             `json:"operation_ref,omitempty"`
	ActorRef           string             `json:"actor_ref,omitempty"`
	ContentState       string             `json:"content_state"`
	UnavailableFields  []UnavailableField `json:"unavailable_fields,omitempty"`
}

type Run struct {
	SchemaVersion      string   `json:"schema_version"`
	ProfileID          string   `json:"profile_id"`
	HarnessFamily      string   `json:"harness_family"`
	EventSchemaVersion string   `json:"event_schema_version"`
	SourcePath         string   `json:"source_path"`
	SourceDigest       string   `json:"source_digest"`
	EventCount         int      `json:"event_count"`
	EventRefs          []string `json:"event_refs"`
	CreatedAt          string   `json:"created_at"`
}
