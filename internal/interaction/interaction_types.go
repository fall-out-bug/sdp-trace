package interaction

import "time"

type Actor struct {
	ID        string `json:"id"`
	ActorType string `json:"actor_type"`
}

type Source struct {
	SourceType    string `json:"source_type"`
	SourceID      string `json:"source_id"`
	SourceVersion string `json:"source_version,omitempty"`
	SourceRef     string `json:"source_ref,omitempty"`
}

type Redaction struct {
	PolicyRef    string `json:"policy_ref"`
	FindingCount int    `json:"finding_count"`
}

type LLMRef struct {
	GatewayTraceID string `json:"gateway_trace_id,omitempty"`
	LLMRequestID   string `json:"llm_request_id,omitempty"`
	ProviderID     string `json:"provider_request_id,omitempty"`
	ModelFamily    string `json:"model_family,omitempty"`
	ModelVersion   string `json:"model_version,omitempty"`
	GatewayID      string `json:"llm_gateway_id,omitempty"`
	LinkageState   string `json:"linkage_state"`
	RetentionState string `json:"retention_state,omitempty"`
}

type Event struct {
	SchemaVersion          string    `json:"schema_version"`
	InteractionID          string    `json:"interaction_id"`
	TaskID                 string    `json:"task_id"`
	OperationID            string    `json:"operation_id,omitempty"`
	StageID                string    `json:"stage_id,omitempty"`
	SourceID               string    `json:"source_id"`
	SourceSequence         int       `json:"source_sequence"`
	EventType              string    `json:"event_type"`
	FrictionClass          string    `json:"friction_class"`
	Actor                  Actor     `json:"actor"`
	Target                 string    `json:"target,omitempty"`
	Source                 Source    `json:"source"`
	ContentRef             string    `json:"content_ref,omitempty"`
	ContentDigest          string    `json:"content_digest"`
	DigestAlgorithm        string    `json:"digest_algorithm"`
	Retention              string    `json:"retention"`
	State                  string    `json:"state"`
	ReferenceRefs          []string  `json:"reference_refs,omitempty"`
	LLMRefs                []LLMRef  `json:"llm_refs,omitempty"`
	ObservedBeforeDelivery bool      `json:"observed_before_delivery"`
	ChannelExclusivity     string    `json:"channel_exclusivity_state"`
	CompletenessState      string    `json:"completeness_state"`
	NotRetainedReason      string    `json:"not_retained_reason,omitempty"`
	Redaction              Redaction `json:"redaction"`
	ObservedAt             string    `json:"observed_at"`
	CreatedAt              string    `json:"created_at"`
}

type Trace struct {
	SchemaVersion     string   `json:"schema_version"`
	TraceID           string   `json:"trace_id"`
	TaskID            string   `json:"task_id"`
	SourceType        string   `json:"source_type"`
	CompletenessState string   `json:"completeness_state"`
	Events            []Event  `json:"events"`
	AssessmentState   string   `json:"assessment_state"`
	NotAssessed       []string `json:"not_assessed,omitempty"`
	CannotVerify      []string `json:"cannot_verify,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

type RelayOptions struct {
	TaskID      string
	ActorType   string
	ActorID     string
	Target      string
	EventType   string
	OperationID string
	StageID     string
	Out         string
	Command     []string
	Now         time.Time
}
type ImportOptions struct {
	TaskID      string
	Source      string
	SourceRef   string
	EventsJSONL string
	Out         string
	Now         time.Time
}

type Summary struct {
	SchemaVersion          string         `json:"schema_version"`
	TaskID                 string         `json:"task_id"`
	TraceID                string         `json:"trace_id,omitempty"`
	EnvelopeID             string         `json:"envelope_id,omitempty"`
	AssessmentState        string         `json:"assessment_state"`
	EventCount             int            `json:"event_count,omitempty"`
	FrictionCounts         map[string]int `json:"friction_counts,omitempty"`
	CorrectionAfterTask    int            `json:"correction_after_assignment_count,omitempty"`
	PlanRejectionCount     int            `json:"plan_rejection_count,omitempty"`
	ClarificationTurnCount int            `json:"clarification_turn_count,omitempty"`
	UnreferencedEventCount int            `json:"unreferenced_event_count,omitempty"`
	RunRefCount            int            `json:"run_ref_count,omitempty"`
	SourceRefCount         int            `json:"source_ref_count,omitempty"`
	TaskRefCount           int            `json:"task_ref_count,omitempty"`
	PromiseRefCount        int            `json:"promise_ref_count,omitempty"`
	InteractionRefCount    int            `json:"interaction_ref_count,omitempty"`
	OperationRefCount      int            `json:"operation_ref_count,omitempty"`
	LLMRefCount            int            `json:"llm_ref_count,omitempty"`
	ToolRefCount           int            `json:"tool_ref_count,omitempty"`
	MutationRefCount       int            `json:"mutation_ref_count,omitempty"`
	StageRefCount          int            `json:"stage_ref_count,omitempty"`
	FrictionRefCount       int            `json:"friction_ref_count,omitempty"`
	NotAssessed            []string       `json:"not_assessed,omitempty"`
	CannotVerify           []string       `json:"cannot_verify,omitempty"`
}

type Envelope struct {
	SchemaVersion   string   `json:"schema_version"`
	TaskID          string   `json:"task_id"`
	EnvelopeID      string   `json:"envelope_id"`
	RunRefs         []string `json:"run_refs"`
	SourceRefs      []string `json:"source_refs"`
	TaskRefs        []string `json:"task_refs"`
	PromiseRefs     []string `json:"promise_refs"`
	InteractionRefs []string `json:"interaction_refs"`
	OperationRefs   []string `json:"operation_refs"`
	LLMRefs         []string `json:"llm_refs"`
	ToolRefs        []string `json:"tool_refs"`
	MutationRefs    []string `json:"mutation_refs"`
	StageRefs       []string `json:"stage_refs"`
	FrictionRefs    []string `json:"friction_refs"`
	AssessmentState string   `json:"assessment_state"`
	NotAssessed     []string `json:"not_assessed,omitempty"`
	CannotVerify    []string `json:"cannot_verify,omitempty"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type traceAssessment struct {
	completeness string
	assessment   string
	notAssessed  []string
	cannotVerify []string
}

type traceCounts struct {
	events         int
	friction       map[string]int
	corrections    int
	planRejections int
	clarifications int
	unreferenced   int
}

type traceSummaryCounter struct {
	frictionCounts             map[string]int
	correctionsAfterAssignment int
	planRejectionCount         int
	clarificationCount         int
	unreferencedCount          int
	assignmentObserved         bool
}

type envelopeRefCounts struct {
	run, source, task, promise, interaction int
	operation, llm, tool, mutation          int
	stage, friction                         int
}
