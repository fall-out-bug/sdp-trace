package witness

// Record is the portable witness artifact. Its fields are deliberately plain
// JSON values so other harnesses can replay trust decisions without importing a
// runtime-specific SDK.
type Record struct {
	SchemaVersion         string           `json:"schema_version,omitempty"`
	Kind                  string           `json:"kind"`
	ProfileID             string           `json:"profile_id,omitempty"`
	ProfileVersion        string           `json:"profile_version,omitempty"`
	ProviderKind          string           `json:"provider_kind,omitempty"`
	Status                string           `json:"status"`
	TrustScope            string           `json:"trust_scope"`
	RequestedTrustScope   string           `json:"requested_trust_scope,omitempty"`
	EstablishedTrustScope string           `json:"established_trust_scope,omitempty"`
	Reason                string           `json:"reason"`
	ReasonCodes           []string         `json:"reason_codes,omitempty"`
	GeneratedAt           string           `json:"generated_at"`
	MissingIdentityFields []string         `json:"missing_identity_fields,omitempty"`
	Source                SourceIdentity   `json:"source"`
	CI                    CIIdentity       `json:"ci"`
	OIDC                  *OIDCClaims      `json:"oidc,omitempty"`
	RunArtifacts          []ArtifactDigest `json:"run_artifacts"`
	ReportArtifacts       []ArtifactDigest `json:"report_artifacts"`
	ProfileStates         *ProfileStates   `json:"profile_states"`
	OutputSafety          *OutputSafety    `json:"output_safety"`
}
