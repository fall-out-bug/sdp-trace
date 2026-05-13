package harnessobs

const (
	ProfileSchemaVersion    = "harness-observation-profile-v1"
	EventSchemaVersion      = "harness-event-v1"
	RunSchemaVersion        = "harness-observation-run-v1"
	ValidationSchemaVersion = "harness-observation-validation-v1"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"

	ContentRedacted      = "redacted"
	ContentDigestOnly    = "digest_only"
	ContentRetainedSafe  = "retained_safe"
	ContentNotApplicable = "not_applicable"

	DefaultMaxLineBytes = 1024 * 1024
	DefaultMaxEvents    = 100000

	SessionProfileSchemaVersion = "harness-session-profile-v1"
	SessionRunSchemaVersion     = "harness-session-run-v1"
	OpenCodeJSONLRawFormat      = "opencode-jsonl-v1"

	safeTokenRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.:-"
)
