package authority

type Package struct {
	SchemaVersion      string                 `json:"schema_version"`
	SelectedPolicyID   string                 `json:"selected_policy_id"`
	Actors             []ActorDeclaration     `json:"actors"`
	AuthorityEnvelopes []AuthorityEnvelope    `json:"authority_envelopes"`
	ObservedActions    []ObservedAction       `json:"observed_actions"`
	EvidenceBindings   []EvidenceBindingInput `json:"evidence_bindings,omitempty"`
	EvidenceResolution EvidenceResolution     `json:"evidence_resolution,omitempty"`
}

type AuthorityEnvelope struct {
	SchemaVersion        string                `json:"schema_version"`
	TaskID               string                `json:"task_id"`
	PolicyID             string                `json:"policy_id"`
	AuthorityScope       string                `json:"authority_scope"`
	ActorRef             string                `json:"actor_ref"`
	AllowedEvents        []string              `json:"allowed_events"`
	DeniedEvents         []string              `json:"denied_events"`
	TargetRules          []TargetRule          `json:"target_rules"`
	ApprovalRequirements []ApprovalRequirement `json:"approval_requirements,omitempty"`
	EffectiveFromEventID string                `json:"effective_from_event_id,omitempty"`
	SupersedesPolicyID   string                `json:"supersedes_policy_id,omitempty"`
}

type TargetRule struct {
	RuleID        string   `json:"rule_id"`
	TargetPattern string   `json:"target_pattern"`
	AllowedEvents []string `json:"allowed_events"`
	DeniedEvents  []string `json:"denied_events"`
}

type ApprovalRequirement struct {
	RequirementID       string `json:"requirement_id"`
	EventType           string `json:"event_type,omitempty"`
	TargetRuleRef       string `json:"target_rule_ref,omitempty"`
	ApprovalEvidenceRef string `json:"approval_evidence_ref"`
}

type ActorDeclaration struct {
	ActorID                string `json:"actor_id"`
	ActorType              string `json:"actor_type"`
	DeclaredRole           string `json:"declared_role"`
	Harness                string `json:"harness,omitempty"`
	Model                  string `json:"model,omitempty"`
	ModelAttributionSource string `json:"model_attribution_source,omitempty"`
	OperationID            string `json:"operation_id,omitempty"`
}

type ObservedAction struct {
	EventID               string   `json:"event_id"`
	TaskID                string   `json:"task_id,omitempty"`
	EventType             string   `json:"event_type"`
	Target                string   `json:"target,omitempty"`
	SourceType            string   `json:"source_type"`
	EvidenceRefs          []string `json:"evidence_refs"`
	ActorID               string   `json:"actor_id,omitempty"`
	OperationID           string   `json:"operation_id,omitempty"`
	ObservedAt            string   `json:"observed_at"`
	ObservationConfidence string   `json:"observation_confidence"`
}

type EvidenceBindingInput struct {
	BindingID     string   `json:"binding_id"`
	LeftEventID   string   `json:"left_event_id"`
	RightEventID  string   `json:"right_event_id"`
	BindingType   string   `json:"binding_type"`
	BindingState  string   `json:"binding_state"`
	MatchedFields []string `json:"matched_fields"`
	EvidenceRef   string   `json:"evidence_ref"`
}

type EvidenceResolution struct {
	ResolvedExternalRefs []string `json:"resolved_external_refs,omitempty"`
	InaccessibleRefs     []string `json:"inaccessible_refs,omitempty"`
	MalformedRefs        []string `json:"malformed_refs,omitempty"`
	StaleRefs            []string `json:"stale_refs,omitempty"`
}

type Result struct {
	SchemaVersion            string                `json:"schema_version"`
	SelectedProfile          string                `json:"selected_profile"`
	SelectedPolicyID         string                `json:"selected_policy_id"`
	AuthorityEvaluationState string                `json:"authority_evaluation_state"`
	Evaluations              []AuthorityEvaluation `json:"evaluations"`
	BindingEvaluations       []EvidenceBinding     `json:"binding_evaluations"`
	SourceCoverage           []string              `json:"source_coverage"`
	Reasons                  []string              `json:"reasons"`
	NextActions              []string              `json:"next_actions"`
}

type AuthorityEvaluation struct {
	EvaluationID      string   `json:"evaluation_id"`
	EventID           string   `json:"event_id"`
	PolicyID          string   `json:"policy_id"`
	State             string   `json:"state"`
	ReasonCode        string   `json:"reason_code"`
	MatchedRuleRef    string   `json:"matched_rule_ref,omitempty"`
	ActorAttribution  string   `json:"actor_attribution"`
	ToolAttribution   string   `json:"tool_attribution"`
	ModelAttribution  string   `json:"model_attribution"`
	SourceCoverage    []string `json:"source_coverage"`
	EvidenceRefs      []string `json:"evidence_refs"`
	ActorID           string   `json:"actor_id,omitempty"`
	OperationID       string   `json:"operation_id,omitempty"`
	MissingAttributes []string `json:"missing_attributes,omitempty"`
}

type EvidenceBinding struct {
	BindingID     string   `json:"binding_id"`
	LeftEventID   string   `json:"left_event_id"`
	RightEventID  string   `json:"right_event_id"`
	BindingType   string   `json:"binding_type"`
	BindingState  string   `json:"binding_state"`
	MatchedFields []string `json:"matched_fields"`
	EvidenceRef   string   `json:"evidence_ref"`
	ReasonCode    string   `json:"reason_code"`
}
