package main

type adapterCapturePreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ExpectedEvidence map[string]string `json:"expected_evidence"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}
