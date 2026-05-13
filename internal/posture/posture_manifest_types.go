package posture

type SelectionManifest struct {
	SchemaVersion           string             `json:"schema_version"`
	ProfileID               string             `json:"profile_id"`
	ProfileVersion          string             `json:"profile_version,omitempty"`
	GroupingSetID           string             `json:"grouping_set_id"`
	FreshnessBoundary       string             `json:"freshness_boundary"`
	DimensionExposurePolicy []string           `json:"dimension_exposure_policy"`
	CurrentWindow           string             `json:"current_window"`
	PreviousWindow          string             `json:"previous_window"`
	Repositories            []RepositoryWindow `json:"repositories"`
	Handoff                 map[string]string  `json:"handoff,omitempty"`
}

// RepositoryWindow binds one selected repository window to its replayable
// query, digest, and optional posture-signal inputs.
type RepositoryWindow struct {
	InputID                string `json:"input_id"`
	Repo                   string `json:"repo"`
	Team                   string `json:"team"`
	Service                string `json:"service"`
	Harness                string `json:"harness"`
	ChangeType             string `json:"change_type"`
	TimeWindow             string `json:"time_window"`
	InputObservedAt        string `json:"input_observed_at"`
	QueryPackResult        string `json:"query_pack_result"`
	ArtifactDigestManifest string `json:"artifact_digest_manifest"`
	PostureSignalManifest  string `json:"posture_signal_manifest,omitempty"`
}

// DigestManifest records artifact digests that must match before selected
// query-pack evidence can enter posture aggregation.
type DigestManifest struct {
	SchemaVersion string           `json:"schema_version"`
	Artifacts     []DigestArtifact `json:"artifacts"`
}

type DigestArtifact struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// SignalManifest carries optional posture signals that may refine metric counts
// only after signal payload safety checks pass.
type SignalManifest struct {
	SchemaVersion string          `json:"schema_version"`
	Signals       []PostureSignal `json:"signals"`
}

// PostureSignal contains row-bound posture evidence that is never trusted as
// free text; every field is screened before aggregation.
type PostureSignal struct {
	RowRef               string `json:"row_ref"`
	WitnessScope         string `json:"witness_scope,omitempty"`
	ObserverState        string `json:"observer_state,omitempty"`
	OverrideMarker       string `json:"override_marker,omitempty"`
	LateAttachMarker     string `json:"late_attach_marker,omitempty"`
	ContractChangeMarker string `json:"contract_change_marker,omitempty"`
}

// ExportResult is the portable posture export surface produced from trusted
