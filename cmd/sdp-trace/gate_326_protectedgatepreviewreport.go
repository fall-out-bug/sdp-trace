package main

type protectedGatePreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	TrustCap        string            `json:"trust_cap"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}
