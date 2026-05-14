package main

type gatePreviewReport struct {
	Command            string   `json:"command"`
	GateMode           string   `json:"gate_mode"`
	TrustCap           string   `json:"trust_cap"`
	RequiredRuns       []string `json:"required_runs"`
	RequiredEvidence   []string `json:"required_evidence"`
	WitnessInspectable bool     `json:"witness_inspectable"`
	WitnessMismatches  []string `json:"witness_mismatches,omitempty"`
	Claim              string   `json:"claim"`
}
