package interaction

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
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

func Relay(ctx context.Context, opts RelayOptions, stdin io.Reader, stdout, stderr io.Writer) (Trace, int, error) {

	opts = normalizeRelay(opts)
	if len(opts.Command) == 0 {
		return Trace{}, 0, errors.New("interaction relay requires forward command after --")
	}
	return relayWithCommand(ctx, opts, stdin, stdout, stderr)
}

func relayWithCommand(ctx context.Context, opts RelayOptions, stdin io.Reader, stdout, stderr io.Writer) (Trace, int, error) {

	body, err := readBody(stdin)
	if err != nil {
		return Trace{}, 0, err
	}
	event, err := NewObservedEvent(opts, body, 0)
	if err != nil {
		return Trace{}, 0, err
	}
	trace := NewTrace(opts.TaskID, SourceObservedControlChannel, []Event{event}, opts.Now)
	if err := writeRelayTrace(opts.Out, trace, event, body); err != nil {
		return Trace{}, 0, err
	}
	exitCode, err := runForward(ctx, opts.Command, body, stdout, stderr)
	return trace, exitCode, err
}

func writeRelayTrace(path string, trace Trace, event Event, body []byte) error {

	if err := WriteContentBlobs(path, trace, map[string][]byte{event.InteractionID: body}); err != nil {
		return err
	}
	return WriteTrace(path, trace)
}
func ImportTranscript(opts ImportOptions) (Trace, error) {

	opts = normalizeImport(opts)
	if err := validateImportOptions(opts); err != nil {
		return Trace{}, err
	}
	events, err := importTranscriptEvents(opts)
	if err != nil {
		return Trace{}, err
	}
	trace := NewTrace(opts.TaskID, SourcePreclassifiedTranscript, events, opts.Now)
	if err := WriteTrace(opts.Out, trace); err != nil {
		return Trace{}, err
	}
	return trace, nil
}

func validateImportOptions(opts ImportOptions) error {

	if opts.Source != SourcePreclassifiedTranscript {
		return errors.New("interaction import-transcript requires --source preclassified-transcript-import")
	}
	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return err
	}
	if strings.TrimSpace(opts.EventsJSONL) == "" {
		return errors.New("interaction import-transcript requires --events-jsonl")
	}
	return nil
}

func importTranscriptEvents(opts ImportOptions) ([]Event, error) {

	events, err := readJSONLEvents(opts.EventsJSONL)
	if err != nil {
		return nil, err
	}
	if err := normalizeTranscriptEvents(events, opts); err != nil {
		return nil, err
	}
	return events, nil
}

func normalizeTranscriptEvents(events []Event, opts ImportOptions) error {

	if len(events) == 0 {
		return errors.New("interaction import-transcript requires at least one event")
	}
	for i := range events {
		if err := normalizeTranscriptEvent(&events[i], opts); err != nil {
			return err
		}
	}
	return validateOrdering(events)
}

func normalizeTranscriptEvent(event *Event, opts ImportOptions) error {

	if err := validateTranscriptEventTask(*event, opts); err != nil {
		return err
	}
	if err := validateTranscriptEventSource(*event); err != nil {
		return err
	}
	event.Source.SourceType = SourcePreclassifiedTranscript
	if opts.SourceRef != "" {
		event.Source.SourceRef = opts.SourceRef
	}
	return ValidateEvent(*event)
}

func validateTranscriptEventTask(event Event, opts ImportOptions) error {
	if event.TaskID != opts.TaskID {

		return fmt.Errorf("event task_id %q does not match import task_id", event.TaskID)
	}
	return nil
}

func validateTranscriptEventSource(event Event) error {

	if event.Source.SourceType == SourceAgentReported {
		return errors.New("agent-reported interaction is not accepted as event evidence")
	}
	if event.Source.SourceType != "" && event.Source.SourceType != SourcePreclassifiedTranscript {
		return fmt.Errorf("unsupported source_type %q", event.Source.SourceType)
	}
	return nil
}

func NewObservedEvent(opts RelayOptions, body []byte, sequence int) (Event, error) {

	if err := validateObservedEventOptions(opts); err != nil {
		return Event{}, err
	}
	if unsafeCount(body) > 0 {
		return Event{}, errors.New("interaction content contains unsafe material and cannot be retained")
	}
	return observedEvent(opts, body, sequence), nil
}

func validateObservedEventOptions(opts RelayOptions) error {
	if err := validateObservedEventIDs(opts); err != nil {
		return err
	}

	return validateObservedEventCatalog(opts)
}

func validateObservedEventIDs(opts RelayOptions) error {

	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return err
	}
	if err := validateSafeID("target", opts.Target); err != nil {
		return err
	}
	return validateObservedActorID(opts)
}
func validateObservedActorID(opts RelayOptions) error {
	if opts.ActorID == "" {

		opts.ActorID = opts.ActorType
	}
	return validateSafeID("actor_id", opts.ActorID)
}
func validateObservedEventCatalog(opts RelayOptions) error {

	if !validActorType(opts.ActorType) {
		return fmt.Errorf("unsupported actor_type %q", opts.ActorType)
	}
	if !validEventType(opts.EventType) {
		return fmt.Errorf("unsupported event_type %q", opts.EventType)
	}
	return nil
}

func observedEvent(opts RelayOptions, body []byte, sequence int) Event {

	if opts.ActorID == "" {
		opts.ActorID = opts.ActorType
	}
	id := "ix-" + randomHex(12)
	now := opts.Now.UTC().Format(time.RFC3339)

	return observedEventRecord(opts, body, id, sequence, now)
}

func observedEventRecord(opts RelayOptions, body []byte, id string, sequence int, now string) Event {

	event := observedEventIdentity(opts, id, sequence)
	event.Source = observedRelaySource()
	applyObservedEventContent(&event, opts.TaskID, id, body)
	applyObservedEventAssessment(&event, now)
	return event
}

func observedEventIdentity(opts RelayOptions, id string, sequence int) Event {

	return Event{
		SchemaVersion: SchemaVersion,
		InteractionID: id,

		TaskID:      opts.TaskID,
		OperationID: opts.OperationID,
		StageID:     opts.StageID,

		SourceID:       DefaultRelaySourceID,
		SourceSequence: sequence,

		EventType:     opts.EventType,
		FrictionClass: frictionClass(opts.EventType),

		Actor:  Actor{ID: opts.ActorID, ActorType: opts.ActorType},
		Target: opts.Target,
	}
}

func applyObservedEventContent(event *Event, taskID, id string, body []byte) {

	event.ContentRef = interactionContentRef(taskID, id)
	event.ContentDigest = interactionSHA256Hex(body)
	event.DigestAlgorithm = DigestAlgorithmSHA256
}

func applyObservedEventAssessment(event *Event, now string) {

	event.Retention = RetentionSanitizedExcerpt
	event.State = StateUnreferenced
	event.ObservedBeforeDelivery = true
	event.ChannelExclusivity = ChannelExclusivityNotAssessed
	event.CompletenessState = CompletenessComplete
	event.Redaction = defaultRedaction()

	event.ObservedAt = now
	event.CreatedAt = now
}

func observedRelaySource() Source {

	return Source{SourceType: SourceObservedControlChannel, SourceID: DefaultRelaySourceID, SourceVersion: SchemaVersion}
}

func interactionContentRef(taskID, interactionID string) string {

	return fmt.Sprintf("sdp://interaction/%s/%s", taskID, interactionID)
}

func defaultRedaction() Redaction {

	return Redaction{PolicyRef: DefaultRedactionPolicyRef, FindingCount: 0}
}

func interactionSHA256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func NewTrace(taskID, sourceType string, events []Event, now time.Time) Trace {

	if now.IsZero() {
		now = time.Now().UTC()
	}
	state := newTraceAssessment()
	for _, event := range events {
		state.applyCompleteness(event.CompletenessState)
	}
	return traceFromAssessment(taskID, sourceType, events, now, state)
}

func newTraceAssessment() *traceAssessment {
	return &traceAssessment{completeness: CompletenessComplete, assessment: "assessed"}
}
func (state *traceAssessment) applyCompleteness(completeness string) {

	switch completeness {
	case CompletenessCannotVerify:
		state.markCannotVerify()
	case CompletenessNotAssessed:
		state.markNotAssessed()
	case CompletenessPartial:
		state.markPartial()
	}
}

func (state *traceAssessment) markCannotVerify() {

	state.assessment = StateNotAssessed
	state.completeness = CompletenessCannotVerify
	state.cannotVerify = append(state.cannotVerify, "source completeness cannot be verified")
}

func (state *traceAssessment) markNotAssessed() {

	if state.completeness != CompletenessCannotVerify {
		state.assessment = StateNotAssessed
		state.completeness = CompletenessNotAssessed
	}
	state.notAssessed = append(state.notAssessed, "source completeness was not assessed")
}

func (state *traceAssessment) markPartial() {
	if state.completeness == CompletenessComplete {

		state.assessment = "partial"
		state.completeness = CompletenessPartial
	}
}
func traceFromAssessment(taskID, sourceType string, events []Event, now time.Time, state *traceAssessment) Trace {

	stamp := now.UTC().Format(time.RFC3339)
	return Trace{
		SchemaVersion:     SchemaVersion,
		TraceID:           "it-" + randomHex(12),
		TaskID:            taskID,
		SourceType:        sourceType,
		CompletenessState: state.completeness,
		Events:            events,
		AssessmentState:   state.assessment,
		NotAssessed:       state.notAssessed,
		CannotVerify:      state.cannotVerify,
		CreatedAt:         stamp,
		UpdatedAt:         stamp,
	}
}

func ValidateEvent(event Event) error {

	return firstValidationError(
		func() error { return validateEventIdentity(event) },
		func() error { return validateEventCatalog(event) },
		func() error { return validateEventSource(event) },
		func() error { return validateEventContent(event) },
		func() error { return validateEventTiming(event) },
		func() error { return validateEventRefs(event) },
	)
}

func firstValidationError(checks ...func() error) error {

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}

func validateEventIdentity(event Event) error {

	if event.SchemaVersion != SchemaVersion {
		return errors.New("interaction event has unsupported schema_version")
	}
	if err := validateEventPrimaryIDs(event); err != nil {
		return err
	}
	return validateEventOptionalIDs(event)
}

func validateEventPrimaryIDs(event Event) error {
	if err := validateSafeID("interaction_id", event.InteractionID); err != nil {

		return err
	}
	return validateSafeID("task_id", event.TaskID)
}

func validateEventOptionalIDs(event Event) error {
	if err := validateOptionalSafeID("operation_id", event.OperationID); err != nil {

		return err
	}
	return validateOptionalSafeID("stage_id", event.StageID)
}

func validateOptionalSafeID(label, value string) error {
	if value == "" {

		return nil
	}
	return validateSafeID(label, value)
}

func validateEventCatalog(event Event) error {

	if err := validateSafeID("actor.id", event.Actor.ID); err != nil {
		return err
	}
	if err := validateEventTypeAndFriction(event); err != nil {
		return err
	}
	if err := validateEventActorAndState(event); err != nil {
		return err
	}
	return validateEventRetentionStates(event)
}
func validateEventTypeAndFriction(event Event) error {

	if !validEventType(event.EventType) {
		return fmt.Errorf("unsupported event_type %q", event.EventType)
	}
	if event.FrictionClass != frictionClass(event.EventType) {
		return fmt.Errorf("event friction_class %q does not match event_type %q", event.FrictionClass, event.EventType)
	}
	return nil
}

func validateEventActorAndState(event Event) error {

	if !validActorType(event.Actor.ActorType) {
		return fmt.Errorf("unsupported actor_type %q", event.Actor.ActorType)
	}
	if !validEventState(event.State) {
		return fmt.Errorf("unsupported state %q", event.State)
	}
	return nil
}

func validateEventRetentionStates(event Event) error {

	if !validRetention(event.Retention) {
		return fmt.Errorf("unsupported retention %q", event.Retention)
	}
	if !validCompleteness(event.CompletenessState) {
		return fmt.Errorf("unsupported completeness_state %q", event.CompletenessState)
	}
	if !validChannelExclusivity(event.ChannelExclusivity) {
		return fmt.Errorf("unsupported channel_exclusivity_state %q", event.ChannelExclusivity)
	}
	return nil
}

func validateEventSource(event Event) error {

	if !validSourceType(event.Source.SourceType) {
		return fmt.Errorf("unsupported source_type %q", event.Source.SourceType)
	}
	if err := validateSafeID("source_id", event.SourceID); err != nil {
		return err
	}
	return nil
}

func validateEventContent(event Event) error {
	if err := validateEventDigest(event); err != nil {
		return err
	}

	return validateEventContentRef(event)
}

func validateEventDigest(event Event) error {
	if event.DigestAlgorithm != DigestAlgorithmSHA256 || !sha256Pattern.MatchString(event.ContentDigest) {

		return errors.New("interaction event requires sha256 content digest")
	}
	return nil
}

func validateEventContentRef(event Event) error {

	if err := validateContentRefFormat(event.ContentRef); err != nil {
		return err
	}
	if event.ContentRef == "" && event.NotRetainedReason == "" {
		return errors.New("interaction event without content_ref requires not_retained_reason")
	}
	return nil
}
func validateContentRefFormat(contentRef string) error {
	if contentRef != "" && !contentRefPattern.MatchString(contentRef) {

		return fmt.Errorf("unsupported content_ref %q", contentRef)
	}
	return nil
}

func validateEventTiming(event Event) error {

	if _, err := time.Parse(time.RFC3339, event.ObservedAt); err != nil {
		return errors.New("observed_at must be RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, event.CreatedAt); err != nil {
		return errors.New("created_at must be RFC3339")
	}
	if event.SourceSequence < 0 {
		return errors.New("source_sequence must be non-negative")
	}
	return nil
}

func validateEventRefs(event Event) error {
	if err := validateReferenceRefs(event.ReferenceRefs); err != nil {

		return err
	}
	return validateLLMRefs(event.LLMRefs)
}

func validateReferenceRefs(refs []string) error {

	for _, ref := range refs {
		if !validReference(ref) {
			return fmt.Errorf("unsupported reference_ref %q", ref)
		}
	}
	return nil
}

func validateLLMRefs(refs []LLMRef) error {

	for _, ref := range refs {
		if !validLLMLinkageState(ref.LinkageState) {
			return fmt.Errorf("unsupported llm linkage_state %q", ref.LinkageState)
		}
	}
	return nil
}

func validLLMLinkageState(state string) bool {
	return state == StateNotAssessed || state == StateCannotVerify || state == StateReferenced
}
func WriteTrace(path string, trace Trace) error {

	if strings.TrimSpace(path) == "" {
		return errors.New("interaction command requires --out")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(trace, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func WriteContentBlobs(tracePath string, trace Trace, bodies map[string][]byte) error {

	for _, event := range trace.Events {
		body, ok := bodies[event.InteractionID]
		if !ok {
			continue
		}
		if err := writeContentBlob(tracePath, event, body); err != nil {
			return err
		}
	}
	return nil
}

func writeContentBlob(tracePath string, event Event, body []byte) error {

	sum := sha256.Sum256(body)
	if event.ContentDigest != hex.EncodeToString(sum[:]) {
		return fmt.Errorf("content digest mismatch for %s", event.InteractionID)
	}
	path := contentBlobPath(tracePath, event)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o600)
}

func ReadTrace(path string) (Trace, error) {
	// Reads validate immediately so callers never receive structurally invalid
	// trace data as trusted input.
	var trace Trace
	data, err := os.ReadFile(path)
	if err != nil {
		return Trace{}, err
	}
	if err := json.Unmarshal(data, &trace); err != nil {
		return Trace{}, err
	}
	if err := ValidateTrace(trace); err != nil {
		return Trace{}, err
	}
	return trace, nil
}

func ValidateTrace(trace Trace) error {

	if trace.SchemaVersion != SchemaVersion {
		return errors.New("interaction trace has unsupported schema_version")
	}
	if err := validateTraceHeader(trace); err != nil {
		return err
	}
	if err := validateTraceEvents(trace); err != nil {
		return err
	}
	return validateOrdering(trace.Events)
}

func validateTraceHeader(trace Trace) error {

	if err := validateSafeID("task_id", trace.TaskID); err != nil {
		return err
	}
	if !validSourceType(trace.SourceType) {
		return fmt.Errorf("unsupported source_type %q", trace.SourceType)
	}
	if len(trace.Events) == 0 {
		return errors.New("interaction trace requires events")
	}
	return nil
}

func validateTraceEvents(trace Trace) error {

	for _, event := range trace.Events {
		if event.TaskID != trace.TaskID {
			return errors.New("interaction trace event task_id mismatch")
		}
		if err := ValidateEvent(event); err != nil {
			return err
		}
	}
	return nil
}
func SummarizeTrace(trace Trace) Summary {

	summary := newTraceSummaryCounter()
	for _, event := range trace.Events {
		summarizeTraceEvent(event, summary)
	}
	notAssessed := append([]string{}, trace.NotAssessed...)
	if !summary.assignmentObserved {
		notAssessed = append(notAssessed, "task_assignment event absent; post-assignment correction count is not assessed")
	}

	return traceSummary(trace, summary, notAssessed)
}
func traceSummary(trace Trace, summary *traceSummaryCounter, notAssessed []string) Summary {

	counts := traceSummaryCounts(trace, summary)
	return Summary{
		SchemaVersion:   SchemaVersion,
		TaskID:          trace.TaskID,
		TraceID:         trace.TraceID,
		AssessmentState: trace.AssessmentState,

		EventCount:             counts.events,
		FrictionCounts:         counts.friction,
		CorrectionAfterTask:    counts.corrections,
		PlanRejectionCount:     counts.planRejections,
		ClarificationTurnCount: counts.clarifications,
		UnreferencedEventCount: counts.unreferenced,

		NotAssessed:  notAssessed,
		CannotVerify: trace.CannotVerify,
	}
}

func traceSummaryCounts(trace Trace, summary *traceSummaryCounter) traceCounts {

	return traceCounts{
		events:         len(trace.Events),
		friction:       summary.frictionCounts,
		corrections:    summary.correctionsAfterAssignment,
		planRejections: summary.planRejectionCount,
		clarifications: summary.clarificationCount,
		unreferenced:   summary.unreferencedCount,
	}
}

func newTraceSummaryCounter() *traceSummaryCounter {
	return &traceSummaryCounter{
		frictionCounts: make(map[string]int),
	}
}

func summarizeTraceEvent(event Event, summary *traceSummaryCounter) {

	summary.frictionCounts[event.FrictionClass]++
	if event.EventType == "task_assignment" {
		summary.assignmentObserved = true
		return
	}
	summarizeTraceEventTypeCounters(event.EventType, summary)
	summarizeTraceReferenceAndCorrection(event, summary)
}

func summarizeTraceEventTypeCounters(eventType string, summary *traceSummaryCounter) {

	switch eventType {
	case "plan_rejected":
		summary.planRejectionCount++
	case "clarification_request":
		summary.clarificationCount++
	case "clarification_answer":
		summary.clarificationCount++
	}
}

func summarizeTraceReferenceAndCorrection(event Event, summary *traceSummaryCounter) {

	if eventIsUnreferenced(event) {
		summary.unreferencedCount++
	}
	if summary.assignmentObserved && isPostAssignmentCorrection(event.EventType) {
		summary.correctionsAfterAssignment++
	}
}

func eventIsUnreferenced(event Event) bool {
	return len(event.ReferenceRefs) == 0 || event.State == StateUnreferenced
}

func isPostAssignmentCorrection(eventType string) bool {

	switch eventType {
	case "corrective_feedback", "boundary_violation", "evidence_correction":
		return true
	default:
		return false
	}
}

func ReadEnvelope(path string) (Envelope, error) {
	// Envelopes are validated on read for the same reason traces are: consumers
	// should not summarize invalid reference bundles.
	var envelope Envelope
	data, err := os.ReadFile(path)
	if err != nil {
		return Envelope{}, err
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := ValidateEnvelope(envelope); err != nil {
		return Envelope{}, err
	}
	return envelope, nil
}

func ValidateEnvelope(envelope Envelope) error {

	if envelope.SchemaVersion != SchemaVersion {
		return errors.New("delivery envelope has unsupported schema_version")
	}
	if err := validateEnvelopeIdentity(envelope); err != nil {
		return err
	}
	if err := validateEnvelopeRunRefs(envelope.RunRefs); err != nil {
		return err
	}
	return validateEnvelopeRefs(envelope)
}

func validateEnvelopeIdentity(envelope Envelope) error {
	if err := validateSafeID("task_id", envelope.TaskID); err != nil {

		return err
	}
	return validateSafeID("envelope_id", envelope.EnvelopeID)
}
func validateEnvelopeRunRefs(refs []string) error {

	for _, ref := range refs {
		if !runRefPattern.MatchString(ref) {
			return fmt.Errorf("unsupported run_ref %q", ref)
		}
	}
	return nil
}

func validateEnvelopeRefs(envelope Envelope) error {

	for _, refs := range envelopeReferenceGroups(envelope) {
		for _, ref := range refs {
			if !validReference(ref) {
				return fmt.Errorf("unsupported envelope ref %q", ref)
			}
		}
	}
	return nil
}
func envelopeReferenceGroups(envelope Envelope) [][]string {
	return [][]string{envelope.OperationRefs, envelope.ToolRefs, envelope.MutationRefs, envelope.LLMRefs, envelope.InteractionRefs, envelope.PromiseRefs, envelope.StageRefs, envelope.FrictionRefs, envelope.TaskRefs, envelope.SourceRefs}
}

func SummarizeEnvelope(envelope Envelope) Summary {

	return envelopeSummary(envelope, envelopeReferenceCounts(envelope))
}

func envelopeSummary(envelope Envelope, counts envelopeRefCounts) Summary {

	return Summary{
		SchemaVersion: SchemaVersion,

		TaskID:     envelope.TaskID,
		EnvelopeID: envelope.EnvelopeID,

		AssessmentState: envelope.AssessmentState,

		RunRefCount:         counts.run,
		SourceRefCount:      counts.source,
		TaskRefCount:        counts.task,
		PromiseRefCount:     counts.promise,
		InteractionRefCount: counts.interaction,
		OperationRefCount:   counts.operation,
		LLMRefCount:         counts.llm,
		ToolRefCount:        counts.tool,
		MutationRefCount:    counts.mutation,
		StageRefCount:       counts.stage,
		FrictionRefCount:    counts.friction,

		NotAssessed:  envelope.NotAssessed,
		CannotVerify: envelope.CannotVerify,
	}
}

func envelopeReferenceCounts(envelope Envelope) envelopeRefCounts {

	counts := primaryEnvelopeReferenceCounts(envelope)
	addExecutionEnvelopeReferenceCounts(&counts, envelope)
	return counts
}

func primaryEnvelopeReferenceCounts(envelope Envelope) envelopeRefCounts {

	run := len(envelope.RunRefs)
	source := len(envelope.SourceRefs)
	task := len(envelope.TaskRefs)
	promise := len(envelope.PromiseRefs)

	return envelopeRefCounts{run: run, source: source, task: task, promise: promise}
}

func addExecutionEnvelopeReferenceCounts(counts *envelopeRefCounts, envelope Envelope) {

	counts.interaction = len(envelope.InteractionRefs)
	counts.operation = len(envelope.OperationRefs)
	counts.llm = len(envelope.LLMRefs)
	counts.tool = len(envelope.ToolRefs)

	counts.mutation = len(envelope.MutationRefs)
	counts.stage = len(envelope.StageRefs)
	counts.friction = len(envelope.FrictionRefs)
}

func normalizeRelay(opts RelayOptions) RelayOptions {

	opts.TaskID = strings.TrimSpace(opts.TaskID)
	opts.ActorType = strings.TrimSpace(opts.ActorType)
	opts.ActorID = strings.TrimSpace(opts.ActorID)
	opts.Target = strings.TrimSpace(opts.Target)
	opts.EventType = strings.TrimSpace(opts.EventType)
	opts.OperationID = strings.TrimSpace(opts.OperationID)
	opts.StageID = strings.TrimSpace(opts.StageID)
	opts.Out = strings.TrimSpace(opts.Out)
	return defaultRelayOptions(opts)
}

func defaultRelayOptions(opts RelayOptions) RelayOptions {

	opts = defaultRelayStrings(opts)
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	return opts
}

func defaultRelayStrings(opts RelayOptions) RelayOptions {

	if opts.ActorType == "" {
		opts.ActorType = "human_user"
	}
	if opts.Target == "" {
		opts.Target = "agent"
	}
	if opts.EventType == "" {
		opts.EventType = "corrective_feedback"
	}
	return opts
}
func normalizeImport(opts ImportOptions) ImportOptions {

	opts.TaskID = strings.TrimSpace(opts.TaskID)
	opts.Source = strings.TrimSpace(opts.Source)
	opts.SourceRef = strings.TrimSpace(opts.SourceRef)
	opts.EventsJSONL = strings.TrimSpace(opts.EventsJSONL)
	opts.Out = strings.TrimSpace(opts.Out)
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	return opts
}

func readBody(stdin io.Reader) ([]byte, error) {
	// Empty stdin is rejected because a relay event without content would have no
	// observed interaction payload to hash.
	// The limit reader detects oversize bodies without buffering unbounded
	// interaction content.
	var buf bytes.Buffer
	limited := io.LimitReader(stdin, MaxBodyBytes+1)
	if _, err := buf.ReadFrom(limited); err != nil {
		return nil, err
	}
	body := buf.Bytes()
	if len(body) == 0 {
		return nil, errors.New("interaction relay requires stdin content")
	}
	if len(body) > MaxBodyBytes {
		return nil, fmt.Errorf("interaction content exceeds %d bytes", MaxBodyBytes)
	}
	return body, nil
}

func readJSONLEvents(path string) ([]Event, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return scanJSONLEvents(file)
}
func scanJSONLEvents(file *os.File) ([]Event, error) {

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), MaxBodyBytes*4)
	events := make([]Event, 0)
	for scanner.Scan() {
		var err error
		events, err = appendJSONLEventLine(events, scanner.Text())
		if err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func appendJSONLEventLine(events []Event, text string) ([]Event, error) {
	// Blank JSONL records are skipped without advancing source sequence because
	// they are not evidence rows.
	var event Event
	keep, err := parseJSONLEventLine(text, &event)
	if err != nil || !keep {
		return events, err
	}
	return append(events, event), nil
}

func parseJSONLEventLine(text string, event *Event) (bool, error) {

	line := strings.TrimSpace(text)
	if line == "" {
		return false, nil
	}
	if err := json.Unmarshal([]byte(line), event); err != nil {
		return false, err
	}
	return true, nil
}

func validateOrdering(events []Event) error {

	last := map[string]int{}
	seen := map[string]bool{}
	for _, event := range events {
		key := event.SourceID
		if seen[key] && event.SourceSequence <= last[key] {
			return fmt.Errorf("non-monotonic source_sequence for source %s", key)
		}
		seen[key] = true
		last[key] = event.SourceSequence
	}
	return nil
}

func runForward(ctx context.Context, command []string, body []byte, stdout, stderr io.Writer) (int, error) {

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdin = bytes.NewReader(body)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err == nil {
		return 0, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}
	return 1, err
}

func validateSafeID(label, value string) error {

	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !safeIDPattern.MatchString(value) {
		return fmt.Errorf("%s must match [A-Za-z0-9_.:-]+", label)
	}
	return nil
}
func unsafeCount(body []byte) int {

	text := string(body)
	count := 0
	for _, pattern := range []*regexp.Regexp{privatePathPattern, tokenPattern, authURLPattern} {
		count += len(pattern.FindAllString(text, -1))
	}
	return count
}

func randomHex(n int) string {

	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		sum := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
		return hex.EncodeToString(sum[:])[:n*2]
	}
	return hex.EncodeToString(buf)
}

func validActorType(value string) bool {

	switch value {
	case "human_user", "human_role", "human_group", "ai_agent", "model", "system", "tool", "organization", "other":
		return true
	default:
		return false
	}
}

func validSourceType(value string) bool {

	switch value {
	case SourceObservedControlChannel, SourcePreclassifiedTranscript:
		return true
	default:
		return false
	}
}

func validEventType(value string) bool {
	return frictionClass(value) != ""
}

func validRetention(value string) bool {

	switch value {
	case RetentionDigestOnly, RetentionSanitizedExcerpt, RetentionEncryptedRawRef, RetentionExternalArtifactRef, RetentionNotAssessed:
		return true
	default:
		return false
	}
}

func validCompleteness(value string) bool {

	switch value {
	case CompletenessComplete, CompletenessPartial, CompletenessNotAssessed, CompletenessCannotVerify:
		return true
	default:
		return false
	}
}

func validChannelExclusivity(value string) bool {

	switch value {
	case ChannelExclusivityNotAssessed, StateReferenced, StateCannotVerify:
		return true
	default:
		return false
	}
}
func validEventState(value string) bool {

	switch value {
	case StateObserved, StateReferenced, StateUnreferenced, StateNotAssessed, StateCannotVerify, "redacted":
		return true
	default:
		return false
	}
}

func frictionClass(eventType string) string {
	return frictionClasses[eventType]
}

func validReference(ref string) bool {

	if strings.TrimSpace(ref) == "" || strings.Contains(ref, " ") {
		return false
	}
	if strings.HasPrefix(ref, "recorder:") {
		return recorderRefPattern.MatchString(ref)
	}
	return contentRefPattern.MatchString(ref)
}

func contentBlobPath(tracePath string, event Event) string {

	return filepath.Join(filepath.Dir(tracePath), "interactions", event.TaskID, event.InteractionID+".txt")
}

func EventTypes() []string {

	values := []string{
		"task_assignment",
		"clarification_request",
		"clarification_answer",
		"plan_proposed",
		"plan_approved",
		"plan_rejected",
		"corrective_feedback",
		"boundary_violation",
		"evidence_correction",
		"tool_or_model_drift",
		"pause_requested",
		"resume_approved",
	}
	sort.Strings(values)
	return values
}
