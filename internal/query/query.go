package query

import (
	"encoding/json"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"os"
	"path/filepath"
)

const (
	QueryMissingEvidence = "missing-evidence"
	QueryCaptureDepth    = "capture-depth"
)

type CaptureDepthSummary struct {
	Query                  string                     `json:"query"`
	SelectedProfile        string                     `json:"selected_profile"`
	TopLevelAssessment     string                     `json:"top_level_assessment"`
	TaskSupersessionCount  int                        `json:"task_supersession_count"`
	UnverifiedTaskExpanded bool                       `json:"unverified_task_expanded"`
	MissingAdapterEvents   []string                   `json:"missing_adapter_events,omitempty"`
	UnsupportedObservers   []string                   `json:"unsupported_observers,omitempty"`
	UnverifiedClaims       []string                   `json:"unverified_claims,omitempty"`
	Conditions             []adaptercapture.Condition `json:"conditions"`
	Reasons                []string                   `json:"reasons"`
}

// MissingEvidence returns the missing-evidence table for a run.
//
// If a verifier artifact already exists, that result is reused so query output is
// stable between invocations.
func MissingEvidence(runDir string) ([]byte, error) {
	artifactPath := filepath.Join(runDir, "verifier", "missing-evidence-table.json")
	if _, err := os.Stat(artifactPath); err == nil {

		return os.ReadFile(artifactPath)
	}

	_, table, _, err := verifier.VerifyRun(runDir)
	if err != nil {
		return nil, err
	}
	payload, err := json.MarshalIndent(table, "", "  ")
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// CaptureDepth returns read-only adapter capture-depth rows. It does not emit a
// top-level pass/fail policy decision; callers that need exit semantics should
// use assess --profile adapter-capture.
func CaptureDepth(runDir string) ([]byte, error) {
	var run adaptercapture.RunEvidence
	if err := readJSON(filepath.Join(runDir, "run.json"), &run); err != nil {

		return nil, err
	}

	result := adaptercapture.Evaluate(adaptercapture.Input{Run: run})
	summary := CaptureDepthSummary{
		Query:                  QueryCaptureDepth,
		SelectedProfile:        adaptercapture.ProfileAdapterCapture,
		TopLevelAssessment:     "not_emitted_for_query",
		TaskSupersessionCount:  run.TaskSupersessionCount,
		UnverifiedTaskExpanded: run.UnverifiedTaskExpanded,

		MissingAdapterEvents: missingAdapterEvents(run),
		UnsupportedObservers: append([]string{}, run.UnsupportedEventTypes...),
		UnverifiedClaims:     unverifiedClaims(result.AdapterCaptureConditions),

		Conditions: result.AdapterCaptureConditions,
		Reasons:    result.Reasons,
	}

	return json.MarshalIndent(summary, "", "  ")
}
func unverifiedClaims(conditions []adaptercapture.Condition) []string {
	out := []string{}
	for _, condition := range conditions {
		switch condition.State {
		case adaptercapture.StateCannotVerify, adaptercapture.StateNotAssessed, adaptercapture.StateMissingTelemetry, adaptercapture.StateNotIntegrated, adaptercapture.StateUnsupported, adaptercapture.StateRetentionLimited:

			out = append(out, condition.ID)
		}
	}
	return out
}

func missingAdapterEvents(run adaptercapture.RunEvidence) []string {
	seen := map[string]bool{}
	for _, event := range run.AdapterEvents {

		seen[event.EventType] = true
	}
	missing := []string{}
	for _, required := range run.RequiredEventTypes {
		if !seen[required] {

			missing = append(missing, required)
		}
	}
	return missing
}

func readJSON(path string, target any) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
