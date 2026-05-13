package ciartifact

type Manifest struct {
	SchemaVersion    string              `json:"schema_version"`
	SelectedProfile  string              `json:"selected_profile"`
	AuthorityScope   string              `json:"authority_scope"`
	SelectedSource   SourceIdentity      `json:"selected_source"`
	SelectedRun      RunIdentity         `json:"selected_run"`
	RequiredFamilies []FamilyRequirement `json:"required_families"`
	ArtifactFamilies []FamilyInput       `json:"artifact_families"`
	ArtifactIndex    ArtifactIndexInput  `json:"artifact_index"`
	OutputSafety     OutputSafetyInput   `json:"output_safety"`
	SafetyRuleset    SafetyRuleset       `json:"safety_ruleset"`
}

type SourceIdentity struct {
	Repository string `json:"repository,omitempty"`
	Ref        string `json:"ref,omitempty"`
	CommitSHA  string `json:"commit_sha,omitempty"`
}

type RunIdentity struct {
	Provider   string `json:"provider,omitempty"`
	RunID      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	WorkflowID string `json:"workflow_id,omitempty"`
	JobID      string `json:"job_id,omitempty"`
}

type FamilyRequirement struct {
	Family                string `json:"family"`
	RequiredProducerScope string `json:"required_producer_scope"`
}

type FamilyInput struct {
	Family              string `json:"family"`
	ProducerScope       string `json:"producer_scope"`
	ArtifactAccessState string `json:"artifact_access_state"`
	BindingState        string `json:"binding_state"`
	ClaimSource         string `json:"claim_source,omitempty"`
}

type ArtifactIndexInput struct {
	State string `json:"state"`
}

type OutputSafetyInput struct {
	State         string   `json:"state"`
	UnsafeClasses []string `json:"unsafe_classes,omitempty"`
}

type SafetyRuleset struct {
	ID     string `json:"id"`
	SHA256 string `json:"sha256"`
}

type ObservationResult struct {
	SchemaVersion            string              `json:"schema_version"`
	SelectedProfile          string              `json:"selected_profile"`
	AuthorityScope           string              `json:"authority_scope"`
	ArtifactObservationState string              `json:"artifact_observation_state"`
	SelectedSource           SourceIdentity      `json:"selected_source"`
	SelectedRun              RunIdentity         `json:"selected_run"`
	ProducerScope            string              `json:"producer_scope"`
	ArtifactAccessState      string              `json:"artifact_access_state"`
	RequiredFamilies         []FamilyRequirement `json:"required_families"`
	ArtifactFamilies         []FamilyObservation `json:"artifact_families"`
	Bindings                 BindingSummary      `json:"bindings"`
	ArtifactIndex            ArtifactIndexResult `json:"artifact_index"`
	OutputSafety             OutputSafetyResult  `json:"output_safety"`
	SafetyRuleset            SafetyRuleset       `json:"safety_ruleset"`
	Reasons                  []string            `json:"reasons"`
	NextActions              []string            `json:"next_actions"`
}

type FamilyObservation struct {
	Family              string `json:"family"`
	Required            bool   `json:"required"`
	RequiredProducer    string `json:"required_producer_scope,omitempty"`
	ProducerScope       string `json:"producer_scope"`
	ArtifactAccessState string `json:"artifact_access_state"`
	BindingState        string `json:"binding_state"`
	FamilyState         string `json:"family_state"`
	ReasonCode          string `json:"reason_code"`
	Reason              string `json:"reason"`
	NextAction          string `json:"next_action,omitempty"`
}

type BindingSummary struct {
	SourceBindingState   string `json:"source_binding_state"`
	RunBindingState      string `json:"run_binding_state"`
	ProducerBindingState string `json:"producer_binding_state"`
}

type ArtifactIndexResult struct {
	State      string `json:"state"`
	Result     string `json:"result"`
	ReasonCode string `json:"reason_code"`
	Reason     string `json:"reason"`
}

type OutputSafetyResult struct {
	State         string   `json:"state"`
	UnsafeClasses []string `json:"unsafe_classes,omitempty"`
	ReasonCode    string   `json:"reason_code"`
	Reason        string   `json:"reason"`
}
