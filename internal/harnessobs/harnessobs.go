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
)

var (
	safeIDPattern       = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)
	safeFileIDPattern   = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
	sha256Pattern       = regexp.MustCompile(`^[a-f0-9]{64}$`)
	base64TokenPattern  = regexp.MustCompile(`(?i)^[A-Za-z0-9+/_-]{43,}={0,2}$`)
	providerTokenPrefix = regexp.MustCompile(`(?i)(sk-[A-Za-z0-9_-]{16,}|xox[baprs]-[A-Za-z0-9-]{16,}|gh[pousr]_[A-Za-z0-9_]{20,})`)
	privatePathPattern  = regexp.MustCompile(`(^|[\s"'])/(Users|home|private|var|tmp)/[^\s"']+`)
	rawFieldNames       = map[string]bool{
		"raw_prompt":         true,
		"prompt":             true,
		"raw_model_response": true,
		"model_response":     true,
		"raw_command":        true,
		"command_body":       true,
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

type ValidateOptions struct {
	ProfilePath string
	RunDir      string
	OutPath     string
}

func Observe(opts ObserveOptions) (Run, error) {
	if strings.TrimSpace(opts.ProfilePath) == "" {
		return Run{}, errors.New("harness observe requires --profile")
	}
	if strings.TrimSpace(opts.SourcePath) == "" {
		return Run{}, errors.New("harness observe requires --source")
	}
	if strings.TrimSpace(opts.OutDir) == "" {
		return Run{}, errors.New("harness observe requires --out")
	}
	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return Run{}, fmt.Errorf("unsafe profile path: %w", err)
	}
	sourcePath, err := safeExistingFile(opts.SourcePath)
	if err != nil {
		return Run{}, fmt.Errorf("unsafe source path: %w", err)
	}
	outDir, err := safeOutDir(opts.OutDir)
	if err != nil {
		return Run{}, err
	}
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Run{}, err
	}
	events, sourceDigest, err := readEvents(profile, sourcePath)
	if err != nil {
		return Run{}, err
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return Run{}, err
	}
	for _, event := range events {
		path := filepath.Join(outDir, "events", event.EventID+".json")
		if err := writeJSON(path, event); err != nil {
			return Run{}, err
		}
	}
	run := Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,
		SourcePath:         filepath.Base(sourcePath),
		SourceDigest:       sourceDigest,
		EventCount:         len(events),
		EventRefs:          eventRefs(events),
		CreatedAt:          opts.Now.Format(time.RFC3339),
	}
	if err := writeJSON(filepath.Join(outDir, "run.json"), run); err != nil {
		return Run{}, err
	}
	return run, nil
}

func Validate(opts ValidateOptions) (Validation, error) {
	if strings.TrimSpace(opts.ProfilePath) == "" {
		return Validation{}, errors.New("harness validate requires --profile")
	}
	if strings.TrimSpace(opts.RunDir) == "" {
		return Validation{}, errors.New("harness validate requires --run")
	}
	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return Validation{}, fmt.Errorf("unsafe profile path: %w", err)
	}
	runDir, err := safeExistingDir(opts.RunDir)
	if err != nil {
		return Validation{}, fmt.Errorf("unsafe run path: %w", err)
	}
	outPath := ""
	if opts.OutPath != "" {
		outPath, err = safeOutFile(opts.OutPath)
		if err != nil {
			return Validation{}, fmt.Errorf("unsafe out path: %w", err)
		}
	}
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Validation{}, err
	}
	run, events, err := LoadRun(runDir)
	if err != nil {
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
		if opts.OutPath != "" {
			if err := writeJSON(outPath, validation); err != nil {
				return Validation{}, err
			}
		}
		return validation, nil
	}
	validation := evaluate(profile, run, events)
	if opts.OutPath != "" {
		if err := writeJSON(outPath, validation); err != nil {
			return Validation{}, err
		}
	}
	return validation, nil
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
	if profile.SchemaVersion != ProfileSchemaVersion {
		return fmt.Errorf("unsupported harness profile schema_version: %s", profile.SchemaVersion)
	}
	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe profile_id")
	}
	if !safeIDPattern.MatchString(profile.HarnessFamily) {
		return errors.New("unsafe harness_family")
	}
	if profile.EventSchemaVersion != EventSchemaVersion {
		return errors.New("unsupported event_schema_version")
	}
	if len(profile.RequiredEventFamilies) == 0 {
		return errors.New("profile requires at least one required_event_family")
	}
	for _, family := range append(profile.RequiredEventFamilies, profile.OptionalEventFamilies...) {
		if !validFamily(family) {
			return fmt.Errorf("unsupported event family: %s", family)
		}
	}
	for key, rule := range profile.DegradationRules {
		if !validRuleKey(key) {
			return fmt.Errorf("unsupported degradation rule: %s", key)
		}
		if !validState(rule.State) || !safeIDPattern.MatchString(rule.ReasonCode) {
			return fmt.Errorf("invalid degradation rule %s", key)
		}
	}
	return nil
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
	if !validContentState(event.ContentState) {
		return errors.New("invalid content_state")
	}
	for _, field := range event.UnavailableFields {
		if !safeIDPattern.MatchString(field.Field) || field.State != StateNotAssessed || !safeIDPattern.MatchString(field.ReasonCode) {
			return errors.New("invalid unavailable_fields")
		}
	}
	return nil
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
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("path must be relative local file without traversal")
	}
	clean := filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(clean)
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
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", errors.New("path escapes working directory")
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("path must be a file")
	}
	return rel, nil
}

func safeExistingDir(path string) (string, error) {
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("path must be relative local directory without traversal")
	}
	clean := filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(clean)
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
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", errors.New("path escapes working directory")
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("path must be a directory")
	}
	return rel, nil
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
	if path == "" {
		path = "."
	}
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("parent path must be relative local directory without traversal")
	}
	clean := filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			parent := filepath.Dir(clean)
			if parent == clean {
				return "", err
			}
			return safeParentDir(parent)
		}
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
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", errors.New("parent path escapes working directory")
	}
	return rel, nil
}

func safeOutDir(path string) (string, error) {
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New("out must be a relative local directory without traversal")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(path)
	if _, err := os.Lstat(clean); err == nil {
		resolved, err := filepath.EvalSymlinks(clean)
		if err != nil {
			return "", err
		}
		abs, err := filepath.Abs(resolved)
		if err != nil {
			return "", err
		}
		rel, err := filepath.Rel(cwd, abs)
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			return "", errors.New("out path escapes working directory")
		}
		clean = rel
	} else {
		parent := filepath.Dir(clean)
		for parent != "." && parent != string(filepath.Separator) {
			if _, statErr := os.Lstat(parent); statErr == nil {
				resolved, err := filepath.EvalSymlinks(parent)
				if err != nil {
					return "", err
				}
				abs, err := filepath.Abs(resolved)
				if err != nil {
					return "", err
				}
				rel, err := filepath.Rel(cwd, abs)
				if err != nil {
					return "", err
				}
				if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
					return "", errors.New("out parent path escapes working directory")
				}
				break
			}
			parent = filepath.Dir(parent)
		}
	}
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

func findUnsafe(value any) (string, string) {
	return findUnsafeAt("", value)
}

func findUnsafeAt(path string, value any) (string, string) {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if rawFieldNames[strings.ToLower(key)] {
				return childPath, "forbidden_raw_field"
			}
			if field, reason := findUnsafeAt(childPath, child); field != "" {
				return field, reason
			}
		}
	case []any:
		for i, child := range v {
			if field, reason := findUnsafeAt(fmt.Sprintf("%s[%d]", path, i), child); field != "" {
				return field, reason
			}
		}
	case string:
		if strings.TrimSpace(v) == "" {
			return "", ""
		}
		if privatePathPattern.MatchString(v) || strings.HasPrefix(v, "/") || strings.Contains(v, "../") {
			return path, "unsafe_path_or_private_path"
		}
		if unsafeURL(v) {
			return path, "authenticated_url"
		}
		if providerTokenPrefix.MatchString(v) {
			return path, "token_like_value"
		}
		if !digestField(path) && base64TokenPattern.MatchString(v) && !sha256Pattern.MatchString(v) {
			return path, "token_like_value"
		}
	}
	return "", ""
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
