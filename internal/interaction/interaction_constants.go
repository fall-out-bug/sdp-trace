package interaction

import (
	"regexp"
)

const (
	SchemaVersion = "1.0.0"
	MaxBodyBytes  = 16 * 1024

	SourceObservedControlChannel  = "observed-control-channel"
	SourcePreclassifiedTranscript = "preclassified-transcript-import"
	SourceAgentReported           = "agent-reported"
	StateObserved                 = "observed"
	StateReferenced               = "referenced"
	StateUnreferenced             = "unreferenced"
	StateNotAssessed              = "not_assessed"
	StateCannotVerify             = "cannot_verify"
	CompletenessComplete          = "complete"
	CompletenessPartial           = "partial"
	CompletenessNotAssessed       = "not_assessed"
	CompletenessCannotVerify      = "cannot_verify"
	RetentionDigestOnly           = "digest_only"
	RetentionSanitizedExcerpt     = "sanitized_excerpt"
	RetentionEncryptedRawRef      = "encrypted_raw_ref"
	RetentionExternalArtifactRef  = "external_artifact_ref"
	RetentionNotAssessed          = "not_assessed"
	ChannelExclusivityNotAssessed = "not_assessed"
	DigestAlgorithmSHA256         = "sha256"
	DefaultRedactionPolicyRef     = "block29-safe-default-v1"
	DefaultRelaySourceID          = "interaction-relay-v1"
)

var (
	safeIDPattern      = regexp.MustCompile(`^[A-Za-z0-9_.:-]+$`)
	sha256Pattern      = regexp.MustCompile(`^[a-f0-9]{64}$`)
	privatePathPattern = regexp.MustCompile(`(^|\s)/(Users|home|private|var|tmp)/[^\s]+`)
	tokenPattern       = regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password|authorization)\s*[:=]\s*[^\s]+`)
	authURLPattern     = regexp.MustCompile(`https?://[^/\s:@]+:[^/\s@]+@`)
	contentRefPattern  = regexp.MustCompile(`^(evidence:[A-Za-z0-9_./:-]+|sdp://interaction/[A-Za-z0-9_.:-]+/[A-Za-z0-9_.:-]+|external:[A-Za-z0-9_.:-]+|recorder:[A-Za-z0-9_.:-]+/event:[0-9]+)$`)
	recorderRefPattern = regexp.MustCompile(`^recorder:[A-Za-z0-9_.:-]+(?:/event:[0-9]+)?$`)
	runRefPattern      = regexp.MustCompile(`^recorder:[A-Za-z0-9_.:-]+$`)
)

var frictionClasses = map[string]string{
	"task_assignment":       "none",
	"plan_approved":         "none",
	"clarification_request": "clarification",
	"clarification_answer":  "clarification",
	"plan_proposed":         "planning",
	"plan_rejected":         "correction",
	"corrective_feedback":   "correction",
	"boundary_violation":    "correction",
	"evidence_correction":   "evidence",
	"tool_or_model_drift":   "drift",
	"pause_requested":       "coordination",
	"resume_approved":       "coordination",
}
