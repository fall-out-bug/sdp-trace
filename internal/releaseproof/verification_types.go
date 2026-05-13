package releaseproof

// Verification is the durable local release-proof output. The flat field set
// intentionally preserves keyed literal compatibility for tests and callers
// while grouping assessed local proof separately from explicit not_assessed
// external trust fields.
type Verification struct {
	ID                         string            `json:"id"`
	SchemaVersion              string            `json:"schema_version"`
	ArtifactRole               string            `json:"artifact_role"`
	TrustScope                 string            `json:"trust_scope"`
	ReleaseVerificationState   string            `json:"release_verification_state"`
	ManifestRef                string            `json:"manifest_ref"`
	ManifestDigest             string            `json:"manifest_digest"`
	ManifestDigestStatus       string            `json:"manifest_digest_status"`
	ArtifactDigestStatus       string            `json:"artifact_digest_status"`
	SignatureProfile           string            `json:"signature_profile"`
	SignatureStatus            string            `json:"signature_status"`
	IdentityPolicyRef          string            `json:"identity_policy_ref"`
	IdentityPolicyStatus       string            `json:"identity_policy_status"`
	SourceCommit               string            `json:"source_commit"`
	SourceCommitStatus         string            `json:"source_commit_status"`
	SourceCommitArtifactStatus string            `json:"source_commit_artifact_status"`
	SourceCommitArtifactCounts ArtifactCounts    `json:"source_commit_artifact_counts"`
	ExternalTrustProfile       string            `json:"external_trust_profile"`
	ExternalAttestationRef     *string           `json:"external_attestation_ref"`
	TransparencyEvidenceRef    *string           `json:"transparency_evidence_ref"`
	SourceURIStatus            string            `json:"source_uri_status"`
	ProtectedRefStatus         string            `json:"protected_ref_status"`
	WorkflowIdentityStatus     string            `json:"workflow_identity_status"`
	ApprovalStatus             string            `json:"approval_status"`
	ProductionReleaseVerified  ProofStateBoolean `json:"production_release_verified"`
	TransparencyLogStatus      string            `json:"transparency_log_status"`
	FreshnessStatus            string            `json:"freshness_status"`
	VerifiedAt                 string            `json:"verified_at"`
	TrustedContractRelease     bool              `json:"trusted_contract_release"`
	PrivateEquivalentProfile   string            `json:"private_equivalent_profile,omitempty"`
	ProvenanceRefs             []string          `json:"provenance_refs,omitempty"`
	Accountability             Accountability    `json:"accountability"`
	SourceCommitReason         string            `json:"source_commit_reason,omitempty"`
	ExternalTrustReason        string            `json:"external_trust_reason,omitempty"`
	ArtifactIssues             []ArtifactIssue   `json:"artifact_issues,omitempty"`
}
