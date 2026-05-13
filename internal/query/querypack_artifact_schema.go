package query

import "encoding/json"

func readArtifactSchemaVersion(payload []byte) string {
	// readArtifactSchemaVersion keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var envelope struct {
		SchemaVersion string `json:"schema_version"`
	}
	if err := json.Unmarshal(payload, &envelope); err == nil {
		return envelope.SchemaVersion
	}
	return ""
}
