package posture

type ExportResult struct {
	SchemaVersion        string            `json:"schema_version"`
	ExportProfileID      string            `json:"export_profile_id"`
	ExportProfileVersion string            `json:"export_profile_version"`
	ExportID             string            `json:"export_id"`
	Producer             string            `json:"producer"`
	GeneratedAt          string            `json:"generated_at"`
	GroupingSetID        string            `json:"grouping_set_id"`
	ActiveGroupingKeys   []string          `json:"active_grouping_keys"`
	InputSelection       []InputSelection  `json:"input_selection"`
	MetricRows           []MetricRow       `json:"metric_rows"`
	MovementRows         []MovementRow     `json:"movement_rows"`
	MovementSummary      MovementSummary   `json:"movement_summary"`
	RefusalRows          []RefusalRow      `json:"refusal_rows"`
	Handoff              map[string]string `json:"handoff"`
	OutputSafety         OutputSafety      `json:"output_safety"`
}

// InputSelection records selected input identity without exposing raw artifact
// paths in downstream posture output.
type InputSelection struct {
	InputID         string `json:"input_id"`
	Repository      string `json:"repo"`
	TimeWindow      string `json:"time_window"`
	PathRedactedID  string `json:"path_redacted_id"`
	SHA256          string `json:"sha256,omitempty"`
	InputTrustState string `json:"input_trust_state"`
}

// MetricRow is an evidence-backed aggregate for one metric, window, and
// dimension key.
type MetricRow struct {
	ID                      string            `json:"id"`
	MetricID                string            `json:"metric_id"`
	MetricVersion           string            `json:"metric_version"`
	Numerator               int               `json:"numerator"`
	Denominator             int               `json:"denominator"`
	Unit                    string            `json:"unit"`
	TimeWindow              string            `json:"time_window"`
	Dimensions              map[string]string `json:"dimensions"`
	DimensionKey            string            `json:"dimension_key"`
	SourceInputRefs         []string          `json:"source_input_refs"`
	SourceArtifactDigestSet string            `json:"source_artifact_digest_set_hash"`
	SourceFieldState        string            `json:"source_field_state"`
	NotAssessedCount        int               `json:"not_assessed_count"`
	InputTrustStateSummary  map[string]int    `json:"input_trust_state_summary"`
}

// MovementRow compares matching metric rows across the selected current and
// previous windows without fabricating missing-window evidence.
type MovementRow struct {
	ID                   string `json:"id"`
	MetricID             string `json:"metric_id"`
	MetricVersion        string `json:"metric_version"`
	DimensionKey         string `json:"dimension_key"`
	CurrentMetricRowRef  string `json:"current_metric_row_ref,omitempty"`
	PreviousMetricRowRef string `json:"previous_metric_row_ref,omitempty"`
	CurrentValue         int    `json:"current_value"`
	PreviousValue        int    `json:"previous_value"`
	Delta                int    `json:"delta"`
	ComparisonBasis      string `json:"comparison_basis"`
	Comparable           bool   `json:"comparable"`
	NonComparableReason  string `json:"non_comparable_reason,omitempty"`
}

// MovementSummary totals comparable and non-comparable movement rows using the
// same reason vocabulary as individual rows.
type MovementSummary struct {
	ComparableCount     int            `json:"comparable_count"`
	NonComparableCount  int            `json:"non_comparable_count"`
	NonComparableReason map[string]int `json:"non_comparable_reason_counts"`
}

// RefusalRow preserves rejected inputs as explicit posture evidence rather than
// silently dropping them from the export.
type RefusalRow struct {
	ID              string `json:"id"`
	InputID         string `json:"input_id"`
	TimeWindow      string `json:"time_window,omitempty"`
	RefusalReason   string `json:"refusal_reason"`
	InputTrustState string `json:"input_trust_state"`
}

// OutputSafety lists sensitive classes this package checked as absent from
// rendered human-readable posture output.
type OutputSafety struct {
	VerifiedAbsentSensitiveClasses []string `json:"verified_absent_sensitive_classes"`
}
