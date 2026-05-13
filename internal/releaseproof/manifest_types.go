package releaseproof

// Manifest is the source-bound release proof input. It names the source commit,
// artifact obligations, signing profile, identity policy, and accountability
// metadata that the local evaluator can render without claiming external trust.
type Manifest struct {
	ID                       string             `json:"id"`
	SigningProfile           string             `json:"signing_profile"`
	TrustedIdentityPolicyRef string             `json:"trusted_identity_policy_ref"`
	SourceCommit             string             `json:"source_commit"`
	Artifacts                []ManifestArtifact `json:"artifacts"`
	Accountability           Accountability     `json:"accountability"`
}

// ManifestArtifact is a manifest obligation checked against the source commit.
type ManifestArtifact struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	SHA256 string `json:"sha256"`
}

// Accountability records the human ownership fields copied into the proof.
type Accountability struct {
	DRI                 Actor  `json:"dri"`
	Approver            Actor  `json:"approver"`
	Escalation          Actor  `json:"escalation"`
	AuthorityScope      string `json:"authority_scope"`
	AccountabilityClaim string `json:"accountability_claim"`
	ApprovalRef         string `json:"approval_ref"`
	RiskOwner           Actor  `json:"risk_owner"`
	LineOfDefense       string `json:"line_of_defense"`
}

// Actor is an identity reference plus its closed actor classification.
type Actor struct {
	IdentityRef string `json:"identity_ref"`
	ActorType   string `json:"actor_type"`
}
