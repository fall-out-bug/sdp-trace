package witness

const (
	// ReasonMissingArtifact records absent artifact digest evidence.
	ReasonMissingArtifact = "witness_artifact_digest_missing"
	// ReasonArtifactMismatch records artifact digests that do not match the
	// expected packet or run evidence.
	ReasonArtifactMismatch = "witness_artifact_digest_mismatch"
	// ReasonSourceMissing records absent source binding evidence.
	ReasonSourceMissing = "witness_source_binding_missing"
	// ReasonSourceMismatch records a source binding conflict.
	ReasonSourceMismatch = "witness_source_mismatch"
	// ReasonRunMissing records absent run binding evidence.
	ReasonRunMissing = "witness_run_binding_missing"
	// ReasonRunMismatch records a run binding conflict.
	ReasonRunMismatch = "witness_run_mismatch"
	// ReasonPolicyMissing records absent policy binding evidence.
	ReasonPolicyMissing = "witness_policy_binding_missing"
	// ReasonPolicyMismatch records a policy binding conflict.
	ReasonPolicyMismatch = "witness_policy_mismatch"
)
