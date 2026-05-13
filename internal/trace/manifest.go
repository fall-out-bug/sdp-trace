package trace

// RunManifest records run-level metadata.
type RunManifest struct {
	SchemaVersion   string `json:"schema_version"`
	RunID           string `json:"run_id"`
	RecorderVersion string `json:"recorder_version"`
	CreatedAt       string `json:"created_at"`
	ClosedAt        string `json:"closed_at,omitempty"`
	Task            string `json:"task,omitempty"`
	ContractID      string `json:"contract_id"`
	ContractPath    string `json:"contract_path,omitempty"`
	ContractDigest  string `json:"contract_digest,omitempty"`
	SourceSnapshot  string `json:"source_snapshot_digest"`
	SourceState     string `json:"source_snapshot_state"`
	EventCount      int    `json:"event_count"`
	EventChainHead  string `json:"event_chain_head"`
	FinalChainHead  string `json:"final_chain_head"`
	ClosureState    string `json:"closure_state"`
}

// Validate checks required manifest fields for shared verifier paths.
func (manifest RunManifest) Validate() error {
	// Manifest validation is intentionally shape-only; chain and event-count
	// consistency are checked against the run directory during replay.
	return firstValidationError(
		requiredString(manifest.SchemaVersion, "run manifest missing schema_version"),
		requiredString(manifest.RunID, "run manifest missing run_id"),
		nonNegative(manifest.EventCount, "run manifest event_count must be >= 0"),
		requiredString(manifest.ContractID, "run manifest missing contract_id"),
	)
}
