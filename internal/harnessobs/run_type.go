package harnessobs

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
