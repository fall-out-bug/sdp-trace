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
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/posture"
	"github.com/fall_out_bug/sdp-trace/internal/query"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

const (
	exitUsage        = 2
	exitCannotVerify = 3
)

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
	ctx := context.Background()

	switch cmd {
	case "wrap":
		return runWrap(ctx, cmdArgs, stdout, stderr)
	case "run":
		return runWrappedCommand(ctx, cmdArgs, stdout, stderr)
	case "dry-run":
		return runDryRun(ctx, cmdArgs, stdout, stderr)
	case "preview":
		return runPreview(ctx, cmdArgs, stdout, stderr)
	case "doctor":
		return runDoctor(ctx, cmdArgs, stdout, stderr)
	case "verify":
		return runVerify(ctx, cmdArgs, stdout, stderr)
	case "explain":
		return runExplain(ctx, cmdArgs, stdout, stderr)
	case "query":
		return runQuery(ctx, cmdArgs, stdout, stderr)
	case "query-pack":
		return runQueryPack(ctx, cmdArgs, stdout, stderr)
	case "export":
		return runExport(ctx, cmdArgs, stdout, stderr)
	case "report":
		return runReport(ctx, cmdArgs, stdout, stderr)
	case "gate":
		return runGate(ctx, cmdArgs, stdout, stderr)
	case "assess":
		return runAssess(ctx, cmdArgs, stdout, stderr)
	case "override":
		return runOverride(ctx, cmdArgs, stdout, stderr)
	case "checkpoint":
		return runCheckpoint(ctx, cmdArgs, stdout, stderr)
	case "witness":
		return runWitness(ctx, cmdArgs, stdout, stderr)
	case "validate-fixtures":
		return runValidateFixtures(ctx, cmdArgs, stdout, stderr)
	case "release-proof":
		return runReleaseProof(ctx, cmdArgs, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return 1
	}
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

func runAssess(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "preview":
			return runAssessPreview(args[1:], stdout, stderr)
		case "explain":
			return runAssessExplain(args[1:], stdout, stderr)
		}
	}
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
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "assess accepts only flags")
		return exitUsage
	}
	switch opts.stringValue("profile") {
	case "adapter-capture":
		return runAdapterCaptureAssess(opts, stdout, stderr)
	case "managed-harness":
		return runManagedAssess(opts, stdout, stderr)
	case "forensic-retention":
		return runForensicAssess(opts, stdout, stderr)
	case "ci-artifact-observation":
		return runCIArtifactAssess(opts, stdout, stderr)
	default:
		fmt.Fprintln(stderr, "assess requires --profile adapter-capture, managed-harness, forensic-retention, or ci-artifact-observation")
		return exitUsage
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
	opts := &flagSet{name: "assess preview"}
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("run", "")
	opts.setString("adapter-registry", "")
	opts.setString("managed-policy", "")
	opts.setString("managed-witness", "")
	opts.setString("redaction-policy", "")
	opts.setString("artifact-manifest", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "assess preview accepts only flags")
		return exitUsage
	}
	switch opts.stringValue("profile") {
	case "adapter-capture":
		return runAdapterCaptureAssessPreview(opts, stdout)
	case "managed-harness":
		return runManagedAssessPreview(opts, stdout)
	case "forensic-retention":
		return runForensicAssessPreview(opts, stdout)
	case "ci-artifact-observation":
		return runCIArtifactAssessPreview(opts, stdout)
	default:
		fmt.Fprintln(stderr, "assess preview requires --profile adapter-capture, managed-harness, forensic-retention, or ci-artifact-observation")
		return exitUsage
	}
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
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "assess explain accepts only flags")
		return exitUsage
	}
	path := opts.stringValue("assessment-result")
	if path == "" {
		fmt.Fprintln(stderr, "assess explain requires --assessment-result <file>")
		return exitUsage
	}
	var envelope struct {
		SchemaVersion   string `json:"schema_version"`
		SelectedProfile string `json:"selected_profile"`
	}
	if err := readJSONFile(path, &envelope); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	switch envelope.SchemaVersion {
	case adaptercapture.SchemaVersion:
		var result adaptercapture.AssessmentResult
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explainAdapterCaptureAssessment(result, stdout)
	case managed.SchemaVersion:
		var result managed.AssessmentResult
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explainManagedAssessment(result, stdout)
	case forensic.SchemaVersion:
		var result forensic.AssessmentResult
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explainForensicAssessment(result, stdout)
	case ciartifact.SchemaVersion:
		var result ciartifact.ObservationResult
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explainCIArtifactObservation(result, stdout)
	default:
		fmt.Fprintf(stderr, "unsupported assessment-result schema_version: %s\n", envelope.SchemaVersion)
		return exitCannotVerify
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
	required := map[string]string{
		"--checkpoint":        opts.stringValue("checkpoint"),
		"--checkpoint-policy": opts.stringValue("checkpoint-policy"),
		"--witness":           opts.stringValue("witness"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "protected gate requires %s\n", flag)
			return exitUsage
		}
	}
	var signed checkpoint.SignedCheckpoint
	if err := readJSONFile(opts.stringValue("checkpoint"), &signed); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	var policy checkpoint.TrustedCheckpointPolicy
	if err := readJSONFile(opts.stringValue("checkpoint-policy"), &policy); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	var witnessSummary demo.WitnessSummary
	if err := readJSONFile(opts.stringValue("witness"), &witnessSummary); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	rows, err := demo.VerifiedRows(target, contract)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	runDir, err := protectedRunDir(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	checkpointResult := checkpoint.Verify(runDir, signed, &policy)
	expected, err := demoWitnessExpectation(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	checkpointResult = protectedCheckpointVerification(checkpointResult, signed, policy, witnessSummary, expected)
	result := demo.EvaluateProtectedGate(rows, contract, demo.ProtectedGateInput{
		Checkpoint:         checkpointResult,
		PolicyProvided:     true,
		Witness:            &witnessSummary,
		WitnessExpectation: expected,
		Now:                time.Now().UTC(),
	})
	if err := writeJSONFile(outPath, result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return gateExitCode(result)
}

func runGateExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "gate explain"}
	opts.setString("gate-result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "gate explain accepts only flags")
		return exitUsage
	}
	path := opts.stringValue("gate-result")
	if path == "" {
		fmt.Fprintln(stderr, "gate explain requires --gate-result <file>")
		return exitUsage
	}
	var result demo.GateResult
	if err := readJSONFile(path, &result); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if result.SchemaVersion != demo.GateSchemaVersion && result.SchemaVersion != demo.GateSchemaVersionBlock16 {
		fmt.Fprintf(stderr, "unsupported gate-result schema_version: %s\n", result.SchemaVersion)
		return exitCannotVerify
	}
	if result.SchemaVersion == demo.GateSchemaVersion {
		fmt.Fprintln(stdout, "Protected profile fields: absent")
	}
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
	if result.CheckpointVerification != nil {
		fmt.Fprintf(stdout, "Checkpoint result: %s\n", result.CheckpointVerification.Result)
		fmt.Fprintf(stdout, "Checkpoint trust scope: %s\n", result.CheckpointVerification.TrustScope)
	}
	for _, condition := range result.ProtectedConditions {
		fmt.Fprintf(stdout, "Protected condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
	for _, requiredRun := range result.RequiredRuns {
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
	for _, binding := range result.WitnessBindings {
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
	for _, missing := range result.MissingAuditEvidence {
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
	for _, override := range result.OverrideRequests {
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
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
	if witnessSummary.Kind != "github-actions" || witnessSummary.Status != demo.GatePass || witnessSummary.TrustScope != "ci_witnessed" {
		return false
	}
	if expected.Repository != "" && witnessSummary.Source.Repository != expected.Repository {
		return false
	}
	if expected.Ref != "" && witnessSummary.Source.Ref != expected.Ref {
		return false
	}
	if expected.CommitSHA != "" && witnessSummary.Source.CommitSHA != expected.CommitSHA {
		return false
	}
	if expected.RunID != "" && witnessSummary.CIIdentity.RunID != expected.RunID {
		return false
	}
	expectedArtifacts := map[string]string{}
	for _, artifact := range expected.RunArtifacts {
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	if len(expectedArtifacts) > 0 && len(witnessSummary.RunArtifacts) == 0 {
		return false
	}
	for _, artifact := range witnessSummary.RunArtifacts {
		if expectedArtifacts[artifact.Path] != artifact.SHA256 {
			return false
		}
		delete(expectedArtifacts, artifact.Path)
	}
	return len(expectedArtifacts) == 0
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
	if result.SelectedProfile == demo.GateProfileProtected {
		switch result.ProtectedGate {
		case demo.GatePass:
			return 0
		case demo.GateFail:
			return 1
		case demo.GateCannotVerify, demo.GateNotAssessed:
			return exitCannotVerify
		}
	}
	for _, requiredRun := range result.RequiredRuns {
		if requiredRun.State == demo.GateFail || requiredRun.State == demo.GateMissingTelemetry {
			return 1
		}
	}
	for _, state := range []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate} {
		if state == demo.GateFail || state == demo.GateMissingTelemetry {
			return 1
		}
	}
	for _, requiredRun := range result.RequiredRuns {
		if requiredRun.State == demo.GateCannotVerify {
			return exitCannotVerify
		}
	}
	for _, state := range []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate} {
		if state == demo.GateCannotVerify {
			return exitCannotVerify
		}
	}
	return 0
}

func runWitness(_ context.Context, args []string, stdout, stderr io.Writer) int {
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
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "witness requires <runs-root-or-run-dir>")
		return exitUsage
	}
	kind := opts.stringValue("kind")
	if !allowedWitnessKind(kind) {
		fmt.Fprintln(stderr, "witness requires --kind github-actions, gitlab-ci, buildkite, or customer-pki")
		return exitUsage
	}
	outPath := opts.stringValue("out")
	if outPath == "" {
		fmt.Fprintln(stderr, "witness requires --out <file>")
		return exitUsage
	}
	var record witness.Record
	var err error
	switch kind {
	case witness.KindGitHubActions:
		record, err = witness.WriteGitHubActions(outPath, targets[0], opts.stringValue("report-dir"), witness.EnvironmentFromOS())
	case witness.KindGitLabCI, witness.KindBuildkite:
		record, err = witness.WriteProfile(kind, outPath, targets[0], opts.stringValue("report-dir"), witness.ProfileOptions{
			EnvelopePath: opts.stringValue("witness-envelope"),
		})
	case witness.KindCustomerPKI:
		if missing := missingCustomerPKIFlags(opts); len(missing) > 0 {
			fmt.Fprintf(stderr, "customer-pki witness requires %s\n", strings.Join(missing, ", "))
			return exitUsage
		}
		record, err = witness.WriteProfile(kind, outPath, targets[0], opts.stringValue("report-dir"), witness.ProfileOptions{
			CustomerPKIAuthorityPolicy: opts.stringValue("customer-pki-authority-policy"),
			CustomerPKIPublicCert:      opts.stringValue("customer-pki-public-cert"),
			CustomerPKIPublicKey:       opts.stringValue("customer-pki-public-key"),
			CustomerPKIPayloadDigest:   opts.stringValue("customer-pki-payload-digest"),
			CustomerPKIFreshness:       opts.stringValue("customer-pki-freshness-evidence"),
		})
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(record, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	if record.Status == witness.StatusCannotVerify {
		return exitCannotVerify
	}
	if record.Status == witness.StatusFail {
		return 1
	}
	if record.Status == witness.StatusNotAssessed {
		return exitCannotVerify
	}
	return 0
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
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return exitUsage
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
	opts := &flagSet{name: "query-pack"}
	opts.setString("pack", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "query-pack accepts only flags")
		return exitUsage
	}
	switch strings.TrimSpace(opts.stringValue("pack")) {
	case "":
		fmt.Fprintln(stderr, "error: ambiguous pack selection; --pack is required")
		return exitUsage
	case query.QueryPackForensicsBasic:
	default:
		fmt.Fprintf(stderr, "error: unknown pack %q\n", opts.stringValue("pack"))
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		fmt.Fprintln(stderr, "query-pack requires --run")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		fmt.Fprintln(stderr, "query-pack requires --out")
		return exitUsage
	}
	result, err := query.ForensicsBasicPack(opts.stringValue("run"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func runQueryPackExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "query-pack explain"}
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "query-pack explain accepts only flags")
		return exitUsage
	}
	if strings.TrimSpace(opts.stringValue("result")) == "" {
		fmt.Fprintln(stderr, "query-pack explain requires --result")
		return exitUsage
	}
	var result query.QueryPackResult
	if err := readJSONFile(opts.stringValue("result"), &result); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
		fmt.Fprintln(stderr, "unsupported query-pack result")
		return exitCannotVerify
	}
	fmt.Fprint(stdout, query.ExplainForensicsPack(result))
	return 0
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
	fixtureRoot := "."
	if len(args) > 0 {
		fixtureRoot = args[0]
	}
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	failed := false
	for _, runDir := range runDirs {
		result, table, audit, verifyErr := verifier.VerifyRun(runDir)
		if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
			fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, err)
			failed = true
			continue
		}
		fmt.Fprintf(stdout, "%s => %s\n", runDir, result.Result)
		if verifyErr != nil {
			fmt.Fprintf(stderr, "%s verification error: %v\n", runDir, verifyErr)
		}
		expectation, err := readFixtureExpectation(fixtureRoot, runDir)
		if err != nil {
			fmt.Fprintf(stderr, "invalid fixture expectation for %s: %v\n", runDir, err)
			failed = true
			continue
		}
		if expectation.ExpectedResult != "" && expectation.ExpectedResult != string(result.Result) {
			fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
			failed = true
			continue
		}
		if expectation.ExpectedResult == "" && result.Result == trace.VerdictFail {
			failed = true
		}
		if expectation.ExpectedResult == "" && result.Result == trace.VerdictCannotVerify {
			failed = true
		}
	}
	if failed {
		return 1
	}
	return 0
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
  sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation> [profile inputs]
  sdp-trace assess explain --assessment-result <file>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>
  sdp-trace release-proof --manifest <file> --out <file>
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
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			key := parts[0]
			_, knownString := f.data[key]
			_, knownBool := f.bools[key]
			if !knownString && !knownBool {
				return fmt.Errorf("unknown flag --%s", key)
			}
			switch len(parts) {
			case 1:
				switch {
				case knownBool:
					if i+1 < len(args) && isBoolLiteral(args[i+1]) {
						f.bools[key] = parseBoolLiteral(args[i+1])
						i++
					} else {
						f.bools[key] = true
					}
				default:
					if i+1 >= len(args) {
						return fmt.Errorf("flag --%s requires value", key)
					}
					val := args[i+1]
					if strings.HasPrefix(val, "--") {
						return fmt.Errorf("flag --%s requires value", key)
					}
					i++
					f.data[key] = val
				}
			case 2:
				if _, ok := f.bools[key]; ok {
					lower := strings.ToLower(parts[1])
					if lower == "false" || lower == "0" {
						f.bools[key] = false
					} else if lower == "true" || lower == "1" || lower == "" {
						f.bools[key] = true
					} else {
						return fmt.Errorf("invalid boolean value for --%s: %s", key, parts[1])
					}
					continue
				}
				f.data[key] = parts[1]
			default:
			}
			continue
		}
		rest = append(rest, arg)
	}
	f.args = rest
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

func parseBoolLiteral(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "1"
}
