package posture

import (
	"regexp"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

const (
	// SchemaVersion identifies the only posture export contract this package
	// can build or validate.
	SchemaVersion               = "block21-cross-repo-posture-export-v1"
	ProfileID                   = "cross-repo-evidence-posture-v1"
	ProfileVer                  = "v1"
	SelectionSchemaVersion      = "block21-cross-repo-selection-v1"
	DigestManifestSchemaVersion = "block21-artifact-digest-manifest-v1"
	SignalManifestSchemaVersion = "block21-posture-signal-manifest-v1"

	GroupingRepoWindow          = "repo_window_v1"
	GroupingTeamServiceWindow   = "team_service_window_v1"
	GroupingHarnessChangeWindow = "harness_change_window_v1"
)

var metricCatalog = []metricDef{
	{"missing_telemetry_rows", "v1", "row_state"},
	{"not_assessed_rows", "v1", "row_state"},
	{"cannot_verify_rows", "v1", "row_state"},
	{"unsupported_observer_rows", "v1", "row_or_signal"},
	{"not_integrated_rows", "v1", "row_state"},
	{"retention_limited_rows", "v1", "row_state"},
	{"local_only_evidence_rows", "v1", "posture_signal"},
	{"ci_witnessed_evidence_rows", "v1", "posture_signal"},
	{"external_witnessed_evidence_rows", "v1", "posture_signal"},
	{"issue_observed_rows", "v1", "row_state"},
	{"override_rows", "v1", "posture_signal"},
	{"late_attach_rows", "v1", "posture_signal"},
	{"contract_change_rows", "v1", "posture_signal"},
}

var rowStateMetrics = map[string]string{
	"missing_telemetry_rows": query.RowStateMissingTelemetry,
	"not_assessed_rows":      query.RowStateNotAssessed,
	"cannot_verify_rows":     query.RowStateCannotVerify,
	"not_integrated_rows":    query.RowStateNotIntegrated,
	"retention_limited_rows": query.RowStateRetentionLimited,
	"issue_observed_rows":    query.RowStateIssueObserved,
}

var signalMetricPredicates = map[string]func(PostureSignal) bool{
	"local_only_evidence_rows":         func(signal PostureSignal) bool { return signal.WitnessScope == "local_only" },
	"ci_witnessed_evidence_rows":       func(signal PostureSignal) bool { return signal.WitnessScope == "ci_witnessed" },
	"external_witnessed_evidence_rows": func(signal PostureSignal) bool { return signal.WitnessScope == "external_witnessed" },
	"override_rows":                    func(signal PostureSignal) bool { return signal.OverrideMarker == "override_present" },
	"late_attach_rows":                 func(signal PostureSignal) bool { return signal.LateAttachMarker == "late_attach_observed" },
	"contract_change_rows":             func(signal PostureSignal) bool { return signal.ContractChangeMarker == "contract_change_observed" },
}

var safeLabelPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)
