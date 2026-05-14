package main

type forensicPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	PolicyEffects   map[string]string `json:"policy_effects"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}
