package trace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SchemaVersion is the local, in-repo Block 10 event schema version.
const SchemaVersion = "block10-event-v1"

// Hashing and canonicalization constants used by the shared trace model.
const (
	NullEventHash = ""
	// CanonicalAlgoV is preserved for pre-existing call-sites.
	CanonicalAlgoV = CanonicalAlgoVersion
)

// Verifier constants used by CLI and test fixtures.
const (
	ClosureStateCompleted      = "completed"
	ClosureStateCommandFailure = "command_failed"
	ClosureStateUnknown        = "unknown"
)

// RecorderVersion is embedded in run metadata.
const RecorderVersion = "0.1.0"

// EventType defines the supported event class names for the first milestone.
type EventType string

// EventType constants.
const (
	EventRecorderAttached        EventType = "recorder_attached"
	EventRunStarted              EventType = "run_started"
	EventCommandStarted          EventType = "command_started"
	EventCommandFinished         EventType = "command_finished"
	EventRunClosed               EventType = "run_closed"
	EventPolicyOverrideRequested EventType = "policy_override_requested"
)

// EvidenceState maps direct missing-evidence states used by first-milestone verification.
type EvidenceState string

// EvidenceState constants.
const (
	EvidenceStatePresent      EvidenceState = "present"
	EvidenceStateMissing      EvidenceState = "missing"
	EvidenceStateNotAssessed  EvidenceState = "not_assessed"
	EvidenceStateCannotVerify EvidenceState = "cannot_verify"
)

// VerifierVerdict tracks the high-level result.
type VerifierVerdict string

// Verifier constants.
const (
	VerdictObserved     VerifierVerdict = "observed"
	VerdictFail         VerifierVerdict = "fail"
	VerdictCannotVerify VerifierVerdict = "cannot_verify"
	VerdictNotAssessed  VerifierVerdict = "not_assessed"
)

// TrustScope identifies evidence source context.
type TrustScope string

// TrustScope constants.
const (
	TrustScopeLocalObserved TrustScope = "local_observed"
)

// Completeness tracks output completeness.
type Completeness string

// Completeness constants.
const (
	CompletenessComplete         Completeness = "complete"
	CompletenessPartial          Completeness = "partial"
	CompletenessMissingTelemetry Completeness = "missing_telemetry"
	CompletenessUnknown          Completeness = "unknown"
)

// Replayability expresses reusability of recorded output.
type Replayability string

// Replayability constants.
const (
	ReplayabilityFull    Replayability = "full"
	ReplayabilityNone    Replayability = "none"
	ReplayabilityPartial Replayability = "partial"
)

// RunManifest records run-level metadata.
type RunManifest struct {
	SchemaVersion   string `json:"schema_version"`
	RunID           string `json:"run_id"`
	RecorderVersion string `json:"recorder_version"`
	CreatedAt       string `json:"created_at"`
	ClosedAt        string `json:"closed_at,omitempty"`
	Task            string `json:"task,omitempty"`
	ContractID      string `json:"contract_id"`
	ContractPath    string `json:"contract_path,omitempty"`
	ContractDigest  string `json:"contract_digest,omitempty"`
	SourceSnapshot  string `json:"source_snapshot_digest"`
	SourceState     string `json:"source_snapshot_state"`
	EventCount      int    `json:"event_count"`
	EventChainHead  string `json:"event_chain_head"`
	FinalChainHead  string `json:"final_chain_head"`
	ClosureState    string `json:"closure_state"`
}

// Validate checks required manifest fields for shared verifier paths.
func (manifest RunManifest) Validate() error {
	if strings.TrimSpace(manifest.SchemaVersion) == "" {
		return errors.New("run manifest missing schema_version")
	}
	if strings.TrimSpace(manifest.RunID) == "" {
		return errors.New("run manifest missing run_id")
	}
	if manifest.EventCount < 0 {
		return errors.New("run manifest event_count must be >= 0")
	}
	if strings.TrimSpace(manifest.ContractID) == "" {
		return errors.New("run manifest missing contract_id")
	}
	return nil
}

// Event captures one append-only record in a run.
type Event struct {
	SchemaVersion    string           `json:"schema_version"`
	RunID            string           `json:"run_id"`
	EventID          string           `json:"event_id"`
	Sequence         int              `json:"sequence"`
	EventType        EventType        `json:"event_type"`
	Timestamp        string           `json:"timestamp"`
	PrevEventHash    string           `json:"prev_event_hash"`
	HashAlgorithm    string           `json:"hash_algorithm"`
	Canonicalization Canonicalization `json:"canonicalization"`
	PayloadDigest    string           `json:"payload_digest"`
	EventPayload     map[string]any   `json:"event_payload"`
	Payload          json.RawMessage  `json:"payload"`
	EventHash        string           `json:"event_hash"`
	ObservedBy       string           `json:"observed_by,omitempty"`
	ClockMonotonic   int64            `json:"clock_monotonic,omitempty"`
}

type Canonicalization struct {
	Algorithm string `json:"algorithm"`
	Version   string `json:"version"`
}

// EnsureDefaults populates event defaults used during hashing and writing.
func (event Event) EnsureDefaults() Event {
	if event.SchemaVersion == "" {
		event.SchemaVersion = SchemaVersion
	}
	if event.HashAlgorithm == "" {
		event.HashAlgorithm = HashAlgSHA256
	}
	if event.Canonicalization.Algorithm == "" {
		event.Canonicalization.Algorithm = CanonicalSchemaAlgo
	}
	if event.Canonicalization.Version == "" {
		event.Canonicalization.Version = CanonicalAlgoVersion
	}
	if event.PrevEventHash == "" {
		event.PrevEventHash = NullEventHash
	}
	if synced, err := event.syncPayloadRepresentation(); err == nil {
		event = synced
	}
	return event
}

// Validate checks local invariants before persistence and replay.
func (event Event) Validate() error {
	if strings.TrimSpace(event.SchemaVersion) == "" {
		return errors.New("missing schema_version")
	}
	if strings.TrimSpace(event.RunID) == "" {
		return errors.New("missing run_id")
	}
	if strings.TrimSpace(event.EventID) == "" {
		return errors.New("missing event_id")
	}
	if strings.TrimSpace(event.Timestamp) == "" {
		return errors.New("missing timestamp")
	}
	if strings.TrimSpace(event.EventHash) == "" {
		return errors.New("missing event_hash")
	}
	if event.Sequence < 0 {
		return fmt.Errorf("invalid sequence %d", event.Sequence)
	}
	if event.HashAlgorithm != "" && event.HashAlgorithm != HashAlgSHA256 {
		return fmt.Errorf("unsupported hash_algorithm %s", event.HashAlgorithm)
	}
	_, err := event.syncPayloadRepresentation()
	return err
}

// WithComputedPayloadDigest computes PayloadDigest from the event payload bytes.
func (event Event) WithComputedPayloadDigest() (Event, error) {
	event = event.EnsureDefaults()
	synced, err := event.syncPayloadRepresentation()
	if err != nil {
		return Event{}, err
	}
	event = synced
	payloadDigest, err := CanonicalEventPayloadDigest(event.Payload)
	if err != nil {
		return Event{}, err
	}
	event.PayloadDigest = payloadDigest
	return event, nil
}

// WithComputedEventHash computes payload and event hashes.
func (event Event) WithComputedEventHash() (Event, error) {
	withDigest, err := event.WithComputedPayloadDigest()
	if err != nil {
		return Event{}, err
	}
	eventHash, err := ComputeEventHash(withDigest)
	if err != nil {
		return Event{}, err
	}
	withDigest.EventHash = eventHash
	return withDigest, nil
}

// VerifyPayloadDigest validates the payload digest for the event.
func (event Event) VerifyPayloadDigest() error {
	if strings.TrimSpace(event.PayloadDigest) == "" {
		return nil
	}
	synced, err := event.syncPayloadRepresentation()
	if err != nil {
		return err
	}
	event = synced
	computed, err := CanonicalEventPayloadDigest(event.Payload)
	if err != nil {
		return err
	}
	if event.PayloadDigest != computed {
		return fmt.Errorf("payload_digest mismatch: expected %s got %s", computed, event.PayloadDigest)
	}
	return nil
}

func (event Event) syncPayloadRepresentation() (Event, error) {
	switch {
	case len(event.Payload) > 0:
		if event.EventPayload == nil {
			var decoded map[string]any
			if err := json.Unmarshal(event.Payload, &decoded); err != nil {
				return Event{}, fmt.Errorf("invalid payload: %w", err)
			}
			event.EventPayload = decoded
		}
		return event, nil
	case event.EventPayload != nil:
		payload, err := json.Marshal(event.EventPayload)
		if err != nil {
			return Event{}, err
		}
		event.Payload = payload
		return event, nil
	default:
		return event, nil
	}
}

// RecordedCommandPayload stores common command lifecycle metadata.
type RecordedCommandPayload struct {
	CommandID string `json:"command_id"`
	Command   string `json:"command"`
}

// CommandStartedPayload captures launch metadata.
type CommandStartedPayload struct {
	RecordedCommandPayload
	CmdPID         int    `json:"cmd_pid"`
	WrapperName    string `json:"wrapper_name"`
	Cwd            string `json:"cwd"`
	ContractDigest string `json:"contract_digest"`
}

// CommandFinishedPayload captures close metadata.
type CommandFinishedPayload struct {
	RecordedCommandPayload
	ExitCode       int    `json:"exit_code"`
	Signal         string `json:"signal,omitempty"`
	StdoutDigest   string `json:"stdout_digest"`
	StderrDigest   string `json:"stderr_digest"`
	ClockMonotonic int64  `json:"clock_monotonic,omitempty"`
}

// RecorderAttachedPayload captures launch metadata for the recorder.
type RecorderAttachedPayload struct {
	RecorderVersion string `json:"recorder_version"`
	RunNonce        string `json:"run_nonce"`
	Pid             int    `json:"pid"`
}

// RunStartedPayload captures command + contract context.
type RunStartedPayload struct {
	TaskRef        string `json:"task_ref"`
	Command        string `json:"command"`
	ContractID     string `json:"contract_id"`
	ContractDigest string `json:"contract_digest"`
	ContractPath   string `json:"contract_path,omitempty"`
}

// RunClosedPayload summarizes terminal state.
type RunClosedPayload struct {
	ChainHead     string `json:"chain_head"`
	ClosureState  string `json:"closure_state"`
	CommandID     string `json:"command_id"`
	MissingCount  int    `json:"missing_evidence_count"`
	ObservedCount int    `json:"observed_event_count"`
}

// Contract controls required events for milestone verification.
type Contract struct {
	ContractID         string                `json:"contract_id"`
	Version            string                `json:"version"`
	RequiredEvents     []string              `json:"required_events"`
	RequiredEvidence   []EvidenceRequirement `json:"required_evidence,omitempty"`
	RequiredRuns       []RequiredRun         `json:"required_runs,omitempty"`
	LockRequiredBefore string                `json:"lock_required_before,omitempty"`
}

// RequiredRun names a contract-declared run that should be observed for an advisory gate.
type RequiredRun struct {
	ID               string   `json:"id"`
	WrapperName      string   `json:"wrapper_name"`
	RequiredEvidence []string `json:"required_evidence,omitempty"`
	Profile          string   `json:"profile,omitempty"`
}

// EvidenceRequirement names a contract-declared observation that can be
// matched against event payload fields without product-specific classifiers.
type EvidenceRequirement struct {
	ID            string `json:"id"`
	EventType     string `json:"event_type"`
	PayloadField  string `json:"payload_field"`
	PayloadEquals string `json:"payload_equals"`
}

// DefaultContract is the minimal local contract for first-milestone local recorder output.
var DefaultContract = Contract{
	ContractID: "local-default-v1",
	Version:    SchemaVersion,
	RequiredEvents: []string{
		string(EventRecorderAttached),
		string(EventRunStarted),
		string(EventCommandStarted),
		string(EventCommandFinished),
		string(EventRunClosed),
	},
	LockRequiredBefore: "run_started",
}

// MissingEvidenceRow describes one required-event mismatch.
type MissingEvidenceRow struct {
	ExpectedEvent       string `json:"expected_event"`
	ObservedState       string `json:"observed_state"`
	Reason              string `json:"reason"`
	PolicyReference     string `json:"policy_reference,omitempty"`
	ReplayabilityImpact string `json:"replayability_impact"`
}

// MissingEvidenceTable summarizes observed-vs-required events.
type MissingEvidenceTable struct {
	ContractID string               `json:"contract_id"`
	Rows       []MissingEvidenceRow `json:"rows"`
}

// VerifierResult is the live machine-readable verification envelope.
type VerifierResult struct {
	RunID         string          `json:"run_id"`
	Result        VerifierVerdict `json:"result"`
	TrustScope    TrustScope      `json:"trust_scope"`
	Completeness  Completeness    `json:"completeness"`
	Replayability Replayability   `json:"replayability"`
	Reason        string          `json:"reason"`
	RunDir        string          `json:"run_dir"`
}

// IntegrityAudit records one structural issue for explainability.
type IntegrityAudit struct {
	RunID   string            `json:"run_id"`
	Issue   string            `json:"issue"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details,omitempty"`
}

// EventHash computes sha256 over a canonicalized payload copy without event_hash.
func EventHash(event Event) (string, error) {
	computed, err := ComputeEventHash(event)
	if err != nil {
		return "", err
	}
	return computed, nil
}

// ReadJSON decodes any JSON file into dst.
func ReadJSON(ctx context.Context, path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_ = ctx
	return json.Unmarshal(data, dst)
}

// EventSeqFromFilename returns the numeric sequence prefix from `<NNNNNN>-<event>.json`.
func EventSeqFromFilename(name string) (int, error) {
	re := regexp.MustCompile(`^(\d{6})-`)
	matches := re.FindStringSubmatch(name)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid event filename: %s", name)
	}
	return strconv.Atoi(matches[1])
}

// SHA256Hex hashes arbitrary text for fixture metadata.
func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return strings.ToLower(hex.EncodeToString(sum[:]))
}

// ResolveContractPath converts empty or relative paths to absolute file system paths.
func ResolveContractPath(baseDir, contractPath string) string {
	if contractPath == "" {
		return ""
	}
	if filepath.IsAbs(contractPath) {
		return contractPath
	}
	if baseDir == "" {
		return contractPath
	}
	return filepath.Join(baseDir, contractPath)
}

// LocalSourceSnapshot returns a local source digest and state disclosure.
func LocalSourceSnapshot(baseDir string) (string, string) {
	cleanBase := filepath.Clean(baseDir)
	tree, treeErr := gitOutput(cleanBase, "rev-parse", "HEAD^{tree}")
	status, statusErr := gitOutput(cleanBase, "status", "--porcelain")
	if treeErr == nil && statusErr == nil {
		state := "git_tree_clean"
		if strings.TrimSpace(status) != "" {
			state = "git_tree_dirty"
		}
		return SHA256Hex("git-tree:" + strings.TrimSpace(tree)), state
	}
	return SHA256Hex("source-not-assessed:" + cleanBase), "not_assessed"
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
