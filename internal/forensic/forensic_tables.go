package forensic

var insufficientCriticalRetentionConditions = map[string]Condition{
	RetentionModeDigestOnly:  failWithCap("critical_evidence_reconstructable", "critical_evidence_digest_only", "critical evidence is digest-only and not reconstructable", RetentionModeDigestOnly, "Retain sanitized excerpts, encrypted raw references, or external artifact references for critical event families."),
	RetentionModeNotAssessed: failWithCap("critical_evidence_reconstructable", "critical_evidence_not_assessed", "critical evidence retention is not assessed", RetentionModeNotAssessed, "Capture critical evidence or keep forensic retention open."),
}

type rawReferenceValidationRule struct {
	invalid   func(*RawReference) bool
	condition Condition
}

var rawReferenceValidationRules = []rawReferenceValidationRule{
	{
		invalid: func(ref *RawReference) bool {

			return ref.Digest.Algorithm != "sha256" || len(ref.Digest.Value) != 64
		},
		condition: fail("raw_reference_bound", "weak_digest", "raw reference digest is weak, unknown, or malformed", "Use SHA-256 or stronger digest binding for raw references."),
	},
	{
		invalid: func(ref *RawReference) bool {

			return ref.ReferenceType != RetentionModeEncryptedRawRef && ref.ReferenceType != RetentionModeExternalArtifactRef
		},
		condition: fail("raw_reference_bound", "raw_reference_type_invalid", "raw reference type is not an accepted FR-054 raw reference mode", "Use encrypted_raw_ref or external_artifact_ref."),
	},
	{
		invalid: func(ref *RawReference) bool {
			return ref.ReferenceURI == ""
		},
		condition: cannotVerify("raw_reference_bound", "missing_reference", "raw reference URI is missing", "Provide a stable encrypted or external raw reference."),
	},
	{
		invalid: func(ref *RawReference) bool {
			return rawReferenceAccessUnverifiable(ref)
		},
		condition: cannotVerify("raw_reference_bound", "access_unverifiable", "raw reference access state is not verifiably available", "Record current access verification state and time."),
	},
	{
		invalid: func(ref *RawReference) bool {

			return ref.AccessStateLastVerified == ""
		},
		condition: cannotVerify("raw_reference_bound", "access_unverifiable", "raw reference access verification time is missing", "Record access_state_last_verified for the assessment."),
	},
	{
		invalid: func(ref *RawReference) bool {
			return encryptedKeyCustodyUnverifiable(ref)
		},
		condition: cannotVerify("raw_reference_bound", "key_custody_unverifiable", "encrypted raw reference key custody is not verifiable", "Record holder_known or escrowed key custody state."),
	},
	{
		invalid: func(ref *RawReference) bool {
			return retentionLifecycleUnverifiable(ref.RetentionLifecycle.State)
		},
		condition: cannotVerify("raw_reference_bound", "retention_lifecycle_unverifiable", "raw reference retention lifecycle is not active", "Record active retention lifecycle evidence."),
	},
}
var defaultCriticalEventTypes = []string{
	"command_started",
	"command_finished",
	"test_output_observed",
	"file_mutation_observed",
	"artifact_captured",
	"model_identity_observed",
	"harness_identity_observed",
	"requirement_superseded",
	"redaction_applied",
	"run_closed",
}

var validRetentionModes = map[string]bool{
	RetentionModeDigestOnly:          true,
	RetentionModeSanitizedExcerpt:    true,
	RetentionModeEncryptedRawRef:     true,
	RetentionModeExternalArtifactRef: true,
	RetentionModeNotAssessed:         true,
}

var validAllowedRetentionModeFixture = []string{
	RetentionModeDigestOnly,
	RetentionModeSanitizedExcerpt,
	RetentionModeEncryptedRawRef,
	RetentionModeExternalArtifactRef,
	RetentionModeNotAssessed,
}

var validRedactionActionFixture = []string{
	RedactionActionApplyRule,
	RedactionActionWithhold,
	RedactionActionMarkUnavailable,
}

var validForbiddenPersistenceClassFixture = []string{
	"credentials",
	"tokens",
	"raw_prompts",
	"raw_model_responses",
	"source_snippets",
	"stdout_stderr_bodies",
	"oidc_tokens",
	"adapter_secrets",
	"gateway_tokens",
	"checkpoint_key_material",
}
