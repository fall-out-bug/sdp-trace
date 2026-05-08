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

func Relay(ctx context.Context, opts RelayOptions, stdin io.Reader, stdout, stderr io.Writer) (Trace, int, error) {
	opts = normalizeRelay(opts)
	if len(opts.Command) == 0 {
		return Trace{}, 0, errors.New("interaction relay requires forward command after --")
	}
	body, err := readBody(stdin)
	if err != nil {
		return Trace{}, 0, err
	}
	event, err := NewObservedEvent(opts, body, 0)
	if err != nil {
		return Trace{}, 0, err
	}
	trace := NewTrace(opts.TaskID, SourceObservedControlChannel, []Event{event}, opts.Now)
	if err := WriteContentBlobs(opts.Out, trace, map[string][]byte{event.InteractionID: body}); err != nil {
		return Trace{}, 0, err
	}
	if err := WriteTrace(opts.Out, trace); err != nil {
		return Trace{}, 0, err
	}
	exitCode, err := runForward(ctx, opts.Command, body, stdout, stderr)
	return trace, exitCode, err
}

func ImportTranscript(opts ImportOptions) (Trace, error) {
	opts = normalizeImport(opts)
	if opts.Source != SourcePreclassifiedTranscript {
		return Trace{}, errors.New("interaction import-transcript requires --source preclassified-transcript-import")
	}
	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return Trace{}, err
	}
	if strings.TrimSpace(opts.EventsJSONL) == "" {
		return Trace{}, errors.New("interaction import-transcript requires --events-jsonl")
	}
	events, err := readJSONLEvents(opts.EventsJSONL)
	if err != nil {
		return Trace{}, err
	}
	if len(events) == 0 {
		return Trace{}, errors.New("interaction import-transcript requires at least one event")
	}
	for i := range events {
		if events[i].TaskID != opts.TaskID {
			return Trace{}, fmt.Errorf("event task_id %q does not match import task_id", events[i].TaskID)
		}
		if events[i].Source.SourceType == SourceAgentReported {
			return Trace{}, errors.New("agent-reported interaction is not accepted as event evidence")
		}
		if events[i].Source.SourceType != "" && events[i].Source.SourceType != SourcePreclassifiedTranscript {
			return Trace{}, fmt.Errorf("unsupported source_type %q", events[i].Source.SourceType)
		}
		events[i].Source.SourceType = SourcePreclassifiedTranscript
		if opts.SourceRef != "" {
			events[i].Source.SourceRef = opts.SourceRef
		}
		if err := ValidateEvent(events[i]); err != nil {
			return Trace{}, err
		}
	}
	if err := validateOrdering(events); err != nil {
		return Trace{}, err
	}
	trace := NewTrace(opts.TaskID, SourcePreclassifiedTranscript, events, opts.Now)
	if err := WriteTrace(opts.Out, trace); err != nil {
		return Trace{}, err
	}
	return trace, nil
}

func NewObservedEvent(opts RelayOptions, body []byte, sequence int) (Event, error) {
	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return Event{}, err
	}
	if err := validateSafeID("target", opts.Target); err != nil {
		return Event{}, err
	}
	if opts.ActorID == "" {
		opts.ActorID = opts.ActorType
	}
	if err := validateSafeID("actor_id", opts.ActorID); err != nil {
		return Event{}, err
	}
	if !validActorType(opts.ActorType) {
		return Event{}, fmt.Errorf("unsupported actor_type %q", opts.ActorType)
	}
	if !validEventType(opts.EventType) {
		return Event{}, fmt.Errorf("unsupported event_type %q", opts.EventType)
	}
	if unsafeCount(body) > 0 {
		return Event{}, errors.New("interaction content contains unsafe material and cannot be retained")
	}
	id := "ix-" + randomHex(12)
	sum := sha256.Sum256(body)
	now := opts.Now.UTC().Format(time.RFC3339)
	return Event{
		SchemaVersion:          SchemaVersion,
		InteractionID:          id,
		TaskID:                 opts.TaskID,
		OperationID:            opts.OperationID,
		StageID:                opts.StageID,
		SourceID:               DefaultRelaySourceID,
		SourceSequence:         sequence,
		EventType:              opts.EventType,
		FrictionClass:          frictionClass(opts.EventType),
		Actor:                  Actor{ID: opts.ActorID, ActorType: opts.ActorType},
		Target:                 opts.Target,
		Source:                 Source{SourceType: SourceObservedControlChannel, SourceID: DefaultRelaySourceID, SourceVersion: SchemaVersion},
		ContentRef:             fmt.Sprintf("sdp://interaction/%s/%s", opts.TaskID, id),
		ContentDigest:          hex.EncodeToString(sum[:]),
		DigestAlgorithm:        DigestAlgorithmSHA256,
		Retention:              RetentionSanitizedExcerpt,
		State:                  StateUnreferenced,
		ObservedBeforeDelivery: true,
		ChannelExclusivity:     ChannelExclusivityNotAssessed,
		CompletenessState:      CompletenessComplete,
		Redaction:              Redaction{PolicyRef: DefaultRedactionPolicyRef, FindingCount: 0},
		ObservedAt:             now,
		CreatedAt:              now,
	}, nil
}

func NewTrace(taskID, sourceType string, events []Event, now time.Time) Trace {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	completeness := CompletenessComplete
	assessment := "assessed"
	notAssessed := []string(nil)
	cannotVerify := []string(nil)
	for _, event := range events {
		switch event.CompletenessState {
		case CompletenessCannotVerify:
			assessment = StateNotAssessed
			completeness = CompletenessCannotVerify
			cannotVerify = append(cannotVerify, "source completeness cannot be verified")
		case CompletenessNotAssessed:
			if completeness != CompletenessCannotVerify {
				assessment = StateNotAssessed
				completeness = CompletenessNotAssessed
			}
			notAssessed = append(notAssessed, "source completeness was not assessed")
		case CompletenessPartial:
			if completeness == CompletenessComplete {
				assessment = "partial"
				completeness = CompletenessPartial
			}
		}
	}
	stamp := now.UTC().Format(time.RFC3339)
	return Trace{
		SchemaVersion:     SchemaVersion,
		TraceID:           "it-" + randomHex(12),
		TaskID:            taskID,
		SourceType:        sourceType,
		CompletenessState: completeness,
		Events:            events,
		AssessmentState:   assessment,
		NotAssessed:       notAssessed,
		CannotVerify:      cannotVerify,
		CreatedAt:         stamp,
		UpdatedAt:         stamp,
	}
}

func ValidateEvent(event Event) error {
	if event.SchemaVersion != SchemaVersion {
		return errors.New("interaction event has unsupported schema_version")
	}
	if err := validateSafeID("interaction_id", event.InteractionID); err != nil {
		return err
	}
	if err := validateSafeID("task_id", event.TaskID); err != nil {
		return err
	}
	if event.OperationID != "" {
		if err := validateSafeID("operation_id", event.OperationID); err != nil {
			return err
		}
	}
	if event.StageID != "" {
		if err := validateSafeID("stage_id", event.StageID); err != nil {
			return err
		}
	}
	if !validEventType(event.EventType) {
		return fmt.Errorf("unsupported event_type %q", event.EventType)
	}
	if event.FrictionClass != frictionClass(event.EventType) {
		return fmt.Errorf("event friction_class %q does not match event_type %q", event.FrictionClass, event.EventType)
	}
	if !validActorType(event.Actor.ActorType) {
		return fmt.Errorf("unsupported actor_type %q", event.Actor.ActorType)
	}
	if err := validateSafeID("actor.id", event.Actor.ID); err != nil {
		return err
	}
	if !validSourceType(event.Source.SourceType) {
		return fmt.Errorf("unsupported source_type %q", event.Source.SourceType)
	}
	if err := validateSafeID("source_id", event.SourceID); err != nil {
		return err
	}
	if !validRetention(event.Retention) {
		return fmt.Errorf("unsupported retention %q", event.Retention)
	}
	if !validCompleteness(event.CompletenessState) {
		return fmt.Errorf("unsupported completeness_state %q", event.CompletenessState)
	}
	if !validChannelExclusivity(event.ChannelExclusivity) {
		return fmt.Errorf("unsupported channel_exclusivity_state %q", event.ChannelExclusivity)
	}
	if !validEventState(event.State) {
		return fmt.Errorf("unsupported state %q", event.State)
	}
	if event.DigestAlgorithm != DigestAlgorithmSHA256 || !sha256Pattern.MatchString(event.ContentDigest) {
		return errors.New("interaction event requires sha256 content digest")
	}
	if event.ContentRef != "" && !contentRefPattern.MatchString(event.ContentRef) {
		return fmt.Errorf("unsupported content_ref %q", event.ContentRef)
	}
	if event.ContentRef == "" && event.NotRetainedReason == "" {
		return errors.New("interaction event without content_ref requires not_retained_reason")
	}
	if _, err := time.Parse(time.RFC3339, event.ObservedAt); err != nil {
		return errors.New("observed_at must be RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, event.CreatedAt); err != nil {
		return errors.New("created_at must be RFC3339")
	}
	if event.SourceSequence < 0 {
		return errors.New("source_sequence must be non-negative")
	}
	for _, ref := range event.ReferenceRefs {
		if !validReference(ref) {
			return fmt.Errorf("unsupported reference_ref %q", ref)
		}
	}
	for _, ref := range event.LLMRefs {
		if ref.LinkageState != StateNotAssessed && ref.LinkageState != StateCannotVerify && ref.LinkageState != StateReferenced {
			return fmt.Errorf("unsupported llm linkage_state %q", ref.LinkageState)
		}
	}
	return nil
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
		sum := sha256.Sum256(body)
		if event.ContentDigest != hex.EncodeToString(sum[:]) {
			return fmt.Errorf("content digest mismatch for %s", event.InteractionID)
		}
		path := contentBlobPath(tracePath, event)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, body, 0o600); err != nil {
			return err
		}
	}
	return nil
}

func ReadTrace(path string) (Trace, error) {
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
	if err := validateSafeID("task_id", trace.TaskID); err != nil {
		return err
	}
	if !validSourceType(trace.SourceType) {
		return fmt.Errorf("unsupported source_type %q", trace.SourceType)
	}
	if len(trace.Events) == 0 {
		return errors.New("interaction trace requires events")
	}
	for _, event := range trace.Events {
		if event.TaskID != trace.TaskID {
			return errors.New("interaction trace event task_id mismatch")
		}
		if err := ValidateEvent(event); err != nil {
			return err
		}
	}
	return validateOrdering(trace.Events)
}

func SummarizeTrace(trace Trace) Summary {
	counts := map[string]int{}
	corrections := 0
	planRejections := 0
	clarifications := 0
	assignmentObserved := false
	unreferenced := 0
	for _, event := range trace.Events {
		counts[event.FrictionClass]++
		switch event.EventType {
		case "plan_rejected":
			planRejections++
		case "clarification_request", "clarification_answer":
			clarifications++
		}
		if event.EventType == "task_assignment" {
			assignmentObserved = true
			continue
		}
		if len(event.ReferenceRefs) == 0 || event.State == StateUnreferenced {
			unreferenced++
		}
		if assignmentObserved && (event.EventType == "corrective_feedback" || event.EventType == "boundary_violation" || event.EventType == "evidence_correction") {
			corrections++
		}
	}
	notAssessed := append([]string{}, trace.NotAssessed...)
	if !assignmentObserved {
		notAssessed = append(notAssessed, "task_assignment event absent; post-assignment correction count is not assessed")
	}
	return Summary{
		SchemaVersion:          SchemaVersion,
		TaskID:                 trace.TaskID,
		TraceID:                trace.TraceID,
		AssessmentState:        trace.AssessmentState,
		EventCount:             len(trace.Events),
		FrictionCounts:         counts,
		CorrectionAfterTask:    corrections,
		PlanRejectionCount:     planRejections,
		ClarificationTurnCount: clarifications,
		UnreferencedEventCount: unreferenced,
		NotAssessed:            notAssessed,
		CannotVerify:           trace.CannotVerify,
	}
}

func ReadEnvelope(path string) (Envelope, error) {
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
	if err := validateSafeID("task_id", envelope.TaskID); err != nil {
		return err
	}
	if err := validateSafeID("envelope_id", envelope.EnvelopeID); err != nil {
		return err
	}
	for _, ref := range envelope.RunRefs {
		if !runRefPattern.MatchString(ref) {
			return fmt.Errorf("unsupported run_ref %q", ref)
		}
	}
	for _, refs := range [][]string{envelope.OperationRefs, envelope.ToolRefs, envelope.MutationRefs, envelope.LLMRefs, envelope.InteractionRefs, envelope.PromiseRefs, envelope.StageRefs, envelope.FrictionRefs, envelope.TaskRefs, envelope.SourceRefs} {
		for _, ref := range refs {
			if !validReference(ref) {
				return fmt.Errorf("unsupported envelope ref %q", ref)
			}
		}
	}
	return nil
}

func SummarizeEnvelope(envelope Envelope) Summary {
	return Summary{
		SchemaVersion:       SchemaVersion,
		TaskID:              envelope.TaskID,
		EnvelopeID:          envelope.EnvelopeID,
		AssessmentState:     envelope.AssessmentState,
		RunRefCount:         len(envelope.RunRefs),
		SourceRefCount:      len(envelope.SourceRefs),
		TaskRefCount:        len(envelope.TaskRefs),
		PromiseRefCount:     len(envelope.PromiseRefs),
		InteractionRefCount: len(envelope.InteractionRefs),
		OperationRefCount:   len(envelope.OperationRefs),
		LLMRefCount:         len(envelope.LLMRefs),
		ToolRefCount:        len(envelope.ToolRefs),
		MutationRefCount:    len(envelope.MutationRefs),
		StageRefCount:       len(envelope.StageRefs),
		FrictionRefCount:    len(envelope.FrictionRefs),
		NotAssessed:         envelope.NotAssessed,
		CannotVerify:        envelope.CannotVerify,
	}
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
	if opts.ActorType == "" {
		opts.ActorType = "human_user"
	}
	if opts.Target == "" {
		opts.Target = "agent"
	}
	if opts.EventType == "" {
		opts.EventType = "corrective_feedback"
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
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
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), MaxBodyBytes*4)
	events := make([]Event, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
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
	switch eventType {
	case "task_assignment", "plan_approved":
		return "none"
	case "clarification_request", "clarification_answer":
		return "clarification"
	case "plan_proposed":
		return "planning"
	case "plan_rejected", "corrective_feedback", "boundary_violation":
		return "correction"
	case "evidence_correction":
		return "evidence"
	case "tool_or_model_drift":
		return "drift"
	case "pause_requested", "resume_approved":
		return "coordination"
	default:
		return ""
	}
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
