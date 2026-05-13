package query

type QueryPackInputArtifact struct {
	Role             string `json:"role"`
	SHA256           string `json:"sha256,omitempty"`
	PathRedactedID   string `json:"path_redacted_id"`
	SchemaVersion    string `json:"schema_version,omitempty"`
	ArtifactRequired bool   `json:"artifact_required"`
}
