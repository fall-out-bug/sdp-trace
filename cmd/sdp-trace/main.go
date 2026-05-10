package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
	"github.com/fall_out_bug/sdp-trace/internal/interaction"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/posture"
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
	"github.com/fall_out_bug/sdp-trace/internal/query"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

const (
	exitUsage        = 2
	exitCannotVerify = 3
)

var cliStdin io.Reader = os.Stdin

type commandHandler func(context.Context, []string, io.Writer, io.Writer) int
type subcommandHandler func([]string, io.Writer, io.Writer) int

var commandHandlers = map[string]commandHandler{
	"wrap":              runWrap,
	"run":               runWrappedCommand,
	"dry-run":           runDryRun,
	"preview":           runPreview,
	"doctor":            runDoctor,
	"install":           runInstall,
	"interaction":       runInteraction,
	"observe":           runObserveCommand,
	"harness":           runHarnessCommand,
	"envelope":          runEnvelope,
	"verify":            runVerify,
	"explain":           runExplain,
	"query":             runQuery,
	"query-pack":        runQueryPack,
	"export":            runExport,
	"report":            runReport,
	"gate":              runGate,
	"assess":            runAssess,
	"override":          runOverride,
	"checkpoint":        runCheckpoint,
	"witness":           runWitness,
	"validate-fixtures": runValidateFixtures,
	"release-proof":     runReleaseProof,
	"pr-review":         runPRReview,
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printUsage(stdout)
		return 0
	}
	cmd := args[0]
	cmdArgs := args[1:]

	handler, ok := commandHandlers[cmd]
	if !ok {
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return 1
	}
	return handler(context.Background(), cmdArgs, stdout, stderr)
}

func runObserveCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runObserve(args, stdout, stderr)
}

func runHarnessCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runHarness(args, stdout, stderr)
}

func runSubcommand(args []string, stdout, stderr io.Writer, label, usage string, handlers map[string]subcommandHandler) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return exitUsage
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintf(stdout, "Usage: sdp-trace %s\n", label)
		return 0
	}
	handler, ok := handlers[args[0]]
	if !ok {
		fmt.Fprintf(stderr, "unknown %s command: %s\n", subcommandName(label), args[0])
		return exitUsage
	}
	return handler(args[1:], stdout, stderr)
}

func subcommandName(label string) string {
	if before, _, ok := strings.Cut(label, " "); ok {
		return before
	}
	return label
}

func runObserve(args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "observe <setup|collect|session> [flags]", "observe requires setup, collect, or session", map[string]subcommandHandler{
		"setup":   runObserveSetup,
		"collect": runObserveCollect,
		"session": runObserveSession,
	})
}

func runHarness(args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "harness <observe|validate|summarize> [flags]", "harness requires observe, validate, or summarize", map[string]subcommandHandler{
		"observe":   runHarnessObserve,
		"validate":  runHarnessValidate,
		"summarize": runHarnessSummarize,
	})
}

func runHarnessObserve(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "harness observe"}
	opts.setString("profile", "")
	opts.setString("source", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "harness observe accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) == "" {
		fmt.Fprintln(stderr, "harness observe requires --profile")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("source")) == "" {
		fmt.Fprintln(stderr, "harness observe requires --source")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "harness observe requires --out")
		return exitUsage
	}
	run, err := harnessobs.Observe(harnessobs.ObserveOptions{
		ProfilePath: opts.stringValue("profile"),
		SourcePath:  opts.stringValue("source"),
		OutDir:      opts.stringValue("out"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "marshal harness run: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runObserveSetup(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "observe setup"}
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("command", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "observe setup accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) == "" {
		fmt.Fprintln(stderr, "observe setup requires --profile")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "observe setup requires --out")
		return exitUsage
	}
	session, err := harnessobs.SetupSession(harnessobs.SessionSetupOptions{
		ProfilePath: opts.stringValue("profile"),
		OutDir:      opts.stringValue("out"),
		Command:     opts.stringValue("command"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "marshal observe setup: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runObserveCollect(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "observe collect"}
	opts.setString("profile", "")
	opts.setString("run", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "observe collect accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) == "" {
		fmt.Fprintln(stderr, "observe collect requires --profile")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		fmt.Fprintln(stderr, "observe collect requires --run")
		return exitUsage
	}
	session, observed, err := harnessobs.CollectSession(harnessobs.SessionCollectOptions{
		ProfilePath: opts.stringValue("profile"),
		RunDir:      opts.stringValue("run"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, err := json.MarshalIndent(struct {
		Session harnessobs.SessionRun `json:"session"`
		Run     harnessobs.Run        `json:"run"`
	}{Session: session, Run: observed}, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "marshal observe collect: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	if session.CollectionState == harnessobs.StateCannotVerify {
		return exitCannotVerify
	}
	return 0
}

func runObserveSession(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "observe session"}
	opts.setString("profile", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) == "" {
		fmt.Fprintln(stderr, "observe session requires --profile")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "observe session requires --out")
		return exitUsage
	}
	if len(opts.rest()) == 0 {
		fmt.Fprintln(stderr, "observe session requires command after --")
		return exitUsage
	}
	session, observed, err := harnessobs.RunSession(harnessobs.SessionOptions{
		ProfilePath: opts.stringValue("profile"),
		OutDir:      opts.stringValue("out"),
		Command:     opts.rest(),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, err := json.MarshalIndent(struct {
		Session harnessobs.SessionRun `json:"session"`
		Run     harnessobs.Run        `json:"run"`
	}{Session: session, Run: observed}, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "marshal observe session: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runHarnessValidate(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "harness validate"}
	opts.setString("profile", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "harness validate accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) == "" {
		fmt.Fprintln(stderr, "harness validate requires --profile")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		fmt.Fprintln(stderr, "harness validate requires --run")
		return exitUsage
	}
	validation, err := harnessobs.Validate(harnessobs.ValidateOptions{
		ProfilePath: opts.stringValue("profile"),
		RunDir:      opts.stringValue("run"),
		OutPath:     opts.stringValue("out"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, err := json.MarshalIndent(validation, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "marshal harness validation: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	switch validation.ValidationState {
	case harnessobs.StatePass:
		return 0
	case harnessobs.StateFail:
		return 1
	case harnessobs.StateNotAssessed, harnessobs.StateCannotVerify:
		return exitCannotVerify
	default:
		fmt.Fprintf(stderr, "unknown harness validation state: %s\n", validation.ValidationState)
		return exitCannotVerify
	}
}

func runHarnessSummarize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "harness summarize"}
	opts.setString("validation", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "harness summarize accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("validation")) == "" {
		fmt.Fprintln(stderr, "harness summarize requires --validation")
		return exitUsage
	}
	validation, err := harnessobs.LoadValidation(opts.stringValue("validation"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprint(stdout, harnessobs.Summarize(validation))
	return 0
}

func runPRReview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "pr-review <packet|run|synthesize|validate|summarize|check> [flags]", "pr-review requires packet, run, synthesize, validate, summarize, or check", map[string]subcommandHandler{
		"packet":     runPRReviewPacket,
		"run":        runPRReviewRun,
		"synthesize": runPRReviewSynthesize,
		"validate":   runPRReviewValidate,
		"summarize":  runPRReviewSummarize,
		"check":      runPRReviewCheck,
	})
}

func runPRReviewPacket(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "pr-review packet"}
	opts.setString("out", "")
	opts.setString("repo-id", "")
	opts.setString("change-ref", "")
	opts.setString("base", "")
	opts.setString("head", "")
	opts.setString("diff", "")
	opts.setString("metadata", "")
	opts.setString("context", "")
	opts.setString("verification", "")
	opts.setString("ci-state", prreview.StateNotAssessed)
	opts.setString("created-by", "sdp-trace-cli")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review packet accepts only flags")
		return exitUsage
	}
	if err := requirePRReviewPacketInputs(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	packet, err := prreview.BuildPacket(prreview.PacketOptions{
		OutDir:            opts.stringValue("out"),
		RepoID:            opts.stringValue("repo-id"),
		ChangeRef:         opts.stringValue("change-ref"),
		BaseCommit:        opts.stringValue("base"),
		HeadCommit:        opts.stringValue("head"),
		DiffPath:          opts.stringValue("diff"),
		MetadataPath:      opts.stringValue("metadata"),
		ContextPaths:      repeatedFlagValues(args, "context", opts.stringValue("context")),
		VerificationPaths: repeatedFlagValues(args, "verification", opts.stringValue("verification")),
		CIState:           opts.stringValue("ci-state"),
		CreatedBy:         opts.stringValue("created-by"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, _ := json.MarshalIndent(packet, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runPRReviewRun(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "pr-review run"}
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	opts.setBool("preview", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review run accepts only flags")
		return exitUsage
	}
	packet, profile, ok := readPRReviewPacketAndProfile(opts, stderr)
	if !ok {
		return exitCannotVerify
	}
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	runs, preview, err := prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:         opts.stringValue("out"),
		AllowedRunners: allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:        opts.boolValue("preview"),
		WorkDir:        opts.stringValue("work-dir"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	var payload []byte
	if preview != nil {
		payload, _ = json.MarshalIndent(preview, "", "  ")
	} else {
		payload, _ = json.MarshalIndent(runs, "", "  ")
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runPRReviewSynthesize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "pr-review synthesize"}
	opts.setString("packet", "")
	opts.setString("runs", "")
	opts.setString("existing-ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review synthesize accepts only flags")
		return exitUsage
	}
	if err := requireOutputFile("pr-review synthesize", opts.stringValue("out")); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	var existing *prreview.Ledger
	if opts.stringValue("existing-ledger") != "" {
		ledger, err := prreview.ReadLedger(opts.stringValue("existing-ledger"))
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		existing = &ledger
	}
	ledger := prreview.SynthesizeLedger(packet, runs, existing)
	if err := prreview.WriteJSON(opts.stringValue("out"), ledger); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(ledger, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runPRReviewValidate(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "pr-review validate"}
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("runs", "")
	opts.setString("ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review validate accepts only flags")
		return exitUsage
	}
	if err := requireOutputFile("pr-review validate", opts.stringValue("out")); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	packet, profile, ok := readPRReviewPacketAndProfile(opts, stderr)
	if !ok {
		return exitCannotVerify
	}
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	validation := prreview.Validate(packet, profile, runs, ledger)
	if err := prreview.WriteJSON(opts.stringValue("out"), validation); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(validation, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	if reviewValidationExitCode(validation) != 0 {
		return exitCannotVerify
	}
	return 0
}

func runPRReviewSummarize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "pr-review summarize"}
	opts.setString("validation", "")
	opts.setString("ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review summarize accepts only flags")
		return exitUsage
	}
	validation, err := prreview.ReadValidation(opts.stringValue("validation"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	summary := prreview.Summarize(validation, ledger)
	if opts.stringValue("out") != "" {
		if err := refuseExistingFile(opts.stringValue("out")); err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
		if err := os.WriteFile(opts.stringValue("out"), []byte(summary), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	fmt.Fprint(stdout, summary)
	return 0
}

func runPRReviewCheck(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewCheckArgs(args, stderr)
	if !ok {
		return code
	}
	packet, profile, code, ok := preparePRReviewCheck(opts, args, stderr)
	if !ok {
		return code
	}
	runs, preview, code, ok := executePRReviewCheck(packet, profile, opts, args, stderr)
	if !ok {
		return code
	}
	return finishPRReviewCheck(opts.stringValue("out"), packet, profile, runs, preview, stdout, stderr)
}

func parsePRReviewCheckArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review check"}
	opts.setString("out", "")
	opts.setString("repo-id", "")
	opts.setString("change-ref", "")
	opts.setString("base", "")
	opts.setString("head", "")
	opts.setString("diff", "")
	opts.setString("metadata", "")
	opts.setString("context", "")
	opts.setString("verification", "")
	opts.setString("profile", "")
	opts.setString("ci-state", prreview.StateNotAssessed)
	opts.setString("created-by", "sdp-trace-cli")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	opts.setBool("preview", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "pr-review check accepts only flags")
		return nil, exitUsage, false
	}
	if err := requirePRReviewCheckInputs(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func requirePRReviewCheckInputs(opts *flagSet) error {
	outDir := opts.stringValue("out")
	if strings.TrimSpace(outDir) == "" {
		return errors.New("pr-review check requires --out")
	}
	return requirePRReviewPacketInputs(opts)
}

func preparePRReviewCheck(opts *flagSet, args []string, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, int, bool) {
	outDir := opts.stringValue("out")
	packet, err := prreview.BuildPacket(prreview.PacketOptions{
		OutDir:            filepath.Join(outDir, "packet"),
		RepoID:            opts.stringValue("repo-id"),
		ChangeRef:         opts.stringValue("change-ref"),
		BaseCommit:        opts.stringValue("base"),
		HeadCommit:        opts.stringValue("head"),
		DiffPath:          opts.stringValue("diff"),
		MetadataPath:      opts.stringValue("metadata"),
		ContextPaths:      repeatedFlagValues(args, "context", opts.stringValue("context")),
		VerificationPaths: repeatedFlagValues(args, "verification", opts.stringValue("verification")),
		CIState:           opts.stringValue("ci-state"),
		CreatedBy:         opts.stringValue("created-by"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	return packet, profile, 0, true
}

func executePRReviewCheck(packet prreview.Packet, profile prreview.ReviewProfile, opts *flagSet, args []string, stderr io.Writer) (prreview.RunSet, *prreview.RunPreview, int, bool) {
	outDir := opts.stringValue("out")
	runs, preview, err := prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:         filepath.Join(outDir, "runs"),
		AllowedRunners: allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:        opts.boolValue("preview"),
		WorkDir:        opts.stringValue("work-dir"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.RunSet{}, nil, exitCannotVerify, false
	}
	return runs, preview, 0, true
}

func finishPRReviewCheck(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, preview *prreview.RunPreview, stdout, stderr io.Writer) int {
	if writePRReviewCheckPreview(stdout, preview) {
		return 0
	}
	ledger, validation, code, ok := writePRReviewCheckArtifacts(outDir, packet, profile, runs, stderr)
	if !ok {
		return code
	}
	fmt.Fprint(stdout, prreview.Summarize(validation, ledger))
	return reviewValidationExit(validation)
}

func writePRReviewCheckPreview(stdout io.Writer, preview *prreview.RunPreview) bool {
	if preview == nil {
		return false
	}
	writeIndentedPayload(stdout, preview)
	return true
}

func writePRReviewCheckArtifacts(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, stderr io.Writer) (prreview.Ledger, prreview.Validation, int, bool) {
	if !writePRReviewJSON(filepath.Join(outDir, "runs", "results.json"), runs, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	ledger := prreview.SynthesizeLedger(packet, runs, nil)
	validation := prreview.Validate(packet, profile, runs, ledger)
	if !writePRReviewJSON(filepath.Join(outDir, "ledger.json"), ledger, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	if !writePRReviewJSON(filepath.Join(outDir, "validation.json"), validation, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	return ledger, validation, 0, true
}

func writePRReviewJSON(path string, value any, stderr io.Writer) bool {
	if err := prreview.WriteJSON(path, value); err != nil {
		fmt.Fprintln(stderr, err)
		return false
	}
	return true
}

func reviewValidationExit(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		return exitCannotVerify
	}
	return 0
}

func writeIndentedPayload(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
}

func requirePRReviewPacketInputs(opts *flagSet) error {
	required := map[string]string{
		"--out":        opts.stringValue("out"),
		"--repo-id":    opts.stringValue("repo-id"),
		"--change-ref": opts.stringValue("change-ref"),
		"--base":       opts.stringValue("base"),
		"--head":       opts.stringValue("head"),
		"--diff":       opts.stringValue("diff"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("pr-review packet requires %s", flag)
		}
	}
	return nil
}

func requireOutputFile(command, path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("%s requires --out", command)
	}
	return refuseExistingFile(path)
}

func refuseExistingFile(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("output path is a directory: %s", path)
		}
		return fmt.Errorf("output file exists: %s", path)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func requireDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("work-dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("work-dir is not a directory: %s", path)
	}
	return nil
}

func reviewValidationExitCode(validation prreview.Validation) int {
	switch validation.ReviewCoverageState {
	case prreview.CoverageCannotVerify, prreview.CoverageUnresolved:
		return exitCannotVerify
	default:
		return 0
	}
}

func readPRReviewPacketAndProfile(opts *flagSet, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, bool) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, false
	}
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, false
	}
	return packet, profile, true
}

func repeatedFlagValues(args []string, key, parsedFallback string) []string {
	prefix := "--" + key + "="
	values := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, prefix) {
			values = append(values, strings.TrimPrefix(arg, prefix))
			continue
		}
		if arg == "--"+key && i+1 < len(args) {
			values = append(values, args[i+1])
			i++
		}
	}
	if len(values) == 0 && strings.TrimSpace(parsedFallback) != "" {
		values = append(values, parsedFallback)
	}
	return values
}

func allowedRunnerSet(values []string) map[string]bool {
	allowed := map[string]bool{}
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				allowed[item] = true
			}
		}
	}
	return allowed
}

func runReleaseProof(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "release-proof"}
	opts.setString("manifest", "examples/contract-foundation/contract-manifest.example.json")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "release-proof accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "release-proof requires --out")
		return exitUsage
	}
	repoRoot, err := releaseproof.RepoRoot(".")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	result, err := releaseproof.Evaluate(repoRoot, opts.stringValue("manifest"), time.Now())
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := releaseproof.Write(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	switch result.ReleaseVerificationState {
	case releaseproof.StatePass:
		return 0
	case releaseproof.StateFail:
		return 1
	default:
		return exitCannotVerify
	}
}

func runInteraction(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "interaction requires relay, import-transcript, or summarize")
		return exitUsage
	}
	switch args[0] {
	case "relay":
		return runInteractionRelay(ctx, args[1:], stdout, stderr)
	case "import-transcript":
		return runInteractionImportTranscript(args[1:], stdout, stderr)
	case "summarize":
		return runInteractionSummarize(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown interaction command: %s\n", args[0])
		return exitUsage
	}
}

func runInteractionRelay(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "interaction relay"}
	opts.setString("task-id", "")
	opts.setString("actor-type", "human_user")
	opts.setString("actor-id", "")
	opts.setString("target", "agent")
	opts.setString("event-type", "corrective_feedback")
	opts.setString("operation-id", "")
	opts.setString("stage-id", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("task-id")) == "" {
		fmt.Fprintln(stderr, "interaction relay requires --task-id")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "interaction relay requires --out")
		return exitUsage
	}
	if len(opts.rest()) == 0 {
		fmt.Fprintln(stderr, "interaction relay requires forward command after --")
		return exitUsage
	}
	_, exitCode, err := interaction.Relay(ctx, interaction.RelayOptions{
		TaskID:      opts.stringValue("task-id"),
		ActorType:   opts.stringValue("actor-type"),
		ActorID:     opts.stringValue("actor-id"),
		Target:      opts.stringValue("target"),
		EventType:   opts.stringValue("event-type"),
		OperationID: opts.stringValue("operation-id"),
		StageID:     opts.stringValue("stage-id"),
		Out:         opts.stringValue("out"),
		Command:     opts.rest(),
	}, cliStdin, stdout, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return exitCode
}

func runInteractionImportTranscript(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "interaction import-transcript"}
	opts.setString("source", "")
	opts.setString("source-ref", "")
	opts.setString("task-id", "")
	opts.setString("events-jsonl", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "interaction import-transcript accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("task-id")) == "" {
		fmt.Fprintln(stderr, "interaction import-transcript requires --task-id")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("events-jsonl")) == "" {
		fmt.Fprintln(stderr, "interaction import-transcript requires --events-jsonl")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "interaction import-transcript requires --out")
		return exitUsage
	}
	trace, err := interaction.ImportTranscript(interaction.ImportOptions{
		TaskID:      opts.stringValue("task-id"),
		Source:      opts.stringValue("source"),
		SourceRef:   opts.stringValue("source-ref"),
		EventsJSONL: opts.stringValue("events-jsonl"),
		Out:         opts.stringValue("out"),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	payload, _ := json.MarshalIndent(trace, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runInteractionSummarize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "interaction summarize"}
	opts.setString("trace", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "interaction summarize accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("trace")) == "" {
		fmt.Fprintln(stderr, "interaction summarize requires --trace")
		return exitUsage
	}
	trace, err := interaction.ReadTrace(opts.stringValue("trace"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	summary := interaction.SummarizeTrace(trace)
	if strings.TrimSpace(opts.stringValue("out")) != "" {
		if err := writeJSONFile(opts.stringValue("out"), summary); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	payload, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runEnvelope(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "summarize" {
		fmt.Fprintln(stderr, "envelope requires summarize")
		return exitUsage
	}
	opts := &flagSet{name: "envelope summarize"}
	opts.setString("envelope", "")
	opts.setString("out", "")
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "envelope summarize accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("envelope")) == "" {
		fmt.Fprintln(stderr, "envelope summarize requires --envelope")
		return exitUsage
	}
	envelope, err := interaction.ReadEnvelope(opts.stringValue("envelope"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	summary := interaction.SummarizeEnvelope(envelope)
	if strings.TrimSpace(opts.stringValue("out")) != "" {
		if err := writeJSONFile(opts.stringValue("out"), summary); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	payload, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runAssess(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if code, ok := runAssessSubcommand(args, stdout, stderr); ok {
		return code
	}
	opts, ok := parseAssessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	handler, ok := assessHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout, stderr)
}

func runAssessSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	if len(args) == 0 {
		return 0, false
	}
	switch args[0] {
	case "preview":
		return runAssessPreview(args[1:], stdout, stderr), true
	case "explain":
		return runAssessExplain(args[1:], stdout, stderr), true
	default:
		return 0, false
	}
}

func parseAssessOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := &flagSet{name: "assess"}
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("contract", "")
	opts.setString("run", "")
	opts.setString("adapter-registry", "")
	opts.setString("managed-policy", "")
	opts.setString("managed-witness", "")
	opts.setString("redaction-policy", "")
	opts.setString("artifact-manifest", "")
	opts.setString("authority-package", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "assess accepts only flags")
		return nil, false
	}
	return opts, true
}

type assessHandler func(*flagSet, io.Writer, io.Writer) int

func assessHandlers() map[string]assessHandler {
	return map[string]assessHandler{
		"adapter-capture":         runAdapterCaptureAssess,
		"managed-harness":         runManagedAssess,
		"forensic-retention":      runForensicAssess,
		"ci-artifact-observation": runCIArtifactAssess,
		"authority-envelope":      runAuthorityAssess,
	}
}

func runAdapterCaptureAssess(opts *flagSet, stdout, stderr io.Writer) int {
	required := map[string]string{
		"--out": opts.stringValue("out"),
		"--run": opts.stringValue("run"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "adapter capture assess requires %s\n", flag)
			return exitUsage
		}
	}
	input, err := loadAdapterCaptureInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := adaptercapture.Evaluate(input)
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return adapterCaptureExitCode(result)
}

func runManagedAssess(opts *flagSet, stdout, stderr io.Writer) int {
	required := map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--adapter-registry": opts.stringValue("adapter-registry"),
		"--managed-policy":   opts.stringValue("managed-policy"),
		"--managed-witness":  opts.stringValue("managed-witness"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "managed assess requires %s\n", flag)
			return exitUsage
		}
	}
	input, err := loadManagedInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := managed.Evaluate(input)
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return managedExitCode(result)
}

func runForensicAssess(opts *flagSet, stdout, stderr io.Writer) int {
	required := map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--redaction-policy": opts.stringValue("redaction-policy"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "forensic assess requires %s\n", flag)
			return exitUsage
		}
	}
	input, err := loadForensicInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := forensic.Evaluate(input)
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return forensicExitCode(result)
}

func runCIArtifactAssess(opts *flagSet, stdout, stderr io.Writer) int {
	required := map[string]string{
		"--out":               opts.stringValue("out"),
		"--artifact-manifest": opts.stringValue("artifact-manifest"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "ci artifact observation assess requires %s\n", flag)
			return exitUsage
		}
	}
	var manifest ciartifact.Manifest
	if err := readJSONFile(opts.stringValue("artifact-manifest"), &manifest); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := ciartifact.Evaluate(manifest)
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return ciArtifactExitCode(result)
}

func loadManagedInput(opts *flagSet) (managed.Input, error) {
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		return managed.Input{}, err
	}
	var policy managed.Policy
	if err := readJSONFile(opts.stringValue("managed-policy"), &policy); err != nil {
		return managed.Input{}, err
	}
	var registry managed.Registry
	if err := readJSONFile(opts.stringValue("adapter-registry"), &registry); err != nil {
		return managed.Input{}, err
	}
	var runEvidence managed.RunEvidence
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return managed.Input{}, err
	}
	var witness managed.Witness
	if err := readJSONFile(opts.stringValue("managed-witness"), &witness); err != nil {
		return managed.Input{}, err
	}
	return managed.Input{
		Contract: managed.Contract{RequiredEventTypes: append([]string(nil), contract.RequiredEvents...)},
		Policy:   policy,
		Registry: registry,
		Run:      runEvidence,
		Witness:  witness,
	}, nil
}

func loadForensicInput(opts *flagSet) (forensic.Input, error) {
	var policy forensic.Policy
	if err := readJSONFile(opts.stringValue("redaction-policy"), &policy); err != nil {
		return forensic.Input{}, err
	}
	var runEvidence forensic.RunEvidence
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return forensic.Input{}, err
	}
	return forensic.Input{Policy: policy, Run: runEvidence}, nil
}

func loadAdapterCaptureInput(opts *flagSet) (adaptercapture.Input, error) {
	var runEvidence adaptercapture.RunEvidence
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return adaptercapture.Input{}, err
	}
	return adaptercapture.Input{Run: runEvidence}, nil
}

type managedPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runAssessPreview(args []string, stdout, stderr io.Writer) int {
	opts, ok := parseAssessPreviewOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	handler, ok := assessPreviewHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess preview requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout)
}

func parseAssessPreviewOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := &flagSet{name: "assess preview"}
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("run", "")
	opts.setString("adapter-registry", "")
	opts.setString("managed-policy", "")
	opts.setString("managed-witness", "")
	opts.setString("redaction-policy", "")
	opts.setString("artifact-manifest", "")
	opts.setString("authority-package", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "assess preview accepts only flags")
		return nil, false
	}
	return opts, true
}

type assessPreviewHandler func(*flagSet, io.Writer) int

func assessPreviewHandlers() map[string]assessPreviewHandler {
	return map[string]assessPreviewHandler{
		"adapter-capture":         runAdapterCaptureAssessPreview,
		"managed-harness":         runManagedAssessPreview,
		"forensic-retention":      runForensicAssessPreview,
		"ci-artifact-observation": runCIArtifactAssessPreview,
		"authority-envelope":      runAuthorityAssessPreview,
	}
}

func runAuthorityAssess(opts *flagSet, stdout, stderr io.Writer) int {
	required := map[string]string{
		"--out":               opts.stringValue("out"),
		"--authority-package": opts.stringValue("authority-package"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "authority-envelope assess requires %s\n", flag)
			return exitUsage
		}
	}
	pkg, err := authority.ReadPackage(opts.stringValue("authority-package"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	result := authority.Evaluate(pkg)
	if err := authority.Write(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return authorityExitCode(result)
}

type adapterCapturePreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ExpectedEvidence map[string]string `json:"expected_evidence"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

func runAdapterCaptureAssessPreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"run": managedInputStatus(opts.stringValue("run")),
	}
	report := adapterCapturePreviewReport{
		Command:         "assess preview",
		SelectedProfile: adaptercapture.ProfileAdapterCapture,
		Inputs:          inputs,
		ExpectedEvidence: map[string]string{
			"binding_modes":        "same_chain,adapter_bundle",
			"test_provenance":      "ci_executed,wrapper_executed,harness_observed,agent_reported,cannot_verify",
			"capture_depth_states": "captured,missing_telemetry,not_integrated,unsupported,retention_limited,not_assessed,cannot_verify",
		},
		Safety: map[string]string{
			"raw_payloads":    "not_rendered",
			"adapter_secrets": "not_rendered",
			"gateway_refs":    "token_free_refs_only",
			"model_payloads":  "digest_or_block18_reference_only",
		},
		NextActions: adapterCapturePreviewActions(inputs),
		Claim:       "preview is read-only and does not emit an adapter capture verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

func runManagedAssessPreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"adapter_registry": managedInputStatus(opts.stringValue("adapter-registry")),
		"managed_policy":   managedInputStatus(opts.stringValue("managed-policy")),
		"managed_witness":  managedInputStatus(opts.stringValue("managed-witness")),
	}
	report := managedPreviewReport{
		Command:         "assess preview",
		SelectedProfile: managed.ProfileManagedHarness,
		Inputs:          inputs,
		NextActions:     managedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a managed verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

type forensicPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	PolicyEffects   map[string]string `json:"policy_effects"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runForensicAssessPreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"redaction_policy": managedInputStatus(opts.stringValue("redaction-policy")),
	}
	report := forensicPreviewReport{
		Command:         "assess preview",
		SelectedProfile: forensic.ProfileForensicRetention,
		Inputs:          inputs,
		PolicyEffects: map[string]string{
			"redaction_engine": "not_executed_in_preview",
			"matched_values":   "not_rendered",
			"rule_refs":        "shown_when_present_in_policy_or_run_metadata",
			"retention_modes":  "digest_only,sanitized_excerpt,encrypted_raw_ref,external_artifact_ref,not_assessed",
		},
		NextActions: forensicPreviewActions(inputs),
		Claim:       "preview is read-only and does not emit a forensic verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

type ciArtifactPreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ObservedFamilies []string          `json:"observed_families"`
	StateModel       map[string]string `json:"state_model"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

type authorityPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	StateModel      map[string]string `json:"state_model"`
	Safety          map[string]string `json:"safety"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runAuthorityAssessPreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"authority_package": managedInputStatus(opts.stringValue("authority-package")),
	}
	report := authorityPreviewReport{
		Command:         "assess preview",
		SelectedProfile: authority.Profile,
		Inputs:          inputs,
		StateModel: map[string]string{
			"authority":   "within_authority,outside_authority,not_assessed,cannot_verify",
			"attribution": "verified,not_assessed,cannot_verify",
			"binding":     "verified,not_assessed,cannot_verify",
		},
		Safety: map[string]string{
			"raw_prompts":       "not_accepted",
			"raw_model_outputs": "not_accepted",
			"credential_refs":   "rejected_as_malformed",
			"policy_effects":    "not_emitted",
		},
		NextActions: authorityPreviewActions(inputs),
		Claim:       "preview is read-only and does not emit an authority or policy verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

func runCIArtifactAssessPreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"artifact_manifest": managedInputStatus(opts.stringValue("artifact-manifest")),
	}
	report := ciArtifactPreviewReport{
		Command:         "assess preview",
		SelectedProfile: ciartifact.ProfileCIArtifactObservation,
		Inputs:          inputs,
		ObservedFamilies: []string{
			"run", "report", "witness", "provenance", "evidence",
			"trace", "artifact_index", "redaction_scan", "review", "change_ci",
		},
		StateModel: map[string]string{
			"top_level": "pass,fail,cannot_verify,not_assessed",
			"producer":  "ci_uploaded,checked_in,local_generated,agent_reported,harness_observed,external_artifact_ref,not_assessed",
			"access":    "present,absent,partial,expired,inaccessible,malformed,unsafe,not_assessed,cannot_verify",
		},
		Safety: map[string]string{
			"raw_artifact_content": "not_rendered",
			"reason_payloads":      "safe_reason_codes_only",
			"network_fetch":        "not_performed",
		},
		NextActions: ciArtifactPreviewActions(inputs),
		Claim:       "preview is read-only and does not emit a CI artifact observation verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

func runAssessExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "assess explain"}
	opts.setString("assessment-result", "")
	path, err := parseAssessExplainArgs(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	return explainAssessmentResult(path, stdout, stderr)
}

func parseAssessExplainArgs(opts *flagSet, args []string) (string, error) {
	if err := opts.parse(args); err != nil {
		return "", err
	}
	if len(opts.rest()) != 0 {
		return "", errors.New("assess explain accepts only flags")
	}
	path := opts.stringValue("assessment-result")
	if path == "" {
		return "", errors.New("assess explain requires --assessment-result <file>")
	}
	return path, nil
}

func explainAssessmentResult(path string, stdout, stderr io.Writer) int {
	var envelope struct {
		SchemaVersion   string `json:"schema_version"`
		SelectedProfile string `json:"selected_profile"`
	}
	if err := readJSONFile(path, &envelope); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return dispatchAssessmentExplanation(path, envelope.SchemaVersion, stdout, stderr)
}

func dispatchAssessmentExplanation(path, schemaVersion string, stdout, stderr io.Writer) int {
	handler, ok := assessmentExplainHandlers[schemaVersion]
	if !ok {
		fmt.Fprintf(stderr, "unsupported assessment-result schema_version: %s\n", schemaVersion)
		return exitCannotVerify
	}
	return handler(path, stdout, stderr)
}

type assessmentExplainHandler func(string, io.Writer, io.Writer) int

var assessmentExplainHandlers = map[string]assessmentExplainHandler{
	adaptercapture.SchemaVersion:  explainTypedAssessment[adaptercapture.AssessmentResult](explainAdapterCaptureAssessment),
	managed.SchemaVersion:         explainTypedAssessment[managed.AssessmentResult](explainManagedAssessment),
	forensic.SchemaVersion:        explainTypedAssessment[forensic.AssessmentResult](explainForensicAssessment),
	ciartifact.SchemaVersion:      explainTypedAssessment[ciartifact.ObservationResult](explainCIArtifactObservation),
	authority.ResultSchemaVersion: explainTypedAssessment[authority.Result](explainAuthorityEvaluation),
}

func explainTypedAssessment[T any](explain func(T, io.Writer) int) assessmentExplainHandler {
	return func(path string, stdout, stderr io.Writer) int {
		var result T
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explain(result, stdout)
	}
}

func explainAdapterCaptureAssessment(result adaptercapture.AssessmentResult, stdout io.Writer) int {
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Adapter capture assessment: %s\n", result.AdapterCaptureAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.AdapterCaptureConditions {
		fmt.Fprintf(stdout, "Adapter condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
		if condition.CappedToRetentionMode != "" {
			fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
		}
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
}

func explainManagedAssessment(result managed.AssessmentResult, stdout io.Writer) int {
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Managed harness assessment: %s\n", result.ManagedHarnessAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ManagedConditions {
		fmt.Fprintf(stdout, "Managed condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
}

func explainForensicAssessment(result forensic.AssessmentResult, stdout io.Writer) int {
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Forensic retention assessment: %s\n", result.ForensicRetentionAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ForensicConditions {
		fmt.Fprintf(stdout, "Forensic condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
		if condition.CappedToRetentionMode != "" {
			fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
		}
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
}

func explainCIArtifactObservation(result ciartifact.ObservationResult, stdout io.Writer) int {
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "CI artifact observation: %s\n", result.ArtifactObservationState)
	fmt.Fprintf(stdout, "Authority scope: %s\n", result.AuthorityScope)
	fmt.Fprintf(stdout, "Producer scope: %s\n", result.ProducerScope)
	fmt.Fprintf(stdout, "Artifact access state: %s\n", result.ArtifactAccessState)
	for _, family := range result.ArtifactFamilies {
		fmt.Fprintf(stdout, "Artifact family %s: %s (%s)\n", family.Family, family.FamilyState, family.ReasonCode)
		fmt.Fprintf(stdout, "  Producer scope: %s\n", family.ProducerScope)
		fmt.Fprintf(stdout, "  Artifact access: %s\n", family.ArtifactAccessState)
		fmt.Fprintf(stdout, "  Binding: %s\n", family.BindingState)
	}
	fmt.Fprintf(stdout, "Artifact index: %s (%s)\n", result.ArtifactIndex.Result, result.ArtifactIndex.ReasonCode)
	fmt.Fprintf(stdout, "Output safety: %s (%s)\n", result.OutputSafety.State, result.OutputSafety.ReasonCode)
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
}

func explainAuthorityEvaluation(result authority.Result, stdout io.Writer) int {
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Authority evaluation: %s\n", result.AuthorityEvaluationState)
	fmt.Fprintf(stdout, "Selected policy: %s\n", result.SelectedPolicyID)
	for _, eval := range result.Evaluations {
		fmt.Fprintf(stdout, "Observed action %s: %s (%s)\n", eval.EventID, eval.State, eval.ReasonCode)
		fmt.Fprintf(stdout, "  Actor attribution: %s\n", eval.ActorAttribution)
		fmt.Fprintf(stdout, "  Tool attribution: %s\n", eval.ToolAttribution)
		fmt.Fprintf(stdout, "  Model attribution: %s\n", eval.ModelAttribution)
		if eval.MatchedRuleRef != "" {
			fmt.Fprintf(stdout, "  Matched rule: %s\n", eval.MatchedRuleRef)
		}
	}
	for _, binding := range result.BindingEvaluations {
		fmt.Fprintf(stdout, "Binding %s: %s (%s)\n", binding.BindingID, binding.BindingState, binding.ReasonCode)
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
}

func managedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		return "absent"
	}
	info, err := os.Stat(path)
	if err != nil {
		return "present_unreadable"
	}
	if info.IsDir() {
		data, err := os.ReadFile(filepath.Join(path, "run.json"))
		if err != nil {
			return "present_unreadable"
		}
		var raw any
		if err := json.Unmarshal(data, &raw); err != nil {
			return "present_malformed"
		}
		return "present_readable"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "present_unreadable"
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return "present_malformed"
	}
	return "present_readable"
}

func managedPreviewActions(inputs map[string]string) []string {
	order := []string{"run", "adapter_registry", "managed_policy", "managed_witness"}
	var actions []string
	for _, key := range order {
		switch inputs[key] {
		case "absent":
			actions = append(actions, "Supply "+key+" before managed assessment.")
		case "present_unreadable", "present_malformed":
			actions = append(actions, "Fix "+key+" so it is readable JSON or a run directory.")
		}
	}
	return actions
}

func forensicPreviewActions(inputs map[string]string) []string {
	order := []string{"run", "redaction_policy"}
	var actions []string
	for _, key := range order {
		switch inputs[key] {
		case "absent":
			actions = append(actions, "Supply "+key+" before forensic retention assessment.")
		case "present_unreadable", "present_malformed":
			actions = append(actions, "Fix "+key+" so it is readable JSON or a run directory.")
		}
	}
	return actions
}

func adapterCapturePreviewActions(inputs map[string]string) []string {
	switch inputs["run"] {
	case "absent":
		return []string{"Supply run before adapter capture assessment."}
	case "present_unreadable", "present_malformed":
		return []string{"Fix run so it is a readable JSON run directory."}
	default:
		return nil
	}
}

func ciArtifactPreviewActions(inputs map[string]string) []string {
	switch inputs["artifact_manifest"] {
	case "absent":
		return []string{"Supply artifact manifest before CI artifact observation assessment."}
	case "present_unreadable", "present_malformed":
		return []string{"Fix artifact manifest so it is readable JSON."}
	default:
		return nil
	}
}

func authorityPreviewActions(inputs map[string]string) []string {
	switch inputs["authority_package"] {
	case "absent":
		return []string{"Supply authority package before authority envelope assessment."}
	case "present_unreadable", "present_malformed":
		return []string{"Fix authority package so it is readable JSON."}
	default:
		return nil
	}
}

func adapterCaptureExitCode(result adaptercapture.AssessmentResult) int {
	switch result.AdapterCaptureAssessment {
	case adaptercapture.StatePass:
		return 0
	case adaptercapture.StateFail:
		return 1
	default:
		return exitCannotVerify
	}
}

func managedExitCode(result managed.AssessmentResult) int {
	switch result.ManagedHarnessAssessment {
	case managed.StatePass:
		return 0
	case managed.StateFail:
		return 1
	default:
		return exitCannotVerify
	}
}

func forensicExitCode(result forensic.AssessmentResult) int {
	switch result.ForensicRetentionAssessment {
	case forensic.StatePass:
		return 0
	case forensic.StateFail:
		return 1
	default:
		return exitCannotVerify
	}
}

func ciArtifactExitCode(result ciartifact.ObservationResult) int {
	switch result.ArtifactObservationState {
	case ciartifact.StatePass:
		return 0
	case ciartifact.StateFail:
		return 1
	default:
		return exitCannotVerify
	}
}

func authorityExitCode(result authority.Result) int {
	switch result.AuthorityEvaluationState {
	case authority.StateWithinAuthority:
		return 0
	case authority.StateOutsideAuthority:
		return 1
	case authority.StateNotAssessed:
		return exitCannotVerify
	default:
		return exitCannotVerify
	}
}

func runCheckpoint(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return exitUsage
	}
	switch args[0] {
	case "create":
		return runCheckpointCreate(args[1:], stdout, stderr)
	case "verify":
		return runCheckpointVerify(args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return exitUsage
	}
}

func runCheckpointCreate(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "checkpoint create"}
	opts.setString("run", "")
	opts.setString("out", "")
	opts.setString("private-key", "")
	opts.setString("signer-id", "")
	opts.setString("id", "checkpoint-001")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "checkpoint create accepts only flags")
		return exitUsage
	}
	required := map[string]string{
		"--run":         opts.stringValue("run"),
		"--out":         opts.stringValue("out"),
		"--private-key": opts.stringValue("private-key"),
		"--signer-id":   opts.stringValue("signer-id"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "checkpoint create requires %s\n", flag)
			return exitUsage
		}
	}
	var key checkpoint.KeyPair
	if err := readJSONFile(opts.stringValue("private-key"), &key); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	created, err := checkpoint.Create(opts.stringValue("run"), checkpoint.CreateOptions{
		CheckpointID: opts.stringValue("id"),
		SignerID:     opts.stringValue("signer-id"),
		Key:          key,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := writeJSONFile(opts.stringValue("out"), created); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "checkpoint: %s\n", created.CheckpointID)
	return 0
}

func runCheckpointVerify(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "checkpoint verify"}
	opts.setString("run", "")
	opts.setString("checkpoint", "")
	opts.setString("policy", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "checkpoint verify accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		fmt.Fprintln(stderr, "checkpoint verify requires --run")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("checkpoint")) == "" {
		fmt.Fprintln(stderr, "checkpoint verify requires --checkpoint")
		return exitUsage
	}
	var signed checkpoint.SignedCheckpoint
	if err := readJSONFile(opts.stringValue("checkpoint"), &signed); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	var policy *checkpoint.TrustedCheckpointPolicy
	if opts.stringValue("policy") != "" {
		var loaded checkpoint.TrustedCheckpointPolicy
		if err := readJSONFile(opts.stringValue("policy"), &loaded); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		policy = &loaded
	}
	result := checkpoint.Verify(opts.stringValue("run"), signed, policy)
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	switch result.Result {
	case checkpoint.StatePass:
		return 0
	case checkpoint.StateCannotVerify:
		return exitCannotVerify
	default:
		return 1
	}
}

func runReport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "report"}
	opts.setString("out", "")
	opts.setString("contract", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "report requires <runs-root-or-run-dir>")
		return exitUsage
	}
	outDir := opts.stringValue("out")
	if outDir == "" {
		fmt.Fprintln(stderr, "report requires --out <dir>")
		return exitUsage
	}
	artifacts, err := demo.WriteReport(targets[0], outDir, opts.stringValue("contract"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(artifacts.Summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	if artifacts.Summary.CannotVerifyCount > 0 {
		return exitCannotVerify
	}
	if artifacts.Summary.FailedCount > 0 {
		return 1
	}
	return 0
}

func runGate(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "explain":
			return runGateExplain(args[1:], stdout, stderr)
		case "preview":
			return runGatePreview(args[1:], stdout, stderr)
		}
	}
	opts := &flagSet{name: "gate"}
	opts.setString("out", "")
	opts.setString("contract", "")
	opts.setString("witness", "")
	opts.setString("profile", "")
	opts.setString("checkpoint", "")
	opts.setString("checkpoint-policy", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "gate requires <runs-root-or-run-dir>")
		return exitUsage
	}
	outPath := opts.stringValue("out")
	if outPath == "" {
		fmt.Fprintln(stderr, "gate requires --out <file>")
		return exitUsage
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		return runProtectedGate(targets[0], outPath, opts, stdout, stderr)
	}
	result, err := demo.WriteGate(targets[0], outPath, opts.stringValue("contract"), opts.stringValue("witness"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return gateExitCode(result)
}

func runProtectedGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	result, code := resolveProtectedGate(target, opts, stderr)
	if code != 0 {
		return code
	}
	return writeProtectedGateResult(outPath, result, stdout, stderr)
}

func resolveProtectedGate(target string, opts *flagSet, stderr io.Writer) (demo.GateResult, int) {
	signed, policy, witnessSummary, code, ok := readProtectedGateInputs(opts, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	contract, rows, runDir, code, ok := loadProtectedGateRows(target, opts.stringValue("contract"), stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	expected, code, ok := loadProtectedWitnessExpectation(target, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	checkpointResult := protectedCheckpointVerification(
		checkpoint.Verify(runDir, signed, &policy),
		signed,
		policy,
		witnessSummary,
		expected,
	)
	return demo.EvaluateProtectedGate(rows, contract, demo.ProtectedGateInput{
		Checkpoint:         checkpointResult,
		PolicyProvided:     true,
		Witness:            &witnessSummary,
		WitnessExpectation: expected,
		Now:                time.Now().UTC(),
	}), 0
}

func writeProtectedGateResult(path string, result demo.GateResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(path, result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeIndentedPayload(stdout, result)
	return gateExitCode(result)
}

func readProtectedGateInputs(opts *flagSet, stderr io.Writer) (checkpoint.SignedCheckpoint, checkpoint.TrustedCheckpointPolicy, demo.WitnessSummary, int, bool) {
	var signed checkpoint.SignedCheckpoint
	var policy checkpoint.TrustedCheckpointPolicy
	var witness demo.WitnessSummary
	inputs := []struct {
		flag  string
		path  string
		value any
	}{
		{"--checkpoint", opts.stringValue("checkpoint"), &signed},
		{"--checkpoint-policy", opts.stringValue("checkpoint-policy"), &policy},
		{"--witness", opts.stringValue("witness"), &witness},
	}
	for _, input := range inputs {
		if strings.TrimSpace(input.path) == "" {
			fmt.Fprintf(stderr, "protected gate requires %s\n", input.flag)
			return checkpoint.SignedCheckpoint{}, checkpoint.TrustedCheckpointPolicy{}, demo.WitnessSummary{}, exitUsage, false
		}
		if err := readJSONFile(input.path, input.value); err != nil {
			fmt.Fprintln(stderr, err)
			return checkpoint.SignedCheckpoint{}, checkpoint.TrustedCheckpointPolicy{}, demo.WitnessSummary{}, exitUsage, false
		}
	}
	return signed, policy, witness, 0, true
}

func loadProtectedGateRows(target, contractPath string, stderr io.Writer) (trace.Contract, []demo.RunRow, string, int, bool) {
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	rows, err := demo.VerifiedRows(target, contract)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	runDir, err := protectedRunDir(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", exitCannotVerify, false
	}
	return contract, rows, runDir, 0, true
}

func loadProtectedWitnessExpectation(target string, stderr io.Writer) (demo.WitnessExpectation, int, bool) {
	expected, err := demoWitnessExpectation(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return demo.WitnessExpectation{}, exitCannotVerify, false
	}
	return expected, 0, true
}

func runGateExplain(args []string, stdout, stderr io.Writer) int {
	path, code, ok := parseGateExplainArgs(args, stderr)
	if !ok {
		return code
	}
	result, code, ok := readGateExplainResult(path, stderr)
	if !ok {
		return code
	}
	explainGateResult(result, stdout)
	return 0
}

func parseGateExplainArgs(args []string, stderr io.Writer) (string, int, bool) {
	opts := &flagSet{name: "gate explain"}
	opts.setString("gate-result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return "", exitUsage, false
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "gate explain accepts only flags")
		return "", exitUsage, false
	}
	path := opts.stringValue("gate-result")
	if path == "" {
		fmt.Fprintln(stderr, "gate explain requires --gate-result <file>")
		return "", exitUsage, false
	}
	return path, 0, true
}

func readGateExplainResult(path string, stderr io.Writer) (demo.GateResult, int, bool) {
	var result demo.GateResult
	if err := readJSONFile(path, &result); err != nil {
		fmt.Fprintln(stderr, err)
		return demo.GateResult{}, exitCannotVerify, false
	}
	if result.SchemaVersion != demo.GateSchemaVersion && result.SchemaVersion != demo.GateSchemaVersionBlock16 {
		fmt.Fprintf(stderr, "unsupported gate-result schema_version: %s\n", result.SchemaVersion)
		return demo.GateResult{}, exitCannotVerify, false
	}
	return result, 0, true
}

func explainGateResult(result demo.GateResult, stdout io.Writer) {
	if result.SchemaVersion == demo.GateSchemaVersion {
		fmt.Fprintln(stdout, "Protected profile fields: absent")
	}
	explainGateSummary(result, stdout)
	explainProtectedGateDetails(result, stdout)
	explainGateCollections(result, stdout)
}

func explainGateSummary(result demo.GateResult, stdout io.Writer) {
	fmt.Fprintf(stdout, "Gate mode: %s\n", result.GateMode)
	fmt.Fprintf(stdout, "Trust cap: %s\n", result.TrustCap)
	if result.SelectedProfile != "" {
		fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	}
	fmt.Fprintf(stdout, "Local gate: %s\n", result.LocalGate)
	fmt.Fprintf(stdout, "CI witness gate: %s\n", result.CIWitnessGate)
	fmt.Fprintf(stdout, "Audit-grade gate: %s\n", result.AuditGradeGate)
	if result.ProtectedGate != "" {
		fmt.Fprintf(stdout, "Protected gate: %s\n", result.ProtectedGate)
	}
}

func explainProtectedGateDetails(result demo.GateResult, stdout io.Writer) {
	if result.CheckpointVerification != nil {
		fmt.Fprintf(stdout, "Checkpoint result: %s\n", result.CheckpointVerification.Result)
		fmt.Fprintf(stdout, "Checkpoint trust scope: %s\n", result.CheckpointVerification.TrustScope)
	}
	for _, condition := range result.ProtectedConditions {
		fmt.Fprintf(stdout, "Protected condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
}

func explainGateCollections(result demo.GateResult, stdout io.Writer) {
	explainRequiredRuns(result.RequiredRuns, stdout)
	explainWitnessBindings(result.WitnessBindings, stdout)
	explainMissingAuditEvidence(result.MissingAuditEvidence, stdout)
	explainOverrideRequests(result.OverrideRequests, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
}

func explainRequiredRuns(requiredRuns []demo.RequiredRunResult, stdout io.Writer) {
	for _, requiredRun := range requiredRuns {
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
}

func explainWitnessBindings(bindings []demo.WitnessBinding, stdout io.Writer) {
	for _, binding := range bindings {
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
}

func explainMissingAuditEvidence(missingEvidence []string, stdout io.Writer) {
	for _, missing := range missingEvidence {
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
}

func explainOverrideRequests(overrides []demo.OverrideRequest, stdout io.Writer) {
	for _, override := range overrides {
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
}

func explainReasons(reasons []string, stdout io.Writer) {
	for _, reason := range reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
}

func explainNextActions(actions []string, stdout io.Writer) {
	for _, action := range actions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
}

type gatePreviewReport struct {
	Command            string   `json:"command"`
	GateMode           string   `json:"gate_mode"`
	TrustCap           string   `json:"trust_cap"`
	RequiredRuns       []string `json:"required_runs"`
	RequiredEvidence   []string `json:"required_evidence"`
	WitnessInspectable bool     `json:"witness_inspectable"`
	WitnessMismatches  []string `json:"witness_mismatches,omitempty"`
	Claim              string   `json:"claim"`
}

type protectedGatePreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	TrustCap        string            `json:"trust_cap"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runGatePreview(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "gate preview"}
	opts.setString("contract", "")
	opts.setString("witness", "")
	opts.setString("profile", "")
	opts.setString("checkpoint", "")
	opts.setString("checkpoint-policy", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "gate preview requires <runs-root-or-run-dir>")
		return exitUsage
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		return runProtectedGatePreview(opts, stdout)
	}
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	report := gatePreviewReport{
		Command:          "gate preview",
		GateMode:         previewGateMode(contract),
		TrustCap:         string(trace.TrustScopeLocalObserved),
		RequiredRuns:     requiredRunIDs(contract),
		RequiredEvidence: requiredEvidenceIDsForCLI(contract),
		Claim:            "preview is read-only and does not claim the gate will pass",
	}
	witnessPath := opts.stringValue("witness")
	if witnessPath != "" {
		report.WitnessInspectable, report.WitnessMismatches = demo.PreviewWitnessBinding(witnessPath, targets[0])
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	_ = targets[0]
	return 0
}

func runProtectedGatePreview(opts *flagSet, stdout io.Writer) int {
	inputs := map[string]string{
		"checkpoint":        protectedInputStatus(opts.stringValue("checkpoint")),
		"checkpoint_policy": protectedInputStatus(opts.stringValue("checkpoint-policy")),
		"witness":           protectedInputStatus(opts.stringValue("witness")),
	}
	report := protectedGatePreviewReport{
		Command:         "gate preview",
		SelectedProfile: demo.GateProfileProtected,
		TrustCap:        string(trace.TrustScopeLocalObserved),
		Inputs:          inputs,
		NextActions:     protectedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a protected verdict",
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			return exitCannotVerify
		}
	}
	return 0
}

func protectedRunDir(target string) (string, error) {
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return "", err
	}
	if len(runDirs) != 1 {
		return "", fmt.Errorf("protected gate requires one selected run, got %d", len(runDirs))
	}
	return runDirs[0], nil
}

func protectedCheckpointVerification(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) checkpoint.VerificationResult {
	if result.Result == checkpoint.StateFail {
		return result
	}
	if signed.Signer.Authority != checkpoint.AuthorityCIIsolatedJob {
		return result
	}
	if !policyAllowsSigner(policy, signed) {
		return result
	}
	if !witnessMatchesProtectedInput(witnessSummary, expected) {
		return result
	}
	result.SignerAuthorityState = checkpoint.StatePass
	result.TrustScope = checkpoint.TrustScopeCISigned
	result.Result = checkpoint.StatePass
	return result
}

func policyAllowsSigner(policy checkpoint.TrustedCheckpointPolicy, signed checkpoint.SignedCheckpoint) bool {
	for _, signer := range policy.AllowedSigners {
		if signer.SignerID == signed.Signer.SignerID &&
			signer.Authority == signed.Signer.Authority &&
			signer.PublicKey == signed.Signature.PublicKey {
			return true
		}
	}
	return false
}

func witnessMatchesProtectedInput(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	if !witnessHasProtectedTrust(witnessSummary) || !witnessSourceMatches(witnessSummary, expected) {
		return false
	}
	return witnessArtifactsMatch(witnessSummary.RunArtifacts, expected.RunArtifacts)
}

func witnessArtifactsMatch(runArtifacts, expectedRunArtifacts []demo.WitnessArtifactDigest) bool {
	expectedArtifacts := expectedArtifactDigests(expectedRunArtifacts)
	if len(runArtifacts) != len(expectedArtifacts) {
		return false
	}
	for _, artifact := range runArtifacts {
		if !witnessArtifactMatchesExpectation(artifact, expectedArtifacts) {
			return false
		}
	}
	return true
}

func expectedArtifactDigests(expectedRunArtifacts []demo.WitnessArtifactDigest) map[string]string {
	expectedArtifacts := map[string]string{}
	for _, artifact := range expectedRunArtifacts {
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	return expectedArtifacts
}

func witnessArtifactMatchesExpectation(artifact demo.WitnessArtifactDigest, expectedArtifacts map[string]string) bool {
	expectedSHA, ok := expectedArtifacts[artifact.Path]
	return ok && expectedSHA == artifact.SHA256
}

func witnessHasProtectedTrust(witnessSummary demo.WitnessSummary) bool {
	return witnessSummary.Kind == "github-actions" &&
		witnessSummary.Status == demo.GatePass &&
		witnessSummary.TrustScope == "ci_witnessed"
}

func witnessSourceMatches(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	return optionalStringMatches(expected.Repository, witnessSummary.Source.Repository) &&
		optionalStringMatches(expected.Ref, witnessSummary.Source.Ref) &&
		optionalStringMatches(expected.CommitSHA, witnessSummary.Source.CommitSHA) &&
		optionalStringMatches(expected.RunID, witnessSummary.CIIdentity.RunID)
}

func optionalStringMatches(expected, actual string) bool {
	return expected == "" || actual == expected
}

func demoWitnessExpectation(target string) (demo.WitnessExpectation, error) {
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return demo.WitnessExpectation{}, err
	}
	artifacts := make([]demo.WitnessArtifactDigest, 0, len(runDirs))
	runID := ""
	for _, runDir := range runDirs {
		artifact, err := trace.OpenRunArtifact(runDir)
		if err != nil {
			return demo.WitnessExpectation{}, err
		}
		if runID == "" {
			runID = artifact.Manifest.RunID
		}
		digest, err := sha256File(runDir, "run.json")
		if err != nil {
			return demo.WitnessExpectation{}, err
		}
		artifacts = append(artifacts, demo.WitnessArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	return demo.WitnessExpectation{RunID: runID, RunArtifacts: artifacts}, nil
}

func sha256File(dir, name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func protectedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		return "absent"
	}
	var value any
	if err := readJSONFile(path, &value); err != nil {
		if os.IsNotExist(err) || errors.Is(err, os.ErrPermission) {
			return "present_unreadable"
		}
		return "present_malformed"
	}
	return "present_readable"
}

func protectedPreviewActions(inputs map[string]string) []string {
	names := []string{"checkpoint", "checkpoint_policy", "witness"}
	actions := make([]string, 0)
	for _, name := range names {
		switch inputs[name] {
		case "absent":
			actions = append(actions, fmt.Sprintf("Supply %s input before running protected gate.", name))
		case "present_unreadable", "present_malformed":
			actions = append(actions, fmt.Sprintf("Replace %s input with readable JSON.", name))
		}
	}
	return actions
}

func runOverride(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "request" {
		fmt.Fprintln(stderr, "override requires request")
		return exitUsage
	}
	opts := &flagSet{name: "override request"}
	opts.setString("out", "")
	opts.setString("id", "")
	opts.setString("by", "")
	opts.setString("reason", "")
	opts.setString("source-ref", "")
	opts.setString("scope", "")
	opts.setString("external-reference", "")
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "override request accepts only flags")
		return exitUsage
	}
	runDir := opts.stringValue("out")
	required := map[string]string{
		"--out":        runDir,
		"--id":         opts.stringValue("id"),
		"--by":         opts.stringValue("by"),
		"--reason":     opts.stringValue("reason"),
		"--source-ref": opts.stringValue("source-ref"),
		"--scope":      opts.stringValue("scope"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "override request requires %s\n", flag)
			return exitUsage
		}
	}
	payload := map[string]any{
		"override_id":  opts.stringValue("id"),
		"producer":     "sdp-trace-cli",
		"origin":       "native_cli",
		"requested_by": opts.stringValue("by"),
		"reason":       opts.stringValue("reason"),
		"source_ref":   opts.stringValue("source-ref"),
		"scope":        opts.stringValue("scope"),
		"created_at":   time.Now().UTC().Format(time.RFC3339Nano),
	}
	if external := opts.stringValue("external-reference"); external != "" {
		payload["external_reference"] = external
	}
	event, err := trace.AppendRunEvent(runDir, trace.EventPolicyOverrideRequested, payload, "sdp-trace-cli")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "override_event: %s\n", event.EventID)
	return 0
}

func readJSONFile(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func writeTextFileAtomic(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.WriteString(value); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func previewGateMode(contract trace.Contract) string {
	mode := demo.GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case demo.GateModeProtectedFuture:
			return demo.GateModeProtectedFuture
		case demo.GateModeAdvisoryCI:
			mode = demo.GateModeAdvisoryCI
		}
	}
	return mode
}

func requiredRunIDs(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredRuns))
	for _, required := range contract.RequiredRuns {
		if required.ID != "" {
			ids = append(ids, required.ID)
		}
	}
	return ids
}

func requiredEvidenceIDsForCLI(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}

func gateExitCode(result demo.GateResult) int {
	if code, ok := protectedGateExitCode(result); ok {
		return code
	}
	return gateStateExitCode(gateExitStates(result))
}

func protectedGateExitCode(result demo.GateResult) (int, bool) {
	if result.SelectedProfile != demo.GateProfileProtected {
		return 0, false
	}
	code, ok := protectedGateExitCodes[result.ProtectedGate]
	if !ok {
		return 0, false
	}
	return code, true
}

var protectedGateExitCodes = map[string]int{
	demo.GatePass:         0,
	demo.GateFail:         1,
	demo.GateCannotVerify: exitCannotVerify,
	demo.GateNotAssessed:  exitCannotVerify,
}

func gateExitStates(result demo.GateResult) []string {
	states := []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate}
	for _, requiredRun := range result.RequiredRuns {
		states = append(states, requiredRun.State)
	}
	return states
}

func gateStateExitCode(states []string) int {
	if hasGateState(states, demo.GateFail, demo.GateMissingTelemetry) {
		return 1
	}
	if hasGateState(states, demo.GateCannotVerify) {
		return exitCannotVerify
	}
	return 0
}

func hasGateState(states []string, targets ...string) bool {
	for _, state := range states {
		if slices.Contains(targets, state) {
			return true
		}
	}
	return false
}

func runWitness(_ context.Context, args []string, stdout, stderr io.Writer) int {
	options, ok := parseWitnessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	record, err := buildWitnessRecord(options)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return writeWitnessRecordOutput(stdout, record)
}

type witnessOptions struct {
	kind                      string
	target                    string
	out                       string
	reportDir                 string
	witnessEnvelope           string
	customerPKIAuthorityPath  string
	customerPKIPublicCertPath string
	customerPKIPublicKeyPath  string
	customerPKIPayloadDigest  string
	customerPKIFreshnessPath  string
}

func parseWitnessOptions(args []string, stderr io.Writer) (witnessOptions, bool) {
	opts, ok := parseWitnessFlagSet(args, stderr)
	if !ok {
		return witnessOptions{}, false
	}
	options, message, ok := witnessOptionsFromFlags(opts)
	if !ok {
		fmt.Fprintln(stderr, message)
		return witnessOptions{}, false
	}
	return options, true
}

func parseWitnessFlagSet(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := &flagSet{name: "witness"}
	opts.setString("kind", "")
	opts.setString("out", "")
	opts.setString("report-dir", "")
	opts.setString("witness-envelope", "")
	opts.setString("customer-pki-authority-policy", "")
	opts.setString("customer-pki-public-cert", "")
	opts.setString("customer-pki-public-key", "")
	opts.setString("customer-pki-payload-digest", "")
	opts.setString("customer-pki-freshness-evidence", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	return opts, true
}

func witnessOptionsFromFlags(opts *flagSet) (witnessOptions, string, bool) {
	fields, message, ok := witnessRequiredFieldsFromFlags(opts)
	if !ok {
		return witnessOptions{}, message, false
	}
	return witnessOptionsFromRequiredFields(fields, opts), "", true
}

type witnessRequiredFields struct {
	target string
	kind   string
	out    string
}

func witnessRequiredFieldsFromFlags(opts *flagSet) (witnessRequiredFields, string, bool) {
	target, message, ok := witnessTargetFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessKindOutFromFlags(opts, target)
}

func witnessKindOutFromFlags(opts *flagSet, target string) (witnessRequiredFields, string, bool) {
	kind, message, ok := witnessKindFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	out, message, ok := witnessOutFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	if message, ok := validateWitnessKindFlags(kind, opts); !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessRequiredFields{target: target, kind: kind, out: out}, "", true
}

func witnessOptionsFromRequiredFields(fields witnessRequiredFields, opts *flagSet) witnessOptions {
	return witnessOptions{
		kind:                      fields.kind,
		target:                    fields.target,
		out:                       fields.out,
		reportDir:                 opts.stringValue("report-dir"),
		witnessEnvelope:           opts.stringValue("witness-envelope"),
		customerPKIAuthorityPath:  opts.stringValue("customer-pki-authority-policy"),
		customerPKIPublicCertPath: opts.stringValue("customer-pki-public-cert"),
		customerPKIPublicKeyPath:  opts.stringValue("customer-pki-public-key"),
		customerPKIPayloadDigest:  opts.stringValue("customer-pki-payload-digest"),
		customerPKIFreshnessPath:  opts.stringValue("customer-pki-freshness-evidence"),
	}
}

func witnessTargetFromFlags(opts *flagSet) (string, string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		return "", "witness requires <runs-root-or-run-dir>", false
	}
	return targets[0], "", true
}

func witnessKindFromFlags(opts *flagSet) (string, string, bool) {
	kind := opts.stringValue("kind")
	if !allowedWitnessKind(kind) {
		return "", "witness requires --kind github-actions, gitlab-ci, buildkite, or customer-pki", false
	}
	return kind, "", true
}

func witnessOutFromFlags(opts *flagSet) (string, string, bool) {
	out := opts.stringValue("out")
	if out == "" {
		return "", "witness requires --out <file>", false
	}
	return out, "", true
}

func validateWitnessKindFlags(kind string, opts *flagSet) (string, bool) {
	missing := missingWitnessKindFlags(kind, opts)
	if len(missing) > 0 {
		return fmt.Sprintf("customer-pki witness requires %s", strings.Join(missing, ", ")), false
	}
	return "", true
}

func missingWitnessKindFlags(kind string, opts *flagSet) []string {
	if kind != witness.KindCustomerPKI {
		return nil
	}
	return missingCustomerPKIFlags(opts)
}

func buildWitnessRecord(opts witnessOptions) (witness.Record, error) {
	builder, ok := witnessRecordBuilders()[opts.kind]
	if !ok {
		return witness.Record{}, fmt.Errorf("unsupported witness kind %q", opts.kind)
	}
	return builder(opts)
}

type witnessRecordBuilder func(witnessOptions) (witness.Record, error)

func witnessRecordBuilders() map[string]witnessRecordBuilder {
	return map[string]witnessRecordBuilder{
		witness.KindGitHubActions: buildGitHubActionsWitness,
		witness.KindGitLabCI:      buildEnvelopeWitness,
		witness.KindBuildkite:     buildEnvelopeWitness,
		witness.KindCustomerPKI:   buildCustomerPKIWitness,
	}
}

func buildGitHubActionsWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteGitHubActions(opts.out, opts.target, opts.reportDir, witness.EnvironmentFromOS())
}

func buildEnvelopeWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteProfile(opts.kind, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		EnvelopePath: opts.witnessEnvelope,
	})
}

func buildCustomerPKIWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteProfile(witness.KindCustomerPKI, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		CustomerPKIAuthorityPolicy: opts.customerPKIAuthorityPath,
		CustomerPKIPublicCert:      opts.customerPKIPublicCertPath,
		CustomerPKIPublicKey:       opts.customerPKIPublicKeyPath,
		CustomerPKIPayloadDigest:   opts.customerPKIPayloadDigest,
		CustomerPKIFreshness:       opts.customerPKIFreshnessPath,
	})
}

func writeWitnessRecordOutput(stdout io.Writer, record witness.Record) int {
	payload, _ := json.MarshalIndent(record, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	switch record.Status {
	case witness.StatusCannotVerify, witness.StatusNotAssessed:
		return exitCannotVerify
	case witness.StatusFail:
		return 1
	default:
		return 0
	}
}

func missingCustomerPKIFlags(opts *flagSet) []string {
	missing := []string{}
	required := map[string]string{
		"--customer-pki-authority-policy":   opts.stringValue("customer-pki-authority-policy"),
		"--customer-pki-payload-digest":     opts.stringValue("customer-pki-payload-digest"),
		"--customer-pki-freshness-evidence": opts.stringValue("customer-pki-freshness-evidence"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, flag)
		}
	}
	if strings.TrimSpace(opts.stringValue("customer-pki-public-cert")) == "" && strings.TrimSpace(opts.stringValue("customer-pki-public-key")) == "" {
		missing = append(missing, "--customer-pki-public-cert or --customer-pki-public-key")
	}
	sort.Strings(missing)
	return missing
}

func allowedWitnessKind(kind string) bool {
	switch kind {
	case witness.KindGitHubActions, witness.KindGitLabCI, witness.KindBuildkite, witness.KindCustomerPKI:
		return true
	default:
		return false
	}
}

func runWrap(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "wrap"}
	opts.setString("name", "")
	opts.setString("contract", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintln(stderr, "wrap requires a command")
		return exitUsage
	}
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: true,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
}

func runWrappedCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "run"}
	opts.setString("task", "")
	opts.setString("contract", "")
	opts.setBool("use-default-contract", false)
	opts.setString("name", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintln(stderr, "run requires a command")
		return exitUsage
	}
	useDefault := opts.boolValue("use-default-contract")
	task := opts.stringValue("task")
	if task == "" {
		fmt.Fprintln(stderr, "run requires --task")
		return exitUsage
	}
	if opts.stringValue("contract") == "" && !useDefault {
		fmt.Fprintln(stderr, "run requires --contract unless --use-default-contract is set")
		return exitUsage
	}
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		Task:               task,
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: useDefault,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
}

func runDryRun(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("dry-run", "simulation", args, stdout, stderr)
}

func runPreview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("preview", "preview", args, stdout, stderr)
}

func runPreviewCommand(commandName, mode string, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: commandName}
	opts.setString("contract", "")
	opts.setBool("use-default-contract", true)
	opts.setString("name", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintf(stderr, "%s requires a command\n", commandName)
		return exitUsage
	}
	contractPath := opts.stringValue("contract")
	useDefault := opts.boolValue("use-default-contract")
	if contractPath == "" && !useDefault {
		fmt.Fprintf(stderr, "%s requires --contract unless --use-default-contract is set\n", commandName)
		return exitUsage
	}
	contract := trace.DefaultContract
	if contractPath != "" {
		loaded, err := trace.LoadContract(contractPath)
		if err != nil {
			fmt.Fprintf(stderr, "failed to load contract: %v\n", err)
			return exitCannotVerify
		}
		contract = loaded
	}
	payload := map[string]any{
		"mode":                 mode,
		"command_descriptor":   trace.NewCommandDescriptor(command),
		"contract":             contract,
		"boundaries":           previewBoundaries(),
		"offline_implications": previewOfflineImplications(),
		"writes_artifacts":     false,
		"safe_retention_modes": safeRetentionModes(),
		"warning":              "no run artifacts were written",
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	return 0
}

func runDoctor(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "doctor"}
	opts.setString("contract", "")
	opts.setString("output-dir", defaultRunRoot)
	opts.setString("report-dir", defaultReportDir)
	opts.setString("profile", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) != "" {
		if opts.stringValue("profile") != repoobserver.ProfileGithubActionsGitHooksV1 {
			fmt.Fprintf(stderr, "doctor --profile requires %s\n", repoobserver.ProfileGithubActionsGitHooksV1)
			return exitUsage
		}
		status, err := repoobserver.Doctor(repoobserver.Options{Profile: opts.stringValue("profile")})
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		if err := repoobserver.WriteJSON(opts.stringValue("out"), status); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprint(stdout, repoobserver.HumanTable(status))
		return repoObserverExitCode(status)
	}
	report, exitCode := buildDoctorReport(doctorOptions{
		ContractPath: opts.stringValue("contract"),
		OutputDir:    opts.stringValue("output-dir"),
		ReportDir:    opts.stringValue("report-dir"),
		Env:          witness.EnvironmentFromOS(),
	})
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", data)
	return exitCode
}

func runInstall(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	if len(args) == 0 || args[0] != "repo-observer" {
		fmt.Fprintln(stderr, "install requires repo-observer")
		return exitUsage
	}
	opts := &flagSet{name: "install repo-observer"}
	opts.setString("profile", repoobserver.ProfileGithubActionsGitHooksV1)
	opts.setString("repository-id", "")
	opts.setString("out", "")
	opts.setBool("write", false)
	opts.setBool("force", false)
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "install repo-observer accepts only flags")
		return exitUsage
	}
	status, err := repoobserver.Install(repoobserver.Options{
		Profile:      opts.stringValue("profile"),
		RepositoryID: opts.stringValue("repository-id"),
		Write:        opts.boolValue("write"),
		Force:        opts.boolValue("force"),
	})
	if writeErr := repoobserver.WriteJSON(opts.stringValue("out"), status); writeErr != nil {
		fmt.Fprintln(stderr, writeErr)
		return 1
	}
	if err != nil {
		if status.SchemaVersion != "" {
			fmt.Fprint(stdout, repoobserver.HumanTable(status))
		}
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprint(stdout, repoobserver.HumanTable(status))
	if !opts.boolValue("write") {
		return 0
	}
	return repoObserverExitCode(status)
}

func repoObserverExitCode(status repoobserver.Status) int {
	if status.InstallState == repoobserver.StateCannotVerify || status.ProofState == repoobserver.StateCannotVerify {
		return exitCannotVerify
	}
	if status.InstallState == repoobserver.StateFail {
		return 1
	}
	return 0
}

func runVerify(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "verify requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
	if info, err := os.Stat(runDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "run directory does not exist: %s\n", runDir)
		return exitCannotVerify
	}
	result, table, audit, err := verifier.VerifyRun(runDir)
	if writeErr := verifier.WriteVerifierArtifacts(runDir, result, table, audit); writeErr != nil {
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, writeErr)
		return 1
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
	}
	switch result.Result {
	case trace.VerdictObserved, trace.VerdictNotAssessed:
		return 0
	case trace.VerdictFail:
		return 1
	case trace.VerdictCannotVerify:
		return exitCannotVerify
	default:
		return 0
	}
}

func runExplain(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "explain requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
	explanation, err := verifier.ExplainRun(runDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, explanation)
	return 0
}

func runQuery(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "query"}
	opts.setString("query", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	queryName := opts.stringValue("query")
	runDirs := opts.rest()
	if len(runDirs) == 0 {
		fmt.Fprintln(stderr, "query requires <run-dir>")
		return exitUsage
	}
	if queryName == query.QueryCaptureDepth {
		payload, err := query.CaptureDepth(runDirs[0])
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		fmt.Fprintf(stdout, "%s\n", payload)
		return 0
	}
	if queryName != query.QueryMissingEvidence {
		fmt.Fprintf(stderr, "unsupported query: %s\n", queryName)
		return exitUsage
	}
	payload, err := query.MissingEvidence(runDirs[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runQueryPack(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "explain" {
		return runQueryPackExplain(args[1:], stdout, stderr)
	}
	return runQueryPackBuild(args, stderr)
}

func runQueryPackBuild(args []string, stderr io.Writer) int {
	opts, err := parseQueryPackArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if err := validateQueryPackOptions(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	code, err := writeQueryPackArtifact(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return code
	}
	return 0
}

func writeQueryPackArtifact(opts *queryPackOptions) (int, error) {
	result, err := query.ForensicsBasicPack(opts.runPath)
	if err != nil {
		return exitCannotVerify, err
	}
	if err := writeJSONFile(opts.outPath, result); err != nil {
		return 1, err
	}
	return 0, nil
}

func runQueryPackExplain(args []string, stdout, stderr io.Writer) int {
	opts, err := parseQueryPackExplainArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result, err := readQueryPackResult(opts.resultPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := validateQueryPackExplainResult(result); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprint(stdout, query.ExplainForensicsPack(result))
	return 0
}

type queryPackOptions struct {
	pack    string
	runPath string
	outPath string
}

type queryPackExplainOptions struct {
	resultPath string
}

func parseQueryPackArgs(args []string) (*queryPackOptions, error) {
	opts := &flagSet{name: "query-pack"}
	opts.setString("pack", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		return nil, fmt.Errorf("query-pack accepts only flags")
	}
	return &queryPackOptions{
		pack:    strings.TrimSpace(opts.stringValue("pack")),
		runPath: strings.TrimSpace(opts.stringValue("run")),
		outPath: strings.TrimSpace(opts.stringValue("out")),
	}, nil
}

func parseQueryPackExplainArgs(args []string) (*queryPackExplainOptions, error) {
	opts := &flagSet{name: "query-pack explain"}
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		return nil, fmt.Errorf("query-pack explain accepts only flags")
	}
	resultPath := strings.TrimSpace(opts.stringValue("result"))
	if resultPath == "" {
		return nil, fmt.Errorf("query-pack explain requires --result")
	}
	return &queryPackExplainOptions{resultPath: resultPath}, nil
}

func validateQueryPackOptions(opts *queryPackOptions) error {
	if opts.pack == "" {
		return fmt.Errorf("error: ambiguous pack selection; --pack is required")
	}
	if opts.pack != query.QueryPackForensicsBasic {
		return fmt.Errorf("error: unknown pack %q", opts.pack)
	}
	return requireQueryPackRequiredInputs(opts.runPath, opts.outPath)
}

func requireQueryPackRequiredInputs(runPath, outPath string) error {
	if runPath == "" {
		return fmt.Errorf("query-pack requires --run")
	}
	if outPath == "" {
		return fmt.Errorf("query-pack requires --out")
	}
	return nil
}

func readQueryPackResult(path string) (query.QueryPackResult, error) {
	var result query.QueryPackResult
	if err := readJSONFile(path, &result); err != nil {
		return query.QueryPackResult{}, err
	}
	return result, nil
}

func validateQueryPackExplainResult(result query.QueryPackResult) error {
	if result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
		return fmt.Errorf("unsupported query-pack result")
	}
	return nil
}

func runExport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "telemetry" {
		return runTelemetryExport(args[1:], stdout, stderr)
	}
	if len(args) > 1 && args[0] == "cross-repo-posture" && args[1] == "explain" {
		return runCrossRepoPostureExplain(args[2:], stdout, stderr)
	}
	if len(args) > 0 && args[0] == "cross-repo-posture" {
		return runCrossRepoPostureExport(args[1:], stdout, stderr)
	}
	fmt.Fprintln(stderr, "export requires cross-repo-posture or telemetry")
	return exitUsage
}

func runTelemetryExport(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "export telemetry"}
	opts.setString("profile", "")
	opts.setString("cross-repo-posture", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "export telemetry accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) != telemetry.ProfilePrometheusTextV1 {
		fmt.Fprintln(stderr, "export telemetry requires --profile prometheus-text-v1")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("cross-repo-posture")) == "" {
		fmt.Fprintln(stderr, "export telemetry requires --cross-repo-posture")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "export telemetry requires --out")
		return exitUsage
	}
	var result posture.ExportResult
	if err := readJSONFile(opts.stringValue("cross-repo-posture"), &result); err != nil {
		fmt.Fprintln(stderr, "posture_unreadable")
		return exitCannotVerify
	}
	rendered, err := telemetry.RenderPrometheus(result)
	if err != nil {
		fmt.Fprintln(stderr, "telemetry_cannot_verify")
		return exitCannotVerify
	}
	if opts.stringValue("out") == "-" {
		fmt.Fprint(stdout, rendered)
		return 0
	}
	if err := writeTextFileAtomic(opts.stringValue("out"), rendered); err != nil {
		fmt.Fprintln(stderr, "out_unwritable")
		return 1
	}
	return 0
}

func runCrossRepoPostureExport(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "export cross-repo-posture"}
	opts.setString("profile", "")
	opts.setString("selection", "")
	opts.setString("out", "")
	opts.setBool("validate-only", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "export cross-repo-posture accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("profile")) != posture.ProfileID {
		fmt.Fprintln(stderr, "export cross-repo-posture requires --profile cross-repo-evidence-posture-v1")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("selection")) == "" {
		fmt.Fprintln(stderr, "export cross-repo-posture requires --selection")
		return exitUsage
	}
	result, err := posture.Build(opts.stringValue("selection"), time.Now())
	if err != nil {
		fmt.Fprintln(stderr, "no_export_artifact")
		return exitCannotVerify
	}
	if opts.boolValue("validate-only") {
		return 0
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "export cross-repo-posture requires --out")
		return exitUsage
	}
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, "out_unwritable")
		return 1
	}
	_ = stdout
	return 0
}

func runCrossRepoPostureExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "export cross-repo-posture-explain"}
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "export cross-repo-posture-explain accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("result")) == "" {
		fmt.Fprintln(stderr, "export cross-repo-posture-explain requires --result")
		return exitUsage
	}
	var result posture.ExportResult
	if err := readJSONFile(opts.stringValue("result"), &result); err != nil {
		fmt.Fprintln(stderr, "result_unreadable")
		return exitCannotVerify
	}
	if result.SchemaVersion != posture.SchemaVersion || result.ExportProfileID != posture.ProfileID {
		fmt.Fprintln(stderr, "unsupported cross-repo posture export")
		return exitCannotVerify
	}
	rendered, err := posture.Explain(result)
	if err != nil {
		fmt.Fprintln(stderr, "output_safety_violation")
		return exitCannotVerify
	}
	fmt.Fprint(stdout, rendered)
	return 0
}

func runValidateFixtures(_ context.Context, args []string, stdout, stderr io.Writer) int {
	fixtureRoot := fixtureRootArg(args)
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if validateFixtureRuns(fixtureRoot, runDirs, stdout, stderr) {
		return 1
	}
	return 0
}

func fixtureRootArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}

func validateFixtureRuns(fixtureRoot string, runDirs []string, stdout, stderr io.Writer) bool {
	failed := false
	for _, runDir := range runDirs {
		if validateFixtureRun(fixtureRoot, runDir, stdout, stderr) {
			failed = true
		}
	}
	return failed
}

func validateFixtureRun(fixtureRoot, runDir string, stdout, stderr io.Writer) bool {
	result, table, audit, verifyErr := verifier.VerifyRun(runDir)
	if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, err)
		return true
	}
	fmt.Fprintf(stdout, "%s => %s\n", runDir, result.Result)
	if verifyErr != nil {
		fmt.Fprintf(stderr, "%s verification error: %v\n", runDir, verifyErr)
	}
	return fixtureExpectationFailed(fixtureRoot, runDir, result, stderr)
}

func fixtureExpectationFailed(fixtureRoot, runDir string, result trace.VerifierResult, stderr io.Writer) bool {
	expectation, err := readFixtureExpectation(fixtureRoot, runDir)
	if err != nil {
		fmt.Fprintf(stderr, "invalid fixture expectation for %s: %v\n", runDir, err)
		return true
	}
	if expectation.ExpectedResult != "" {
		return expectedFixtureResultFailed(runDir, result, expectation, stderr)
	}
	return unexpectedFixtureResultFailed(result)
}

func expectedFixtureResultFailed(runDir string, result trace.VerifierResult, expectation fixtureExpectation, stderr io.Writer) bool {
	if expectation.ExpectedResult == string(result.Result) {
		return false
	}
	fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
	return true
}

func unexpectedFixtureResultFailed(result trace.VerifierResult) bool {
	return result.Result == trace.VerdictFail || result.Result == trace.VerdictCannotVerify
}

type doctorReport struct {
	Command            string        `json:"command"`
	Result             string        `json:"result"`
	Environment        []doctorCheck `json:"environment"`
	ControlPoints      []doctorCheck `json:"control_points"`
	SafeRetentionModes []string      `json:"safe_retention_modes"`
}

type doctorCheck struct {
	ID        string   `json:"id"`
	State     string   `json:"state"`
	Reason    string   `json:"reason"`
	Contract  string   `json:"contract_id,omitempty"`
	Missing   []string `json:"missing,omitempty"`
	Reference string   `json:"reference,omitempty"`
}

type doctorOptions struct {
	ContractPath string
	OutputDir    string
	ReportDir    string
	Env          map[string]string
}

type previewBoundary struct {
	Boundary string `json:"boundary"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
}

type previewOfflineImplication struct {
	Requirement string `json:"requirement"`
	State       string `json:"state"`
	Reason      string `json:"reason"`
}

const (
	defaultRunRoot   = ".sdp-trace-runs"
	defaultReportDir = ".sdp-trace-report"
)

func buildDoctorReport(options doctorOptions) (doctorReport, int) {
	defaultContract := trace.DefaultContract
	result := "offline_dev"
	exitCode := 0
	contract := defaultContract
	contractCheck := doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "default contract is available",
		Contract:  defaultContract.ContractID,
		Reference: "local-default-v1",
	}
	if options.ContractPath != "" {
		loaded, err := trace.LoadContract(options.ContractPath)
		if err != nil {
			result = string(trace.VerdictCannotVerify)
			exitCode = exitCannotVerify
			contractCheck = doctorCheck{
				ID:        "contract",
				State:     string(trace.VerdictCannotVerify),
				Reason:    "contract cannot be loaded",
				Reference: options.ContractPath,
			}
		} else {
			contract = loaded
			contractCheck = doctorCheck{
				ID:        "contract",
				State:     "pass",
				Reason:    "contract can be loaded",
				Contract:  contract.ContractID,
				Reference: options.ContractPath,
			}
		}
	}
	ciCheck := ciWitnessPrerequisiteCheck(options.Env)
	outputDirCheck := writablePathCheck("output_directory", options.OutputDir, "run artifact output directory is writable")
	reportDirCheck := writablePathCheck("report_directory", options.ReportDir, "report artifact directory is writable")
	expectedEvidenceCheck := expectedEvidenceReferenceCheck(contract)
	for _, check := range []doctorCheck{outputDirCheck, reportDirCheck, expectedEvidenceCheck} {
		if check.State == string(trace.VerdictCannotVerify) {
			result = string(trace.VerdictCannotVerify)
			exitCode = exitCannotVerify
		}
	}
	report := doctorReport{
		Command: "doctor",
		Result:  result,
		Environment: []doctorCheck{
			{
				ID:     "local_process",
				State:  "pass",
				Reason: "current process can inspect local environment",
			},
			{
				ID:     "offline_development",
				State:  "offline_dev",
				Reason: "external CI identity is not required for local preview or wrapper readiness",
			},
		},
		ControlPoints: []doctorCheck{
			{
				ID:     "local_wrapper",
				State:  "pass",
				Reason: "wrap and run commands are registered in this binary",
			},
			outputDirCheck,
			reportDirCheck,
			contractCheck,
			expectedEvidenceCheck,
			{
				ID:        "default_contract",
				State:     "pass",
				Reason:    "built-in contract is available for local development",
				Contract:  defaultContract.ContractID,
				Reference: defaultContract.Version,
			},
			ciCheck,
		},
		SafeRetentionModes: safeRetentionModes(),
	}
	return report, exitCode
}

func writablePathCheck(id, path, okReason string) doctorCheck {
	if strings.TrimSpace(path) == "" {
		return doctorCheck{
			ID:     id,
			State:  string(trace.VerdictCannotVerify),
			Reason: "path is empty",
		}
	}
	target := path
	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "path exists but is not a directory",
			Reference: path,
		}
	}
	if os.IsNotExist(err) {
		target = filepath.Dir(path)
		if target == "" {
			target = "."
		}
	}
	probe, err := os.CreateTemp(target, ".sdp-trace-doctor-")
	if err != nil {
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "directory is not writable",
			Reference: path,
		}
	}
	probeName := probe.Name()
	_ = probe.Close()
	_ = os.Remove(probeName)
	return doctorCheck{
		ID:        id,
		State:     "pass",
		Reason:    okReason,
		Reference: path,
	}
}

func expectedEvidenceReferenceCheck(contract trace.Contract) doctorCheck {
	if len(contract.RequiredEvents) == 0 {
		return doctorCheck{
			ID:       "expected_evidence_references",
			State:    string(trace.VerdictCannotVerify),
			Reason:   "contract has no required_events",
			Contract: contract.ContractID,
		}
	}
	missing := make([]string, 0)
	for _, eventType := range contract.RequiredEvents {
		if !knownEventType(eventType) {
			missing = append(missing, "required_events:"+eventType)
		}
	}
	for _, evidence := range contract.RequiredEvidence {
		if strings.TrimSpace(evidence.ID) == "" {
			missing = append(missing, "required_evidence:<missing_id>")
		}
		if strings.TrimSpace(evidence.EventType) == "" {
			missing = append(missing, "required_evidence:"+evidence.ID+":<missing_event_type>")
			continue
		}
		if !knownEventType(evidence.EventType) {
			missing = append(missing, "required_evidence:"+evidence.ID+":"+evidence.EventType)
		}
	}
	if len(missing) > 0 {
		return doctorCheck{
			ID:       "expected_evidence_references",
			State:    string(trace.VerdictCannotVerify),
			Reason:   "contract references unsupported event types",
			Contract: contract.ContractID,
			Missing:  missing,
		}
	}
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    "pass",
		Reason:   "contract required events and evidence references are supported by the current local event model",
		Contract: contract.ContractID,
	}
}

func knownEventType(eventType string) bool {
	switch trace.EventType(eventType) {
	case trace.EventRecorderAttached,
		trace.EventRunStarted,
		trace.EventCommandStarted,
		trace.EventCommandFinished,
		trace.EventRunClosed,
		trace.EventPolicyOverrideRequested:
		return true
	default:
		return false
	}
}

func ciWitnessPrerequisiteCheck(env map[string]string) doctorCheck {
	missing := missingCIWitnessFields(env)
	if len(missing) > 0 {
		return doctorCheck{
			ID:      "ci_witness_prerequisites",
			State:   string(trace.VerdictCannotVerify),
			Reason:  "GitHub Actions identity or OIDC prerequisite is unavailable in this environment",
			Missing: missing,
		}
	}
	return doctorCheck{
		ID:     "ci_witness_prerequisites",
		State:  "pass",
		Reason: "GitHub Actions identity and OIDC prerequisites are present",
	}
}

func missingCIWitnessFields(env map[string]string) []string {
	required := []string{
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_URL",
		"GITHUB_ACTIONS",
		"GITHUB_ACTOR",
		"GITHUB_JOB",
		"GITHUB_REF",
		"GITHUB_REPOSITORY",
		"GITHUB_RUN_ATTEMPT",
		"GITHUB_RUN_ID",
		"GITHUB_SERVER_URL",
		"GITHUB_SHA",
		"GITHUB_WORKFLOW",
	}
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	return missing
}

func safeRetentionModes() []string {
	return []string{
		string(trace.RetentionModeDigestOnly),
		string(trace.RetentionModeSanitizedExcerpt),
		string(trace.RetentionModeEncryptedRawRef),
		string(trace.RetentionModeExternalArtifactRef),
		string(trace.RetentionModeNotAssessed),
	}
}

func previewBoundaries() []previewBoundary {
	return []previewBoundary{
		{
			Boundary: string(trace.ObservationBoundaryProcessWrapper),
			State:    "pass",
			Reason:   "preview covers local process-wrapper capture only",
		},
		{
			Boundary: string(trace.ObservationBoundaryAdapterSocket),
			State:    string(trace.ObservationStateNotIntegrated),
			Reason:   "adapter socket/API capture is not configured in Block 13B",
		},
		{
			Boundary: string(trace.ObservationBoundaryToolWrapper),
			State:    string(trace.ObservationStateUnsupported),
			Reason:   "tool-level wrapping is a future observation boundary",
		},
		{
			Boundary: string(trace.ObservationBoundaryVCSPRObserver),
			State:    string(trace.ObservationStateNotIntegrated),
			Reason:   "VCS/PR observer is not configured in Block 13B",
		},
		{
			Boundary: string(trace.ObservationBoundaryCIObserver),
			State:    string(trace.ObservationStateOfflineDev),
			Reason:   "CI witness cannot be produced by local preview",
		},
		{
			Boundary: string(trace.ObservationBoundaryExternalWitness),
			State:    string(trace.ObservationStateNotIntegrated),
			Reason:   "external witness profile is not implemented in Block 13B",
		},
	}
}

func previewOfflineImplications() []previewOfflineImplication {
	return []previewOfflineImplication{
		{
			Requirement: "ci_witnessed",
			State:       string(trace.ObservationStateOfflineDev),
			Reason:      "rerun in CI with OIDC before using CI witness evidence",
		},
		{
			Requirement: "external_witnessed",
			State:       string(trace.ObservationStateNotIntegrated),
			Reason:      "external witness profile is not implemented in Block 13B",
		},
	}
}

func printUsage(w io.Writer) {
	const usage = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace doctor [--contract <file>]
  sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]
  sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]
  sdp-trace interaction relay --task-id <safe-id> --event-type <type> --out <file> -- <forward-command...>
  sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id <safe-id> --events-jsonl <file> --out <file>
  sdp-trace interaction summarize --trace <file> [--out <file>]
  sdp-trace observe setup --profile <session-profile.json> --out <run-dir> [--command <harness-command-preview>]
  sdp-trace observe collect --profile <session-profile.json> --run <run-dir>
  sdp-trace observe session --profile <session-profile.json> --out <run-dir> -- <harness-command...>
  sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>
  sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>
  sdp-trace harness summarize --validation <validation.json>
  sdp-trace envelope summarize --envelope <file> [--out <file>]
  sdp-trace verify <run-dir>
  sdp-trace explain <run-dir>
  sdp-trace query --query <missing-evidence|capture-depth> <run-dir>
  sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>
  sdp-trace query-pack explain --result <file>
  sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection <file> --out <file>
  sdp-trace export cross-repo-posture explain --result <file>
  sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file|->
  sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>
  sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>
  sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>
  sdp-trace assess --profile ci-artifact-observation --out <file> --artifact-manifest <file>
  sdp-trace assess --profile authority-envelope --out <file> --authority-package <file>
  sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope> [profile inputs]
  sdp-trace assess explain --assessment-result <file>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>
  sdp-trace release-proof --manifest <file> --out <file>
  sdp-trace pr-review packet --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> [--ci-state <state>] [--created-by <actor>]
  sdp-trace pr-review run --packet <dir> --profile <file> --out <dir> [--preview] [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace pr-review synthesize --packet <dir> --runs <dir> --out <file>
  sdp-trace pr-review validate --packet <dir> --profile <file> --runs <dir> --ledger <file> --out <file>
  sdp-trace pr-review summarize --validation <file> --ledger <file> [--out <file>]
  sdp-trace pr-review check --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> --profile <file> [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace validate-fixtures [root-dir]
`
	fmt.Fprint(w, usage)
}

type fixtureExpectation struct {
	ExpectedResult string `json:"expected_result"`
}

func readFixtureExpectation(root, runDir string) (fixtureExpectation, error) {
	expectations, err := readFixtureExpectations(root)
	if err != nil {
		return fixtureExpectation{}, err
	}
	if len(expectations) == 0 {
		return fixtureExpectation{}, nil
	}
	name := filepath.Base(runDir)
	return fixtureExpectation{ExpectedResult: expectations[name]}, nil
}

func readFixtureExpectations(root string) (map[string]string, error) {
	path := filepath.Join(root, "fixture-expectations.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var expectations map[string]string
	if err := json.Unmarshal(data, &expectations); err != nil {
		return nil, err
	}
	return expectations, nil
}

// flagSet is a tiny local parser for limited CLI needs.
type flagSet struct {
	name  string
	data  map[string]string
	bools map[string]bool
	args  []string
}

func (f *flagSet) setString(key, defaultValue string) {
	if f.data == nil {
		f.data = map[string]string{}
	}
	f.data[key] = defaultValue
}

func (f *flagSet) setBool(key string, defaultValue bool) {
	if f.bools == nil {
		f.bools = map[string]bool{}
	}
	f.bools[key] = defaultValue
}

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}

func (f *flagSet) consumeArg(args []string, idx *int, rest *[]string) (bool, error) {
	arg := args[*idx]
	if arg == "--" {
		*rest = append(*rest, args[*idx+1:]...)
		return true, nil
	}
	if !strings.HasPrefix(arg, "--") {
		*rest = append(*rest, arg)
		return false, nil
	}
	flag, flagValue, hasValue := splitFlag(arg)
	return false, f.consumeFlag(flag, flagValue, hasValue, args, idx)
}

func splitFlag(arg string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
	if len(parts) == 1 {
		return parts[0], "", false
	}
	return parts[0], parts[1], true
}

func (f *flagSet) consumeFlag(flag string, flagValue string, hasValue bool, args []string, idx *int) error {
	isString, isBool := f.isKnownFlag(flag)
	if !isString && !isBool {
		return fmt.Errorf("unknown flag --%s", flag)
	}
	if hasValue {
		return f.consumeValue(flag, flagValue, isBool)
	}
	return f.consumeNoEqualsValue(flag, args, idx, isBool)
}

func (f *flagSet) isKnownFlag(flag string) (bool, bool) {
	_, isString := f.data[flag]
	_, isBool := f.bools[flag]
	return isString, isBool
}
func (f *flagSet) consumeValue(flag, flagValue string, isBool bool) error {
	if !isBool {
		f.data[flag] = flagValue
		return nil
	}
	return f.consumeBoolValue(flag, flagValue)
}

func (f *flagSet) consumeNoEqualsValue(flag string, args []string, idx *int, isBool bool) error {
	if !isBool {
		return f.consumeStringFromNext(flag, args, idx)
	}
	nextIdx := *idx + 1
	if !isBoolValueAt(args, nextIdx) {
		f.bools[flag] = true
		return nil
	}
	*idx = nextIdx
	return f.consumeBoolValue(flag, args[*idx])
}

func (f *flagSet) consumeStringFromNext(flag string, args []string, idx *int) error {
	nextIdx := *idx + 1
	if nextIdx >= len(args) {
		return fmt.Errorf("flag --%s requires value", flag)
	}
	value := args[nextIdx]
	if strings.HasPrefix(value, "--") {
		return fmt.Errorf("flag --%s requires value", flag)
	}
	*idx = nextIdx
	f.data[flag] = value
	return nil
}

func isBoolValueAt(args []string, idx int) bool {
	return idx < len(args) && isBoolLiteral(args[idx])
}

func (f *flagSet) consumeBoolValue(flag, flagValue string) error {
	switch strings.ToLower(flagValue) {
	case "false", "0":
		f.bools[flag] = false
	case "true", "1", "":
		f.bools[flag] = true
	default:
		return fmt.Errorf("invalid boolean value for --%s: %s", flag, flagValue)
	}
	return nil
}

func (f *flagSet) stringValue(key string) string {
	if f.data == nil {
		return ""
	}
	return f.data[key]
}

func (f *flagSet) boolValue(key string) bool {
	if f.bools == nil {
		return false
	}
	return f.bools[key]
}

func (f *flagSet) rest() []string {
	return f.args
}

func isHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help")
}

func isBoolLiteral(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "1" || lower == "0"
}
