package recorder

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

const defaultBaseOutputDir = ".sdp-trace-runs"

// RecorderOptions configures local capture behavior.
type RecorderOptions struct {
	Task               string
	ContractPath       string
	UseDefaultContract bool
	WrapperName        string
	OutputDir          string
	Command            []string
	Env                []string
}

// RecorderResult represents run output for CLI callers and tests.
type RecorderResult struct {
	RunDir   string
	ExitCode int
	Contract trace.Contract
}

// Run executes a command and writes a local chain with minimal milestone fields.
func Run(ctx context.Context, options RecorderOptions) (RecorderResult, error) {
	if len(options.Command) == 0 {
		return RecorderResult{}, errors.New("missing command")
	}
	prepared, err := prepareRun(options)
	if err != nil {
		return RecorderResult{}, err
	}
	exitCode, signal := runCommand(ctx, options.Command, prepared.options.Env, prepared.writer)
	if err := finishRun(prepared.writer, prepared.commandID, exitCode, signal, prepared.options); err != nil {
		return RecorderResult{}, err
	}
	return RecorderResult{
		RunDir:   prepared.runDir,
		ExitCode: exitCode,
		Contract: prepared.contract,
	}, nil
}

type preparedRun struct {
	options   RecorderOptions
	runDir    string
	contract  trace.Contract
	writer    *runWriter
	commandID string
}

func prepareRun(options RecorderOptions) (preparedRun, error) {
	runDir, contract, writer, err := prepareRunFiles(options)
	if err != nil {
		return preparedRun{}, err
	}
	options.Env = ifNilCopy(os.Environ())
	commandID, err := startRecordedCommand(writer, options, contract)
	if err != nil {
		return preparedRun{}, err
	}
	return preparedRun{options: options, runDir: runDir, contract: contract, writer: writer, commandID: commandID}, nil
}

func prepareRunFiles(options RecorderOptions) (string, trace.Contract, *runWriter, error) {
	runDir, err := prepareRunDir(options.OutputDir)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}
	contract, err := resolveRecorderContract(options.ContractPath, options.UseDefaultContract)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}
	writer, err := initializedRunWriter(runDir, contract, options)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}
	return runDir, contract, writer, nil
}

func initializedRunWriter(runDir string, contract trace.Contract, options RecorderOptions) (*runWriter, error) {
	writer, err := newRunWriter(runDir, contract, options.Task)
	if err != nil {
		return nil, err
	}
	if err := initializeRunWriter(writer, options, contract); err != nil {
		return nil, err
	}
	return writer, nil
}

func resolveRecorderContract(contractPath string, useDefault bool) (trace.Contract, error) {
	contract, err := resolveContract(contractPath, useDefault)
	if err != nil {
		return trace.Contract{}, err
	}
	if contract.ContractID == "" {
		return trace.Contract{}, errors.New("contract missing identifier")
	}
	return contract, nil
}

func startRecordedCommand(writer *runWriter, options RecorderOptions, contract trace.Contract) (string, error) {
	if err := appendRunStartEvents(writer, options, contract); err != nil {
		return "", err
	}
	commandID := randomHex(12)
	if err := appendCommandStarted(writer, options, commandID); err != nil {
		return "", err
	}
	return commandID, nil
}

func initializeRunWriter(writer *runWriter, options RecorderOptions, contract trace.Contract) error {
	if options.ContractPath != "" {
		writer.manifest.ContractPath = options.ContractPath
		writer.manifest.ContractDigest = trace.SHA256Hex(string(mustMarshalJSON(contract)))
	}
	return writer.writeManifest()
}

func finishRun(writer *runWriter, commandID string, exitCode int, signal string, options RecorderOptions) error {
	if err := appendCommandFinished(writer, options, commandID, exitCode, signal); err != nil {
		return err
	}
	closureState := trace.ClosureStateCompleted
	if exitCode != 0 {
		closureState = trace.ClosureStateCommandFailure
	}
	if err := appendRunClosed(writer, commandID, closureState); err != nil {
		return err
	}
	return writer.finalize(closureState)
}

func prepareRunDir(outputDir string) (string, error) {
	if outputDir != "" {
		return outputDir, ensureFreshOutputDir(outputDir)
	}
	if err := os.MkdirAll(defaultBaseOutputDir, 0o755); err != nil {
		return "", err
	}
	return os.MkdirTemp(defaultBaseOutputDir, "run-")
}

func appendRunStartEvents(writer *runWriter, options RecorderOptions, contract trace.Contract) error {
	if err := writer.appendEvent(trace.EventRecorderAttached, trace.RecorderAttachedPayload{
		RecorderVersion: trace.RecorderVersion,
		RunNonce:        fmt.Sprintf("run-nonce-%s", writer.manifest.RunID),
		Pid:             os.Getpid(),
	}); err != nil {
		return err
	}
	return writer.appendEvent(trace.EventRunStarted, trace.RunStartedPayload{
		TaskRef:        options.Task,
		Command:        commandName(options.Command),
		ContractID:     contract.ContractID,
		ContractDigest: writer.manifest.ContractDigest,
		ContractPath:   writer.manifest.ContractPath,
	})
}

func appendCommandStarted(writer *runWriter, options RecorderOptions, commandID string) error {
	return writer.appendEvent(trace.EventCommandStarted, trace.CommandStartedPayload{
		RecordedCommandPayload: trace.RecordedCommandPayload{
			CommandID: commandID,
			Command:   commandName(options.Command),
		},
		CmdPID:         0,
		WrapperName:    options.WrapperName,
		Cwd:            currentWorkingDir(),
		ContractDigest: writer.manifest.ContractDigest,
	})
}

func appendCommandFinished(writer *runWriter, options RecorderOptions, commandID string, exitCode int, signal string) error {
	return writer.appendEvent(trace.EventCommandFinished, trace.CommandFinishedPayload{
		RecordedCommandPayload: trace.RecordedCommandPayload{
			CommandID: commandID,
			Command:   commandName(options.Command),
		},
		ExitCode:     exitCode,
		Signal:       signal,
		StdoutDigest: writer.stdoutDigest(),
		StderrDigest: writer.stderrDigest(),
	})
}

func appendRunClosed(writer *runWriter, commandID, closureState string) error {
	return writer.appendEvent(trace.EventRunClosed, trace.RunClosedPayload{
		ChainHead:     writer.lastEventHash(),
		ClosureState:  closureState,
		CommandID:     commandID,
		MissingCount:  0,
		ObservedCount: writer.eventCount(),
	})
}

func resolveContract(contractPath string, useDefault bool) (trace.Contract, error) {
	if contractPath == "" {
		if useDefault {
			return trace.DefaultContract, nil
		}
		return trace.Contract{}, errors.New("contract required unless --use-default-contract is set")
	}
	return trace.LoadContract(contractPath)
}

func ensureFreshOutputDir(runDir string) error {
	entries, err := os.ReadDir(runDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("output directory must be empty: %s", runDir)
	}
	return nil
}

type runWriter struct {
	runDir     string
	contract   trace.Contract
	manifest   trace.RunManifest
	sequence   int
	lastHash   string
	events     []trace.Event
	stdoutHash hashWriter
	stderrHash hashWriter
}

type hashWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (h *hashWriter) Write(p []byte) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.buf.Write(p)
}

func (h *hashWriter) Digest() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	sum := sha256.Sum256(h.buf.Bytes())
	return hex.EncodeToString(sum[:])
}

func newRunWriter(runDir string, contract trace.Contract, task string) (*runWriter, error) {
	if err := os.MkdirAll(filepath.Join(runDir, "events"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(runDir, "artifacts"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(runDir, "verifier"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(runDir, "export"), 0o755); err != nil {
		return nil, err
	}

	runID := randomHex(16)
	manifest := trace.RunManifest{
		SchemaVersion:   trace.SchemaVersion,
		RunID:           runID,
		RecorderVersion: trace.RecorderVersion,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339Nano),
		Task:            task,
		ContractID:      contract.ContractID,
		SourceSnapshot:  "",
		SourceState:     "",
		EventCount:      0,
		ClosureState:    trace.ClosureStateUnknown,
		ContractPath:    "",
		ContractDigest:  trace.SHA256Hex(string(mustMarshalJSON(contract))),
	}

	sourceDigest, sourceState := trace.LocalSourceSnapshot(currentWorkingDir())
	manifest.SourceSnapshot = sourceDigest
	manifest.SourceState = sourceState

	return &runWriter{
		runDir:   runDir,
		contract: contract,
		manifest: manifest,
	}, nil
}

func (w *runWriter) appendEvent(eventType trace.EventType, payload any) error {
	payloadMap, err := toEventPayload(payload)
	if err != nil {
		return err
	}
	event := trace.Event{
		SchemaVersion: trace.SchemaVersion,
		RunID:         w.manifest.RunID,
		EventID:       randomHex(24),
		Sequence:      w.sequence,
		EventType:     eventType,
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: w.lastHash,
		HashAlgorithm: trace.HashAlgSHA256,
		Canonicalization: trace.Canonicalization{
			Algorithm: trace.CanonicalSchemaAlgo,
			Version:   trace.CanonicalAlgoV,
		},
		EventPayload: payloadMap,
		ObservedBy:   "local_recorder",
	}
	event = event.EnsureDefaults()
	eventHash, err := event.WithComputedEventHash()
	if err != nil {
		return err
	}
	event = eventHash

	filename := filepath.Join(w.runDir, "events", eventFilename(event.Sequence, eventType))
	if err := writeIndentedJSON(filename, event); err != nil {
		return err
	}
	w.sequence++
	w.lastHash = event.EventHash
	w.events = append(w.events, event)
	return w.writeManifest()
}

func (w *runWriter) writeManifest() error {
	w.manifest.EventCount = w.sequence
	if w.sequence > 0 {
		w.manifest.EventChainHead = w.lastHash
		w.manifest.FinalChainHead = w.lastHash
	}
	return writeIndentedJSON(filepath.Join(w.runDir, "run.json"), w.manifest)
}

func (w *runWriter) finalize(closureState string) error {
	w.manifest.ClosureState = closureState
	w.manifest.ClosedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if err := w.writeManifest(); err != nil {
		return err
	}
	if err := writeText(filepath.Join(w.runDir, "artifacts", "stdout.digest"), w.stdoutDigest()+"\n"); err != nil {
		return err
	}
	return writeText(filepath.Join(w.runDir, "artifacts", "stderr.digest"), w.stderrDigest()+"\n")
}

func (w *runWriter) stdoutDigest() string {
	return w.stdoutHash.Digest()
}

func (w *runWriter) stderrDigest() string {
	return w.stderrHash.Digest()
}

func (w *runWriter) lastEventHash() string {
	return w.lastHash
}

func (w *runWriter) eventCount() int {
	return len(w.events)
}

func runCommand(ctx context.Context, command []string, env []string, writer *runWriter) (int, string) {
	cmd := recordedCommand(ctx, command, env, writer)
	if err := cmd.Start(); err != nil {
		return 1, "start_failed"
	}
	return waitCommand(ctx, cmd)
}

func recordedCommand(ctx context.Context, command []string, env []string, writer *runWriter) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = io.MultiWriter(os.Stdout, &writer.stdoutHash)
	cmd.Stderr = io.MultiWriter(os.Stderr, &writer.stderrHash)
	return cmd
}

func waitCommand(ctx context.Context, cmd *exec.Cmd) (int, string) {
	err := cmd.Wait()
	if err == nil {
		return 0, ""
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), processSignal(exitErr.ProcessState)
	}
	if ctx.Err() != nil {
		return 1, "context_cancelled"
	}
	return 1, ""
}

func processSignal(processState *os.ProcessState) string {
	if noProcessSignal(processState) {
		return ""
	}
	status, ok := processState.Sys().(syscall.WaitStatus)
	if !ok {
		return ""
	}
	return status.Signal().String()
}

func noProcessSignal(processState *os.ProcessState) bool {
	return processState == nil || processState.Exited()
}

func eventFilename(sequence int, eventType trace.EventType) string {
	return fmt.Sprintf("%06d-%s.json", sequence, eventType)
}

func toEventPayload(payload any) (map[string]any, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	payloadMap, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("event payload must be a JSON object")
	}
	return payloadMap, nil
}

func commandName(command []string) string {
	if len(command) == 0 {
		return ""
	}
	return filepath.Base(command[0])
}

func currentWorkingDir() string {
	cwd, _ := os.Getwd()
	return cwd
}

func mustMarshalJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func writeIndentedJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func writeText(path string, value string) error {
	return os.WriteFile(path, []byte(value), 0o644)
}

func ifNilCopy(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string{}, values...)
}

func randomHex(length int) string {
	const alphabet = "0123456789abcdef"
	out := make([]byte, length)
	for i := range out {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return fmt.Sprintf("fallback-%x", time.Now().UnixNano())
		}
		out[i] = alphabet[v.Int64()]
	}
	return string(out)
}
