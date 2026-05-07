package adaptercapture

import "sort"

const (
	SchemaVersion = "block19-adapter-capture-assessment-v1"

	ProfileAdapterCapture = "adapter_capture"
	TrustScopeAdapter     = "adapter_capture_observed"
	TrustScopeLocal       = "local_observed"

	StatePass             = "pass"
	StateFail             = "fail"
	StateCannotVerify     = "cannot_verify"
	StateNotAssessed      = "not_assessed"
	StateMissingTelemetry = "missing_telemetry"
	StateNotIntegrated    = "not_integrated"
	StateUnsupported      = "unsupported"
	StateRetentionLimited = "retention_limited"

	BindingSameChain     = "same_chain"
	BindingAdapterBundle = "adapter_bundle"

	IdentitySelfAsserted = "self_asserted"
	IdentityBound        = "bound"

	RetentionDigestOnly          = "digest_only"
	RetentionSanitizedExcerpt    = "sanitized_excerpt"
	RetentionEncryptedRawRef     = "encrypted_raw_ref"
	RetentionExternalArtifactRef = "external_artifact_ref"
	RetentionNotAssessed         = "not_assessed"
)

var adapterConditionIDs = []string{
	"adapter_event_contract_valid",
	"adapter_identity_visible",
	"run_binding_established",
	"task_drift_visible",
	"tool_call_depth_visible",
	"file_mutation_correlated",
	"model_identity_not_overclaimed",
	"test_provenance_not_overclaimed",
	"provider_refs_portable",
	"redaction_metadata_consistent",
	"capture_depth_not_overclaimed",
}

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

func Evaluate(input Input) AssessmentResult {
	conditions := []Condition{
		contractCondition(input.Run),
		identityCondition(input.Run),
		runBindingCondition(input.Run),
		taskDriftCondition(input.Run),
		toolDepthCondition(input.Run),
		fileMutationCondition(input.Run),
		modelIdentityCondition(input.Run),
		testProvenanceCondition(input.Run),
		providerRefsCondition(input.Run),
		redactionMetadataCondition(input.Run),
		overclaimCondition(input.Run),
	}
	result := AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileAdapterCapture,
		AdapterCaptureAssessment: topLevel(conditions),
		TrustScope:               TrustScopeAdapter,
		AdapterCaptureConditions: conditions,
	}
	if result.AdapterCaptureAssessment != StatePass {
		result.TrustScope = TrustScopeLocal
	}
	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}

func contractCondition(run RunEvidence) Condition {
	if len(run.AdapterEvents) == 0 {
		return cannotVerify("adapter_event_contract_valid", "adapter_events_missing", "adapter event evidence is missing", "Supply same-chain adapter events or an adapter bundle.")
	}
	seen := map[string]bool{}
	for _, event := range run.AdapterEvents {
		if event.EventID == "" || event.EventType == "" || event.ProducerIdentity == "" || event.AdapterIdentity == "" || event.EventPayloadDigest == "" {
			return fail("adapter_event_contract_valid", "adapter_event_malformed", "adapter event is missing required contract fields", "Emit schema-valid adapter events with producer, adapter, type, and digest fields.")
		}
		key := event.EventType + "\x00" + event.CorrelationRef
		if seen[key] && event.CorrelationRef != "" {
			return cannotVerify("adapter_event_contract_valid", "conflicting_adapter_events", "multiple adapter events share a correlation key", "Deduplicate adapter events or make conflicts explicit.")
		}
		seen[key] = true
	}
	return pass("adapter_event_contract_valid", "adapter_event_contract_valid", "adapter events match required contract fields")
}

func identityCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if event.ProducerIdentity == "" || event.AdapterIdentity == "" {
			return cannotVerify("adapter_identity_visible", "adapter_identity_missing", "adapter or producer identity is missing", "Record adapter and producer identity.")
		}
		if event.IdentityBinding != IdentitySelfAsserted && event.IdentityBinding != IdentityBound {
			return cannotVerify("adapter_identity_visible", "adapter_identity_unclassified", "adapter identity binding state is not classified", "Classify adapter identity as self_asserted or bound.")
		}
	}
	return pass("adapter_identity_visible", "adapter_identity_visible", "adapter and producer identity are visible with binding classification")
}

func runBindingCondition(run RunEvidence) Condition {
	if run.RunID == "" || run.RunNonce == "" {
		return cannotVerify("run_binding_established", "run_binding_missing", "run id or nonce is missing", "Record run id and run nonce before assessing adapter capture.")
	}
	for _, event := range run.AdapterEvents {
		if event.RunID != run.RunID || event.RunNonce != run.RunNonce {
			return fail("run_binding_established", "run_binding_mismatch", "adapter event contradicts run id or nonce", "Use adapter events bound to the selected run.")
		}
		if run.RunClosedSequence > 0 && event.BindingMode == BindingSameChain && event.Sequence > run.RunClosedSequence {
			return cannotVerify("run_binding_established", "late_adapter_event", "adapter event appears after run closure", "Do not use late adapter events to satisfy capture-depth assessment.")
		}
		switch event.BindingMode {
		case BindingSameChain:
			if event.EventHash == "" || event.PrevEventHash == "" {
				return cannotVerify("run_binding_established", "same_chain_digest_missing", "same-chain adapter event lacks hash linkage", "Record prev_event_hash and event_hash.")
			}
		case BindingAdapterBundle:
			if run.AdapterBundle == nil || event.AdapterBundleHeadDigest == "" || event.AdapterBundleHeadDigest != run.AdapterBundle.HeadDigest || event.AdapterBundleID != run.AdapterBundle.BundleID {
				return cannotVerify("run_binding_established", "adapter_bundle_unbound", "adapter bundle is not bound to the selected run", "Bind the adapter bundle head digest into the run artifact.")
			}
			if run.RunClosedSequence > 0 && run.AdapterBundle.ReferencedSequence > run.RunClosedSequence {
				return cannotVerify("run_binding_established", "late_adapter_bundle", "adapter bundle was first referenced after run closure", "Reference adapter bundles before run closure.")
			}
		default:
			return cannotVerify("run_binding_established", "binding_mode_missing", "adapter event binding mode is missing or unsupported", "Use same_chain or adapter_bundle binding.")
		}
	}
	return pass("run_binding_established", "run_binding_established", "adapter events are bound to run id, nonce, and chain or bundle context")
}

func taskDriftCondition(run RunEvidence) Condition {
	if !run.TaskDriftAssessed {
		return Condition{ID: "task_drift_visible", State: StateNotAssessed, ReasonCode: "task_drift_not_assessed", Reason: "task drift assessment was not selected", NextAction: "Assess task locks and task_superseded events."}
	}
	for _, event := range run.AdapterEvents {
		if event.EventType == "task_superseded" && event.ActorAttributionState == "" {
			return cannotVerify("task_drift_visible", "task_supersession_actor_missing", "task supersession lacks actor attribution state", "Record actor attribution state and task digest refs.")
		}
	}
	if run.TaskSupersessionCount == 0 {
		return pass("task_drift_visible", "no_supersessions_observed", "task drift was assessed and no supersessions were observed")
	}
	return pass("task_drift_visible", "task_supersessions_visible", "task supersessions include visible attribution and digest evidence")
}

func toolDepthCondition(run RunEvidence) Condition {
	if hasRequired(run, "tool_call") && !hasEvent(run.AdapterEvents, "tool_call") {
		if unsupported(run, "tool_call") {
			return Condition{ID: "tool_call_depth_visible", State: StateUnsupported, ReasonCode: "tool_event_unsupported", Reason: "adapter declares no tool-call capability", NextAction: "Use an adapter with tool-call capture capability."}
		}
		return Condition{ID: "tool_call_depth_visible", State: StateMissingTelemetry, ReasonCode: "tool_event_missing", Reason: "required tool-call adapter event is missing", NextAction: "Capture tool_call adapter events or mark the observer unsupported."}
	}
	return pass("tool_call_depth_visible", "tool_call_depth_visible", "required tool-call families are captured or not required")
}

func fileMutationCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if event.EventType == "file_mutation" {
			if event.SourceBaseline == "" || event.RunID == "" {
				return cannotVerify("file_mutation_correlated", "file_mutation_source_missing", "file mutation is not correlated with source baseline and run id", "Record source baseline and run id correlation for file mutation events.")
			}
		}
	}
	return pass("file_mutation_correlated", "file_mutation_correlated", "file mutation evidence is correlated with source baseline and run id")
}

func modelIdentityCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if event.EventType == "model_call_observed" {
			if event.ModelIdentityProvenance == "gateway_observed" && (!run.GatewayIntegrated || !run.GatewayEvidenceBound || event.IdentityBinding != IdentityBound) {
				return fail("model_identity_not_overclaimed", "gateway_identity_overclaimed", "model identity is claimed as gateway-observed without bound gateway evidence", "Keep model identity harness_observed or bind gateway evidence.")
			}
		}
	}
	if !run.GatewayIntegrated {
		return Condition{ID: "model_identity_not_overclaimed", State: StateNotIntegrated, ReasonCode: "gateway_not_integrated", Reason: "gateway provenance is not integrated", NextAction: "Integrate gateway evidence before claiming gateway-observed model identity."}
	}
	return pass("model_identity_not_overclaimed", "model_identity_not_overclaimed", "model identity provenance stays within available gateway evidence")
}

func testProvenanceCondition(run RunEvidence) Condition {
	observed := false
	for _, event := range run.AdapterEvents {
		if event.EventType != "test_observed" {
			continue
		}
		observed = true
		switch event.TestProvenance {
		case "ci_executed", "wrapper_executed":
			return pass("test_provenance_not_overclaimed", "test_provenance_executed", "test evidence is bound to CI or wrapper execution")
		case "agent_reported":
			if event.ExecutedEvidenceClaimed {
				return fail("test_provenance_not_overclaimed", "agent_reported_test_not_executed", "agent-reported tests are claimed as executed evidence", "Bind test evidence to CI or wrapper execution.")
			}
			return cannotVerify("test_provenance_not_overclaimed", "test_execution_unverified", "agent-reported test evidence is visible but non-executed", "Capture CI or wrapper-executed test evidence.")
		case "harness_observed":
			if event.ExecutedEvidenceClaimed {
				return fail("test_provenance_not_overclaimed", "harness_observed_test_not_executed", "harness-observed test intent is claimed as executed evidence", "Bind test evidence to CI or wrapper execution.")
			}
			return cannotVerify("test_provenance_not_overclaimed", "test_execution_unverified", "harness-observed test evidence is correlation-only", "Capture CI or wrapper-executed test evidence.")
		default:
			return cannotVerify("test_provenance_not_overclaimed", "test_provenance_missing", "test provenance is missing or unverifiable", "Record ci_executed or wrapper_executed test provenance.")
		}
	}
	if hasRequired(run, "test_observed") && !observed {
		return Condition{ID: "test_provenance_not_overclaimed", State: StateMissingTelemetry, ReasonCode: "test_event_missing", Reason: "required test adapter event is missing", NextAction: "Capture test_observed adapter evidence."}
	}
	return pass("test_provenance_not_overclaimed", "test_provenance_not_required", "test provenance was not required")
}

func providerRefsCondition(run RunEvidence) Condition {
	for _, ref := range run.ProviderRefs {
		if containsSecret(ref.SourceRef) || containsSecret(ref.ChangeRef) || containsSecret(ref.ReviewRef) {
			return fail("provider_refs_portable", "provider_ref_contains_secret", "provider-neutral reference contains credential-like material", "Persist canonical token-free provider references.")
		}
	}
	for _, event := range run.AdapterEvents {
		for _, ref := range event.ProviderRefs {
			if containsSecret(ref) {
				return fail("provider_refs_portable", "provider_ref_contains_secret", "event-level provider reference contains credential-like material", "Persist canonical token-free provider references.")
			}
		}
	}
	return pass("provider_refs_portable", "provider_refs_portable", "provider references are portable and token-free")
}

func redactionMetadataCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if event.SensitiveMetadataPersisted || containsSecret(event.GatewayEvidenceRef) || stringSliceContainsSecret(event.ProviderRefs) {
			return fail("redaction_metadata_consistent", "forbidden_adapter_metadata_persisted", "adapter metadata contains forbidden raw or credential-like material", "Redact adapter metadata before persistence.")
		}
		if sensitiveEvent(event.EventType) && (event.RedactionPolicyDigest == "" || event.RetentionMode == "") {
			return cannotVerify("redaction_metadata_consistent", "redaction_metadata_missing", "sensitive adapter event lacks redaction policy or retention metadata", "Record Block 18 redaction policy and retention mode metadata.")
		}
		if event.RetentionMode != "" && !validRetentionMode(event.RetentionMode) {
			return fail("redaction_metadata_consistent", "invalid_retention_mode", "adapter event declares a non-FR-054 retention mode", "Use FR-054 retention modes.")
		}
	}
	return pass("redaction_metadata_consistent", "redaction_metadata_consistent", "sensitive adapter fields carry safe redaction and retention metadata")
}

func overclaimCondition(run RunEvidence) Condition {
	for _, summary := range run.EventFamilySummaries {
		insufficient := summary.State == StateMissingTelemetry || summary.State == StateUnsupported || summary.State == StateNotIntegrated || summary.State == StateNotAssessed || summary.State == StateCannotVerify || summary.State == StateRetentionLimited || summary.RetentionMode == RetentionDigestOnly || summary.RetentionMode == RetentionNotAssessed
		if summary.Reconstructable && insufficient && summary.CapAnnotation == "" {
			return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "capture-depth output claims reconstruction without sufficient evidence", "Emit a visible capture-depth cap for insufficient evidence.")
		}
	}
	for _, event := range run.AdapterEvents {
		if event.ReconstructableClaimed && (event.CaptureState != "captured" || event.RetentionMode == RetentionDigestOnly || event.RetentionMode == RetentionNotAssessed) && event.CapAnnotation == "" {
			return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "adapter event claims reconstruction beyond captured and retained evidence", "Emit a visible cap annotation or lower the claim.")
		}
	}
	return pass("capture_depth_not_overclaimed", "capture_depth_not_overclaimed", "capture-depth output does not exceed available evidence")
}

func topLevel(conditions []Condition) string {
	highest := StatePass
	for _, condition := range conditions {
		if condition.State == StateFail {
			return StateFail
		}
		switch condition.State {
		case StateCannotVerify, StateNotAssessed, StateMissingTelemetry, StateNotIntegrated, StateUnsupported, StateRetentionLimited:
			highest = StateCannotVerify
		}
	}
	return highest
}

func reasons(conditions []Condition) []string {
	out := []string{}
	for _, condition := range conditions {
		if condition.State != StatePass {
			out = append(out, condition.ReasonCode+": "+condition.Reason)
		}
	}
	sort.Strings(out)
	return out
}

func nextActions(conditions []Condition) []string {
	set := map[string]bool{}
	for _, condition := range conditions {
		if condition.State != StatePass && condition.NextAction != "" {
			set[condition.NextAction] = true
		}
	}
	out := []string{}
	for action := range set {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

func hasRequired(run RunEvidence, eventType string) bool {
	for _, required := range run.RequiredEventTypes {
		if required == eventType {
			return true
		}
	}
	return false
}

func unsupported(run RunEvidence, eventType string) bool {
	for _, unsupported := range run.UnsupportedEventTypes {
		if unsupported == eventType {
			return true
		}
	}
	return false
}

func hasEvent(events []AdapterEvent, eventType string) bool {
	for _, event := range events {
		if event.EventType == eventType {
			return true
		}
	}
	return false
}

func sensitiveEvent(eventType string) bool {
	switch eventType {
	case "tool_call", "command_started", "model_call_observed", "test_observed":
		return true
	default:
		return false
	}
}

func validRetentionMode(mode string) bool {
	switch mode {
	case RetentionDigestOnly, RetentionSanitizedExcerpt, RetentionEncryptedRawRef, RetentionExternalArtifactRef, RetentionNotAssessed:
		return true
	default:
		return false
	}
}

func containsSecret(value string) bool {
	if value == "" {
		return false
	}
	needles := []string{
		"secret-token",
		"password=",
		"token=",
		"bearer ",
		"access_token=",
		"credential=",
		"oidc_token",
		"session_id=",
		"raw prompt",
		"raw_prompt",
		"raw response",
		"raw_response",
		"raw_review_body",
		"tool_input_body",
		"tool_output_body",
		"model_request_payload",
		"model_response_payload",
		"adapter_config_raw",
	}
	for _, needle := range needles {
		if containsFold(value, needle) {
			return true
		}
	}
	return false
}

func stringSliceContainsSecret(values []string) bool {
	for _, value := range values {
		if containsSecret(value) {
			return true
		}
	}
	return false
}

func containsFold(value, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(value); i++ {
		if equalFold(value[i:i+len(needle)], needle) {
			return true
		}
	}
	return false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

func pass(id, code, reason string) Condition {
	return Condition{ID: id, State: StatePass, ReasonCode: code, Reason: reason}
}

func fail(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateFail, ReasonCode: code, Reason: reason, NextAction: action}
}

func cannotVerify(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateCannotVerify, ReasonCode: code, Reason: reason, NextAction: action}
}

// ValidTestInput is used by CLI and fixture tests to avoid duplicating a large
// representative adapter capture run outside this package.
func ValidTestInput() Input {
	return validInput()
}

func validInput() Input {
	runID := "adapter-run-1"
	nonce := "nonce-1"
	source := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	policy := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	return Input{Run: RunEvidence{
		RunID:                 runID,
		RunNonce:              nonce,
		SourceBaseline:        source,
		RunClosedSequence:     20,
		RequiredEventTypes:    []string{"run_started", "task_locked", "tool_call", "command_started", "file_mutation", "model_call_observed", "test_observed", "run_closed"},
		RedactionPolicyDigest: policy,
		GatewayIntegrated:     true,
		GatewayEvidenceBound:  true,
		TaskDriftAssessed:     true,
		AdapterEvents: []AdapterEvent{
			validEvent("evt-run", "run_started", 1, runID, nonce, source, policy),
			validEvent("evt-task", "task_locked", 2, runID, nonce, source, policy),
			validEvent("evt-tool", "tool_call", 3, runID, nonce, source, policy),
			validEvent("evt-command", "command_started", 4, runID, nonce, source, policy),
			validEvent("evt-file", "file_mutation", 5, runID, nonce, source, policy),
			validEvent("evt-model", "model_call_observed", 6, runID, nonce, source, policy),
			validEvent("evt-test", "test_observed", 7, runID, nonce, source, policy),
			validEvent("evt-close", "run_closed", 8, runID, nonce, source, policy),
		},
		ProviderRefs:         []ProviderRef{{SourceRef: "repo:generic/source", SourceCommit: source, ChangeRef: "change:42", ReviewRef: "review:7", Producer: "generic_git_host", ObservedAt: "2026-05-07T10:00:00Z"}},
		EventFamilySummaries: []EventFamilyState{{EventFamily: "tool_call", State: StatePass, RetentionMode: RetentionSanitizedExcerpt, Reconstructable: true}},
	}}
}

func validEvent(id, eventType string, sequence int, runID, nonce, source, policy string) AdapterEvent {
	event := AdapterEvent{
		EventID:                 id,
		EventType:               eventType,
		ProducerIdentity:        "adapter:generic",
		AdapterIdentity:         "adapter:generic",
		IdentityBinding:         IdentityBound,
		ProvenanceScope:         "adapter_observed",
		CaptureState:            "captured",
		RunID:                   runID,
		RunNonce:                nonce,
		SourceBaseline:          source,
		BindingMode:             BindingSameChain,
		Sequence:                sequence,
		PrevEventHash:           "1111111111111111111111111111111111111111111111111111111111111111",
		EventHash:               "2222222222222222222222222222222222222222222222222222222222222222",
		CorrelationRef:          "corr:" + id,
		EventPayloadDigest:      "3333333333333333333333333333333333333333333333333333333333333333",
		RedactionPolicyDigest:   policy,
		RetentionMode:           RetentionSanitizedExcerpt,
		ActorAttributionState:   "bound",
		ModelIdentityProvenance: "gateway_observed",
		TestProvenance:          "ci_executed",
		ExecutedEvidenceClaimed: eventType == "test_observed",
		ToolFamily:              "edit",
	}
	if eventType == "run_started" || eventType == "task_locked" || eventType == "run_closed" || eventType == "file_mutation" {
		event.RetentionMode = RetentionDigestOnly
	}
	return event
}
