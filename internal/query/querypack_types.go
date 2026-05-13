package query

type QueryPackResult struct {
	SchemaVersion      string                    `json:"schema_version"`
	QueryPackID        string                    `json:"query_pack_id"`
	QueryPackVersion   string                    `json:"query_pack_version"`
	RunID              string                    `json:"run_id,omitempty"`
	RunNonce           string                    `json:"run_nonce,omitempty"`
	SourceBaseline     string                    `json:"source_baseline,omitempty"`
	TopLevelAssessment string                    `json:"top_level_assessment,omitempty"`
	InputArtifacts     []QueryPackInputArtifact  `json:"input_artifacts"`
	QueryRows          map[string][]QueryPackRow `json:"query_rows"`
	OutputSafety       QueryPackOutputSafety     `json:"output_safety"`
}
