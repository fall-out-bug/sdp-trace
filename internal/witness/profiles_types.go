package witness

type EnvelopeInput struct {
	SchemaVersion       string           `json:"schema_version"`
	ProfileID           string           `json:"profile_id"`
	ProfileVersion      string           `json:"profile_version"`
	ProviderKind        string           `json:"provider_kind"`
	RequestedTrustScope string           `json:"requested_trust_scope"`
	GeneratedAt         string           `json:"generated_at"`
	Source              SourceIdentity   `json:"source"`
	CI                  CIIdentity       `json:"ci"`
	RunArtifacts        []ArtifactDigest `json:"run_artifacts"`
	ReportArtifacts     []ArtifactDigest `json:"report_artifacts"`
	ProfileStates       ProfileStates    `json:"profile_states"`
}

type ProfileOptions struct {
	EnvelopePath               string
	CustomerPKIAuthorityPolicy string
	CustomerPKIPublicCert      string
	CustomerPKIPublicKey       string
	CustomerPKIPayloadDigest   string
	CustomerPKIFreshness       string
}

type CustomerPKIAuthorityPolicy struct {
	SchemaVersion      string `json:"schema_version"`
	ProfileID          string `json:"profile_id"`
	AllowedSignerID    string `json:"allowed_signer_id"`
	PublicKeySHA256    string `json:"public_key_sha256"`
	PolicyDigest       string `json:"policy_digest"`
	KeyCustodyState    string `json:"key_custody_state"`
	RevocationRequired bool   `json:"revocation_required"`
	RevocationState    string `json:"revocation_state"`
}

type CustomerPKIFreshnessEvidence struct {
	SchemaVersion string `json:"schema_version"`
	SignerID      string `json:"signer_id"`
	PayloadDigest string `json:"payload_digest"`
	RunID         string `json:"run_id"`
	PolicyDigest  string `json:"policy_digest"`
	IssuedAt      string `json:"issued_at"`
	ValidUntil    string `json:"valid_until"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
}
