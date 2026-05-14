package main

type authorityPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	StateModel      map[string]string `json:"state_model"`
	Safety          map[string]string `json:"safety"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}
