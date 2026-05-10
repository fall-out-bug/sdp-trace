package harnessobs

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	SchemaVersion      string               `json:"schema_version"`
	ProfileID          string               `json:"profile_id"`
	HarnessProfilePath string               `json:"harness_profile_path"`
	EventSourcePath    string               `json:"event_source_path"`
	RawEventSourcePath string               `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string               `json:"raw_event_format,omitempty"`
	SetupActions       []SessionSetupAction `json:"setup_actions,omitempty"`
	StreamCapture      string               `json:"stream_capture"`
}

type SessionSetupAction struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
}

type SessionRun struct {
	SchemaVersion      string   `json:"schema_version"`
	ProfileID          string   `json:"profile_id"`
	HarnessProfilePath string   `json:"harness_profile_path"`
	EventSourcePath    string   `json:"event_source_path"`
	RawEventSourcePath string   `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string   `json:"raw_event_format,omitempty"`
	SetupActionIDs     []string `json:"setup_action_ids,omitempty"`
	CommandDigest      string   `json:"command_digest,omitempty"`
	CommandDigestState string   `json:"command_digest_state,omitempty"`
	CommandModel       string   `json:"command_model,omitempty"`
	CommandModelState  string   `json:"command_model_state,omitempty"`
	ProcessID          int      `json:"process_id,omitempty"`
	ProcessIDState     string   `json:"process_id_state,omitempty"`
	StartTime          string   `json:"start_time,omitempty"`
	EndTime            string   `json:"end_time,omitempty"`
	SourceCommit       string   `json:"source_commit,omitempty"`
	SourceCommitState  string   `json:"source_commit_state,omitempty"`
	ObservedRunDir     string   `json:"observed_run_dir,omitempty"`
	OutputDigest       string   `json:"output_digest,omitempty"`
	NormalizedDigest   string   `json:"normalized_digest,omitempty"`
	CollectionState    string   `json:"collection_state,omitempty"`
	CollectionReason   string   `json:"collection_reason,omitempty"`
	CreatedAt          string   `json:"created_at"`
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

func Observe(opts ObserveOptions) (Run, error) {
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

type observationContext struct {
	outDir       string
	sourcePath   string
	sourceDigest string
	now          time.Time
	profile      Profile
	events       []Event
}

func prepareObservation(opts ObserveOptions) (observationContext, error) {
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
	return observationContext{
		outDir:       outDir,
		sourcePath:   sourcePath,
		sourceDigest: sourceDigest,
		now:          observationTime(opts.Now),
		profile:      profile,
		events:       events,
	}, nil
}

func validateObserveOptions(opts ObserveOptions) (string, string, string, error) {
	if err := requireObserveOptions(opts); err != nil {
		return "", "", "", err
	}
	return resolveObservePaths(opts)
}

func requireObserveOptions(opts ObserveOptions) error {
	if strings.TrimSpace(opts.ProfilePath) == "" {
		return errors.New("harness observe requires --profile")
	}
	if strings.TrimSpace(opts.SourcePath) == "" {
		return errors.New("harness observe requires --source")
	}
	if strings.TrimSpace(opts.OutDir) == "" {
		return errors.New("harness observe requires --out")
	}
	return nil
}

func resolveObservePaths(opts ObserveOptions) (string, string, string, error) {
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
	if now.IsZero() {
		return time.Now().UTC()
	}
	return now
}

func loadObservationSource(profilePath, sourcePath string) (Profile, []Event, string, error) {
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
	for _, event := range events {
		path := filepath.Join(outDir, "events", event.EventID+".json")
		if err := writeJSON(path, event); err != nil {
			return err
		}
	}
	return nil
}

func newObservedRun(ctx observationContext) Run {
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

func resolveSessionSetupProfilePath(profilePath string) (string, error) {
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
	if strings.TrimSpace(outDir) == "" {
		return "", errors.New("observe setup requires --out")
	}
	return safeOutDir(outDir)
}

func setupSessionRun(profilePath, outDir string, now time.Time, rawCommand string) (SessionRun, error) {
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionRun{}, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return SessionRun{}, err
	}
	run := newSessionRunWithCommand(profile, now, rawCommand)
	if err := writeSessionJSON(filepath.Join(outDir, "session.json"), run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

func newSessionRunWithCommand(profile SessionProfile, now time.Time, rawCommand string) SessionRun {
	run := newSessionRun(profile, sessionRunTime(now))
	setSessionCommand(&run, rawCommand)
	return run
}

func setSessionCommand(run *SessionRun, rawCommand string) {
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
	if now.IsZero() {
		return time.Now().UTC()
	}
	return now
}

func writeSessionJSON(path string, run SessionRun) error {
	return writeJSON(path, run)
}

func CollectSession(opts SessionCollectOptions) (SessionRun, Run, error) {
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

type sessionCollectionContext struct {
	profilePath        string
	runDir             string
	now                time.Time
	profile            SessionProfile
	session            SessionRun
	harnessProfile     Profile
	harnessProfilePath string
}

func prepareSessionCollection(opts SessionCollectOptions) (sessionCollectionContext, error) {
	profilePath, runDir, err := validateSessionCollectOptions(opts)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	session, err := LoadSessionRun(filepath.Join(runDir, "session.json"))
	if err != nil {
		return sessionCollectionContext{}, err
	}
	if session.ProfileID != profile.ProfileID {
		return sessionCollectionContext{}, errors.New("session profile mismatch")
	}
	harnessProfilePath, harnessProfile, err := loadHarnessProfile(profilePath, profile)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return sessionCollectionContext{
		profilePath:        profilePath,
		runDir:             runDir,
		now:                now,
		profile:            profile,
		session:            session,
		harnessProfile:     harnessProfile,
		harnessProfilePath: harnessProfilePath,
	}, nil
}

func validateSessionCollectOptions(opts SessionCollectOptions) (string, string, error) {
	if strings.TrimSpace(opts.ProfilePath) == "" {
		return "", "", errors.New("observe collect requires --profile")
	}
	if strings.TrimSpace(opts.RunDir) == "" {
		return "", "", errors.New("observe collect requires --run")
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

func loadHarnessProfile(profilePath string, profile SessionProfile) (string, Profile, error) {
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
	sourcePath, err := safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err == nil {
		return sourcePath, nil
	}
	if ctx.profile.RawEventFormat == "" {
		return "", errSessionSourceUnavailable
	}
	if normalizeErr := normalizeSessionRawEvents(ctx); normalizeErr != nil {
		return "", normalizeErr
	}
	sourcePath, err = safeProfileRelativeFile(ctx.profilePath, ctx.profile.EventSourcePath)
	if err != nil {
		return "", errSessionSourceUnavailable
	}
	return sourcePath, nil
}

func normalizeSessionRawEvents(ctx *sessionCollectionContext) error {
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
	session := ctx.session
	session.CollectionState = StateCannotVerify
	session.CollectionReason = "source_unavailable"
	session.EndTime = ctx.now.Format(time.RFC3339)
	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}
	return session, Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.harnessProfile.ProfileID,
		HarnessFamily:      ctx.harnessProfile.HarnessFamily,
		EventSchemaVersion: ctx.harnessProfile.EventSchemaVersion,
		SourcePath:         filepath.Base(ctx.profile.EventSourcePath),
		EventCount:         0,
		CreatedAt:          ctx.now.Format(time.RFC3339),
	}, nil
}

func collectSessionSource(ctx sessionCollectionContext, sourcePath string) (SessionRun, Run, error) {
	observedDir := filepath.Join(ctx.runDir, "observed")
	if err := os.MkdirAll(observedDir, 0o755); err != nil {
		return SessionRun{}, Run{}, err
	}
	events, sourceDigest, err := readEventsFromPath(ctx.harnessProfilePath, sourcePath)
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	if err := writeObservedEvents(observedDir, events); err != nil {
		return SessionRun{}, Run{}, err
	}
	observed := observedRun(ctx, sourcePath, sourceDigest, events)
	if err := writeJSON(filepath.Join(observedDir, "run.json"), observed); err != nil {
		return SessionRun{}, Run{}, err
	}
	return finalizeCollectedSession(ctx, observedDir, observed)
}

func writeObservedEvents(observedDir string, events []Event) error {
	for _, event := range events {
		if err := writeJSON(filepath.Join(observedDir, "events", event.EventID+".json"), event); err != nil {
			return err
		}
	}
	return nil
}

func observedRun(ctx sessionCollectionContext, sourcePath, sourceDigest string, events []Event) Run {
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
	if len(opts.Command) == 0 {
		return SessionRun{}, Run{}, errors.New("observe session requires command after --")
	}
	session, err := SetupSession(SessionSetupOptions{ProfilePath: opts.ProfilePath, OutDir: opts.OutDir, Now: opts.Now})
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	start := time.Now().UTC()
	cmd := exec.Command(opts.Command[0], opts.Command[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	session.CommandDigest = digestCommand(opts.Command)
	session.CommandDigestState = StatePass
	if model := extractCommandModel(opts.Command); model != "" {
		session.CommandModel = model
		session.CommandModelState = StatePass
	}
	session.StartTime = start.Format(time.RFC3339)
	if err := cmd.Start(); err != nil {
		return SessionRun{}, Run{}, err
	}
	session.ProcessID = cmd.Process.Pid
	session.ProcessIDState = StatePass
	waitErr := cmd.Wait()
	end := time.Now().UTC()
	session.EndTime = end.Format(time.RFC3339)
	if err := writeJSON(filepath.Join(opts.OutDir, "session.json"), session); err != nil {
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

func Validate(opts ValidateOptions) (Validation, error) {
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
	if err := requireValidateOptions(opts); err != nil {
		return "", "", "", err
	}
	return resolveValidateInputs(opts)
}

func requireValidateOptions(opts ValidateOptions) error {
	if strings.TrimSpace(opts.ProfilePath) == "" {
		return errors.New("harness validate requires --profile")
	}
	if strings.TrimSpace(opts.RunDir) == "" {
		return errors.New("harness validate requires --run")
	}
	return nil
}

func resolveValidateInputs(opts ValidateOptions) (string, string, string, error) {
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
	run, events, err := LoadRun(runDir)
	if err != nil {
		return fallbackSourceUnavailable(profile)
	}
	return evaluate(profile, run, events)
}

func fallbackSourceUnavailable(profile Profile) Validation {
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
	if outPath == "" {
		return nil
	}
	return writeJSON(outPath, validation)
}

func LoadProfile(path string) (Profile, error) {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return Profile{}, err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return Profile{}, err
	}
	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return Profile{}, err
	}
	if err := validateProfile(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

func LoadSessionProfile(path string) (SessionProfile, error) {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return SessionProfile{}, err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return SessionProfile{}, err
	}
	var profile SessionProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return SessionProfile{}, err
	}
	if err := validateSessionProfile(&profile); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}

func validateSessionProfile(profile *SessionProfile) error {
	if profile.SchemaVersion != SessionProfileSchemaVersion {
		return fmt.Errorf("unsupported session profile schema_version %q", profile.SchemaVersion)
	}
	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	if err := validateSessionProfilePaths(*profile); err != nil {
		return err
	}
	if err := normalizeSessionStreamCapture(profile); err != nil {
		return err
	}
	return validateSessionSetupActions(profile.SetupActions)
}

func validateSessionProfilePaths(profile SessionProfile) error {
	if err := validateRequiredSessionPaths(profile); err != nil {
		return err
	}
	return validateRawEventConfig(profile)
}

func validateRequiredSessionPaths(profile SessionProfile) error {
	if strings.TrimSpace(profile.HarnessProfilePath) == "" {
		return errors.New("session profile requires harness_profile_path")
	}
	if strings.TrimSpace(profile.EventSourcePath) == "" {
		return errors.New("session profile requires event_source_path")
	}
	return nil
}

func validateRawEventConfig(profile SessionProfile) error {
	if profile.RawEventFormat != "" && profile.RawEventFormat != OpenCodeJSONLRawFormat {
		return errors.New("unsupported raw_event_format")
	}
	if profile.RawEventFormat != "" && strings.TrimSpace(profile.RawEventSourcePath) == "" {
		return errors.New("raw_event_source_path required for raw_event_format")
	}
	if profile.RawEventFormat == "" && strings.TrimSpace(profile.RawEventSourcePath) != "" {
		return errors.New("raw_event_format required for raw_event_source_path")
	}
	return nil
}

func normalizeSessionStreamCapture(profile *SessionProfile) error {
	if profile.StreamCapture == "" {
		profile.StreamCapture = "disabled"
	}
	switch profile.StreamCapture {
	case "disabled":
		return nil
	case ContentDigestOnly, ContentRetainedSafe:
		return errors.New("stream_capture mode not implemented")
	default:
		return errors.New("unsupported stream_capture")
	}
}

func validateSessionSetupActions(actions []SessionSetupAction) error {
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
	if !safeIDPattern.MatchString(action.ID) {
		return errors.New("unsafe setup action id")
	}
	switch action.Kind {
	case "init", "profile", "wrapper", "hook":
		return nil
	default:
		return errors.New("unsupported setup action kind")
	}
}

func LoadSessionRun(path string) (SessionRun, error) {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return SessionRun{}, err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return SessionRun{}, err
	}
	var run SessionRun
	if err := json.Unmarshal(data, &run); err != nil {
		return SessionRun{}, err
	}
	if run.SchemaVersion != SessionRunSchemaVersion {
		return SessionRun{}, fmt.Errorf("unsupported session schema_version %q", run.SchemaVersion)
	}
	if !safeIDPattern.MatchString(run.ProfileID) {
		return SessionRun{}, errors.New("unsafe session profile_id")
	}
	return run, nil
}

func newSessionRun(profile SessionProfile, now time.Time) SessionRun {
	actionIDs := make([]string, 0, len(profile.SetupActions))
	for _, action := range profile.SetupActions {
		actionIDs = append(actionIDs, action.ID)
	}
	sort.Strings(actionIDs)
	commit := sourceCommit()
	commitState := StatePass
	if commit == "" {
		commitState = StateCannotVerify
	}
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
		CollectionState:    StateCannotVerify,
		CollectionReason:   "not_collected",
		CreatedAt:          now.Format(time.RFC3339),
	}
}

func readEventsFromPath(profilePath, sourcePath string) ([]Event, string, error) {
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return nil, "", err
	}
	return readEvents(profile, sourcePath)
}

func safeProfileRelativeFile(profilePath, relPath string) (string, error) {
	if filepath.IsAbs(relPath) || strings.Contains(relPath, "://") || strings.Contains(relPath, "..") {
		return "", errors.New("profile relative path must be local without traversal")
	}
	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeExistingFile(relPath)
	}
	return safeExistingFile(filepath.Join(baseDir, relPath))
}

func safeProfileRelativeOutFile(profilePath, relPath string) (string, error) {
	if filepath.IsAbs(relPath) || strings.Contains(relPath, "://") || strings.Contains(relPath, "..") {
		return "", errors.New("profile relative output path must be local without traversal")
	}
	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeOutFile(relPath)
	}
	return safeOutFile(filepath.Join(baseDir, relPath))
}

func normalizeRawEvents(format, rawPath, outPath string, sessionFacts []Event, now time.Time) error {
	if format != OpenCodeJSONLRawFormat {
		return errors.New("unsupported raw_event_format")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if filepath.Clean(rawPath) == filepath.Clean(outPath) {
		return errors.New("raw_event_source_path and event_source_path must be different files")
	}
	events, err := normalizedOpenCodeRawEvents(rawPath, sessionFacts, now)
	if err != nil {
		return err
	}
	return writeNormalizedEvents(outPath, events)
}

func normalizedOpenCodeRawEvents(rawPath string, sessionFacts []Event, now time.Time) ([]Event, error) {
	file, err := os.Open(rawPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), DefaultMaxLineBytes)
	lineNo := 0
	events := append([]Event{}, sessionFacts...)
	for scanner.Scan() {
		lineNo++
		lineEvents, err := normalizeOpenCodeRawLineBytes(scanner.Bytes(), lineNo, now)
		if err != nil {
			return nil, err
		}
		events = append(events, lineEvents...)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func normalizeOpenCodeRawLineBytes(line []byte, lineNo int, now time.Time) ([]Event, error) {
	if len(strings.TrimSpace(string(line))) == 0 {
		return nil, nil
	}
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, fmt.Errorf("raw source line %d: malformed_jsonl", lineNo)
	}
	if unsafeField, reason := findUnsafeRawEvent(raw); unsafeField != "" {
		return nil, fmt.Errorf("raw source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	events := normalizeOpenCodeRawLine(raw, lineNo, now)
	for i := range events {
		data, err := json.Marshal(events[i])
		if err != nil {
			return nil, err
		}
		events[i].SourceDigest = digestLine(data)
	}
	return events, nil
}

func writeNormalizedEvents(outPath string, events []Event) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		if _, err := out.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return nil
}

func normalizeOpenCodeRawLine(raw map[string]any, lineNo int, now time.Time) []Event {
	signals := rawSignals(raw)
	families := openCodeFamilies(raw, signals)
	if len(families) == 0 {
		return nil
	}
	ordered := sortedFamilies(families)
	observedAt := openCodeObservedAt(raw, now)
	actor := openCodeActor(raw)
	sourceRef := fmt.Sprintf("raw-%06d", lineNo)
	events := make([]Event, 0, len(ordered))
	for _, family := range ordered {
		events = append(events, normalizedEvent(
			fmt.Sprintf("%s-%s", sourceRef, family),
			family,
			family+"_observed",
			observedAt,
			sourceRef,
			actor,
		))
	}
	return events
}

func openCodeFamilies(raw map[string]any, signals []string) map[string]bool {
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
	return hasKey(raw, "role") ||
		hasSignal(signals, "message", "response", "text") ||
		hasSignalPrefix(signals, "message.", "response.")
}

func openCodeToolFamily(raw map[string]any, signals []string) bool {
	return hasKey(raw, "tool", "tool_call", "toolcall") ||
		hasSignal(signals, "tool.call", "tool.result", "tool_use") ||
		hasSignalPrefix(signals, "tool.")
}

func openCodeMutationFamily(raw map[string]any, signals []string) bool {
	return hasSignal(signals, "file.write", "file.edit", "file.patch", "file.delete", "mutation") ||
		hasSignalPrefix(signals, "mutation.") ||
		nativeMutationTool(raw)
}

func openCodeTestFamily(signals []string) bool {
	return hasSignal(signals, "test.finished", "test.started", "test.passed", "test.failed") ||
		hasSignalPrefix(signals, "test.")
}

func openCodePhaseFamily(raw map[string]any, signals []string) bool {
	return hasKey(raw, "phase") ||
		hasSignal(signals, "phase") ||
		hasSignalPrefix(signals, "phase.", "gsd.", "gsd_")
}

func sortedFamilies(families map[string]bool) []string {
	ordered := make([]string, 0, len(families))
	for family := range families {
		ordered = append(ordered, family)
	}
	sort.Strings(ordered)
	return ordered
}

func openCodeObservedAt(raw map[string]any, now time.Time) string {
	observedAt := findTimestamp(raw)
	if observedAt == "" {
		return now.Format(time.RFC3339)
	}
	return observedAt
}

func openCodeActor(raw map[string]any) string {
	if model := findStringByKey(raw, "model", "model_id", "modelid"); model != "" {
		return safeToken(model)
	}
	if provider := findStringByKey(raw, "provider"); provider != "" {
		return safeToken(provider)
	}
	return "opencode"
}

func sessionCommandFacts(session SessionRun) []Event {
	if session.CommandModelState != StatePass || strings.TrimSpace(session.CommandModel) == "" {
		return nil
	}
	observedAt := session.StartTime
	if observedAt == "" {
		observedAt = session.CreatedAt
	}
	if _, err := time.Parse(time.RFC3339, observedAt); err != nil {
		observedAt = time.Now().UTC().Format(time.RFC3339)
	}
	event := normalizedEvent(
		"session-command-model",
		"model",
		"model_observed",
		observedAt,
		"session-command",
		safeToken(session.CommandModel),
	)
	data, err := json.Marshal(event)
	if err != nil {
		return nil
	}
	event.SourceDigest = digestLine(data)
	return []Event{event}
}

func normalizedEvent(id, family, eventType, observedAt, sourceRef, actor string) Event {
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
	switch v := value.(type) {
	case map[string]any:
		parts := make([]string, 0, len(v)*2)
		for key, child := range v {
			parts = append(parts, strings.ToLower(key))
			parts = append(parts, rawSignalsAt(key, child)...)
		}
		return parts
	case []any:
		parts := make([]string, 0, len(v))
		for _, child := range v {
			parts = append(parts, rawSignalsAt(parentKey, child)...)
		}
		return parts
	case string:
		if rawSignalValueKey(parentKey) {
			return []string{strings.ToLower(v)}
		}
		return nil
	default:
		return []string{strings.ToLower(fmt.Sprint(v))}
	}
}

func rawSignalValueKey(key string) bool {
	switch strings.ToLower(key) {
	case "type", "kind", "event", "event_type", "name", "phase", "role", "provider", "model", "model_id", "status", "tool", "action", "operation":
		return true
	default:
		return false
	}
}

func hasSignal(signals []string, values ...string) bool {
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
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return hasKeyIn(value, wanted)
}

func hasKeyIn(value any, wanted map[string]bool) bool {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if wanted[strings.ToLower(key)] || hasKeyIn(child, wanted) {
				return true
			}
		}
	case []any:
		for _, child := range v {
			if hasKeyIn(child, wanted) {
				return true
			}
		}
	}
	return false
}

func findStringByKey(value any, keys ...string) string {
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return findStringByKeyIn(value, wanted)
}

func findStringByKeyIn(value any, wanted map[string]bool) string {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if wanted[strings.ToLower(key)] {
				if s, ok := child.(string); ok && strings.TrimSpace(s) != "" {
					return s
				}
			}
			if s := findStringByKeyIn(child, wanted); s != "" {
				return s
			}
		}
	case []any:
		for _, child := range v {
			if s := findStringByKeyIn(child, wanted); s != "" {
				return s
			}
		}
	}
	return ""
}

func findTimestamp(raw map[string]any) string {
	for _, key := range []string{"time", "timestamp", "created_at", "observed_at"} {
		if value := findStringByKey(raw, key); value != "" {
			if parsed, err := time.Parse(time.RFC3339, value); err == nil {
				return parsed.UTC().Format(time.RFC3339)
			}
		}
		if value, ok := findNumberByKey(raw, key); ok {
			if observedAt := unixMillisTimestamp(value); observedAt != "" {
				return observedAt
			}
		}
	}
	return ""
}

func findNumberByKey(value any, keys ...string) (float64, bool) {
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return findNumberByKeyIn(value, wanted)
}

func findNumberByKeyIn(value any, wanted map[string]bool) (float64, bool) {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if wanted[strings.ToLower(key)] {
				switch n := child.(type) {
				case float64:
					return n, true
				case int:
					return float64(n), true
				}
			}
			if n, ok := findNumberByKeyIn(child, wanted); ok {
				return n, true
			}
		}
	case []any:
		for _, child := range v {
			if n, ok := findNumberByKeyIn(child, wanted); ok {
				return n, true
			}
		}
	}
	return 0, false
}

func unixMillisTimestamp(value float64) string {
	if value <= 0 {
		return ""
	}
	if value > 1_000_000_000_000 {
		return time.UnixMilli(int64(value)).UTC().Format(time.RFC3339)
	}
	if value > 1_000_000_000 {
		return time.Unix(int64(value), 0).UTC().Format(time.RFC3339)
	}
	return ""
}

func nativeMutationTool(raw map[string]any) bool {
	tool := strings.ToLower(findStringByKey(raw, "tool"))
	switch tool {
	case "edit", "write", "patch", "apply_patch", "update", "delete":
		return true
	default:
		return false
	}
}

func safeToken(value string) string {
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_', r == '.', r == ':', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
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

func LoadRun(dir string) (Run, []Event, error) {
	data, err := os.ReadFile(filepath.Join(dir, "run.json"))
	if err != nil {
		return Run{}, nil, err
	}
	var run Run
	if err := json.Unmarshal(data, &run); err != nil {
		return Run{}, nil, err
	}
	if run.SchemaVersion != RunSchemaVersion {
		return Run{}, nil, fmt.Errorf("unsupported run schema_version: %s", run.SchemaVersion)
	}
	events := make([]Event, 0, len(run.EventRefs))
	for _, ref := range run.EventRefs {
		if !safeEventRef(ref) {
			return Run{}, nil, fmt.Errorf("unsafe event ref: %s", ref)
		}
		data, err := os.ReadFile(filepath.Join(dir, ref))
		if err != nil {
			return Run{}, nil, err
		}
		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			return Run{}, nil, err
		}
		events = append(events, event)
	}
	return run, events, nil
}

func Summarize(validation Validation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Harness observation: %s (%s)\n", validation.ValidationState, validation.ReasonCode)
	fmt.Fprintf(&b, "Profile: %s\n", validation.ProfileID)
	fmt.Fprintf(&b, "Event schema: %s\n", validation.EventSchemaVersion)
	fmt.Fprintf(&b, "Events: %d\n", validation.EventCount)
	fmt.Fprintln(&b, "Dimensions:")
	for _, dim := range validation.Dimensions {
		required := "optional"
		if dim.Required {
			required = "required"
		}
		fmt.Fprintf(&b, "- %s [%s]: %s (%s), events=%d\n", dim.Family, required, dim.State, dim.ReasonCode, dim.EventCount)
	}
	fmt.Fprintf(&b, "Boundary: %s\n", nonAuthority())
	return b.String()
}

func LoadValidation(path string) (Validation, error) {
	safePath, err := safeExistingFile(path)
	if err != nil {
		return Validation{}, err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return Validation{}, err
	}
	var validation Validation
	if err := json.Unmarshal(data, &validation); err != nil {
		return Validation{}, err
	}
	if validation.SchemaVersion != ValidationSchemaVersion {
		return Validation{}, fmt.Errorf("unsupported validation schema_version: %s", validation.SchemaVersion)
	}
	return validation, nil
}

func validateProfile(profile Profile) error {
	if err := validateProfileMetadata(profile); err != nil {
		return err
	}
	if err := validateProfileEventFamilies(profile.RequiredEventFamilies, profile.OptionalEventFamilies); err != nil {
		return err
	}
	return validateProfileDegradationRules(profile.DegradationRules)
}

func validateProfileMetadata(profile Profile) error {
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
	if err := validateFamilyList(requiredEventFamilies); err != nil {
		return err
	}
	return validateFamilyList(optionalEventFamilies)
}

func validateFamilyList(families []string) error {
	for _, family := range families {
		if !validFamily(family) {
			return fmt.Errorf("unsupported event family: %s", family)
		}
	}
	return nil
}

func validateProfileDegradationRules(rules map[string]Rule) error {
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
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	maxLine := profile.Limits.MaxLineBytes
	if maxLine <= 0 {
		maxLine = DefaultMaxLineBytes
	}
	maxEvents := profile.Limits.MaxEvents
	if maxEvents <= 0 {
		maxEvents = DefaultMaxEvents
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)
	events := []Event{}
	sourceHash := sha256.New()
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Bytes()
		sourceHash.Write(line)
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}
		if len(events) >= maxEvents {
			return nil, "", fmt.Errorf("source line %d: event limit exceeded", lineNo)
		}
		event, err := parseEvent(profile, line, lineNo)
		if err != nil {
			return nil, "", err
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, "", err
	}
	return events, hex.EncodeToString(sourceHash.Sum(nil)), nil
}

func parseEvent(profile Profile, line []byte, lineNo int) (Event, error) {
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return Event{}, fmt.Errorf("source line %d: malformed_jsonl", lineNo)
	}
	if unsafeField, reason := findUnsafe(raw); unsafeField != "" {
		return Event{}, fmt.Errorf("source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	var event Event
	if err := json.Unmarshal(line, &event); err != nil {
		return Event{}, fmt.Errorf("source line %d: malformed_event", lineNo)
	}
	expected := digestLine(line)
	if event.SourceDigest != expected {
		return Event{}, fmt.Errorf("source line %d: source_digest_mismatch:%s", lineNo, safeEvent(event.EventID))
	}
	if err := validateEvent(profile, event); err != nil {
		return Event{}, fmt.Errorf("source line %d: %w", lineNo, err)
	}
	return event, nil
}

func validateEvent(profile Profile, event Event) error {
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
	if !safeFileIDPattern.MatchString(event.EventID) {
		return errors.New("unsafe event_id")
	}
	if event.EventSchemaVersion != profile.EventSchemaVersion {
		return errors.New("schema_version_mismatch")
	}
	if !validFamily(event.EventFamily) {
		return errors.New("unsupported event_family")
	}
	if !safeIDPattern.MatchString(event.EventType) {
		return errors.New("unsafe event_type")
	}
	if event.ObservedAt != "" {
		if _, err := time.Parse(time.RFC3339, event.ObservedAt); err != nil {
			return errors.New("invalid observed_at")
		}
	}
	return nil
}

func validateEventRefs(event Event) error {
	if !safeRef(event.SourceRef) {
		return errors.New("unsafe source_ref")
	}
	if !sha256Pattern.MatchString(event.SourceDigest) {
		return errors.New("invalid source_digest")
	}
	if event.TaskRef != "" && !safeRef(event.TaskRef) {
		return errors.New("unsafe task_ref")
	}
	if event.OperationRef != "" && !safeOperationRef(event.OperationRef) {
		return errors.New("unsafe operation_ref")
	}
	if event.ActorRef != "" && !safeRef(event.ActorRef) {
		return errors.New("unsafe actor_ref")
	}
	return nil
}

func validateEventContent(event Event) error {
	if !validContentState(event.ContentState) {
		return errors.New("invalid content_state")
	}
	return nil
}

func validateUnavailableFields(fields []UnavailableField) error {
	for _, field := range fields {
		if !validUnavailableField(field) {
			return errors.New("invalid unavailable_fields")
		}
	}
	return nil
}

func validUnavailableField(field UnavailableField) bool {
	return safeIDPattern.MatchString(field.Field) &&
		field.State == StateNotAssessed &&
		safeIDPattern.MatchString(field.ReasonCode)
}

func evaluate(profile Profile, run Run, events []Event) Validation {
	dimensions := []Dimension{}
	counts := map[string]int{}
	for _, event := range events {
		counts[event.EventFamily]++
	}
	required := map[string]bool{}
	for _, family := range profile.RequiredEventFamilies {
		required[family] = true
		dimensions = append(dimensions, dimension(family, true, counts[family]))
	}
	for _, family := range profile.OptionalEventFamilies {
		if required[family] {
			continue
		}
		dimensions = append(dimensions, dimension(family, false, counts[family]))
	}
	sort.Slice(dimensions, func(i, j int) bool {
		if dimensions[i].Required != dimensions[j].Required {
			return dimensions[i].Required
		}
		return dimensions[i].Family < dimensions[j].Family
	})
	state, reason := compose(dimensions)
	if run.EventSchemaVersion != profile.EventSchemaVersion {
		state, reason = StateCannotVerify, "schema_version_mismatch"
	}
	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,
		ValidationState:    state,
		ReasonCode:         reason,
		Dimensions:         dimensions,
		EventCount:         len(events),
		NonAuthority:       nonAuthority(),
	}
	validation.ValidationDigest = validationDigest(validation)
	return validation
}

func dimension(family string, required bool, count int) Dimension {
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
	switch state {
	case StateFail:
		return 4
	case StateCannotVerify:
		return 3
	case StateNotAssessed:
		return 2
	case StatePass:
		return 1
	default:
		return 0
	}
}

func safeExistingFile(path string) (string, error) {
	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local file without traversal",
		requireDir:     false,
		typeError:      "path must be a file",
	})
}

func safeExistingDir(path string) (string, error) {
	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local directory without traversal",
		requireDir:     true,
		typeError:      "path must be a directory",
	})
}

type existingPathSpec struct {
	traversalError string
	requireDir     bool
	typeError      string
}

func safeExistingPath(path string, spec existingPathSpec) (string, error) {
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
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return absolutePath(resolved)
}

func sanitizeExistingPath(path, traversalError string) (string, error) {
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New(traversalError)
	}
	return filepath.Clean(path), nil
}

func relativeWorkingDirectoryPath(abs string) (string, error) {
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
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("out must be a relative local file without traversal")
	}
	parent, err := safeParentDir(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	base := filepath.Base(path)
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	if !safeFileIDPattern.MatchString(stem) || strings.ContainsAny(base, `/\`) {
		return "", errors.New("unsafe output filename")
	}
	return filepath.Join(parent, base), nil
}

func safeParentDir(path string) (string, error) {
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
	path = defaultRelativeParentPath(path)
	if err := validatePotentialParentPath(path); err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

func defaultRelativeParentPath(path string) string {
	if path == "" {
		return "."
	}
	return path
}

func validatePotentialParentPath(path string) error {
	if filepath.IsAbs(path) {
		return errors.New("parent path must be relative local directory without traversal")
	}
	if strings.Contains(path, "://") {
		return errors.New("parent path must be relative local directory without traversal")
	}
	if strings.Contains(path, "..") {
		return errors.New("parent path must be relative local directory without traversal")
	}
	return nil
}

func resolveParentPathWithinWorkingDirectory(clean string) (string, error) {
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
	parent := filepath.Dir(clean)
	if parent == clean {
		return "", os.ErrNotExist
	}
	return resolveParentPathWithinWorkingDirectory(parent)
}

func ensurePathInsideWorkingDirectory(path string) (string, error) {
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
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, path)
}

func safeOutDir(path string) (string, error) {
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("out must be a relative local directory without traversal")
	}
	clean := filepath.Clean(path)
	if _, err := os.Lstat(clean); err == nil {
		return safeExistingOutDir(clean)
	} else {
		if parentEscapes, err := outParentEscapes(clean); err != nil {
			return "", err
		} else if parentEscapes {
			return "", errors.New("out parent path escapes working directory")
		}
	}
	return ensureOutDirEmptyOrMissing(clean)
}

func safeExistingOutDir(clean string) (string, error) {
	rel, err := relativeSymlinkTarget(clean)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {
		return "", errors.New("out path escapes working directory")
	}
	// Preserve the original behavior: when --out is a symlink, check the
	// resolved target directory for emptiness, not the symlink entry itself.
	return ensureOutDirEmptyOrMissing(rel)
}

func outParentEscapes(clean string) (bool, error) {
	parent := filepath.Dir(clean)
	for parent != "." && parent != string(filepath.Separator) {
		if _, statErr := os.Lstat(parent); statErr == nil {
			rel, err := relativeSymlinkTarget(parent)
			if err != nil {
				return false, err
			}
			return pathEscapesWorkingDirectory(rel), nil
		}
		parent = filepath.Dir(parent)
	}
	return false, nil
}

func relativeSymlinkTarget(path string) (string, error) {
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
	entries, err := os.ReadDir(clean)
	if err == nil && len(entries) > 0 {
		return "", errors.New("harness observe refuses existing non-empty --out")
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	return clean, nil
}

func writeJSON(path string, value any) error {
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
	refs := make([]string, 0, len(events))
	for _, event := range events {
		refs = append(refs, filepath.ToSlash(filepath.Join("events", event.EventID+".json")))
	}
	sort.Strings(refs)
	return refs
}

func safeEventRef(ref string) bool {
	if strings.Contains(ref, "\\") || strings.Contains(ref, "..") || filepath.IsAbs(ref) {
		return false
	}
	if !strings.HasPrefix(ref, "events/") || !strings.HasSuffix(ref, ".json") {
		return false
	}
	id := strings.TrimSuffix(strings.TrimPrefix(ref, "events/"), ".json")
	return safeFileIDPattern.MatchString(id)
}

func digestLine(line []byte) string {
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
	copy := validation
	copy.ValidationDigest = ""
	data, _ := json.Marshal(copy)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func digestCommand(command []string) string {
	data, _ := json.Marshal(command)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func extractCommandModel(command []string) string {
	if len(command) >= 3 && command[1] == "-c" {
		base := filepath.Base(command[0])
		if base == "sh" || base == "bash" {
			if model := extractCommandModelArgs(shellFields(command[2])); model != "" {
				return model
			}
		}
	}
	return extractCommandModelArgs(command)
}

func extractCommandModelArgs(args []string) string {
	for i, arg := range args {
		if arg == "--model" || arg == "-m" {
			if i+1 < len(args) {
				return safeCommandModel(args[i+1])
			}
			return ""
		}
		for _, prefix := range []string{"--model=", "-m="} {
			if strings.HasPrefix(arg, prefix) {
				return safeCommandModel(strings.TrimPrefix(arg, prefix))
			}
		}
	}
	return ""
}

// shellFields handles the shell field syntax needed to locate --model inside a
// controlled sh -c wrapper. It is not a general shell parser; model values still
// have to pass safeCommandModel before they become retained facts.
func shellFields(command string) []string {
	scanner := shellFieldScanner{}
	for _, r := range command {
		scanner.scan(r)
	}
	return scanner.finish()
}

type shellFieldScanner struct {
	fields  []string
	current strings.Builder
	quote   rune
	escaped bool
}

func (scanner *shellFieldScanner) scan(r rune) {
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
	switch {
	case r == '\'' || r == '"':
		scanner.quote = r
	case shellFieldSeparator(r):
		scanner.flush()
	default:
		scanner.current.WriteRune(r)
	}
}

func (scanner *shellFieldScanner) finish() []string {
	if scanner.escaped {
		scanner.current.WriteRune('\\')
	}
	scanner.flush()
	return scanner.fields
}

func (scanner *shellFieldScanner) flush() {
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
	model = strings.TrimSpace(model)
	if model == "" || strings.Contains(model, "://") || strings.ContainsAny(model, " \t\n\r\"'`$\\") {
		return ""
	}
	if strings.Contains(model, "../") || strings.HasPrefix(model, "/") {
		return ""
	}
	if len(model) > 128 {
		return ""
	}
	return model
}

func digestFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func sourceCommit() string {
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

func findUnsafeMapAt(path string, values map[string]any, rawEvent bool) (string, string) {
	for key, child := range values {
		childPath := childPath(path, key)
		reason, skip := unsafeMapFieldReason(childPath, strings.ToLower(key), child, rawEvent)
		if reason != "" {
			return childPath, reason
		}
		if skip {
			continue
		}
		if field, reason := findUnsafeValueAt(childPath, child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}

func childPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

func unsafeMapFieldReason(path, key string, value any, rawEvent bool) (string, bool) {
	if rawEvent && unretainedRawToolInputField(path, key, value) {
		return "", true
	}
	if rawFieldNames[key] {
		return "forbidden_raw_field", false
	}
	if sensitiveFieldNames[key] {
		return "sensitive_field", false
	}
	if rawEvent && unretainedRawBodyField(key) && !structuredRawBody(value) {
		return "", true
	}
	return "", false
}

func unretainedRawToolInputField(path, key string, value any) bool {
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
	switch key {
	case "text", "content", "input", "output", "stdout", "stderr":
		return true
	default:
		return false
	}
}

func structuredRawBody(value any) bool {
	switch value.(type) {
	case map[string]any:
		return true
	default:
		return false
	}
}

func findUnsafeSliceAt(path string, values []any, rawEvent bool) (string, string) {
	for i, child := range values {
		if field, reason := findUnsafeValueAt(fmt.Sprintf("%s[%d]", path, i), child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}

func findUnsafeStringAt(path, value string, rawEvent bool) (string, string) {
	if strings.TrimSpace(value) == "" {
		return "", ""
	}
	if !rawEvent && unsafePathValue(value) {
		return path, "unsafe_path_or_private_path"
	}
	if unsafeURL(value) {
		return path, "authenticated_url"
	}
	if providerTokenPrefix.MatchString(value) {
		return path, "token_like_value"
	}
	if unsafeEncodedToken(path, value, rawEvent) {
		return path, "token_like_value"
	}
	return "", ""
}

func unsafePathValue(value string) bool {
	return privatePathPattern.MatchString(value) || strings.HasPrefix(value, "/") || strings.Contains(value, "../")
}

func unsafeEncodedToken(path, value string, rawEvent bool) bool {
	if digestField(path) || sha256Pattern.MatchString(value) {
		return false
	}
	if rawEvent && rawPathLikeField(path) {
		return false
	}
	return base64TokenPattern.MatchString(value)
}

func rawPathLikeField(path string) bool {
	field := path
	if idx := strings.LastIndex(field, "."); idx >= 0 {
		field = field[idx+1:]
	}
	if idx := strings.LastIndex(field, "["); idx >= 0 {
		field = field[:idx]
	}
	switch strings.ToLower(field) {
	case "path", "file", "filepath", "file_path", "dir", "directory", "cwd":
		return true
	default:
		return false
	}
}

func unsafeURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return false
	}
	if parsed.User != nil {
		return true
	}
	for key := range parsed.Query() {
		if authQueryKeys[strings.ToLower(key)] {
			return true
		}
	}
	return false
}

func digestField(path string) bool {
	last := path
	if idx := strings.LastIndex(last, "."); idx >= 0 {
		last = last[idx+1:]
	}
	switch last {
	case "source_digest", "validation_digest", "commit_digest", "envelope_digest", "payload_digest", "sha256":
		return true
	default:
		return false
	}
}

func validFamily(family string) bool {
	switch family {
	case "harness", "model", "interaction", "phase", "review", "tool", "mutation", "test", "pr", "merge", "gap":
		return true
	default:
		return false
	}
}

func validContentState(state string) bool {
	switch state {
	case ContentRedacted, ContentDigestOnly, ContentRetainedSafe, ContentNotApplicable:
		return true
	default:
		return false
	}
}

func validState(state string) bool {
	switch state {
	case StatePass, StateFail, StateCannotVerify, StateNotAssessed:
		return true
	default:
		return false
	}
}

func validRuleKey(key string) bool {
	switch key {
	case "missing_required_family", "missing_optional_family", "source_unavailable", "unsafe_input", "digest_mismatch", "schema_version_mismatch", "cross_link_conflict":
		return true
	default:
		return false
	}
}

func safeRef(ref string) bool {
	return safeIDPattern.MatchString(ref) || sha256Pattern.MatchString(ref)
}

func safeOperationRef(ref string) bool {
	if strings.HasPrefix(ref, "adapter-run:") || strings.HasPrefix(ref, "delivery-trace:") {
		return !strings.Contains(ref, "..") && !strings.Contains(ref, "://") && len(ref) <= 256
	}
	return safeRef(ref)
}

func safeEvent(eventID string) string {
	if safeIDPattern.MatchString(eventID) {
		return eventID
	}
	return "unknown_event"
}

func nonAuthority() string {
	return "harness observation is evidence only; no harness compliance, feature delivery, PR approval, merge approval, release readiness, or production trust is claimed"
}

func DecodeValidation(r io.Reader) (Validation, error) {
	var validation Validation
	if err := json.NewDecoder(r).Decode(&validation); err != nil {
		return Validation{}, err
	}
	return validation, nil
}
