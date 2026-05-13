package forensic

const (
	SchemaVersion = "block18-forensic-retention-assessment-v1"

	ProfileForensicRetention = "forensic_retention"
	TrustScopeForensic       = "forensic_retention_observed"
	TrustScopeLocalObserved  = "local_observed"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"

	RetentionModeDigestOnly          = "digest_only"
	RetentionModeSanitizedExcerpt    = "sanitized_excerpt"
	RetentionModeEncryptedRawRef     = "encrypted_raw_ref"
	RetentionModeExternalArtifactRef = "external_artifact_ref"
	RetentionModeNotAssessed         = "not_assessed"

	RedactionActionApplyRule       = "apply_rule"
	RedactionActionWithhold        = "withhold"
	RedactionActionMarkUnavailable = "mark_unavailable"

	AuthorityVerified    = "verified"
	AuthoritySelfClaimed = "self_claimed"
	AuthorityNotAssessed = "not_assessed"

	AccessStateVerifiedAvailable = "verified_available"
	AccessStateRestricted        = "restricted"
	AccessStateUnavailable       = "unavailable"
	AccessStateRevoked           = "revoked"
	AccessStateNotAssessed       = "not_assessed"

	KeyCustodyNotApplicable = "not_applicable"
	KeyCustodyHolderKnown   = "holder_known"
	KeyCustodyEscrowed      = "escrowed"
	KeyCustodyDestroyed     = "destroyed"
	KeyCustodyCompromised   = "compromised"
	KeyCustodyUnknown       = "unknown"
	KeyCustodyNotAssessed   = "not_assessed"

	RetentionLifecycleActive      = "active"
	RetentionLifecycleExpired     = "expired"
	RetentionLifecycleRevoked     = "revoked"
	RetentionLifecycleDeleted     = "deleted"
	RetentionLifecycleRotated     = "rotated"
	RetentionLifecycleNotAssessed = "not_assessed"

	UnavailableReasonNotAssessed        = "not_assessed"
	UnavailableReasonMissingReference   = "missing_reference"
	UnavailableReasonAccessDenied       = "access_denied"
	UnavailableReasonExpired            = "expired"
	UnavailableReasonKeyUnavailable     = "key_unavailable"
	UnavailableReasonStoreUnreachable   = "store_unreachable"
	UnavailableReasonDigestUnverifiable = "digest_unverifiable"
)

var forensicConditionIDs = []string{
	"redaction_policy_bound",
	"redaction_prewrite_applied",
	"redaction_unresolved_visible",
	"retention_mode_declared",
	"critical_evidence_reconstructable",
	"raw_reference_bound",
	"forensic_profile_not_overclaimed",
	"profile_selection_accountable",
}
