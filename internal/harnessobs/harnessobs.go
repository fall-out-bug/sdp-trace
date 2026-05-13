package harnessobs

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"errors"
	"fmt"

	"hash"
	"io"
	"net/url"

	"os"
	"os/exec"
	"path/filepath"

	"regexp"
	"sort"
	"strings"
	"time"
)

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

var (
	errSessionSourceUnavailable = errors.New("session source unavailable")
	safeIDPattern               = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)
	safeFileIDPattern           = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
	sha256Pattern               = regexp.MustCompile(`^[a-f0-9]{64}$`)
	base64TokenPattern          = regexp.MustCompile(`(?i)^[A-Za-z0-9+/_-]{43,}={0,2}$`)
	providerTokenPrefix         = regexp.MustCompile(`(?i)(sk-[A-Za-z0-9_-]{16,}|xox[baprs]-[A-Za-z0-9-]{16,}|gh[pousr]_[A-Za-z0-9_]{20,})`)
	privatePathPattern          = regexp.MustCompile(`(^|[\s"'])/(Users|home|private|var|tmp)/[^\s"']+`)
	rawFieldNames               = map[string]bool{
		"raw_prompt":         true,
		"prompt":             true,
		"raw_model_response": true,
		"model_response":     true,
		"raw_command":        true,
		"command_body":       true,
	}
	sensitiveFieldNames = map[string]bool{
		"access_token":  true,
		"api_key":       true,
		"apikey":        true,
		"authorization": true,
		"auth":          true,
		"token":         true,
	}
	authQueryKeys = map[string]bool{
		"token": true, "access_token": true, "api_key": true, "apikey": true,
		"key": true, "signature": true, "sig": true, "auth": true,
	}
)

type Profile struct {
	SchemaVersion         string          `json:"schema_version"`
	ProfileID             string          `json:"profile_id"`
	HarnessFamily         string          `json:"harness_family"`
	EventSchemaVersion    string          `json:"event_schema_version"`
	RequiredEventFamilies []string        `json:"required_event_families"`
	OptionalEventFamilies []string        `json:"optional_event_families,omitempty"`
	RawRetentionPolicy    string          `json:"raw_retention_policy"`
	UnsupportedFields     []string        `json:"unsupported_fields,omitempty"`
	DegradationRules      map[string]Rule `json:"degradation_rules"`
	Limits                Limits          `json:"limits,omitempty"`
}

type Limits struct {
	MaxLineBytes int `json:"max_line_bytes,omitempty"`
	MaxEvents    int `json:"max_events,omitempty"`
}

type Rule struct {
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
}

type Event struct {
	EventID            string             `json:"event_id"`
	EventSchemaVersion string             `json:"event_schema_version"`
	EventFamily        string             `json:"event_family"`
	EventType          string             `json:"event_type"`
	ObservedAt         string             `json:"observed_at,omitempty"`
	SourceRef          string             `json:"source_ref"`
	SourceDigest       string             `json:"source_digest"`
	TaskRef            string             `json:"task_ref,omitempty"`
	OperationRef       string             `json:"operation_ref,omitempty"`
	ActorRef           string             `json:"actor_ref,omitempty"`
	ContentState       string             `json:"content_state"`
	UnavailableFields  []UnavailableField `json:"unavailable_fields,omitempty"`
}

type UnavailableField struct {
	Field      string `json:"field"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
}

type Run struct {
	SchemaVersion      string   `json:"schema_version"`
	ProfileID          string   `json:"profile_id"`
	HarnessFamily      string   `json:"harness_family"`
	EventSchemaVersion string   `json:"event_schema_version"`
	SourcePath         string   `json:"source_path"`
	SourceDigest       string   `json:"source_digest"`
	EventCount         int      `json:"event_count"`
	EventRefs          []string `json:"event_refs"`
	CreatedAt          string   `json:"created_at"`
}

type Validation struct {
	SchemaVersion      string      `json:"schema_version"`
	ProfileID          string      `json:"profile_id"`
	HarnessFamily      string      `json:"harness_family"`
	EventSchemaVersion string      `json:"event_schema_version"`
	ValidationState    string      `json:"validation_state"`
	ReasonCode         string      `json:"reason_code"`
	Dimensions         []Dimension `json:"dimensions"`
	EventCount         int         `json:"event_count"`
	ValidationDigest   string      `json:"validation_digest,omitempty"`
	NonAuthority       string      `json:"non_authority"`
}

type Dimension struct {
	Family     string `json:"family"`
	Required   bool   `json:"required"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	EventCount int    `json:"event_count"`
}

type ObserveOptions struct {
	ProfilePath string
	SourcePath  string
	OutDir      string
	Now         time.Time
}
type SessionProfile struct {
	SchemaVersion      string                 `json:"schema_version"`
	ProfileID          string                 `json:"profile_id"`
	HarnessProfilePath string                 `json:"harness_profile_path"`
	EventSourcePath    string                 `json:"event_source_path"`
	RawEventSourcePath string                 `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string                 `json:"raw_event_format,omitempty"`
	SetupActions       []SessionSetupAction   `json:"setup_actions,omitempty"`
	IsolationRules     []SessionIsolationRule `json:"isolation_rules,omitempty"`
	StreamCapture      string                 `json:"stream_capture"`
}

type SessionSetupAction struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
}

type SessionIsolationRule struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	Required   bool   `json:"required"`
}

type SessionRun struct {
	SchemaVersion      string                   `json:"schema_version"`
	ProfileID          string                   `json:"profile_id"`
	HarnessProfilePath string                   `json:"harness_profile_path"`
	EventSourcePath    string                   `json:"event_source_path"`
	RawEventSourcePath string                   `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string                   `json:"raw_event_format,omitempty"`
	SetupActionIDs     []string                 `json:"setup_action_ids,omitempty"`
	IsolationResults   []SessionIsolationResult `json:"isolation_results,omitempty"`
	CommandDigest      string                   `json:"command_digest,omitempty"`
	CommandDigestState string                   `json:"command_digest_state,omitempty"`
	CommandModel       string                   `json:"command_model,omitempty"`
	CommandModelState  string                   `json:"command_model_state,omitempty"`
	ProcessID          int                      `json:"process_id,omitempty"`
	ProcessIDState     string                   `json:"process_id_state,omitempty"`
	StartTime          string                   `json:"start_time,omitempty"`
	EndTime            string                   `json:"end_time,omitempty"`
	SourceCommit       string                   `json:"source_commit,omitempty"`
	SourceCommitState  string                   `json:"source_commit_state,omitempty"`
	ObservedRunDir     string                   `json:"observed_run_dir,omitempty"`
	OutputDigest       string                   `json:"output_digest,omitempty"`
	NormalizedDigest   string                   `json:"normalized_digest,omitempty"`
	CollectionState    string                   `json:"collection_state,omitempty"`
	CollectionReason   string                   `json:"collection_reason,omitempty"`
	CreatedAt          string                   `json:"created_at"`
}

// SessionIsolationResult records the live readback result of a setup rule; it
// is evidence about the local setup artifact, not proof that a harness obeyed it.
type SessionIsolationResult struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	SHA256     string `json:"sha256,omitempty"`
}

type SessionSetupOptions struct {
	ProfilePath string
	OutDir      string
	Command     string
	Now         time.Time
}

type SessionCollectOptions struct {
	ProfilePath string
	RunDir      string
	Now         time.Time
}

type SessionOptions struct {
	ProfilePath string
	OutDir      string
	Command     []string
	Now         time.Time
}

type ValidateOptions struct {
	ProfilePath string
	RunDir      string
	OutPath     string
}

type observationContext struct {
	outDir       string
	sourcePath   string
	sourceDigest string
	now          time.Time
	profile      Profile
	events       []Event
}

type sessionCollectionContext struct {
	profilePath        string
	runDir             string
	now                time.Time
	profile            SessionProfile
	session            SessionRun
	harnessProfile     Profile
	harnessProfilePath string
}

type observedCommandResult struct {
	waitErr error
	end     time.Time
}

var isolationRuleInstallers = map[string]func(string, string) error{
	"ignore_line":    ensureLineFileRule,
	"json_read_deny": ensureJSONReadDenyRule,
}

type eventRefCheck struct {
	ok  bool
	err string
}

var stateRank = map[string]int{
	StateFail:         4,
	StateCannotVerify: 3,
	StateNotAssessed:  2,
	StatePass:         1,
}

type existingPathSpec struct {
	traversalError string
	requireDir     bool
	typeError      string
}
type shellFieldScanner struct {
	fields  []string
	current strings.Builder
	quote   rune
	escaped bool
}

var digestFieldNames = map[string]bool{
	"source_digest":     true,
	"validation_digest": true,
	"commit_digest":     true,
	"envelope_digest":   true,
	"payload_digest":    true,
	"sha256":            true,
}

var validFamilies = map[string]bool{

	"harness":     true,
	"model":       true,
	"interaction": true,
	"phase":       true,
	"review":      true,
	"tool":        true,
	"mutation":    true,
	"test":        true,
	"pr":          true,
	"merge":       true,
	"gap":         true,
}

var validContentStates = map[string]bool{
	ContentRedacted:      true,
	ContentDigestOnly:    true,
	ContentRetainedSafe:  true,
	ContentNotApplicable: true,
}

var validStates = map[string]bool{

	StatePass:         true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
}

var validRuleKeys = map[string]bool{
	"missing_required_family": true,
	"missing_optional_family": true,
	"source_unavailable":      true,
	"unsafe_input":            true,
	"digest_mismatch":         true,
	"schema_version_mismatch": true,
	"cross_link_conflict":     true,
}

func Observe(opts ObserveOptions) (Run, error) {
	// Observe keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ctx, err := prepareObservation(opts)
	if err != nil {
		return Run{}, err
	}

	if err := writeObservationEvents(ctx.outDir, ctx.events); err != nil {
		return Run{}, err
	}
	run := newObservedRun(ctx)

	if err := writeJSON(filepath.Join(ctx.outDir, "run.json"), run); err != nil {
		return Run{}, err
	}
	return run, nil
}

func prepareObservation(opts ObserveOptions) (observationContext, error) {
	// prepareObservation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, sourcePath, outDir, err := validateObserveOptions(opts)
	if err != nil {
		return observationContext{}, err
	}

	profile, events, sourceDigest, err := loadObservationSource(profilePath, sourcePath)
	if err != nil {
		return observationContext{}, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return observationContext{}, err
	}
	return newObservationContext(opts, outDir, sourcePath, sourceDigest, profile, events), nil
}

func newObservationContext(opts ObserveOptions, outDir, sourcePath, sourceDigest string, profile Profile, events []Event) observationContext {
	// newObservationContext keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	return observationContext{
		outDir:       outDir,
		sourcePath:   sourcePath,
		sourceDigest: sourceDigest,
		now:          observationTime(opts.Now),

		profile: profile,
		events:  events,
	}
}

func validateObserveOptions(opts ObserveOptions) (string, string, string, error) {
	// validateObserveOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireObserveOptions(opts); err != nil {
		return "", "", "", err
	}

	return resolveObservePaths(opts)
}

func requireObserveOptions(opts ObserveOptions) error {
	// requireObserveOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "harness observe requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.SourcePath, "harness observe requires --source"); err != nil {
		return err
	}
	if err := requireNonBlank(opts.OutDir, "harness observe requires --out"); err != nil {
		return err
	}
	return nil
}
func resolveObservePaths(opts ObserveOptions) (string, string, string, error) {
	// resolveObservePaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return "", "", "", fmt.Errorf("unsafe profile path: %w", err)
	}

	sourcePath, err := safeExistingFile(opts.SourcePath)
	if err != nil {
		return "", "", "", fmt.Errorf("unsafe source path: %w", err)
	}

	outDir, err := safeOutDir(opts.OutDir)
	if err != nil {
		return "", "", "", err
	}
	return profilePath, sourcePath, outDir, nil
}

func observationTime(now time.Time) time.Time {
	// observationTime keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if now.IsZero() {

		return time.Now().UTC()
	}
	return now
}

func loadObservationSource(profilePath, sourcePath string) (Profile, []Event, string, error) {
	// loadObservationSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Profile{}, nil, "", err
	}

	events, sourceDigest, err := readEvents(profile, sourcePath)
	if err != nil {
		return Profile{}, nil, "", err
	}
	return profile, events, sourceDigest, nil
}

func writeObservationEvents(outDir string, events []Event) error {
	// writeObservationEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, event := range events {

		path := filepath.Join(outDir, "events", event.EventID+".json")
		if err := writeJSON(path, event); err != nil {
			return err
		}
	}
	return nil
}

func newObservedRun(ctx observationContext) Run {
	// newObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.profile.ProfileID,
		HarnessFamily:      ctx.profile.HarnessFamily,
		EventSchemaVersion: ctx.profile.EventSchemaVersion,
		SourcePath:         filepath.Base(ctx.sourcePath),
		SourceDigest:       ctx.sourceDigest,
		EventCount:         len(ctx.events),
		EventRefs:          eventRefs(ctx.events),
		CreatedAt:          ctx.now.Format(time.RFC3339),
	}
}

func SetupSession(opts SessionSetupOptions) (SessionRun, error) {
	// SetupSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, outDir, err := validateSessionSetupOptions(opts)
	if err != nil {
		return SessionRun{}, err
	}

	run, err := setupSessionRun(profilePath, outDir, opts.Now, opts.Command)
	if err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

func validateSessionSetupOptions(opts SessionSetupOptions) (string, string, error) {
	// validateSessionSetupOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	profilePath, err := resolveSessionSetupProfilePath(opts.ProfilePath)
	if err != nil {
		return "", "", err
	}
	outDir, err := resolveSessionSetupOutDir(opts.OutDir)
	if err != nil {
		return "", "", err
	}
	return profilePath, outDir, nil
}
func normalizeOpenCodeRawLine(raw map[string]any, lineNo int, now time.Time) []Event {
	// normalizeOpenCodeRawLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	signals := rawSignals(raw)
	families := openCodeFamilies(raw, signals)
	if len(families) == 0 {

		return nil
	}

	ordered := sortedFamilies(families)
	observedAt := openCodeObservedAt(raw, now)
	actor := openCodeActor(raw)

	sourceRef := fmt.Sprintf("raw-%06d", lineNo)
	return normalizedOpenCodeEvents(ordered, observedAt, sourceRef, actor)
}

func normalizedOpenCodeEvents(ordered []string, observedAt, sourceRef, actor string) []Event {
	// normalizedOpenCodeEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	events := make([]Event, 0, len(ordered))
	for _, family := range ordered {

		events = append(events, normalizedOpenCodeEvent(family, observedAt, sourceRef, actor))
	}
	return events
}

func normalizedOpenCodeEvent(family, observedAt, sourceRef, actor string) Event {
	// normalizedOpenCodeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return normalizedEvent(
		fmt.Sprintf("%s-%s", sourceRef, family),
		family,
		family+"_observed",
		observedAt,
		sourceRef,
		actor,
	)
}
func openCodeFamilies(raw map[string]any, signals []string) map[string]bool {
	// openCodeFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	families := map[string]bool{}

	setFamily(families, "harness", openCodeHarnessFamily(signals))
	setFamily(families, "model", hasKey(raw, "model", "model_id", "modelid", "provider"))
	setFamily(families, "interaction", openCodeInteractionFamily(raw, signals))
	setFamily(families, "tool", openCodeToolFamily(raw, signals))
	setFamily(families, "mutation", openCodeMutationFamily(raw, signals))
	setFamily(families, "test", openCodeTestFamily(signals))
	setFamily(families, "phase", openCodePhaseFamily(raw, signals))
	setFamily(families, "review", hasSignal(signals, "review") || hasSignalPrefix(signals, "review."))
	setFamily(families, "pr", hasSignal(signals, "pull_request", "pull request") || hasSignalPrefix(signals, "pr.", "pr_"))

	setFamily(families, "merge", hasSignal(signals, "merge") || hasSignalPrefix(signals, "merge."))

	return families
}

func setFamily(families map[string]bool, family string, observed bool) {
	if observed {
		families[family] = true
	}
}

func openCodeHarnessFamily(signals []string) bool {
	return hasSignal(signals, "session.started", "session.completed", "run.started", "run.completed", "step_start", "step-start", "step_finish", "step-finish") ||
		hasSignalPrefix(signals, "session.", "run.")
}

func openCodeInteractionFamily(raw map[string]any, signals []string) bool {
	// openCodeInteractionFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "role") ||
		hasSignal(signals, "message", "response", "text") ||
		hasSignalPrefix(signals, "message.", "response.")
}

func openCodeToolFamily(raw map[string]any, signals []string) bool {
	// openCodeToolFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "tool", "tool_call", "toolcall") ||
		hasSignal(signals, "tool.call", "tool.result", "tool_use") ||
		hasSignalPrefix(signals, "tool.")
}

func openCodeMutationFamily(raw map[string]any, signals []string) bool {
	// openCodeMutationFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasSignal(signals, "file.write", "file.edit", "file.patch", "file.delete", "mutation") ||
		hasSignalPrefix(signals, "mutation.") ||
		nativeMutationTool(raw)
}

func openCodeTestFamily(signals []string) bool {
	return hasSignal(signals, "test.finished", "test.started", "test.passed", "test.failed") ||
		hasSignalPrefix(signals, "test.")
}

func openCodePhaseFamily(raw map[string]any, signals []string) bool {
	// openCodePhaseFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "phase") ||
		hasSignal(signals, "phase") ||
		hasSignalPrefix(signals, "phase.", "gsd.", "gsd_")
}

func sortedFamilies(families map[string]bool) []string {
	// sortedFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ordered := make([]string, 0, len(families))
	for family := range families {
		ordered = append(ordered, family)
	}

	sort.Strings(ordered)
	return ordered
}

func openCodeObservedAt(raw map[string]any, now time.Time) string {
	// openCodeObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedAt := findTimestamp(raw)
	if observedAt == "" {

		return now.Format(time.RFC3339)
	}
	return observedAt
}

func openCodeActor(raw map[string]any) string {
	// openCodeActor keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if model := findStringByKey(raw, "model", "model_id", "modelid"); model != "" {

		return safeToken(model)
	}
	if provider := findStringByKey(raw, "provider"); provider != "" {
		return safeToken(provider)
	}
	return "opencode"
}

func sessionCommandFacts(session SessionRun) []Event {
	// sessionCommandFacts keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !hasSessionCommandModel(session) {

		return nil
	}
	event := sessionCommandModelEvent(session)
	data, err := json.Marshal(event)
	if err != nil {
		return nil
	}

	event.SourceDigest = digestLine(data)

	return []Event{event}
}

func hasSessionCommandModel(session SessionRun) bool {
	return session.CommandModelState == StatePass && strings.TrimSpace(session.CommandModel) != ""
}

func sessionCommandModelEvent(session SessionRun) Event {
	// sessionCommandModelEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return normalizedEvent(
		"session-command-model",
		"model",
		"model_observed",
		sessionCommandObservedAt(session),
		"session-command",
		safeToken(session.CommandModel),
	)
}
func sessionCommandObservedAt(session SessionRun) string {
	// sessionCommandObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedAt := session.StartTime
	if observedAt == "" {
		observedAt = session.CreatedAt
	}
	if _, err := time.Parse(time.RFC3339, observedAt); err != nil {

		return time.Now().UTC().Format(time.RFC3339)
	}
	return observedAt
}

func normalizedEvent(id, family, eventType, observedAt, sourceRef, actor string) Event {
	// normalizedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return Event{
		EventID:            id,
		EventSchemaVersion: EventSchemaVersion,
		EventFamily:        family,
		EventType:          eventType,
		ObservedAt:         observedAt,
		SourceRef:          sourceRef,
		SourceDigest:       "",
		ActorRef:           actor,
		ContentState:       ContentDigestOnly,
	}
}

func rawSignals(value any) []string {
	return rawSignalsAt("", value)
}

func rawSignalsAt(parentKey string, value any) []string {
	// rawSignalsAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if signals, ok := rawStructuredSignals(parentKey, value); ok {
		return signals
	}

	return rawLeafSignals(parentKey, value)
}

func rawStructuredSignals(parentKey string, value any) ([]string, bool) {
	// rawStructuredSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:

		return rawMapSignals(v), true
	case []any:

		return rawSliceSignals(parentKey, v), true
	default:
		return nil, false
	}
}

func rawLeafSignals(parentKey string, value any) []string {
	// rawLeafSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case string:
		return rawStringSignals(parentKey, v)
	default:

		return rawScalarSignals(v)
	}
}

func rawScalarSignals(value any) []string {
	return []string{strings.ToLower(fmt.Sprint(value))}
}

func rawStringSignals(parentKey, value string) []string {
	// rawStringSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if rawSignalValueKey(parentKey) {

		return []string{strings.ToLower(value)}
	}
	return nil
}

func rawMapSignals(values map[string]any) []string {
	// rawMapSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parts := make([]string, 0, len(values)*2)
	for key, child := range values {

		parts = append(parts, strings.ToLower(key))
		parts = append(parts, rawSignalsAt(key, child)...)
	}
	return parts
}

func rawSliceSignals(parentKey string, values []any) []string {
	// rawSliceSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parts := make([]string, 0, len(values))
	for _, child := range values {

		parts = append(parts, rawSignalsAt(parentKey, child)...)
	}
	return parts
}

func rawSignalValueKey(key string) bool {
	// rawSignalValueKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch strings.ToLower(key) {
	case "type", "kind", "event", "event_type", "name", "phase", "role", "provider", "model", "model_id", "status", "tool", "action", "operation":

		return true
	default:
		return false
	}
}

func hasSignal(signals []string, values ...string) bool {
	// hasSignal keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, value := range values {

		wanted[strings.ToLower(value)] = true
	}

	for _, signal := range signals {
		if wanted[signal] {
			return true
		}
	}
	return false
}

func hasSignalPrefix(signals []string, prefixes ...string) bool {
	// hasSignalPrefix keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, signal := range signals {
		for _, prefix := range prefixes {

			if strings.HasPrefix(signal, strings.ToLower(prefix)) {
				return true
			}
		}
	}
	return false
}
func hasKey(value any, keys ...string) bool {
	// hasKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}

	return hasKeyIn(value, wanted)
}

func hasKeyIn(value any, wanted map[string]bool) bool {
	// hasKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:
		return hasKeyInMap(v, wanted)
	case []any:
		return hasKeyInSlice(v, wanted)
	}

	return false
}

func hasKeyInMap(values map[string]any, wanted map[string]bool) bool {
	// hasKeyInMap keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, child := range values {

		if wanted[strings.ToLower(key)] || hasKeyIn(child, wanted) {
			return true
		}
	}
	return false
}

func hasKeyInSlice(values []any, wanted map[string]bool) bool {
	// hasKeyInSlice keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, child := range values {

		if hasKeyIn(child, wanted) {
			return true
		}
	}
	return false
}
func findStringByKey(value any, keys ...string) string {
	// findStringByKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}

	return findStringByKeyIn(value, wanted)
}

func findByKeyIn[T any](value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	switch v := value.(type) {
	case map[string]any:
		return findByKeyInMap(v, wanted, match)
	case []any:
		return findByKeyInSlice(v, wanted, match)
	}

	return zero, false
}

func findByKeyInMap[T any](value map[string]any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyInMap keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	for key, child := range value {

		if found, ok := matchWantedKey(key, child, wanted, match); ok {
			return found, true
		}
		if found, ok := findByKeyIn(child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}

func matchWantedKey[T any](key string, value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// matchWantedKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	if !wanted[strings.ToLower(key)] {

		return zero, false
	}
	return match(value)
}

func findByKeyInSlice[T any](value []any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyInSlice keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	for _, child := range value {

		if found, ok := findByKeyIn(child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}

func findStringByKeyIn(value any, wanted map[string]bool) string {
	// findStringByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	matchingString := func(value any) (string, bool) {
		s, ok := value.(string)
		return s, ok && strings.TrimSpace(s) != ""
	}

	s, ok := findByKeyIn(value, wanted, matchingString)
	if !ok {
		return ""
	}
	return s
}

func findTimestamp(raw map[string]any) string {
	// findTimestamp keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, key := range []string{"time", "timestamp", "created_at", "observed_at"} {

		if observedAt := timestampForKey(raw, key); observedAt != "" {
			return observedAt
		}
	}
	return ""
}

func timestampForKey(raw map[string]any, key string) string {
	// timestampForKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if observedAt := stringTimestampForKey(raw, key); observedAt != "" {
		return observedAt
	}

	if value, ok := findNumberByKey(raw, key); ok {
		return unixMillisTimestamp(value)
	}
	return ""
}
func stringTimestampForKey(raw map[string]any, key string) string {
	// stringTimestampForKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	value := findStringByKey(raw, key)
	if value == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {

		return ""
	}
	return parsed.UTC().Format(time.RFC3339)
}

func findNumberByKey(value any, keys ...string) (float64, bool) {
	// findNumberByKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}

	return findNumberByKeyIn(value, wanted)
}

func findNumberByKeyIn(value any, wanted map[string]bool) (float64, bool) {
	// findNumberByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	matchNumber := func(value any) (float64, bool) {
		switch n := value.(type) {
		case float64:

			return n, true
		case int:
			return float64(n), true
		default:
			return 0, false
		}
	}

	return findByKeyIn(value, wanted, matchNumber)
}

func unixMillisTimestamp(value float64) string {
	// unixMillisTimestamp keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if value <= 0 || value <= 1_000_000_000 {

		return ""
	}
	if value > 1_000_000_000_000 {
		return time.UnixMilli(int64(value)).UTC().Format(time.RFC3339)
	}
	return time.Unix(int64(value), 0).UTC().Format(time.RFC3339)
}

func nativeMutationTool(raw map[string]any) bool {
	// nativeMutationTool keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	tool := strings.ToLower(findStringByKey(raw, "tool"))
	switch tool {
	case "edit", "write", "patch", "apply_patch", "update", "delete":

		return true
	default:
		return false
	}
}

func safeToken(value string) string {
	// safeToken keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var b strings.Builder
	for _, r := range value {
		writeSafeTokenRune(&b, r)
		if b.Len() >= 128 {

			break
		}
	}
	token := strings.Trim(b.String(), "-_.:")
	if token == "" {

		return "opencode"
	}
	return token
}
func writeSafeTokenRune(b *strings.Builder, r rune) {
	// writeSafeTokenRune keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeTokenRune(r) {
		b.WriteRune(r)
		return
	}

	b.WriteByte('-')
}

func safeTokenRune(r rune) bool {
	return strings.ContainsRune(safeTokenRunes, r)
}

func LoadRun(dir string) (Run, []Event, error) {
	// LoadRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var run Run
	if err := readJSON(filepath.Join(dir, "run.json"), &run); err != nil {
		return Run{}, nil, err
	}

	if err := validateLoadedRun(run); err != nil {
		return Run{}, nil, err
	}

	events, err := loadRunEvents(dir, run.EventRefs)
	if err != nil {
		return Run{}, nil, err
	}
	return run, events, nil
}

func readJSON(path string, target any) error {
	// readJSON keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func validateLoadedRun(run Run) error {
	// validateLoadedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if run.SchemaVersion != RunSchemaVersion {

		return fmt.Errorf("unsupported run schema_version: %s", run.SchemaVersion)
	}
	return nil
}
func loadRunEvents(dir string, refs []string) ([]Event, error) {
	// loadRunEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	events := make([]Event, 0, len(refs))
	for _, ref := range refs {

		event, err := loadRunEvent(dir, ref)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func loadRunEvent(dir, ref string) (Event, error) {
	// loadRunEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeEventRef(ref) {

		return Event{}, fmt.Errorf("unsafe event ref: %s", ref)
	}

	data, err := os.ReadFile(filepath.Join(dir, ref))
	if err != nil {
		return Event{}, err
	}
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, err
	}

	return event, nil
}

func Summarize(validation Validation) string {
	// Summarize keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var b strings.Builder

	fmt.Fprintf(&b, "Harness observation: %s (%s)\n", validation.ValidationState, validation.ReasonCode)
	fmt.Fprintf(&b, "Profile: %s\n", validation.ProfileID)
	fmt.Fprintf(&b, "Event schema: %s\n", validation.EventSchemaVersion)
	fmt.Fprintf(&b, "Events: %d\n", validation.EventCount)
	fmt.Fprintln(&b, "Dimensions:")
	for _, dim := range validation.Dimensions {
		writeSummaryDimension(&b, dim)
	}
	fmt.Fprintf(&b, "Boundary: %s\n", nonAuthority())
	return b.String()
}

func writeSummaryDimension(b *strings.Builder, dim Dimension) {
	// writeSummaryDimension keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	required := "optional"
	if dim.Required {

		required = "required"
	}

	fmt.Fprintf(b, "- %s [%s]: %s (%s), events=%d\n", dim.Family, required, dim.State, dim.ReasonCode, dim.EventCount)
}

func LoadValidation(path string) (Validation, error) {
	// LoadValidation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var validation Validation
	if err := readExistingJSON(path, &validation); err != nil {
		return Validation{}, err
	}

	if validation.SchemaVersion != ValidationSchemaVersion {
		return Validation{}, fmt.Errorf("unsupported validation schema_version: %s", validation.SchemaVersion)
	}
	return validation, nil
}

func validateProfile(profile Profile) error {
	// validateProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateProfileMetadata(profile); err != nil {
		return err
	}

	if err := validateProfileEventFamilies(profile.RequiredEventFamilies, profile.OptionalEventFamilies); err != nil {
		return err
	}
	return validateProfileDegradationRules(profile.DegradationRules)
}

func validateProfileMetadata(profile Profile) error {
	// validateProfileMetadata keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateProfileIdentity(profile); err != nil {
		return err
	}
	if profile.EventSchemaVersion != EventSchemaVersion {
		return errors.New("unsupported event_schema_version")
	}

	if len(profile.RequiredEventFamilies) == 0 {
		return errors.New("profile requires at least one required_event_family")
	}
	return nil
}

func validateProfileIdentity(profile Profile) error {
	// validateProfileIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.SchemaVersion != ProfileSchemaVersion {
		return fmt.Errorf("unsupported harness profile schema_version: %s", profile.SchemaVersion)
	}

	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe profile_id")
	}

	if !safeIDPattern.MatchString(profile.HarnessFamily) {
		return errors.New("unsafe harness_family")
	}
	return nil
}

func validateProfileEventFamilies(requiredEventFamilies []string, optionalEventFamilies []string) error {
	// validateProfileEventFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateFamilyList(requiredEventFamilies); err != nil {
		return err
	}

	return validateFamilyList(optionalEventFamilies)
}
func validateFamilyList(families []string) error {
	// validateFamilyList keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, family := range families {

		if !validFamily(family) {
			return fmt.Errorf("unsupported event family: %s", family)
		}
	}
	return nil
}
func validateProfileDegradationRules(rules map[string]Rule) error {
	// validateProfileDegradationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, rule := range rules {

		if !validRuleKey(key) {
			return fmt.Errorf("unsupported degradation rule: %s", key)
		}
		if !validDegradationRule(rule) {
			return fmt.Errorf("invalid degradation rule %s", key)
		}
	}
	return nil
}

func validDegradationRule(rule Rule) bool {
	return validState(rule.State) && safeIDPattern.MatchString(rule.ReasonCode)
}

func readEvents(profile Profile, sourcePath string) ([]Event, string, error) {
	// readEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	maxLine, maxEvents := effectiveEventLimits(profile.Limits)
	return scanEvents(profile, file, maxLine, maxEvents)
}

func scanEvents(profile Profile, file io.Reader, maxLine, maxEvents int) ([]Event, string, error) {
	// scanEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)
	events := []Event{}
	sourceHash := sha256.New()

	lineNo := 0
	return scanEventLines(profile, scanner, sourceHash, events, lineNo, maxEvents)
}

func scanEventLines(profile Profile, scanner *bufio.Scanner, sourceHash hash.Hash, events []Event, lineNo, maxEvents int) ([]Event, string, error) {
	// scanEventLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for scanner.Scan() {

		lineNo++
		line := scanner.Bytes()
		event, ok, err := readEventLine(profile, line, lineNo, len(events), maxEvents, sourceHash)
		if err != nil {
			return nil, "", err
		}
		events = appendScannedEvent(events, event, ok)
	}
	return scannedEvents(events, sourceHash, scanner.Err())
}

func appendScannedEvent(events []Event, event Event, ok bool) []Event {
	// appendScannedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if ok {

		return append(events, event)
	}
	return events
}

func scannedEvents(events []Event, sourceHash hash.Hash, scanErr error) ([]Event, string, error) {
	// scannedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanErr != nil {
		return nil, "", scanErr
	}

	return events, hex.EncodeToString(sourceHash.Sum(nil)), nil
}

func readEventLine(profile Profile, line []byte, lineNo, eventCount, maxEvents int, sourceHash io.Writer) (Event, bool, error) {
	// readEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, err := sourceHash.Write(line); err != nil {
		return Event{}, false, err
	}
	if blankJSONLLine(line) {

		return Event{}, false, nil
	}
	if eventCount >= maxEvents {

		return Event{}, false, fmt.Errorf("source line %d: event limit exceeded", lineNo)
	}
	event, err := parseEvent(profile, line, lineNo)
	return event, err == nil, err
}

func effectiveEventLimits(limits Limits) (int, int) {
	// effectiveEventLimits keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	maxLine := limits.MaxLineBytes
	if maxLine <= 0 {
		maxLine = DefaultMaxLineBytes
	}
	maxEvents := limits.MaxEvents
	if maxEvents <= 0 {

		maxEvents = DefaultMaxEvents
	}
	return maxLine, maxEvents
}

func parseEvent(profile Profile, line []byte, lineNo int) (Event, error) {
	// parseEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	event, err := decodeSafeEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}

	return event, validateParsedEvent(profile, event, line, lineNo)
}

func decodeSafeEventLine(line []byte, lineNo int) (Event, error) {
	// decodeSafeEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	raw, err := decodeRawEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}

	if err := rejectUnsafeEvent(raw, lineNo); err != nil {
		return Event{}, err
	}
	return decodeEventLine(line, lineNo)
}

func decodeRawEventLine(line []byte, lineNo int) (map[string]any, error) {
	// decodeRawEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {

		return nil, fmt.Errorf("source line %d: malformed_jsonl", lineNo)
	}
	return raw, nil
}

func decodeEventLine(line []byte, lineNo int) (Event, error) {
	// decodeEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var event Event
	if err := json.Unmarshal(line, &event); err != nil {

		return Event{}, fmt.Errorf("source line %d: malformed_event", lineNo)
	}
	return event, nil
}
func rejectUnsafeEvent(raw map[string]any, lineNo int) error {
	// rejectUnsafeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeField, reason := findUnsafe(raw); unsafeField != "" {

		return fmt.Errorf("source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}

func validateParsedEvent(profile Profile, event Event, line []byte, lineNo int) error {
	// validateParsedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	expected := digestLine(line)
	if event.SourceDigest != expected {
		return fmt.Errorf("source line %d: source_digest_mismatch:%s", lineNo, safeEvent(event.EventID))
	}

	if err := validateEvent(profile, event); err != nil {
		return fmt.Errorf("source line %d: %w", lineNo, err)
	}
	return nil
}
func validateEvent(profile Profile, event Event) error {
	// validateEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateEventIdentity(profile, event); err != nil {
		return err
	}

	if err := validateEventRefs(event); err != nil {
		return err
	}
	if err := validateEventContent(event); err != nil {
		return err
	}
	return validateUnavailableFields(event.UnavailableFields)
}

func validateEventIdentity(profile Profile, event Event) error {
	// validateEventIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, check := range eventIdentityChecks(profile, event) {

		if !check.ok {
			return errors.New(check.err)
		}
	}
	return validateObservedAt(event.ObservedAt)
}

func eventIdentityChecks(profile Profile, event Event) []eventRefCheck {
	// eventIdentityChecks keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return []eventRefCheck{
		{safeFileIDPattern.MatchString(event.EventID), "unsafe event_id"},
		{event.EventSchemaVersion == profile.EventSchemaVersion, "schema_version_mismatch"},
		{validFamily(event.EventFamily), "unsupported event_family"},
		{safeIDPattern.MatchString(event.EventType), "unsafe event_type"},
	}
}

func validateObservedAt(value string) error {
	// validateObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if value == "" {

		return nil
	}
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return errors.New("invalid observed_at")
	}
	return nil
}

func validateEventRefs(event Event) error {
	// validateEventRefs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, check := range eventRefChecks(event) {

		if !check.ok {
			return errors.New(check.err)
		}
	}
	return nil
}

func eventRefChecks(event Event) []eventRefCheck {
	// eventRefChecks keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return []eventRefCheck{
		{safeRef(event.SourceRef), "unsafe source_ref"},
		{sha256Pattern.MatchString(event.SourceDigest), "invalid source_digest"},
		{event.TaskRef == "" || safeRef(event.TaskRef), "unsafe task_ref"},
		{event.OperationRef == "" || safeOperationRef(event.OperationRef), "unsafe operation_ref"},
		{event.ActorRef == "" || safeRef(event.ActorRef), "unsafe actor_ref"},
	}
}

func validateEventContent(event Event) error {
	// validateEventContent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !validContentState(event.ContentState) {

		return errors.New("invalid content_state")
	}
	return nil
}

func validateUnavailableFields(fields []UnavailableField) error {
	// validateUnavailableFields keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, field := range fields {

		if !validUnavailableField(field) {
			return errors.New("invalid unavailable_fields")
		}
	}
	return nil
}

func validUnavailableField(field UnavailableField) bool {
	// validUnavailableField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeIDPattern.MatchString(field.Field) &&
		field.State == StateNotAssessed &&
		safeIDPattern.MatchString(field.ReasonCode)
}

func evaluate(profile Profile, run Run, events []Event) Validation {
	// evaluate keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	dimensions := evaluationDimensions(profile, events)
	state, reason := compose(dimensions)
	if run.EventSchemaVersion != profile.EventSchemaVersion {

		state, reason = StateCannotVerify, "schema_version_mismatch"
	}
	return validationFromEvaluation(profile, dimensions, len(events), state, reason)
}
func validationFromEvaluation(profile Profile, dimensions []Dimension, eventCount int, state, reason string) Validation {
	// validationFromEvaluation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,

		ValidationState: state,
		ReasonCode:      reason,
		Dimensions:      dimensions,
		EventCount:      eventCount,
		NonAuthority:    nonAuthority(),
	}

	validation.ValidationDigest = validationDigest(validation)
	return validation
}

func evaluationDimensions(profile Profile, events []Event) []Dimension {
	// evaluationDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	counts := eventFamilyCounts(events)
	dimensions, required := requiredDimensions(profile.RequiredEventFamilies, counts)
	dimensions = append(dimensions, optionalDimensions(profile.OptionalEventFamilies, required, counts)...)

	sortDimensions(dimensions)
	return dimensions
}

func eventFamilyCounts(events []Event) map[string]int {
	// eventFamilyCounts keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	counts := map[string]int{}
	for _, event := range events {

		counts[event.EventFamily]++
	}
	return counts
}

func requiredDimensions(families []string, counts map[string]int) ([]Dimension, map[string]bool) {
	// requiredDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	required := map[string]bool{}
	dimensions := make([]Dimension, 0, len(families))
	for _, family := range families {

		required[family] = true
		dimensions = append(dimensions, dimension(family, true, counts[family]))
	}
	return dimensions, required
}

func optionalDimensions(families []string, required map[string]bool, counts map[string]int) []Dimension {
	// optionalDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	dimensions := []Dimension{}
	for _, family := range families {
		if !required[family] {

			dimensions = append(dimensions, dimension(family, false, counts[family]))
		}
	}
	return dimensions
}

func sortDimensions(dimensions []Dimension) {
	// sortDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sort.Slice(dimensions, func(i, j int) bool {
		if dimensions[i].Required != dimensions[j].Required {
			return dimensions[i].Required
		}
		return dimensions[i].Family < dimensions[j].Family
	})
}
func dimension(family string, required bool, count int) Dimension {
	// dimension keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if count > 0 {
		return Dimension{Family: family, Required: required, State: StatePass, ReasonCode: "event_family_observed", EventCount: count}
	}

	reason := "optional_event_family_absent"
	if required {
		reason = "required_event_family_absent"
	}
	return Dimension{Family: family, Required: required, State: StateNotAssessed, ReasonCode: reason}
}

func compose(dimensions []Dimension) (string, string) {
	// compose keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	state := StatePass
	reason := "all_required_dimensions_observed"
	for _, dim := range dimensions {
		if !dim.Required {

			continue
		}

		if rank(dim.State) > rank(state) {
			state = dim.State
			reason = dim.ReasonCode
		}
	}
	return state, reason
}

func rank(state string) int {

	return stateRank[state]
}

func safeExistingFile(path string) (string, error) {
	// safeExistingFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local file without traversal",
		requireDir:     false,
		typeError:      "path must be a file",
	})
}

func safeExistingDir(path string) (string, error) {
	// safeExistingDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local directory without traversal",
		requireDir:     true,
		typeError:      "path must be a directory",
	})
}
func safeExistingPath(path string, spec existingPathSpec) (string, error) {
	// safeExistingPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	cleanPath, err := sanitizeExistingPath(path, spec.traversalError)
	if err != nil {
		return "", err
	}

	abs, err := resolveExistingAbsolutePath(cleanPath)
	if err != nil {
		return "", err
	}
	rel, err := relativeWorkingDirectoryPath(abs)
	if err != nil {
		return "", err
	}
	return rel, ensureExpectedPathType(abs, spec)
}

func resolveExistingAbsolutePath(path string) (string, error) {
	// resolveExistingAbsolutePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	return absolutePath(resolved)
}

func sanitizeExistingPath(path, traversalError string) (string, error) {
	// sanitizeExistingPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New(traversalError)
	}

	return filepath.Clean(path), nil
}

func relativeWorkingDirectoryPath(abs string) (string, error) {
	// relativeWorkingDirectoryPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	rel, err := relativePathFromWorkingDirectory(abs)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {

		return "", errors.New("path escapes working directory")
	}
	return rel, nil
}

func ensureExpectedPathType(path string, spec existingPathSpec) error {
	// ensureExpectedPathType keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() == spec.requireDir {
		return nil
	}

	return errors.New(spec.typeError)
}

func safeOutFile(path string) (string, error) {
	// safeOutFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local file without traversal")
	}

	parent, err := safeParentDir(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	base := filepath.Base(path)
	if !safeOutputBaseName(base) {
		return "", errors.New("unsafe output filename")
	}

	return filepath.Join(parent, base), nil
}

func unsafeOutputPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}

func safeOutputBaseName(base string) bool {
	// safeOutputBaseName keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	stem := strings.TrimSuffix(base, filepath.Ext(base))

	return safeFileIDPattern.MatchString(stem) && !strings.ContainsAny(base, `/\`)
}

func safeParentDir(path string) (string, error) {
	// safeParentDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	clean, err := normalizePotentialParentPath(path)
	if err != nil {
		return "", err
	}

	resolved, err := resolveParentPathWithinWorkingDirectory(clean)
	if err != nil {
		return "", err
	}

	return ensurePathInsideWorkingDirectory(resolved)
}

func normalizePotentialParentPath(path string) (string, error) {
	// normalizePotentialParentPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if path == "" {

		path = "."
	}
	if err := validatePotentialParentPath(path); err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

func validatePotentialParentPath(path string) error {
	// validatePotentialParentPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafePotentialParentPath(path) {

		return errors.New("parent path must be relative local directory without traversal")
	}
	return nil
}
func unsafePotentialParentPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}
func resolveParentPathWithinWorkingDirectory(clean string) (string, error) {
	// resolveParentPathWithinWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		return resolveMissingParent(clean)
	}
	return resolved, nil
}

func resolveMissingParent(clean string) (string, error) {
	// resolveMissingParent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	if parent == clean {

		return "", os.ErrNotExist
	}
	return resolveParentPathWithinWorkingDirectory(parent)
}

func ensurePathInsideWorkingDirectory(path string) (string, error) {
	// ensurePathInsideWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	absPath, err := absolutePath(path)
	if err != nil {
		return "", err
	}
	rel, err := relativePathFromWorkingDirectory(absPath)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {

		return "", errors.New("parent path escapes working directory")
	}

	return rel, nil
}

func absolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

func relativePathFromWorkingDirectory(path string) (string, error) {
	// relativePathFromWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, path)
}

func safeOutDir(path string) (string, error) {
	// safeOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local directory without traversal")
	}

	clean := filepath.Clean(path)
	return safeCleanOutDir(clean)
}

func safeCleanOutDir(clean string) (string, error) {
	// safeCleanOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	exists, err := pathExistsForLstat(clean)
	if err != nil {
		return "", err
	}
	if exists {
		return safeExistingOutDir(clean)
	}

	if err := ensureOutParentInsideWorkingDirectory(clean); err != nil {
		return "", err
	}
	return ensureOutDirEmptyOrMissing(clean)
}

func pathExistsForLstat(path string) (bool, error) {
	// pathExistsForLstat keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, err := os.Lstat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {

		return false, err
	}
}

func ensureOutParentInsideWorkingDirectory(clean string) error {
	// ensureOutParentInsideWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parentEscapes, err := outParentEscapes(clean)
	if err != nil {
		return err
	}
	if parentEscapes {

		return errors.New("out parent path escapes working directory")
	}
	return nil
}

func safeExistingOutDir(clean string) (string, error) {
	// safeExistingOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	rel, err := relativeSymlinkTarget(clean)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {

		return "", errors.New("out path escapes working directory")
	}

	return ensureOutDirEmptyOrMissing(rel)
}

func outParentEscapes(clean string) (bool, error) {
	// outParentEscapes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	for parent != "." && parent != string(filepath.Separator) {

		found, escapes, err := existingParentEscapes(parent)
		if found {
			return escapes, err
		}
		parent = filepath.Dir(parent)
	}
	return false, nil
}
func existingParentEscapes(parent string) (bool, bool, error) {
	// existingParentEscapes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, statErr := os.Lstat(parent); statErr != nil {

		return false, false, nil
	}
	rel, err := relativeSymlinkTarget(parent)
	if err != nil {
		return true, false, err
	}
	return true, pathEscapesWorkingDirectory(rel), nil
}

func relativeSymlinkTarget(path string) (string, error) {
	// relativeSymlinkTarget keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}

func pathEscapesWorkingDirectory(rel string) bool {
	return strings.HasPrefix(rel, "..") || filepath.IsAbs(rel)
}

func ensureOutDirEmptyOrMissing(clean string) (string, error) {
	// ensureOutDirEmptyOrMissing keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	entries, err := os.ReadDir(clean)
	if err := validateOutDirEntries(entries, err); err != nil {
		return "", err
	}

	return clean, nil
}
func validateOutDirEntries(entries []os.DirEntry, err error) error {
	// validateOutDirEntries keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err == nil && len(entries) > 0 {
		return errors.New("harness observe refuses existing non-empty --out")
	}

	return ignorableMissingOutDir(err)
}

func ignorableMissingOutDir(err error) error {
	// ignorableMissingOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err == nil || errors.Is(err, os.ErrNotExist) {

		return nil
	}
	return err
}

func writeJSON(path string, value any) error {
	// writeJSON keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func eventRefs(events []Event) []string {
	// eventRefs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	refs := make([]string, 0, len(events))
	for _, event := range events {

		refs = append(refs, filepath.ToSlash(filepath.Join("events", event.EventID+".json")))
	}
	sort.Strings(refs)
	return refs
}

func safeEventRef(ref string) bool {
	// safeEventRef keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeEventRefPath(ref) {
		return false
	}
	if !strings.HasPrefix(ref, "events/") || !strings.HasSuffix(ref, ".json") {
		return false
	}

	id := strings.TrimSuffix(strings.TrimPrefix(ref, "events/"), ".json")
	return safeFileIDPattern.MatchString(id)
}

func unsafeEventRefPath(ref string) bool {

	return strings.Contains(ref, "\\") || strings.Contains(ref, "..") || filepath.IsAbs(ref)
}

func digestLine(line []byte) string {
	// digestLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err == nil {
		raw["source_digest"] = ""
		canonical, err := json.Marshal(raw)
		if err == nil {

			sum := sha256.Sum256(canonical)
			return hex.EncodeToString(sum[:])
		}
	}

	sum := sha256.Sum256(line)
	return hex.EncodeToString(sum[:])
}

func validationDigest(validation Validation) string {
	// validationDigest keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	copy := validation
	copy.ValidationDigest = ""

	data, _ := json.Marshal(copy)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func digestCommand(command []string) string {
	// digestCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	data, _ := json.Marshal(command)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
func extractCommandModel(command []string) string {
	// extractCommandModel keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if shellCommand := shellCommandString(command); shellCommand != "" {

		if model := extractCommandModelArgs(shellFields(shellCommand)); model != "" {
			return model
		}
	}
	return extractCommandModelArgs(command)
}

func shellCommandString(command []string) string {
	// shellCommandString keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !shellCommandShape(command) {
		return ""
	}
	base := filepath.Base(command[0])
	if base != "sh" && base != "bash" {
		return ""
	}

	return command[2]
}

func shellCommandShape(command []string) bool {
	return len(command) >= 3 && command[1] == "-c"
}

func extractCommandModelArgs(args []string) string {
	// extractCommandModelArgs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i, arg := range args {

		if model, matched := commandModelArg(args, i, arg); matched {
			return model
		}
	}
	return ""
}

func commandModelArg(args []string, i int, arg string) (string, bool) {
	// commandModelArg keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if arg == "--model" || arg == "-m" {

		return nextCommandModelArg(args, i), true
	}
	if model, ok := prefixedCommandModelArg(arg); ok {
		return model, true
	}
	return "", false
}

func nextCommandModelArg(args []string, i int) string {
	// nextCommandModelArg keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if i+1 >= len(args) {
		return ""
	}

	return safeCommandModel(args[i+1])
}

func prefixedCommandModelArg(arg string) (string, bool) {
	// prefixedCommandModelArg keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, prefix := range []string{"--model=", "-m="} {
		if strings.HasPrefix(arg, prefix) {

			return safeCommandModel(strings.TrimPrefix(arg, prefix)), true
		}
	}
	return "", false
}

// shellFields handles the shell field syntax needed to locate --model inside a
// controlled sh -c wrapper. It is not a general shell parser; model values still
// have to pass safeCommandModel before they become retained facts.
func shellFields(command string) []string {
	// shellFields keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	scanner := shellFieldScanner{}
	for _, r := range command {
		scanner.scan(r)
	}
	return scanner.finish()
}
func (scanner *shellFieldScanner) scan(r rune) {
	// scan keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.consumeEscaped(r) {
		return
	}
	if scanner.startsEscape(r) {

		scanner.escaped = true
		return
	}
	if scanner.consumeQuoted(r) {
		return
	}

	scanner.consumeUnquoted(r)
}

func (scanner *shellFieldScanner) consumeEscaped(r rune) bool {
	// consumeEscaped keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !scanner.escaped {
		return false
	}
	if r != '\n' {

		scanner.current.WriteRune('\\')
		scanner.current.WriteRune(r)
	}
	scanner.escaped = false
	return true
}

func (scanner *shellFieldScanner) startsEscape(r rune) bool {
	return scanner.quote != '\'' && r == '\\'
}

func (scanner *shellFieldScanner) consumeQuoted(r rune) bool {
	// consumeQuoted keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.quote == 0 {
		return false
	}
	if r == scanner.quote {

		scanner.quote = 0
		return true
	}
	scanner.current.WriteRune(r)
	return true
}

func (scanner *shellFieldScanner) consumeUnquoted(r rune) {
	// consumeUnquoted keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch {
	case shellQuote(r):

		scanner.quote = r
	case shellFieldSeparator(r):
		scanner.flush()
	default:
		scanner.current.WriteRune(r)
	}
}
func shellQuote(r rune) bool {
	return r == '\'' || r == '"'
}

func (scanner *shellFieldScanner) finish() []string {
	// finish keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.escaped {

		scanner.current.WriteRune('\\')
	}
	scanner.flush()
	return scanner.fields
}

func (scanner *shellFieldScanner) flush() {
	// flush keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.current.Len() == 0 {
		return
	}

	scanner.fields = append(scanner.fields, scanner.current.String())
	scanner.current.Reset()
}

func shellFieldSeparator(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func safeCommandModel(model string) string {
	// safeCommandModel keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	model = strings.TrimSpace(model)
	if unsafeCommandModelIdentity(model) {
		return ""
	}
	if unsafeCommandModelPath(model) {
		return ""
	}
	if len(model) > 128 {

		return ""
	}
	return model
}

func unsafeCommandModelChars(model string) bool {
	return strings.Contains(model, "://") || strings.ContainsAny(model, " \t\n\r\"'`$\\")
}

func unsafeCommandModelIdentity(model string) bool {
	return model == "" || unsafeCommandModelChars(model)
}

func unsafeCommandModelPath(model string) bool {
	return strings.Contains(model, "../") || strings.HasPrefix(model, "/")
}

func digestFile(path string) string {
	// digestFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func sourceCommit() string {
	// sourceCommit keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	data, err := cmd.Output()
	if err != nil {
		return ""
	}

	commit := strings.TrimSpace(string(data))
	if regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(commit) {

		return commit
	}
	return ""
}

func findUnsafe(value any) (string, string) {

	return findUnsafeAt("", value)
}

func findUnsafeRawEvent(value any) (string, string) {

	return findUnsafeRawEventAt("", value)
}

func findUnsafeRawEventAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, true)
}

func findUnsafeAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, false)
}

func findUnsafeValueAt(path string, value any, rawEvent bool) (string, string) {
	// findUnsafeValueAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:
		return findUnsafeMapAt(path, v, rawEvent)
	case []any:
		return findUnsafeSliceAt(path, v, rawEvent)
	case string:
		return findUnsafeStringAt(path, v, rawEvent)
	}

	return "", ""
}
func resolveSessionSetupProfilePath(profilePath string) (string, error) {
	// resolveSessionSetupProfilePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(profilePath) == "" {
		return "", errors.New("observe setup requires --profile")
	}

	safePath, err := safeExistingFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("unsafe profile path: %w", err)
	}
	return safePath, nil
}

func resolveSessionSetupOutDir(outDir string) (string, error) {
	// resolveSessionSetupOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(outDir) == "" {
		return "", errors.New("observe setup requires --out")
	}

	return safeOutDir(outDir)
}
func setupSessionRun(profilePath, outDir string, now time.Time, rawCommand string) (SessionRun, error) {
	// setupSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := prepareSessionRun(profilePath, outDir)
	if err != nil {
		return SessionRun{}, err
	}

	run := newSessionRunWithCommand(profile, now, rawCommand)
	results, err := installIsolationRules(profilePath, profile.IsolationRules)
	if err != nil {
		return SessionRun{}, err
	}
	run.IsolationResults = results

	if err := writeSessionJSON(filepath.Join(outDir, "session.json"), run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

func prepareSessionRun(profilePath, outDir string) (SessionProfile, error) {
	// prepareSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionProfile{}, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}

func newSessionRunWithCommand(profile SessionProfile, now time.Time, rawCommand string) SessionRun {
	// newSessionRunWithCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	run := newSessionRun(profile, sessionRunTime(now))
	setSessionCommand(&run, rawCommand)
	return run
}

func setSessionCommand(run *SessionRun, rawCommand string) {
	// setSessionCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(rawCommand) == "" {
		return
	}

	command := []string{rawCommand}
	run.CommandDigest = digestCommand(command)
	run.CommandDigestState = StatePass
	if model := extractCommandModel(command); model != "" {
		run.CommandModel = model
		run.CommandModelState = StatePass
	}
}

func sessionRunTime(now time.Time) time.Time {
	// sessionRunTime keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if now.IsZero() {

		return time.Now().UTC()
	}
	return now
}

func writeSessionJSON(path string, run SessionRun) error {
	return writeJSON(path, run)
}

func CollectSession(opts SessionCollectOptions) (SessionRun, Run, error) {
	// CollectSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ctx, err := prepareSessionCollection(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	sourcePath, err := resolveSessionEventSource(&ctx)
	if err != nil {
		if !errors.Is(err, errSessionSourceUnavailable) {
			return SessionRun{}, Run{}, err
		}
		return markSessionSourceUnavailable(ctx)
	}
	return collectSessionSource(ctx, sourcePath)
}

func prepareSessionCollection(opts SessionCollectOptions) (sessionCollectionContext, error) {
	// prepareSessionCollection keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, err := validateSessionCollectOptions(opts)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	return loadSessionCollection(profilePath, runDir, opts.Now)
}

func loadSessionCollection(profilePath, runDir string, now time.Time) (sessionCollectionContext, error) {
	// loadSessionCollection keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, session, err := loadSessionCollectionInputs(profilePath, runDir)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	harnessProfilePath, harnessProfile, err := loadHarnessProfile(profilePath, profile)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	now = sessionCollectionTime(now)
	return newSessionCollectionContext(profilePath, runDir, now, profile, session, harnessProfilePath, harnessProfile), nil
}

func newSessionCollectionContext(profilePath, runDir string, now time.Time, profile SessionProfile, session SessionRun, harnessProfilePath string, harnessProfile Profile) sessionCollectionContext {
	// newSessionCollectionContext keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return sessionCollectionContext{
		profilePath:        profilePath,
		runDir:             runDir,
		now:                now,
		profile:            profile,
		session:            session,
		harnessProfile:     harnessProfile,
		harnessProfilePath: harnessProfilePath,
	}
}

func sessionCollectionTime(now time.Time) time.Time {
	// sessionCollectionTime keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if now.IsZero() {

		return time.Now().UTC()
	}
	return now
}
func loadSessionCollectionInputs(profilePath, runDir string) (SessionProfile, SessionRun, error) {
	// loadSessionCollectionInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionProfile{}, SessionRun{}, err
	}

	session, err := LoadSessionRun(filepath.Join(runDir, "session.json"))
	if err != nil {
		return SessionProfile{}, SessionRun{}, err
	}
	if session.ProfileID != profile.ProfileID {

		return SessionProfile{}, SessionRun{}, errors.New("session profile mismatch")
	}

	return profile, session, nil
}
func findUnsafeMapAt(path string, values map[string]any, rawEvent bool) (string, string) {
	// findUnsafeMapAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, child := range values {

		if field, reason := findUnsafeMapChild(path, key, child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}

func findUnsafeMapChild(path, key string, child any, rawEvent bool) (string, string) {
	// findUnsafeMapChild keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	childPath := childPath(path, key)
	reason, skip := unsafeMapFieldReason(childPath, strings.ToLower(key), child, rawEvent)
	if reason != "" {
		return childPath, reason
	}
	if skip {

		return "", ""
	}
	return findUnsafeValueAt(childPath, child, rawEvent)
}

func childPath(parent, key string) string {
	// childPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if parent == "" {

		return key
	}
	return parent + "." + key
}

func unsafeMapFieldReason(path, key string, value any, rawEvent bool) (string, bool) {
	// unsafeMapFieldReason keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if skippableRawEventField(path, key, value, rawEvent) {
		return "", true
	}
	if rawFieldNames[key] {

		return "forbidden_raw_field", false
	}
	if sensitiveFieldNames[key] {
		return "sensitive_field", false
	}
	return "", false
}

func skippableRawEventField(path, key string, value any, rawEvent bool) bool {
	// skippableRawEventField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return rawEvent &&
		(unretainedRawToolInputField(path, key, value) ||
			(unretainedRawBodyField(key) && !structuredRawBody(value)))
}

func unretainedRawToolInputField(path, key string, value any) bool {
	// unretainedRawToolInputField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if key != "prompt" {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}

	segments := strings.Split(path, ".")
	if len(segments) < 3 {
		return false
	}
	return path == "part.state.input.prompt"
}

func unretainedRawBodyField(key string) bool {
	// unretainedRawBodyField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch key {
	case "text", "content", "input", "output", "stdout", "stderr":

		return true
	default:
		return false
	}
}

func structuredRawBody(value any) bool {
	// structuredRawBody keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch value.(type) {
	case map[string]any:

		return true
	default:
		return false
	}
}

func findUnsafeSliceAt(path string, values []any, rawEvent bool) (string, string) {
	// findUnsafeSliceAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i, child := range values {

		if field, reason := findUnsafeValueAt(fmt.Sprintf("%s[%d]", path, i), child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}

func findUnsafeStringAt(path, value string, rawEvent bool) (string, string) {
	// findUnsafeStringAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return "", ""
	}
	if reason := unsafeStringReason(path, value, rawEvent); reason != "" {
		return path, reason
	}
	return "", ""
}
func unsafeStringReason(path, value string, rawEvent bool) string {
	// unsafeStringReason keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeStringPath(value, rawEvent) {
		return "unsafe_path_or_private_path"
	}
	if unsafeURL(value) {

		return "authenticated_url"
	}
	if unsafeStringToken(path, value, rawEvent) {
		return "token_like_value"
	}
	return ""
}

func unsafeStringToken(path, value string, rawEvent bool) bool {
	return providerTokenPrefix.MatchString(value) || unsafeEncodedToken(path, value, rawEvent)
}

func unsafeStringPath(value string, rawEvent bool) bool {
	return !rawEvent && unsafePathValue(value)
}

func unsafePathValue(value string) bool {
	return privatePathPattern.MatchString(value) || strings.HasPrefix(value, "/") || strings.Contains(value, "../")
}

func unsafeEncodedToken(path, value string, rawEvent bool) bool {
	// unsafeEncodedToken keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeEncodedTokenExemption(path, value, rawEvent) {
		return false
	}

	return base64TokenPattern.MatchString(value)
}

func safeEncodedTokenExemption(path, value string, rawEvent bool) bool {
	return digestField(path) || sha256Pattern.MatchString(value) || rawEventPathLikeField(path, rawEvent)
}

func rawEventPathLikeField(path string, rawEvent bool) bool {
	return rawEvent && rawPathLikeField(path)
}

func rawPathLikeField(path string) bool {
	// rawPathLikeField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch rawPathFieldName(path) {
	case "path", "file", "filepath", "file_path", "dir", "directory", "cwd":

		return true
	default:
		return false
	}
}

func rawPathFieldName(path string) string {
	// rawPathFieldName keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	field := path
	if idx := strings.LastIndex(field, "."); idx >= 0 {
		field = field[idx+1:]
	}
	if idx := strings.LastIndex(field, "["); idx >= 0 {

		field = field[:idx]
	}
	return strings.ToLower(field)
}
func unsafeURL(raw string) bool {
	// unsafeURL keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return false
	}
	if parsed.User != nil {

		return true
	}
	return queryHasAuthKey(parsed.Query())
}

func queryHasAuthKey(values url.Values) bool {
	// queryHasAuthKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key := range values {

		if authQueryKeys[strings.ToLower(key)] {
			return true
		}
	}
	return false
}

func digestField(path string) bool {
	// digestField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	last := path
	if idx := strings.LastIndex(last, "."); idx >= 0 {

		last = last[idx+1:]
	}
	return digestFieldNames[last]
}

func validFamily(family string) bool {
	return validFamilies[family]
}

func validState(state string) bool {
	return validStates[state]
}

func validContentState(state string) bool {
	return validContentStates[state]
}

func validRuleKey(key string) bool {
	return validRuleKeys[key]
}

func safeRef(ref string) bool {
	return safeIDPattern.MatchString(ref) || sha256Pattern.MatchString(ref)
}

func safeOperationRef(ref string) bool {
	// safeOperationRef keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if operationRefPrefix(ref) {

		return safePrefixedOperationRef(ref)
	}
	return safeRef(ref)
}

func operationRefPrefix(ref string) bool {
	return strings.HasPrefix(ref, "adapter-run:") || strings.HasPrefix(ref, "delivery-trace:")
}

func safePrefixedOperationRef(ref string) bool {
	return !strings.Contains(ref, "..") && !strings.Contains(ref, "://") && len(ref) <= 256
}
func safeEvent(eventID string) string {
	// safeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeIDPattern.MatchString(eventID) {
		return eventID
	}

	return "unknown_event"
}

func nonAuthority() string {
	return "harness observation is evidence only; no harness compliance, feature delivery, PR approval, merge approval, release readiness, or production trust is claimed"
}

func DecodeValidation(r io.Reader) (Validation, error) {
	// DecodeValidation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var validation Validation
	if err := json.NewDecoder(r).Decode(&validation); err != nil {

		return Validation{}, err
	}
	return validation, nil
}
func validateSessionCollectOptions(opts SessionCollectOptions) (string, string, error) {
	// validateSessionCollectOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireSessionCollectOptions(opts); err != nil {
		return "", "", err
	}

	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return "", "", fmt.Errorf("unsafe profile path: %w", err)
	}

	runDir, err := safeExistingDir(opts.RunDir)
	if err != nil {
		return "", "", fmt.Errorf("unsafe run path: %w", err)
	}
	return profilePath, runDir, nil
}

func requireSessionCollectOptions(opts SessionCollectOptions) error {
	// requireSessionCollectOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "observe collect requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "observe collect requires --run"); err != nil {
		return err
	}
	return nil
}

func loadHarnessProfile(profilePath string, profile SessionProfile) (string, Profile, error) {
	// loadHarnessProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	harnessProfilePath, err := safeProfileRelativeFile(profilePath, profile.HarnessProfilePath)
	if err != nil {
		return "", Profile{}, fmt.Errorf("unsafe harness profile path: %w", err)
	}
	harnessProfile, err := LoadProfile(harnessProfilePath)
	if err != nil {
		return "", Profile{}, err
	}
	return harnessProfilePath, harnessProfile, nil
}

func resolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// resolveSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err == nil {
		return sourcePath, nil
	}
	return resolveMissingSessionEventSource(ctx)
}

func resolveMissingSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// resolveMissingSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if ctx.profile.RawEventFormat == "" {

		return "", errSessionSourceUnavailable
	}
	return normalizeAndResolveSessionEventSource(ctx)
}

func normalizeAndResolveSessionEventSource(ctx *sessionCollectionContext) (string, error) {
	// normalizeAndResolveSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if normalizeErr := normalizeSessionRawEvents(ctx); normalizeErr != nil {
		return "", normalizeErr
	}
	sourcePath, ok := resolvedSessionEventSource(ctx)
	if !ok {
		return "", errSessionSourceUnavailable
	}
	return sourcePath, nil
}

func resolvedSessionEventSource(ctx *sessionCollectionContext) (string, bool) {
	// resolvedSessionEventSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	return sourcePath, err == nil
}

func normalizeSessionRawEvents(ctx *sessionCollectionContext) error {
	// normalizeSessionRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	rawPath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.RawEventSourcePath)
	if err != nil {
		return fmt.Errorf("raw_event_source_path invalid: %w", err)
	}

	normalizedPath, err := safeProfileRelativeOutFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err != nil {
		return err
	}
	if err := normalizeRawEvents(ctx.profile.RawEventFormat, rawPath, normalizedPath, sessionCommandFacts(ctx.session), ctx.now); err != nil {
		return err
	}

	ctx.session.NormalizedDigest = digestFile(normalizedPath)
	return nil
}

func markSessionSourceUnavailable(ctx sessionCollectionContext) (SessionRun, Run, error) {
	// markSessionSourceUnavailable keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := unavailableSession(ctx)

	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}

	return session, unavailableObservedRun(ctx), nil
}
func unavailableSession(ctx sessionCollectionContext) SessionRun {
	// unavailableSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := ctx.session

	session.CollectionState = StateCannotVerify
	session.CollectionReason = "source_unavailable"
	session.EndTime = ctx.now.Format(time.RFC3339)
	return session
}

func unavailableObservedRun(ctx sessionCollectionContext) Run {
	// unavailableObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.harnessProfile.ProfileID,
		HarnessFamily:      ctx.harnessProfile.HarnessFamily,
		EventSchemaVersion: ctx.harnessProfile.EventSchemaVersion,

		SourcePath: filepath.Base(ctx.profile.EventSourcePath),
		EventCount: 0,
		CreatedAt:  ctx.now.Format(time.RFC3339),
	}
}

func collectSessionSource(ctx sessionCollectionContext, sourcePath string) (SessionRun, Run, error) {
	// collectSessionSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedDir := filepath.Join(ctx.runDir, "observed")
	if err := os.MkdirAll(observedDir, 0o755); err != nil {
		return SessionRun{}, Run{}, err
	}

	events, sourceDigest, err := readEventsFromPath(ctx.harnessProfilePath, sourcePath)
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	observed := observedRun(ctx, sourcePath, sourceDigest, events)
	if err := writeObservedRun(observedDir, events, observed); err != nil {
		return SessionRun{}, Run{}, err
	}
	return finalizeCollectedSession(ctx, observedDir, observed)
}

func writeObservedRun(observedDir string, events []Event, observed Run) error {
	// writeObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := writeObservedEvents(observedDir, events); err != nil {
		return err
	}
	return writeJSON(filepath.Join(observedDir, "run.json"), observed)
}

func writeObservedEvents(observedDir string, events []Event) error {
	// writeObservedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, event := range events {

		if err := writeJSON(filepath.Join(observedDir, "events", event.EventID+".json"), event); err != nil {
			return err
		}
	}
	return nil
}
func observedRun(ctx sessionCollectionContext, sourcePath, sourceDigest string, events []Event) Run {
	// observedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.harnessProfile.ProfileID,
		HarnessFamily:      ctx.harnessProfile.HarnessFamily,
		EventSchemaVersion: ctx.harnessProfile.EventSchemaVersion,
		SourcePath:         filepath.Base(sourcePath),
		SourceDigest:       sourceDigest,
		EventCount:         len(events),
		EventRefs:          eventRefs(events),
		CreatedAt:          ctx.now.Format(time.RFC3339),
	}
}

func finalizeCollectedSession(ctx sessionCollectionContext, observedDir string, observed Run) (SessionRun, Run, error) {
	// finalizeCollectedSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := ctx.session
	session.ObservedRunDir = filepath.ToSlash("observed")
	session.OutputDigest = digestFile(filepath.Join(observedDir, "run.json"))

	session.CollectionState = StatePass
	session.CollectionReason = "source_collected"

	if session.EndTime == "" {
		session.EndTime = ctx.now.Format(time.RFC3339)
	}
	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}
	return session, observed, nil
}

func RunSession(opts SessionOptions) (SessionRun, Run, error) {
	// RunSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session, err := setupRunnableSession(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	commandResult, err := runObservedCommand(opts.Command, &session)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	return collectFinishedSession(opts, session, commandResult.waitErr, commandResult.end)
}

func setupRunnableSession(opts SessionOptions) (SessionRun, error) {
	// setupRunnableSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireSessionCommand(opts.Command); err != nil {
		return SessionRun{}, err
	}

	return SetupSession(SessionSetupOptions{ProfilePath: opts.ProfilePath, OutDir: opts.OutDir, Now: opts.Now})
}

func runObservedCommand(command []string, session *SessionRun) (observedCommandResult, error) {
	// runObservedCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cmd := discardedCommand(command)
	setSessionProcessCommand(session, command, time.Now().UTC())

	if err := startSessionProcess(cmd, session); err != nil {
		return observedCommandResult{}, err
	}
	waitErr := cmd.Wait()

	end := time.Now().UTC()
	return observedCommandResult{waitErr: waitErr, end: end}, nil
}

func requireSessionCommand(command []string) error {
	// requireSessionCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if len(command) == 0 {

		return errors.New("observe session requires command after --")
	}
	return nil
}
func collectFinishedSession(opts SessionOptions, session SessionRun, waitErr error, end time.Time) (SessionRun, Run, error) {
	// collectFinishedSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := writeFinishedSession(opts.OutDir, &session, end); err != nil {
		return SessionRun{}, Run{}, err
	}

	collected, observed, err := CollectSession(SessionCollectOptions{ProfilePath: opts.ProfilePath, RunDir: opts.OutDir, Now: end})
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	if waitErr != nil {
		return collected, observed, waitErr
	}
	return collected, observed, nil
}

func writeFinishedSession(outDir string, session *SessionRun, end time.Time) error {
	session.EndTime = end.Format(time.RFC3339)
	return writeJSON(filepath.Join(outDir, "session.json"), *session)
}

func discardedCommand(command []string) *exec.Cmd {
	// discardedCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd
}

func setSessionProcessCommand(session *SessionRun, command []string, start time.Time) {
	// setSessionProcessCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	session.CommandDigest = digestCommand(command)
	session.CommandDigestState = StatePass
	if model := extractCommandModel(command); model != "" {
		session.CommandModel = model
		session.CommandModelState = StatePass
	}
	session.StartTime = start.Format(time.RFC3339)
}

func startSessionProcess(cmd *exec.Cmd, session *SessionRun) error {
	// startSessionProcess keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := cmd.Start(); err != nil {
		return err
	}

	session.ProcessID = cmd.Process.Pid
	session.ProcessIDState = StatePass
	return nil
}

func Validate(opts ValidateOptions) (Validation, error) {
	// Validate keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, outPath, err := validateValidateInputs(opts)
	if err != nil {
		return Validation{}, err
	}

	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Validation{}, err
	}
	validation := evaluationFromRun(profile, runDir)

	if err := writeValidationIfRequested(outPath, validation); err != nil {
		return Validation{}, err
	}
	return validation, nil
}

func validateValidateInputs(opts ValidateOptions) (string, string, string, error) {
	// validateValidateInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireValidateOptions(opts); err != nil {
		return "", "", "", err
	}

	return resolveValidateInputs(opts)
}

func requireValidateOptions(opts ValidateOptions) error {
	// requireValidateOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "harness validate requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "harness validate requires --run"); err != nil {
		return err
	}
	return nil
}

func requireNonBlank(value, message string) error {
	// requireNonBlank keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return errors.New(message)
	}
	return nil
}
func resolveValidateInputs(opts ValidateOptions) (string, string, string, error) {
	// resolveValidateInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, err := resolveValidateSourcePaths(opts)
	if err != nil {
		return "", "", "", err
	}

	outPath, err := resolveValidateOutPath(opts.OutPath)
	if err != nil {
		return "", "", "", err
	}
	return profilePath, runDir, outPath, nil
}

func resolveValidateSourcePaths(opts ValidateOptions) (string, string, error) {
	// resolveValidateSourcePaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return "", "", fmt.Errorf("unsafe profile path: %w", err)
	}
	runDir, err := safeExistingDir(opts.RunDir)
	if err != nil {
		return "", "", fmt.Errorf("unsafe run path: %w", err)
	}
	return profilePath, runDir, nil
}

func resolveValidateOutPath(outPath string) (string, error) {
	// resolveValidateOutPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if outPath == "" {

		return "", nil
	}
	safeOut, err := safeOutFile(outPath)
	if err != nil {
		return "", fmt.Errorf("unsafe out path: %w", err)
	}
	return safeOut, nil
}
func evaluationFromRun(profile Profile, runDir string) Validation {
	// evaluationFromRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	run, events, err := LoadRun(runDir)
	if err != nil {

		return fallbackSourceUnavailable(profile)
	}
	return evaluate(profile, run, events)
}

func fallbackSourceUnavailable(profile Profile) Validation {
	// fallbackSourceUnavailable keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,
		ValidationState:    StateCannotVerify,
		ReasonCode:         "source_unavailable",
		NonAuthority:       nonAuthority(),
	}
	validation.ValidationDigest = validationDigest(validation)

	return validation
}

func writeValidationIfRequested(outPath string, validation Validation) error {
	// writeValidationIfRequested keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if outPath == "" {

		return nil
	}
	return writeJSON(outPath, validation)
}

func LoadProfile(path string) (Profile, error) {
	// LoadProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var profile Profile

	if err := readExistingJSON(path, &profile); err != nil {
		return Profile{}, err
	}
	if err := validateProfile(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

func LoadSessionProfile(path string) (SessionProfile, error) {
	// LoadSessionProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var profile SessionProfile

	if err := readExistingJSONStrict(path, &profile); err != nil {
		return SessionProfile{}, err
	}
	if err := validateSessionProfile(&profile); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}

func validateSessionProfile(profile *SessionProfile) error {
	// validateSessionProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := validateSessionProfileIdentity(*profile); err != nil {
		return err
	}

	if err := normalizeSessionStreamCapture(profile); err != nil {
		return err
	}
	if err := validateSessionSetupActions(profile.SetupActions); err != nil {
		return err
	}
	return validateSessionIsolationRules(profile.IsolationRules)
}

func validateSessionProfileIdentity(profile SessionProfile) error {
	// validateSessionProfileIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.SchemaVersion != SessionProfileSchemaVersion {
		return fmt.Errorf("unsupported session profile schema_version %q", profile.SchemaVersion)
	}

	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return validateSessionProfilePaths(profile)
}

func validateSessionProfilePaths(profile SessionProfile) error {
	// validateSessionProfilePaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateRequiredSessionPaths(profile); err != nil {
		return err
	}

	return validateRawEventConfig(profile)
}

func validateRequiredSessionPaths(profile SessionProfile) error {
	// validateRequiredSessionPaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(profile.HarnessProfilePath, "session profile requires harness_profile_path"); err != nil {
		return err
	}

	if err := requireNonBlank(profile.EventSourcePath, "session profile requires event_source_path"); err != nil {
		return err
	}
	return nil
}

func validateRawEventConfig(profile SessionProfile) error {
	// validateRawEventConfig keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	hasFormat := profile.RawEventFormat != ""
	hasSource := strings.TrimSpace(profile.RawEventSourcePath) != ""

	if unsupportedRawEventFormat(profile.RawEventFormat) {
		return errors.New("unsupported raw_event_format")
	}

	return validateRawEventPair(hasFormat, hasSource)
}

func validateRawEventPair(hasFormat, hasSource bool) error {
	// validateRawEventPair keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	switch {
	case hasFormat == hasSource:
		return nil
	case hasFormat:
		return errors.New("raw_event_source_path required for raw_event_format")
	default:
		return errors.New("raw_event_format required for raw_event_source_path")
	}
}

func unsupportedRawEventFormat(format string) bool {
	return format != "" && format != OpenCodeJSONLRawFormat
}
func normalizeSessionStreamCapture(profile *SessionProfile) error {
	// normalizeSessionStreamCapture keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	defaultSessionStreamCapture(profile)
	if profile.StreamCapture == "disabled" {
		return nil
	}

	return unsupportedSessionStreamCapture(profile.StreamCapture)
}

func defaultSessionStreamCapture(profile *SessionProfile) {
	// defaultSessionStreamCapture keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.StreamCapture == "" {

		profile.StreamCapture = "disabled"
	}
}

func unsupportedSessionStreamCapture(mode string) error {
	// unsupportedSessionStreamCapture keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch mode {
	case ContentDigestOnly, ContentRetainedSafe:

		return errors.New("stream_capture mode not implemented")
	default:
		return errors.New("unsupported stream_capture")
	}
}

func validateSessionSetupActions(actions []SessionSetupAction) error {
	// validateSessionSetupActions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if len(actions) > 3 {
		return errors.New("too many setup actions")
	}

	for _, action := range actions {
		if err := validateSessionSetupAction(action); err != nil {
			return err
		}
	}
	return nil
}

func validateSessionSetupAction(action SessionSetupAction) error {
	// validateSessionSetupAction keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeIDPattern.MatchString(action.ID) {
		return errors.New("unsafe setup action id")
	}

	switch action.Kind {
	case "init", "profile", "wrapper", "hook", "context_isolation":
		return nil
	default:
		return errors.New("unsupported setup action kind")
	}
}

func validateSessionIsolationRules(rules []SessionIsolationRule) error {
	// validateSessionIsolationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, rule := range rules {

		if err := validateSessionIsolationRule(rule); err != nil {
			return err
		}
	}
	return nil
}

func validateSessionIsolationRule(rule SessionIsolationRule) error {
	// validateSessionIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeIDPattern.MatchString(rule.ID) {
		return errors.New("unsafe isolation rule id")
	}

	if err := validateIsolationRulePattern(rule.Pattern); err != nil {
		return err
	}
	if unsafeProfileRelativePath(rule.TargetPath) {
		return errors.New("unsafe isolation target path")
	}
	return validateIsolationRuleKind(rule.Kind)
}

func validateIsolationRulePattern(pattern string) error {
	// validateIsolationRulePattern keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeIsolationRulePattern(pattern) {

		return errors.New("unsafe isolation rule pattern")
	}
	return nil
}

func unsafeIsolationRulePattern(pattern string) bool {
	return strings.TrimSpace(pattern) == "" || strings.Contains(pattern, "\n") || strings.Contains(pattern, "\r")
}

func validateIsolationRuleKind(kind string) error {
	// validateIsolationRuleKind keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	switch kind {
	case "ignore_line", "json_read_deny":
		return nil
	default:
		return errors.New("unsupported isolation rule kind")
	}
}

func installIsolationRules(profilePath string, rules []SessionIsolationRule) ([]SessionIsolationResult, error) {
	// installIsolationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	results := make([]SessionIsolationResult, 0, len(rules))
	for _, rule := range rules {

		resolvedRule, err := resolveIsolationRuleTarget(profilePath, rule)
		if err != nil {
			return nil, err
		}

		result, err := installIsolationRule(resolvedRule)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func resolveIsolationRuleTarget(profilePath string, rule SessionIsolationRule) (SessionIsolationRule, error) {
	// resolveIsolationRuleTarget keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	targetPath, err := safeProfileRelativeIsolationFile(profilePath, rule.TargetPath)
	if err != nil {
		return SessionIsolationRule{}, err
	}
	rule.TargetPath = targetPath
	return rule, nil
}
func safeProfileRelativeIsolationFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeIsolationFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative isolation path must be local without traversal")
	}

	clean := cleanProfileRelativePath(profilePath, relPath)
	if err := validateIsolationParent(clean); err != nil {
		return "", err
	}
	if err := validateIsolationFilename(filepath.Base(clean)); err != nil {
		return "", err
	}
	return clean, nil
}

func cleanProfileRelativePath(profilePath, relPath string) string {
	// cleanProfileRelativePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		baseDir = ""
	}
	return filepath.Clean(filepath.Join(baseDir, relPath))
}

func validateIsolationParent(clean string) error {
	// validateIsolationParent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	if err := validatePotentialParentPath(parent); err != nil {
		return err
	}

	return ensureOutParentInsideWorkingDirectory(parent)
}

func validateIsolationFilename(base string) error {
	// validateIsolationFilename keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(base) == "" || strings.ContainsAny(base, `/\`) {

		return errors.New("unsafe isolation filename")
	}
	return nil
}
func installIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	// installIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := ensureIsolationRule(rule); err != nil {
		return SessionIsolationResult{}, err
	}

	return verifyIsolationRule(rule)
}

func ensureIsolationRule(rule SessionIsolationRule) error {
	// ensureIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	installer, ok := isolationRuleInstallers[rule.Kind]
	if !ok {

		return errors.New("unsupported isolation rule kind")
	}
	return installer(rule.TargetPath, rule.Pattern)
}

func ensureLineFileRule(path, line string) error {
	// ensureLineFileRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lines, err := readOptionalLines(path)
	if err != nil {
		return err
	}

	for _, existing := range lines {
		if existing == line {
			return nil
		}
	}
	lines = append(lines, line)
	return writeLines(path, lines)
}

func readOptionalLines(path string) ([]string, error) {
	// readOptionalLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {

		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

func writeLines(path string, lines []string) error {
	// writeLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func ensureJSONReadDenyRule(path, pattern string) error {
	// ensureJSONReadDenyRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	config, err := readOptionalJSONObject(path)
	if err != nil {
		return err
	}
	setJSONReadDeny(config, pattern)
	return writeJSON(path, config)
}

func readOptionalJSONObject(path string) (map[string]any, error) {
	// readOptionalJSONObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := os.ReadFile(path)
	if err != nil {
		return optionalJSONObjectReadError(err)
	}

	return parseOptionalJSONObject(data)
}

func optionalJSONObjectReadError(err error) (map[string]any, error) {
	// optionalJSONObjectReadError keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if errors.Is(err, os.ErrNotExist) {

		return map[string]any{}, nil
	}
	return nil, err
}
func parseOptionalJSONObject(data []byte) (map[string]any, error) {
	// parseOptionalJSONObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	config := map[string]any{}
	if blankJSON(data) {

		return config, nil
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func blankJSON(data []byte) bool {
	return strings.TrimSpace(string(data)) == ""
}

func setJSONReadDeny(config map[string]any, pattern string) {
	// setJSONReadDeny keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	permission := ensureObject(config, "permission")
	read := ensureObject(permission, "read")

	read[pattern] = "deny"
}

func ensureObject(parent map[string]any, key string) map[string]any {
	// ensureObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if child, ok := parent[key].(map[string]any); ok {
		return child
	}

	child := map[string]any{}
	parent[key] = child
	return child
}

func verifyIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	// verifyIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	result := initialIsolationResult(rule)
	ok, err := isolationRulePresent(rule)
	if err != nil {
		return SessionIsolationResult{}, err
	}
	applyIsolationReadback(&result, ok)
	setIsolationDigest(&result, rule.TargetPath)
	return result, nil
}

func applyIsolationReadback(result *SessionIsolationResult, ok bool) {
	// applyIsolationReadback keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !ok {

		result.State = StateCannotVerify
		result.ReasonCode = "isolation_rule_absent"
	}
}

func setIsolationDigest(result *SessionIsolationResult, path string) {
	// setIsolationDigest keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if digest := digestFile(path); digest != "" {

		result.SHA256 = digest
	}
}

func initialIsolationResult(rule SessionIsolationRule) SessionIsolationResult {
	// initialIsolationResult keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return SessionIsolationResult{
		ID:         rule.ID,
		Kind:       rule.Kind,
		TargetPath: filepath.ToSlash(rule.TargetPath),
		Pattern:    rule.Pattern,
		State:      StatePass,
		ReasonCode: "isolation_rule_verified",
	}
}

func isolationRulePresent(rule SessionIsolationRule) (bool, error) {
	// isolationRulePresent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch rule.Kind {
	case "ignore_line":
		return lineIsolationRulePresent(rule.TargetPath, rule.Pattern)
	case "json_read_deny":
		return jsonReadDenyRulePresent(rule.TargetPath, rule.Pattern)
	default:

		return false, errors.New("unsupported isolation rule kind")
	}
}
func lineIsolationRulePresent(path, pattern string) (bool, error) {
	// lineIsolationRulePresent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lines, err := readOptionalLines(path)
	if err != nil {
		return false, err
	}

	for _, line := range lines {
		if line == pattern {
			return true, nil
		}
	}
	return false, nil
}

func jsonReadDenyRulePresent(path, pattern string) (bool, error) {
	// Decode only the subtree needed for this isolation proof.
	var config struct {
		Permission struct {
			Read map[string]string `json:"read"`
		} `json:"permission"`
	}
	if err := readExistingJSON(path, &config); err != nil {
		return false, err
	}
	return config.Permission.Read[pattern] == "deny", nil
}

func LoadSessionRun(path string) (SessionRun, error) {
	// LoadSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var run SessionRun

	if err := readExistingJSON(path, &run); err != nil {
		return SessionRun{}, err
	}
	if err := validateLoadedSessionRun(run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

func validateLoadedSessionRun(run SessionRun) error {
	// validateLoadedSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if run.SchemaVersion != SessionRunSchemaVersion {
		return fmt.Errorf("unsupported session schema_version %q", run.SchemaVersion)
	}

	if !safeIDPattern.MatchString(run.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return nil
}
func readExistingJSON(path string, target any) error {
	// readExistingJSON keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	safePath, err := safeExistingFile(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func readExistingJSONStrict(path string, target any) error {
	// readExistingJSONStrict keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	safePath, err := safeExistingFile(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}
	return decodeStrictJSON(data, target)
}

func decodeStrictJSON(data []byte, target any) error {
	// decodeStrictJSON keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	// A second decode distinguishes trailing JSON from ordinary EOF.
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("strict JSON input contains trailing data")
	}
	return nil
}

func newSessionRun(profile SessionProfile, now time.Time) SessionRun {
	// newSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	actionIDs := sessionSetupActionIDs(profile)
	commit, commitState := currentSourceCommitState()
	return newSessionRunRecord(profile, now, actionIDs, commit, commitState)
}

func newSessionRunRecord(profile SessionProfile, now time.Time, actionIDs []string, commit, commitState string) SessionRun {
	// newSessionRunRecord keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessProfilePath: profile.HarnessProfilePath,
		EventSourcePath:    profile.EventSourcePath,

		RawEventSourcePath: profile.RawEventSourcePath,
		RawEventFormat:     profile.RawEventFormat,
		SetupActionIDs:     actionIDs,

		CommandDigestState: StateCannotVerify,
		ProcessIDState:     StateCannotVerify,
		SourceCommit:       commit,
		SourceCommitState:  commitState,

		CollectionState:  StateCannotVerify,
		CollectionReason: "not_collected",
		CreatedAt:        now.Format(time.RFC3339),
	}
}

func sessionSetupActionIDs(profile SessionProfile) []string {
	// sessionSetupActionIDs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	actionIDs := make([]string, 0, len(profile.SetupActions))
	for _, action := range profile.SetupActions {

		actionIDs = append(actionIDs, action.ID)
	}

	sort.Strings(actionIDs)
	return actionIDs
}

func currentSourceCommitState() (string, string) {
	// currentSourceCommitState keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	commit := sourceCommit()
	if commit == "" {

		return "", StateCannotVerify
	}
	return commit, StatePass
}

func readEventsFromPath(profilePath, sourcePath string) ([]Event, string, error) {
	// readEventsFromPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	profile, err := LoadProfile(profilePath)
	if err != nil {
		return nil, "", err
	}
	return readEvents(profile, sourcePath)
}

func safeProfileRelativeFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeExistingFile(relPath)
	}
	return safeExistingFile(filepath.Join(baseDir, relPath))
}
func safeProfileRelativeOutFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeOutFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative output path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeOutFile(relPath)
	}
	return safeOutFile(filepath.Join(baseDir, relPath))
}

func unsafeProfileRelativePath(path string) bool {

	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}
func normalizeRawEvents(format, rawPath, outPath string, sessionFacts []Event, now time.Time) error {
	// normalizeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := validateRawNormalizationInputs(format, rawPath, outPath); err != nil {
		return err
	}
	events, err := normalizedOpenCodeRawEvents(rawPath, sessionFacts, rawNormalizationTime(now))
	if err != nil {
		return err
	}
	return writeNormalizedEvents(outPath, events)
}

func validateRawNormalizationInputs(format, rawPath, outPath string) error {
	// validateRawNormalizationInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if format != OpenCodeJSONLRawFormat {
		return errors.New("unsupported raw_event_format")
	}

	if filepath.Clean(rawPath) == filepath.Clean(outPath) {
		return errors.New("raw_event_source_path and event_source_path must be different files")
	}
	return nil
}

func rawNormalizationTime(now time.Time) time.Time {
	// rawNormalizationTime keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if now.IsZero() {

		return time.Now().UTC()
	}
	return now
}

func normalizedOpenCodeRawEvents(rawPath string, sessionFacts []Event, now time.Time) ([]Event, error) {
	// normalizedOpenCodeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	file, err := os.Open(rawPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return scanOpenCodeRawEvents(file, sessionFacts, now)
}

func scanOpenCodeRawEvents(file io.Reader, sessionFacts []Event, now time.Time) ([]Event, error) {
	// scanOpenCodeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), DefaultMaxLineBytes)
	lineNo := 0

	events := append([]Event{}, sessionFacts...)
	for scanner.Scan() {
		lineNo++
		var err error

		events, err = appendNormalizedRawLine(events, scanner.Bytes(), lineNo, now)
		if err != nil {
			return nil, err
		}
	}

	return events, scanner.Err()
}

func appendNormalizedRawLine(events []Event, line []byte, lineNo int, now time.Time) ([]Event, error) {
	// appendNormalizedRawLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lineEvents, err := normalizeOpenCodeRawLineBytes(line, lineNo, now)
	if err != nil {
		return nil, err
	}

	return append(events, lineEvents...), nil
}

func normalizeOpenCodeRawLineBytes(line []byte, lineNo int, now time.Time) ([]Event, error) {
	// normalizeOpenCodeRawLineBytes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if blankJSONLLine(line) {
		return nil, nil
	}
	// Decode as a generic map first so unsafe provider fields can be rejected
	// before typed event construction drops unknown data.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, fmt.Errorf("raw source line %d: malformed_jsonl", lineNo)
	}

	if err := rejectUnsafeRawEvent(raw, lineNo); err != nil {
		return nil, err
	}
	events := normalizeOpenCodeRawLine(raw, lineNo, now)
	return addNormalizedSourceDigests(events)
}

func rejectUnsafeRawEvent(raw map[string]any, lineNo int) error {
	// rejectUnsafeRawEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeField, reason := findUnsafeRawEvent(raw); unsafeField != "" {

		return fmt.Errorf("raw source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}

func addNormalizedSourceDigests(events []Event) ([]Event, error) {
	// addNormalizedSourceDigests keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i := range events {
		data, err := json.Marshal(events[i])
		if err != nil {
			return nil, err
		}

		events[i].SourceDigest = digestLine(data)
	}
	return events, nil
}

func blankJSONLLine(line []byte) bool {
	return len(strings.TrimSpace(string(line))) == 0
}

func writeNormalizedEvents(outPath string, events []Event) error {
	// writeNormalizedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	out, err := createNormalizedEventsFile(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, event := range events {

		if err := writeNormalizedEvent(out, event); err != nil {
			return err
		}
	}
	return nil
}
func createNormalizedEventsFile(outPath string) (*os.File, error) {
	// createNormalizedEventsFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return nil, err
	}
	return os.Create(outPath)
}

func writeNormalizedEvent(out io.Writer, event Event) error {
	// writeNormalizedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = out.Write(append(data, '\n'))
	return err
}
