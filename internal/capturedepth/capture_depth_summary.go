package capturedepth

import "github.com/fall_out_bug/sdp-trace/internal/adaptercapture"

// CaptureDepthSummary is the query-only adapter-capture view. It intentionally
// names top-level assessment as not emitted so callers do not treat query JSON
// as a gate verdict.
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

func newCaptureDepthSummary(run adaptercapture.RunEvidence, result adaptercapture.AssessmentResult) CaptureDepthSummary {
	// Keep the top-level field explicit so downstream JSON consumers cannot
	// mistake a query response for a verifier verdict.
	sum := CaptureDepthSummary{
		Query:                  QueryName,
		SelectedProfile:        adaptercapture.ProfileAdapterCapture,
		TopLevelAssessment:     "not_emitted_for_query",
		TaskSupersessionCount:  run.TaskSupersessionCount,
		UnverifiedTaskExpanded: run.UnverifiedTaskExpanded,
	}
	fillSummary(&sum, run, result)
	return sum
}

func fillSummary(s *CaptureDepthSummary, run adaptercapture.RunEvidence, result adaptercapture.AssessmentResult) {
	// Full condition records remain in the response; compact ID slices are only
	// navigation aids for humans and downstream forensic query packs.
	s.MissingAdapterEvents = missingAdapterEvents(run)
	s.UnsupportedObservers = append([]string{}, run.UnsupportedEventTypes...)
	s.UnverifiedClaims = unverifiedClaims(result.AdapterCaptureConditions)
	s.Conditions = result.AdapterCaptureConditions
	s.Reasons = result.Reasons
}
