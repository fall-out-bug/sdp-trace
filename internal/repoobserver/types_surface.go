package repoobserver

// Surface records one observable setup or proof dimension for the selected
// repository profile.
type Surface struct {
	SurfaceID      string `json:"surface_id"`
	InstallState   string `json:"install_state"`
	ProofState     string `json:"proof_state"`
	TrustScope     string `json:"trust_scope"`
	EvidenceSource string `json:"evidence_source"`
	ObservedRef    string `json:"observed_ref,omitempty"`
	ReasonCode     string `json:"reason_code"`
	NextAction     string `json:"next_action,omitempty"`
}
