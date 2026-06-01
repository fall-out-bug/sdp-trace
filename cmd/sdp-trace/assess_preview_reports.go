package main

type managedPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

type adapterCapturePreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ExpectedEvidence map[string]string `json:"expected_evidence"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

type forensicPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	PolicyEffects   map[string]string `json:"policy_effects"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

type ciArtifactPreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ObservedFamilies []string          `json:"observed_families"`
	StateModel       map[string]string `json:"state_model"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

type authorityPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	StateModel      map[string]string `json:"state_model"`
	Safety          map[string]string `json:"safety"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

var ciArtifactPreviewObservedFamilies = []string{
	"run", "report", "witness", "provenance", "evidence",
	"trace", "artifact_index", "redaction_scan", "review", "change_ci",
}

var ciArtifactPreviewStateModel = map[string]string{
	"top_level": "pass,fail,cannot_verify,not_assessed",
	"producer":  "ci_uploaded,checked_in,local_generated,agent_reported,harness_observed,external_artifact_ref,not_assessed",
	"access":    "present,absent,partial,expired,inaccessible,malformed,unsafe,not_assessed,cannot_verify",
}

var ciArtifactPreviewSafety = map[string]string{
	"raw_artifact_content": "not_rendered",
	"reason_payloads":      "safe_reason_codes_only",
	"network_fetch":        "not_performed",
}
