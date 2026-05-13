package policy

// AdapterAuthorityEntry declares which adapter identity may emit which portable
// trace event types. It is policy data, not live adapter authentication.
type AdapterAuthorityEntry struct {
	AdapterID         string   `json:"adapter_id"`
	Provider          string   `json:"provider"`
	IdentityState     string   `json:"identity_state"`
	AllowedEventTypes []string `json:"allowed_event_types"`
	AllowedByPolicy   bool     `json:"allowed_by_policy"`
}

// SignerAuthorityEntry binds signing identities to profile scopes.
type SignerAuthorityEntry struct {
	SignerID             string   `json:"signer_id"`
	ProfileID            string   `json:"profile_id"`
	AllowedScopes        []string `json:"allowed_scopes"`
	IndependenceRequired bool     `json:"independence_required"`
	Environment          string   `json:"environment"`
}

// TrustBoundaryDefaults names labels used when no stronger external trust
// boundary is present.
type TrustBoundaryDefaults struct {
	DefaultWitnessIndependence string `json:"default_witness_independence"`
	LocalProfileLabel          string `json:"local_profile_label"`
}

// SigningProfileSpec describes a signing profile document referenced by
// authority policy fixtures.
type SigningProfileSpec struct {
	SchemaVersion string `json:"schema_version"`
	ProfileID     string `json:"profile_id"`
	Format        string `json:"format"`
	Description   string `json:"description"`
	VerifierHost  string `json:"verifier_host"`
}

// RedactionProfileSpec describes default retention behavior for evidence
// redaction profiles.
type RedactionProfileSpec struct {
	SchemaVersion string `json:"schema_version"`
	ProfileID     string `json:"profile_id"`
	Description   string `json:"description"`
	DefaultMode   string `json:"default_retention_mode"`
}
