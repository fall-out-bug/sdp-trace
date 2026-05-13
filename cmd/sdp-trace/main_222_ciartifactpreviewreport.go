package main

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
