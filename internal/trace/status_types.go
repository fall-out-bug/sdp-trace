package trace

// EventType defines the supported event class names for the first milestone.
type EventType string

// EventType constants.
const (
	EventRecorderAttached        EventType = "recorder_attached"
	EventRunStarted              EventType = "run_started"
	EventCommandStarted          EventType = "command_started"
	EventCommandFinished         EventType = "command_finished"
	EventRunClosed               EventType = "run_closed"
	EventPolicyOverrideRequested EventType = "policy_override_requested"
)

// EvidenceState maps direct missing-evidence states used by first-milestone verification.
type EvidenceState string

// EvidenceState constants.
const (
	EvidenceStatePresent      EvidenceState = "present"
	EvidenceStateMissing      EvidenceState = "missing"
	EvidenceStateNotAssessed  EvidenceState = "not_assessed"
	EvidenceStateCannotVerify EvidenceState = "cannot_verify"
)

// VerifierVerdict tracks the high-level result.
type VerifierVerdict string

// Verifier constants.
const (
	VerdictObserved     VerifierVerdict = "observed"
	VerdictFail         VerifierVerdict = "fail"
	VerdictCannotVerify VerifierVerdict = "cannot_verify"
	VerdictNotAssessed  VerifierVerdict = "not_assessed"
)

// TrustScope identifies evidence source context.
type TrustScope string

// TrustScope constants.
const (
	TrustScopeLocalObserved TrustScope = "local_observed"
)

// Completeness tracks output completeness.
type Completeness string

// Completeness constants.
const (
	CompletenessComplete         Completeness = "complete"
	CompletenessPartial          Completeness = "partial"
	CompletenessMissingTelemetry Completeness = "missing_telemetry"
	CompletenessUnknown          Completeness = "unknown"
)

// Replayability expresses reusability of recorded output.
type Replayability string

// Replayability constants.
const (
	ReplayabilityFull    Replayability = "full"
	ReplayabilityNone    Replayability = "none"
	ReplayabilityPartial Replayability = "partial"
)
