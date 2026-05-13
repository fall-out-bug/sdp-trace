package adaptercapture

import (
	"sort"
	"strings"
)

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
	conditions := adapterCaptureConditions(input.Run)
	result := adapterCaptureAssessmentResult(conditions)
	if result.AdapterCaptureAssessment != StatePass {

		result.TrustScope = TrustScopeLocal
	}

	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}

func adapterCaptureConditions(run RunEvidence) []Condition {

	return []Condition{
		contractCondition(run),
		identityCondition(run),
		runBindingCondition(run),
		taskDriftCondition(run),
		toolDepthCondition(run),
		fileMutationCondition(run),
		modelIdentityCondition(run),
		testProvenanceCondition(run),
		providerRefsCondition(run),
		redactionMetadataCondition(run),
		overclaimCondition(run),
	}
}

func adapterCaptureAssessmentResult(conditions []Condition) AssessmentResult {

	return AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileAdapterCapture,
		AdapterCaptureAssessment: topLevel(conditions),
		TrustScope:               TrustScopeAdapter,
		AdapterCaptureConditions: conditions,
	}
}
func contractCondition(run RunEvidence) Condition {
	if len(run.AdapterEvents) == 0 {
		return cannotVerify("adapter_event_contract_valid", "adapter_events_missing", "adapter event evidence is missing", "Supply same-chain adapter events or an adapter bundle.")
	}

	return contractConditionFromEvents(run.AdapterEvents)
}

func contractConditionFromEvents(events []AdapterEvent) Condition {
	seen := map[string]bool{}
	for _, event := range events {
		if adapterEventIsMalformed(event) {

			return fail("adapter_event_contract_valid", "adapter_event_malformed", "adapter event is missing required contract fields", "Emit schema-valid adapter events with producer, adapter, type, and digest fields.")
		}
		if hasDuplicateCorrelationKey(seen, event) {

			return cannotVerify("adapter_event_contract_valid", "conflicting_adapter_events", "multiple adapter events share a correlation key", "Deduplicate adapter events or make conflicts explicit.")
		}
	}
	return pass("adapter_event_contract_valid", "adapter_event_contract_valid", "adapter events match required contract fields")
}

func adapterEventIsMalformed(event AdapterEvent) bool {
	return missingAdapterEventIdentity(event) || missingAdapterEventPayload(event)
}

func missingAdapterEventIdentity(event AdapterEvent) bool {
	return event.EventID == "" || event.EventType == "" || event.ProducerIdentity == "" || event.AdapterIdentity == ""
}

func missingAdapterEventPayload(event AdapterEvent) bool {
	return event.EventPayloadDigest == ""
}

func hasDuplicateCorrelationKey(seen map[string]bool, event AdapterEvent) bool {
	if event.CorrelationRef == "" {

		return false
	}
	key := contractCorrelationKey(event)
	if seen[key] {

		return true
	}
	seen[key] = true
	return false
}

func contractCorrelationKey(event AdapterEvent) string {
	return event.EventType + "\x00" + event.CorrelationRef
}

func identityCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if condition, ok := identityConditionForEvent(event); ok {

			return condition
		}
	}
	return pass("adapter_identity_visible", "adapter_identity_visible", "adapter and producer identity are visible with binding classification")
}

func identityConditionForEvent(event AdapterEvent) (Condition, bool) {
	if adapterIdentityMissing(event) {

		return cannotVerify("adapter_identity_visible", "adapter_identity_missing", "adapter or producer identity is missing", "Record adapter and producer identity."), true
	}
	if !validIdentityBinding(event.IdentityBinding) {

		return cannotVerify("adapter_identity_visible", "adapter_identity_unclassified", "adapter identity binding state is not classified", "Classify adapter identity as self_asserted or bound."), true
	}
	return Condition{}, false
}

func adapterIdentityMissing(event AdapterEvent) bool {
	return event.ProducerIdentity == "" || event.AdapterIdentity == ""
}

func validIdentityBinding(binding string) bool {
	return binding == IdentitySelfAsserted || binding == IdentityBound
}

func runBindingCondition(run RunEvidence) Condition {
	if runIdentityMissing(run) {

		return cannotVerify("run_binding_established", "run_binding_missing", "run id or nonce is missing", "Record run id and run nonce before assessing adapter capture.")
	}
	for _, event := range run.AdapterEvents {
		if condition := adapterEventRunBindingCondition(run, event); condition.ID != "" {
			return condition
		}
	}
	return pass("run_binding_established", "run_binding_established", "adapter events are bound to run id, nonce, and chain or bundle context")
}

func runIdentityMissing(run RunEvidence) bool {
	return run.RunID == "" || run.RunNonce == ""
}

func adapterEventRunBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	if adapterEventRunIdentityMismatch(run, event) {
		return fail("run_binding_established", "run_binding_mismatch", "adapter event contradicts run id or nonce", "Use adapter events bound to the selected run.")
	}

	return adapterEventBindingModeCondition(run, event)
}

func adapterEventBindingModeCondition(run RunEvidence, event AdapterEvent) Condition {

	switch event.BindingMode {
	case BindingSameChain:
		return sameChainBindingCondition(run, event)
	case BindingAdapterBundle:
		return adapterBundleBindingCondition(run, event)
	default:
		return cannotVerify("run_binding_established", "binding_mode_missing", "adapter event binding mode is missing or unsupported", "Use same_chain or adapter_bundle binding.")
	}
}

func adapterEventRunIdentityMismatch(run RunEvidence, event AdapterEvent) bool {
	return event.RunID != run.RunID || event.RunNonce != run.RunNonce
}

func sameChainBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	if eventAfterRunClosure(run, event.Sequence) {

		return cannotVerify("run_binding_established", "late_adapter_event", "adapter event appears after run closure", "Do not use late adapter events to satisfy capture-depth assessment.")
	}
	if sameChainDigestMissing(event) {

		return cannotVerify("run_binding_established", "same_chain_digest_missing", "same-chain adapter event lacks hash linkage", "Record prev_event_hash and event_hash.")
	}
	return Condition{}
}

func eventAfterRunClosure(run RunEvidence, sequence int) bool {
	return run.RunClosedSequence > 0 && sequence > run.RunClosedSequence
}

func sameChainDigestMissing(event AdapterEvent) bool {
	return event.EventHash == "" || event.PrevEventHash == ""
}

func adapterBundleBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	if adapterBundleUnbound(run.AdapterBundle, event) {

		return cannotVerify("run_binding_established", "adapter_bundle_unbound", "adapter bundle is not bound to the selected run", "Bind the adapter bundle head digest into the run artifact.")
	}
	if run.RunClosedSequence > 0 && run.AdapterBundle.ReferencedSequence > run.RunClosedSequence {

		return cannotVerify("run_binding_established", "late_adapter_bundle", "adapter bundle was first referenced after run closure", "Reference adapter bundles before run closure.")
	}
	return Condition{}
}

func adapterBundleUnbound(bundle *AdapterBundle, event AdapterEvent) bool {

	return bundle == nil ||
		event.AdapterBundleHeadDigest == "" ||
		event.AdapterBundleHeadDigest != bundle.HeadDigest ||
		event.AdapterBundleID != bundle.BundleID
}

func taskDriftCondition(run RunEvidence) Condition {
	if !run.TaskDriftAssessed {

		return Condition{ID: "task_drift_visible", State: StateNotAssessed, ReasonCode: "task_drift_not_assessed", Reason: "task drift assessment was not selected", NextAction: "Assess task locks and task_superseded events."}
	}
	if taskSupersessionActorMissing(run.AdapterEvents) {
		return cannotVerify("task_drift_visible", "task_supersession_actor_missing", "task supersession lacks actor attribution state", "Record actor attribution state and task digest refs.")
	}
	return taskDriftPassCondition(run.TaskSupersessionCount)
}

func taskDriftPassCondition(supersessionCount int) Condition {
	if supersessionCount == 0 {
		return pass("task_drift_visible", "no_supersessions_observed", "task drift was assessed and no supersessions were observed")
	}

	return pass("task_drift_visible", "task_supersessions_visible", "task supersessions include visible attribution and digest evidence")
}

func taskSupersessionActorMissing(events []AdapterEvent) bool {
	for _, event := range events {
		if event.EventType == "task_superseded" && event.ActorAttributionState == "" {

			return true
		}
	}
	return false
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
		if fileMutationCorrelationMissing(event) {

			return cannotVerify("file_mutation_correlated", "file_mutation_source_missing", "file mutation is not correlated with source baseline and run id", "Record source baseline and run id correlation for file mutation events.")
		}
	}
	return pass("file_mutation_correlated", "file_mutation_correlated", "file mutation evidence is correlated with source baseline and run id")
}

func fileMutationCorrelationMissing(event AdapterEvent) bool {
	return event.EventType == "file_mutation" && (event.SourceBaseline == "" || event.RunID == "")
}

func modelIdentityCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if modelIdentityOverclaimed(run, event) {

			return fail("model_identity_not_overclaimed", "gateway_identity_overclaimed", "model identity is claimed as gateway-observed without bound gateway evidence", "Keep model identity harness_observed or bind gateway evidence.")
		}
	}
	if !run.GatewayIntegrated {

		return Condition{ID: "model_identity_not_overclaimed", State: StateNotIntegrated, ReasonCode: "gateway_not_integrated", Reason: "gateway provenance is not integrated", NextAction: "Integrate gateway evidence before claiming gateway-observed model identity."}
	}
	return pass("model_identity_not_overclaimed", "model_identity_not_overclaimed", "model identity provenance stays within available gateway evidence")
}

func modelIdentityOverclaimed(run RunEvidence, event AdapterEvent) bool {
	return event.EventType == "model_call_observed" && event.ModelIdentityProvenance == "gateway_observed" && !gatewayModelIdentityBound(run, event)
}

func gatewayModelIdentityBound(run RunEvidence, event AdapterEvent) bool {
	return run.GatewayIntegrated && run.GatewayEvidenceBound && event.IdentityBinding == IdentityBound
}

func testProvenanceCondition(run RunEvidence) Condition {
	if event, ok := firstEvent(run.AdapterEvents, "test_observed"); ok {

		return testProvenanceEventCondition(event)
	}
	if hasRequired(run, "test_observed") {
		return Condition{ID: "test_provenance_not_overclaimed", State: StateMissingTelemetry, ReasonCode: "test_event_missing", Reason: "required test adapter event is missing", NextAction: "Capture test_observed adapter evidence."}
	}
	return pass("test_provenance_not_overclaimed", "test_provenance_not_required", "test provenance was not required")
}

func firstEvent(events []AdapterEvent, eventType string) (AdapterEvent, bool) {
	for _, event := range events {
		if event.EventType == eventType {

			return event, true
		}
	}
	return AdapterEvent{}, false
}

func testProvenanceEventCondition(event AdapterEvent) Condition {
	if testProvenanceExecuted(event.TestProvenance) {
		return pass("test_provenance_not_overclaimed", "test_provenance_executed", "test evidence is bound to CI or wrapper execution")
	}

	return nonExecutedTestProvenanceCondition(event)
}

func nonExecutedTestProvenanceCondition(event AdapterEvent) Condition {

	switch event.TestProvenance {
	case "agent_reported":
		return reportedTestCondition(event, "agent_reported_test_not_executed", "agent-reported tests are claimed as executed evidence", "agent-reported test evidence is visible but non-executed")
	case "harness_observed":
		return reportedTestCondition(event, "harness_observed_test_not_executed", "harness-observed test intent is claimed as executed evidence", "harness-observed test evidence is correlation-only")
	default:
		return cannotVerify("test_provenance_not_overclaimed", "test_provenance_missing", "test provenance is missing or unverifiable", "Record ci_executed or wrapper_executed test provenance.")
	}
}

func testProvenanceExecuted(provenance string) bool {
	return provenance == "ci_executed" || provenance == "wrapper_executed"
}

func reportedTestCondition(event AdapterEvent, failCode, failReason, cannotReason string) Condition {
	if event.ExecutedEvidenceClaimed {

		return fail("test_provenance_not_overclaimed", failCode, failReason, "Bind test evidence to CI or wrapper execution.")
	}
	return cannotVerify("test_provenance_not_overclaimed", "test_execution_unverified", cannotReason, "Capture CI or wrapper-executed test evidence.")
}

func providerRefsCondition(run RunEvidence) Condition {
	if providerRefsContainSecret(run.ProviderRefs) {

		return fail("provider_refs_portable", "provider_ref_contains_secret", "provider-neutral reference contains credential-like material", "Persist canonical token-free provider references.")
	}
	if adapterEventsProviderRefsContainSecret(run.AdapterEvents) {

		return fail("provider_refs_portable", "provider_ref_contains_secret", "event-level provider reference contains credential-like material", "Persist canonical token-free provider references.")
	}
	return pass("provider_refs_portable", "provider_refs_portable", "provider references are portable and token-free")
}

func providerRefsContainSecret(refs []ProviderRef) bool {
	for _, ref := range refs {
		if providerRefContainsSecret(ref) {

			return true
		}
	}
	return false
}

func adapterEventsProviderRefsContainSecret(events []AdapterEvent) bool {
	for _, event := range events {
		if eventProviderRefsContainSecret(event) {

			return true
		}
	}
	return false
}

func providerRefContainsSecret(ref ProviderRef) bool {
	return containsSecret(ref.SourceRef) || containsSecret(ref.ChangeRef) || containsSecret(ref.ReviewRef)
}

func eventProviderRefsContainSecret(event AdapterEvent) bool {
	return stringSliceContainsSecret(event.ProviderRefs)
}

func redactionMetadataCondition(run RunEvidence) Condition {
	for _, event := range run.AdapterEvents {
		if condition := redactionMetadataConditionForEvent(event); condition.State != "" {

			return condition
		}
	}
	return pass("redaction_metadata_consistent", "redaction_metadata_consistent", "sensitive adapter fields carry safe redaction and retention metadata")
}

func redactionMetadataConditionForEvent(event AdapterEvent) Condition {
	if hasForbiddenRedactionMetadata(event) {

		return fail("redaction_metadata_consistent", "forbidden_adapter_metadata_persisted", "adapter metadata contains forbidden raw or credential-like material", "Redact adapter metadata before persistence.")
	}
	if missingRequiredRedactionMetadata(event) {

		return cannotVerify("redaction_metadata_consistent", "redaction_metadata_missing", "sensitive adapter event lacks redaction policy or retention metadata", "Record Block 18 redaction policy and retention mode metadata.")
	}
	if hasInvalidRetentionMode(event) {

		return fail("redaction_metadata_consistent", "invalid_retention_mode", "adapter event declares a non-FR-054 retention mode", "Use FR-054 retention modes.")
	}
	return Condition{}
}

func hasForbiddenRedactionMetadata(event AdapterEvent) bool {
	return event.SensitiveMetadataPersisted ||
		containsSecret(event.GatewayEvidenceRef) ||
		stringSliceContainsSecret(event.ProviderRefs)
}

func missingRequiredRedactionMetadata(event AdapterEvent) bool {
	return sensitiveEvent(event.EventType) && (event.RedactionPolicyDigest == "" || event.RetentionMode == "")
}

func hasInvalidRetentionMode(event AdapterEvent) bool {
	return event.RetentionMode != "" && !validRetentionMode(event.RetentionMode)
}

func overclaimCondition(run RunEvidence) Condition {
	if eventFamiliesOverclaim(run.EventFamilySummaries) {

		return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "capture-depth output claims reconstruction without sufficient evidence", "Emit a visible capture-depth cap for insufficient evidence.")
	}
	if adapterEventsOverclaim(run.AdapterEvents) {

		return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "adapter event claims reconstruction beyond captured and retained evidence", "Emit a visible cap annotation or lower the claim.")
	}
	return pass("capture_depth_not_overclaimed", "capture_depth_not_overclaimed", "capture-depth output does not exceed available evidence")
}

func eventFamiliesOverclaim(summaries []EventFamilyState) bool {
	for _, summary := range summaries {
		if eventFamilyOverclaims(summary) {

			return true
		}
	}
	return false
}

func adapterEventsOverclaim(events []AdapterEvent) bool {
	for _, event := range events {
		if adapterEventOverclaims(event) {

			return true
		}
	}
	return false
}

func eventFamilyOverclaims(summary EventFamilyState) bool {
	return summary.Reconstructable &&
		eventFamilyInsufficient(summary) &&
		summary.CapAnnotation == ""
}

func eventFamilyInsufficient(summary EventFamilyState) bool {
	return insufficientEventFamilyStates[summary.State] || insufficientRetentionModes[summary.RetentionMode]
}

var insufficientEventFamilyStates = map[string]bool{
	StateMissingTelemetry: true,
	StateUnsupported:      true,
	StateNotIntegrated:    true,
	StateNotAssessed:      true,
	StateCannotVerify:     true,
	StateRetentionLimited: true,
}

var insufficientRetentionModes = map[string]bool{
	RetentionDigestOnly:  true,
	RetentionNotAssessed: true,
}

func adapterEventOverclaims(event AdapterEvent) bool {
	return event.ReconstructableClaimed && adapterEventInsufficient(event) && event.CapAnnotation == ""
}

func adapterEventInsufficient(event AdapterEvent) bool {
	return event.CaptureState != "captured" || insufficientRetentionModes[event.RetentionMode]
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

		addNextAction(set, condition)
	}
	out := []string{}
	for action := range set {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

func addNextAction(set map[string]bool, condition Condition) {
	if condition.State != StatePass && condition.NextAction != "" {

		set[condition.NextAction] = true
	}
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
	return sensitiveEventTypes[eventType]
}

var sensitiveEventTypes = map[string]bool{
	"tool_call":           true,
	"command_started":     true,
	"model_call_observed": true,
	"test_observed":       true,
}

func validRetentionMode(mode string) bool {
	return validRetentionModes[mode]
}

var validRetentionModes = map[string]bool{
	RetentionDigestOnly:          true,
	RetentionSanitizedExcerpt:    true,
	RetentionEncryptedRawRef:     true,
	RetentionExternalArtifactRef: true,
	RetentionNotAssessed:         true,
}

func containsSecret(value string) bool {
	if value == "" {
		return false
	}

	for _, needle := range secretMarkers {
		if containsFold(value, needle) {

			return true
		}
	}
	return false
}

var secretMarkers = []string{
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

func stringSliceContainsSecret(values []string) bool {
	for _, value := range values {
		if containsSecret(value) {

			return true
		}
	}
	return false
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
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
	return Input{Run: validRunEvidence(runID, nonce, source, policy)}
}

func validRunEvidence(runID, nonce, source, policy string) RunEvidence {

	run := validRunHeader(runID, nonce, source, policy)
	run.AdapterEvents = validAdapterEvents(runID, nonce, source, policy)
	run.ProviderRefs = validProviderRefs(source)
	run.EventFamilySummaries = validEventFamilySummaries()
	return run
}

func validRunHeader(runID, nonce, source, policy string) RunEvidence {

	return RunEvidence{
		RunID:                 runID,
		RunNonce:              nonce,
		SourceBaseline:        source,
		RunClosedSequence:     20,
		RequiredEventTypes:    validRequiredEventTypes(),
		RedactionPolicyDigest: policy,
		GatewayIntegrated:     true,
		GatewayEvidenceBound:  true,
		TaskDriftAssessed:     true,
	}
}

func validRequiredEventTypes() []string {
	return []string{"run_started", "task_locked", "tool_call", "command_started", "file_mutation", "model_call_observed", "test_observed", "run_closed"}
}

func validAdapterEvents(runID, nonce, source, policy string) []AdapterEvent {
	out := make([]AdapterEvent, 0, len(validEventSpecs))
	for _, spec := range validEventSpecs {

		out = append(out, validEvent(spec.id, spec.eventType, spec.sequence, runID, nonce, source, policy))
	}
	return out
}

type validEventSpec struct {
	id        string
	eventType string
	sequence  int
}

var validEventSpecs = []validEventSpec{
	{id: "evt-run", eventType: "run_started", sequence: 1},
	{id: "evt-task", eventType: "task_locked", sequence: 2},
	{id: "evt-tool", eventType: "tool_call", sequence: 3},
	{id: "evt-command", eventType: "command_started", sequence: 4},
	{id: "evt-file", eventType: "file_mutation", sequence: 5},
	{id: "evt-model", eventType: "model_call_observed", sequence: 6},
	{id: "evt-test", eventType: "test_observed", sequence: 7},
	{id: "evt-close", eventType: "run_closed", sequence: 8},
}

func validProviderRefs(source string) []ProviderRef {
	return []ProviderRef{{SourceRef: "repo:generic/source", SourceCommit: source, ChangeRef: "change:42", ReviewRef: "review:7", Producer: "generic_git_host", ObservedAt: "2026-05-07T10:00:00Z"}}
}

func validEventFamilySummaries() []EventFamilyState {
	return []EventFamilyState{{EventFamily: "tool_call", State: StatePass, RetentionMode: RetentionSanitizedExcerpt, Reconstructable: true}}
}

func validEvent(id, eventType string, sequence int, runID, nonce, source, policy string) AdapterEvent {

	seed := validEventSeed{id: id, eventType: eventType, sequence: sequence, runID: runID, nonce: nonce, source: source, policy: policy}
	event := baseValidEvent(seed)
	if digestOnlyValidEvent(eventType) {

		event.RetentionMode = RetentionDigestOnly
	}
	return event
}

type validEventSeed struct {
	id        string
	eventType string
	sequence  int
	runID     string
	nonce     string
	source    string
	policy    string
}

func baseValidEvent(seed validEventSeed) AdapterEvent {

	event := AdapterEvent{}
	setValidEventIdentity(&event, seed)
	setValidEventBinding(&event, seed)
	setValidEventEvidence(&event, seed)
	setValidEventClaims(&event, seed.eventType)
	return event
}

func setValidEventIdentity(event *AdapterEvent, seed validEventSeed) {

	event.EventID = seed.id
	event.EventType = seed.eventType
	event.ProducerIdentity = "adapter:generic"
	event.AdapterIdentity = "adapter:generic"
	event.IdentityBinding = IdentityBound
	event.ProvenanceScope = "adapter_observed"
}

func setValidEventBinding(event *AdapterEvent, seed validEventSeed) {

	event.RunID = seed.runID
	event.RunNonce = seed.nonce
	event.SourceBaseline = seed.source
	event.BindingMode = BindingSameChain
	event.Sequence = seed.sequence
	event.PrevEventHash = "1111111111111111111111111111111111111111111111111111111111111111"
	event.EventHash = "2222222222222222222222222222222222222222222222222222222222222222"
}

func setValidEventEvidence(event *AdapterEvent, seed validEventSeed) {

	event.CaptureState = "captured"
	event.CorrelationRef = "corr:" + seed.id
	event.EventPayloadDigest = "3333333333333333333333333333333333333333333333333333333333333333"
	event.RedactionPolicyDigest = seed.policy
	event.RetentionMode = RetentionSanitizedExcerpt
}

func setValidEventClaims(event *AdapterEvent, eventType string) {

	event.ActorAttributionState = "bound"
	event.ModelIdentityProvenance = "gateway_observed"
	event.TestProvenance = "ci_executed"
	event.ExecutedEvidenceClaimed = eventType == "test_observed"
	event.ToolFamily = "edit"
}

func digestOnlyValidEvent(eventType string) bool {
	return digestOnlyValidEvents[eventType]
}

var digestOnlyValidEvents = map[string]bool{
	"run_started":   true,
	"task_locked":   true,
	"run_closed":    true,
	"file_mutation": true,
}
