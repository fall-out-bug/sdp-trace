package adaptercapture

type Input struct {
	Run RunEvidence
}

type RunEvidence struct {
	RunID                  string             `json:"run_id"`
	RunNonce               string             `json:"run_nonce"`
	SourceBaseline         string             `json:"source_baseline"`
	RunClosedSequence      int                `json:"run_closed_sequence"`
	RequiredEventTypes     []string           `json:"required_event_types"`
	UnsupportedEventTypes  []string           `json:"unsupported_event_types,omitempty"`
	AdapterBundle          *AdapterBundle     `json:"adapter_bundle,omitempty"`
	AdapterEvents          []AdapterEvent     `json:"adapter_events"`
	ProviderRefs           []ProviderRef      `json:"provider_refs,omitempty"`
	EventFamilySummaries   []EventFamilyState `json:"event_family_summaries,omitempty"`
	RedactionPolicyDigest  string             `json:"redaction_policy_digest"`
	GatewayIntegrated      bool               `json:"gateway_integrated,omitempty"`
	GatewayEvidenceBound   bool               `json:"gateway_evidence_bound,omitempty"`
	TaskDriftAssessed      bool               `json:"task_drift_assessed"`
	TaskSupersessionCount  int                `json:"task_supersession_count,omitempty"`
	UnverifiedTaskExpanded bool               `json:"unverified_task_expanded,omitempty"`
}

type AdapterBundle struct {
	BundleID           string `json:"bundle_id"`
	HeadDigest         string `json:"head_digest"`
	ReferencedSequence int    `json:"referenced_sequence"`
	EventCount         int    `json:"event_count"`
}

type AdapterEvent struct {
	EventID                    string   `json:"event_id"`
	EventType                  string   `json:"event_type"`
	EventTime                  string   `json:"event_time,omitempty"`
	ProducerIdentity           string   `json:"producer_identity"`
	AdapterIdentity            string   `json:"adapter_identity"`
	IdentityBinding            string   `json:"identity_binding"`
	ProvenanceScope            string   `json:"provenance_scope"`
	CaptureState               string   `json:"capture_state"`
	RunID                      string   `json:"run_id"`
	RunNonce                   string   `json:"run_nonce"`
	SourceBaseline             string   `json:"source_baseline,omitempty"`
	BindingMode                string   `json:"binding_mode"`
	Sequence                   int      `json:"sequence,omitempty"`
	PrevEventHash              string   `json:"prev_event_hash,omitempty"`
	EventHash                  string   `json:"event_hash,omitempty"`
	AdapterBundleID            string   `json:"adapter_bundle_id,omitempty"`
	AdapterBundleHeadDigest    string   `json:"adapter_bundle_head_digest,omitempty"`
	BundleSequence             int      `json:"bundle_sequence,omitempty"`
	CorrelationRef             string   `json:"correlation_ref,omitempty"`
	EventPayloadDigest         string   `json:"event_payload_digest"`
	RedactionPolicyDigest      string   `json:"redaction_policy_digest,omitempty"`
	RetentionMode              string   `json:"retention_mode,omitempty"`
	ToolFamily                 string   `json:"tool_family,omitempty"`
	AdapterLocalToolLabel      string   `json:"adapter_local_tool_label,omitempty"`
	TestProvenance             string   `json:"test_provenance,omitempty"`
	ExecutedEvidenceClaimed    bool     `json:"executed_evidence_claimed,omitempty"`
	ModelIdentityProvenance    string   `json:"model_identity_provenance,omitempty"`
	GatewayEvidenceRef         string   `json:"gateway_evidence_ref,omitempty"`
	ProviderRefs               []string `json:"provider_refs,omitempty"`
	ActorAttributionState      string   `json:"actor_attribution_state,omitempty"`
	SensitiveMetadataPersisted bool     `json:"sensitive_metadata_persisted,omitempty"`
	ReconstructableClaimed     bool     `json:"reconstructable_claimed,omitempty"`
	CapAnnotation              string   `json:"cap_annotation,omitempty"`
}

type ProviderRef struct {
	SourceRef         string `json:"source_ref,omitempty"`
	SourceCommit      string `json:"source_commit,omitempty"`
	SourceTreeDigest  string `json:"source_tree_digest,omitempty"`
	ChangeRef         string `json:"change_ref,omitempty"`
	ReviewRef         string `json:"review_ref,omitempty"`
	ReviewState       string `json:"review_state,omitempty"`
	Producer          string `json:"producer,omitempty"`
	ArtifactDigest    string `json:"artifact_digest,omitempty"`
	ObservedAt        string `json:"observed_at,omitempty"`
	NotAssessedReason string `json:"not_assessed_reason,omitempty"`
}

type EventFamilyState struct {
	EventFamily     string `json:"event_family"`
	State           string `json:"state"`
	RetentionMode   string `json:"retention_mode,omitempty"`
	CapAnnotation   string `json:"cap_annotation,omitempty"`
	Reconstructable bool   `json:"reconstructable,omitempty"`
}

type AssessmentResult struct {
	SchemaVersion            string      `json:"schema_version"`
	SelectedProfile          string      `json:"selected_profile"`
	AdapterCaptureAssessment string      `json:"adapter_capture_assessment"`
	TrustScope               string      `json:"trust_scope"`
	AdapterCaptureConditions []Condition `json:"adapter_capture_conditions"`
	Reasons                  []string    `json:"reasons"`
	NextActions              []string    `json:"next_actions"`
}

type Condition struct {
	ID                    string `json:"id"`
	State                 string `json:"state"`
	ReasonCode            string `json:"reason_code"`
	Reason                string `json:"reason"`
	NextAction            string `json:"next_action,omitempty"`
	CappedToRetentionMode string `json:"capped_to_retention_mode,omitempty"`
}
