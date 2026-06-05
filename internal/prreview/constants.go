package prreview

import (
	"errors"
	"regexp"
)

const (
	SchemaVersionPacket     = "block30-pr-review-packet-v1"
	SchemaVersionProfile    = "block30-pr-review-profile-v1"
	SchemaVersionRunSet     = "block30-pr-review-runs-v1"
	SchemaVersionLedger     = "block30-pr-review-ledger-v1"
	SchemaVersionValidation = "block30-pr-review-validation-v1"

	StatePass         = "pass"
	StateFail         = "fail"
	StatePending      = "pending"
	StateNotAssessed  = "not_assessed"
	StateCannotVerify = "cannot_verify"

	RefKindDiff            = "diff"
	RefKindMetadata        = "metadata"
	RefKindSpec            = "spec"
	RefKindPlan            = "plan"
	RefKindTask            = "task"
	RefKindDoc             = "doc"
	RefKindSchema          = "schema"
	RefKindSourceExcerpt   = "source_excerpt"
	RefKindVerification    = "verification"
	RefKindPrompt          = "prompt"
	RefKindRawOutput       = "raw_output"
	RefKindSanitizedOutput = "sanitized_output"
	RefKindExternal        = "external"

	ContentUnifiedDiff = "unified_diff"
	ContentMarkdown    = "markdown"
	ContentJSON        = "json"
	ContentText        = "text"

	RedactionNone        = "none"
	RedactionRedacted    = "redacted"
	RedactionDigestOnly  = "digest_only"
	RedactionEncrypted   = "encrypted_ref"
	RedactionExternalRef = "external_ref"
	RedactionWithheld    = "withheld"
	RedactionNotAssessed = "not_assessed"

	PlaneCodeCorrectness = "code_correctness"
	PlaneTraceEvidence   = "trace_evidence_provenance"
	PlaneRequirements    = "requirements_vs_implementation"
	PlaneSecurity        = "security_forgery_overclaim"
	PlaneDXReplayability = "dx_replayability"
	PlanePrivacySafety   = "privacy_output_safety"

	RunnerPI             = "pi"
	RunnerOpenCode       = "opencode"
	RunnerManualExternal = "manual_external"

	StatusFindingsReported = "findings_reported"
	StatusNoFindings       = "no_findings"
	StatusNotAssessed      = "not_assessed"
	StatusFailed           = "failed"
	StatusTimedOut         = "timed_out"
	StatusEmptyOutput      = "empty_output"
	StatusOffTask          = "off_task"
	StatusParseFailed      = "parse_failed"
	StatusCannotVerify     = "cannot_verify"

	SeverityCritical      = "critical"
	SeverityMajor         = "major"
	SeverityMinor         = "minor"
	SeverityInformational = "informational"

	DispositionAcceptedFixed           = "accepted_fixed"
	DispositionAcceptedReviewBlocking  = "accepted_review_blocking"
	DispositionAcceptedNarrower        = "accepted_narrower"
	DispositionRejectedFalsePositive   = "rejected_false_positive"
	DispositionDeferredNotAssessed     = "deferred_not_assessed"
	DispositionUnresolvedReviewBlocker = "unresolved_review_blocker"

	CoverageSatisfied    = "coverage_satisfied"
	CoveragePartial      = "coverage_partial"
	CoverageUnresolved   = "coverage_unresolved"
	CoverageNotAssessed  = "not_assessed"
	CoverageCannotVerify = "cannot_verify"

	AuthorityReviewRecordOnly = "review_record_only"
	DecisionNotAuthorized     = "not_authorized_by_sdp_trace"
)

var (
	errPromptEvidenceCannotVerify = errors.New("prompt_evidence_cannot_verify")
	errPromptTemplateCannotVerify = errors.New("prompt_template_cannot_verify")
	repoIDPattern                 = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,62}[a-z0-9]$`)
	changeRefPattern              = regexp.MustCompile(`^(pr|mr|change)-[A-Za-z0-9._-]{1,64}$`)
	sha40Pattern                  = regexp.MustCompile(`^[0-9a-f]{40}$`)
)
