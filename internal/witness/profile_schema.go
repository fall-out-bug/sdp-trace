package witness

// ProfileStates records the per-dimension verifier state that explains why a
// witness did or did not satisfy the requested trust scope.
type ProfileStates struct {
	IdentityState        string `json:"identity_state"`
	SignerAuthorityState string `json:"signer_authority_state"`
	FreshnessState       string `json:"freshness_state"`
	ArtifactBindingState string `json:"artifact_binding_state"`
	SourceBindingState   string `json:"source_binding_state"`
	RunBindingState      string `json:"run_binding_state"`
	PolicyBindingState   string `json:"policy_binding_state"`
	IndependenceState    string `json:"independence_state"`
	KeyCustodyState      string `json:"key_custody_state,omitempty"`
}

// OutputSafety records the safety scan state without copying raw secrets,
// tokens, logs, or customer payloads into the witness artifact.
type OutputSafety struct {
	State                 string   `json:"state"`
	VerifiedAbsentClasses []string `json:"verified_absent_classes"`
}
