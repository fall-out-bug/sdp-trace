package prreview

import (
	"context"
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
	SchemaVersionPacket     = "block30-pr-review-packet-v1"
	SchemaVersionProfile    = "block30-pr-review-profile-v1"
	SchemaVersionRunSet     = "block30-pr-review-runs-v1"
	SchemaVersionLedger     = "block30-pr-review-ledger-v1"
	SchemaVersionValidation = "block30-pr-review-validation-v1"

	StatePass         = "pass"
	StateFail         = "fail"
	StatePending      = "pending"
	StateNotAssessed  = "not_assessed"
	StateCannotVerify = "cannot_verify"

	RefKindDiff            = "diff"
	RefKindMetadata        = "metadata"
	RefKindSpec            = "spec"
	RefKindPlan            = "plan"
	RefKindTask            = "task"
	RefKindDoc             = "doc"
	RefKindSchema          = "schema"
	RefKindSourceExcerpt   = "source_excerpt"
	RefKindVerification    = "verification"
	RefKindPrompt          = "prompt"
	RefKindRawOutput       = "raw_output"
	RefKindSanitizedOutput = "sanitized_output"
	RefKindExternal        = "external"

	ContentUnifiedDiff = "unified_diff"
	ContentMarkdown    = "markdown"
	ContentJSON        = "json"
	ContentText        = "text"

	RedactionNone        = "none"
	RedactionRedacted    = "redacted"
	RedactionDigestOnly  = "digest_only"
	RedactionEncrypted   = "encrypted_ref"
	RedactionExternalRef = "external_ref"
	RedactionWithheld    = "withheld"
	RedactionNotAssessed = "not_assessed"

	PlaneCodeCorrectness = "code_correctness"
	PlaneTraceEvidence   = "trace_evidence_provenance"
	PlaneRequirements    = "requirements_vs_implementation"
	PlaneSecurity        = "security_forgery_overclaim"
	PlaneDXReplayability = "dx_replayability"
	PlanePrivacySafety   = "privacy_output_safety"

	RunnerPI             = "pi"
	RunnerOpenCode       = "opencode"
	RunnerManualExternal = "manual_external"

	StatusFindingsReported = "findings_reported"
	StatusNoFindings       = "no_findings"
	StatusNotAssessed      = "not_assessed"
	StatusFailed           = "failed"
	StatusTimedOut         = "timed_out"
	StatusEmptyOutput      = "empty_output"
	StatusOffTask          = "off_task"
	StatusParseFailed      = "parse_failed"
	StatusCannotVerify     = "cannot_verify"

	SeverityCritical      = "critical"
	SeverityMajor         = "major"
	SeverityMinor         = "minor"
	SeverityInformational = "informational"

	DispositionAcceptedFixed           = "accepted_fixed"
	DispositionAcceptedReviewBlocking  = "accepted_review_blocking"
	DispositionAcceptedNarrower        = "accepted_narrower"
	DispositionRejectedFalsePositive   = "rejected_false_positive"
	DispositionDeferredNotAssessed     = "deferred_not_assessed"
	DispositionUnresolvedReviewBlocker = "unresolved_review_blocker"

	CoverageSatisfied    = "coverage_satisfied"
	CoveragePartial      = "coverage_partial"
	CoverageUnresolved   = "coverage_unresolved"
	CoverageNotAssessed  = "not_assessed"
	CoverageCannotVerify = "cannot_verify"

	AuthorityReviewRecordOnly = "review_record_only"
	DecisionNotAuthorized     = "not_authorized_by_sdp_trace"
)

var (
	repoIDPattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,62}[a-z0-9]$`)
	changeRefPattern = regexp.MustCompile(`^(pr|mr|change)-[A-Za-z0-9._-]{1,64}$`)
	sha40Pattern     = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

type PacketOptions struct {
	OutDir            string
	RepoID            string
	ChangeRef         string
	BaseCommit        string
	HeadCommit        string
	DiffPath          string
	MetadataPath      string
	ContextPaths      []string
	VerificationPaths []string
	CIState           string
	CreatedBy         string
	Now               time.Time
}

// Packet is the frozen review input boundary: reviewers assess the copied refs,
// not the mutable files those refs came from.
type Packet struct {
	SchemaVersion     string             `json:"schema_version"`
	PacketID          string             `json:"packet_id"`
	PacketDigest      string             `json:"packet_digest"`
	RepoID            string             `json:"repo_id"`
	ChangeRef         string             `json:"change_ref"`
	BaseCommit        string             `json:"base_commit"`
	HeadCommit        string             `json:"head_commit"`
	DiffRef           SafeRef            `json:"diff_ref"`
	MetadataRef       *SafeRef           `json:"metadata_ref,omitempty"`
	ContextRefs       []SafeRef          `json:"context_refs"`
	VerificationRefs  []SafeRef          `json:"verification_refs"`
	CIState           string             `json:"ci_state"`
	CreatedAt         string             `json:"created_at"`
	CreatedBy         string             `json:"created_by"`
	RedactionState    string             `json:"redaction_state"`
	UnavailableFields []UnavailableField `json:"unavailable_fields,omitempty"`
}

// SafeRef records enough provenance to replay a packet input without treating
// the referenced content as an approval or proof claim by itself.
type SafeRef struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Ref            string `json:"ref"`
	DigestSHA256   string `json:"digest_sha256"`
	ContentType    string `json:"content_type"`
	RedactionState string `json:"redaction_state"`
}

// UnavailableField makes missing optional packet inputs explicit so absence
// cannot be mistaken for passing evidence.
type UnavailableField struct {
	Field  string `json:"field"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}

// ReviewProfile declares which planes must be reviewed and which runners may
// produce evidence for those planes.
type ReviewProfile struct {
	SchemaVersion  string       `json:"schema_version"`
	ProfileID      string       `json:"profile_id"`
	RequiredPlanes []string     `json:"required_planes"`
	Roles          []ReviewRole `json:"roles"`
}

// ReviewRole binds one reviewer runner to one trust plane and its declared
// execution contract.
type ReviewRole struct {
	RoleID               string   `json:"role_id"`
	Plane                string   `json:"plane"`
	Runner               string   `json:"runner"`
	RequestedModel       string   `json:"requested_model"`
	Command              []string `json:"command,omitempty"`
	TimeoutSeconds       int      `json:"timeout_seconds,omitempty"`
	PromptTemplateRef    string   `json:"prompt_template_ref,omitempty"`
	RequiredOutputSchema string   `json:"required_output_schema,omitempty"`
	RawOutputRetention   string   `json:"raw_output_retention,omitempty"`
	ReadOnlyEnforced     bool     `json:"read_only_enforced,omitempty"`
	WorkingTreeMode      string   `json:"working_tree_mode,omitempty"`
}

// RunOptions contains local execution policy; it is intentionally separate
// from the portable packet so packet evidence stays harness-neutral.
type RunOptions struct {
	OutDir         string
	AllowedRunners map[string]bool
	Preview        bool
	Now            time.Time
	WorkDir        string
}

// RunPreview is a dry-run artifact: it exposes intended commands and prompts
// without creating assessed reviewer evidence.
type RunPreview struct {
	SchemaVersion string        `json:"schema_version"`
	PacketDigest  string        `json:"packet_digest"`
	Roles         []PreviewRole `json:"roles"`
}

// PreviewRole is the per-role portion of a preview, limited to digests and
// runner metadata that are safe to show before execution.
type PreviewRole struct {
	RoleID         string `json:"role_id"`
	Plane          string `json:"plane"`
	Runner         string `json:"runner"`
	RequestedModel string `json:"requested_model"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	CommandDigest  string `json:"command_digest"`
	PromptRef      string `json:"prompt_template_ref,omitempty"`
	PromptDigest   string `json:"prompt_digest,omitempty"`
}

// RunSet records reviewer outputs for a single packet digest.
type RunSet struct {
	SchemaVersion string           `json:"schema_version"`
	PacketDigest  string           `json:"packet_digest"`
	Results       []ReviewerResult `json:"results"`
}

// ReviewerResult is reviewer-authored evidence plus harness-owned execution
// metadata used to decide whether the result is usable.
type ReviewerResult struct {
	ReviewRunID      string    `json:"review_run_id"`
	PacketDigest     string    `json:"packet_digest"`
	Plane            string    `json:"plane"`
	RoleID           string    `json:"role_id"`
	Runner           string    `json:"runner"`
	RequestedModel   string    `json:"requested_model"`
	ObservedModel    string    `json:"observed_model"`
	ModelFamily      string    `json:"model_family"`
	ModelVersion     string    `json:"model_version"`
	FallbackForModel string    `json:"fallback_for_model,omitempty"`
	FallbackReason   string    `json:"fallback_reason,omitempty"`
	Status           string    `json:"status"`
	Findings         []Finding `json:"findings"`
	CommandDigest    string    `json:"command_digest,omitempty"`
	RawOutputRef     *SafeRef  `json:"raw_output_ref,omitempty"`
	PromptRef        *SafeRef  `json:"prompt_ref,omitempty"`
	ContextRefs      []string  `json:"context_refs,omitempty"`
	StartedAt        string    `json:"started_at,omitempty"`
	EndedAt          string    `json:"ended_at,omitempty"`
	Reason           string    `json:"reason,omitempty"`
}

// Finding is a reviewer claim that must be backed by a citation before
// validation can treat it as replayable.
type Finding struct {
	ID           string   `json:"id"`
	Severity     string   `json:"severity"`
	Citation     Citation `json:"citation"`
	Summary      string   `json:"summary"`
	Rationale    string   `json:"rationale,omitempty"`
	SuggestedFix string   `json:"suggested_fix,omitempty"`
	Question     string   `json:"question,omitempty"`
	EvidenceRefs []string `json:"evidence_refs,omitempty"`
}

// Citation points a finding back to packet evidence; validation rejects claims
// that cannot resolve to a packet ref or digest.
type Citation struct {
	ContextRefID string `json:"context_ref_id,omitempty"`
	DiffHunkID   string `json:"diff_hunk_id,omitempty"`
	SourceDigest string `json:"source_digest,omitempty"`
	LineStart    int    `json:"line_start,omitempty"`
	LineEnd      int    `json:"line_end,omitempty"`
}

// Ledger is the durable review disposition record synthesized from run output
// and any prior human decisions.
type Ledger struct {
	SchemaVersion string          `json:"schema_version"`
	PacketDigest  string          `json:"packet_digest"`
	Findings      []LedgerFinding `json:"findings"`
}

// LedgerFinding preserves reviewer findings with human disposition state kept
// outside reviewer control.
type LedgerFinding struct {
	ID                  string   `json:"id"`
	ReviewRunID         string   `json:"review_run_id"`
	Plane               string   `json:"plane"`
	RoleID              string   `json:"role_id"`
	Severity            string   `json:"severity"`
	Summary             string   `json:"summary"`
	Citation            Citation `json:"citation"`
	Disposition         string   `json:"disposition"`
	EvidenceRefs        []string `json:"evidence_refs,omitempty"`
	DispositionEvidence string   `json:"disposition_evidence,omitempty"`
}

// Validation is the portable coverage verdict for review evidence only. It does
// not authorize merge, release, or risk acceptance.
type Validation struct {
	SchemaVersion       string          `json:"schema_version"`
	PacketDigest        string          `json:"packet_digest"`
	ReviewCoverageState string          `json:"review_coverage_state"`
	CIState             string          `json:"ci_state"`
	AuthorityScope      string          `json:"authority_scope"`
	MergeDecision       string          `json:"merge_decision"`
	ReleaseDecision     string          `json:"release_decision"`
	RiskAcceptance      string          `json:"risk_acceptance"`
	PlaneResults        []PlaneResult   `json:"plane_results"`
	Findings            []LedgerFinding `json:"findings"`
	Reasons             []string        `json:"reasons"`
	NextActions         []string        `json:"next_actions"`
}

// PlaneResult states whether one required review plane has usable evidence and
// what action remains when it does not.
type PlaneResult struct {
	Plane      string `json:"plane"`
	Status     string `json:"status"`
	Usable     bool   `json:"usable"`
	RunID      string `json:"review_run_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
	NextAction string `json:"next_action,omitempty"`
}

func BuildPacket(opts PacketOptions) (Packet, error) {
	// BuildPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validatePacketOptions(opts); err != nil {
		return Packet{}, err
	}
	if err := ensureNewDir(opts.OutDir); err != nil {
		return Packet{}, err
	}
	return buildPacketInPreparedDir(opts)
}

func buildPacketInPreparedDir(opts PacketOptions) (Packet, error) {
	// buildPacketInPreparedDir keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	now, createdBy, ciState := packetDefaults(opts)
	refs, err := buildPacketRefs(opts)
	if err != nil {
		return Packet{}, err
	}
	packet := newPacket(opts, refs, now, createdBy, ciState)
	if err := finalizePacket(opts.OutDir, &packet); err != nil {
		return Packet{}, err
	}
	return packet, nil
}

func finalizePacket(outDir string, packet *Packet) error {
	// finalizePacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	digest, err := packetDigest(*packet)
	if err != nil {
		return err
	}
	packet.PacketDigest = "sha256:" + digest
	return WriteJSON(filepath.Join(outDir, "packet.json"), *packet)
}

func newPacket(opts PacketOptions, refs packetRefs, now time.Time, createdBy, ciState string) Packet {
	// newPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet := newPacketIdentity(opts)
	attachPacketRefs(&packet, refs)
	attachPacketProvenance(&packet, now, createdBy, ciState)
	packet.UnavailableFields = unavailablePacketFields(opts)
	return packet
}

func newPacketIdentity(opts PacketOptions) Packet {
	// newPacketIdentity keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return Packet{
		SchemaVersion:  SchemaVersionPacket,
		PacketID:       fmt.Sprintf("%s-%s-%s", opts.RepoID, opts.ChangeRef, opts.HeadCommit[:12]),
		RepoID:         opts.RepoID,
		ChangeRef:      opts.ChangeRef,
		BaseCommit:     opts.BaseCommit,
		HeadCommit:     opts.HeadCommit,
		RedactionState: RedactionNone,
	}
}

func attachPacketRefs(packet *Packet, refs packetRefs) {
	// attachPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.DiffRef = refs.diff
	packet.MetadataRef = refs.metadata
	packet.ContextRefs = refs.context
	packet.VerificationRefs = refs.verification
}

func attachPacketProvenance(packet *Packet, now time.Time, createdBy, ciState string) {
	// attachPacketProvenance keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.CIState = ciState
	packet.CreatedAt = now.Format(time.RFC3339)
	packet.CreatedBy = createdBy
}

type packetRefs struct {
	diff         SafeRef
	metadata     *SafeRef
	context      []SafeRef
	verification []SafeRef
}

func packetDefaults(opts PacketOptions) (time.Time, string, string) {
	// packetDefaults keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	now := opts.Now
	if now.IsZero() {

		now = time.Now().UTC()
	}
	createdBy := strings.TrimSpace(opts.CreatedBy)
	if createdBy == "" {
		createdBy = "sdp-trace"
	}
	ciState := opts.CIState
	if ciState == "" {

		ciState = StateNotAssessed
	}
	return now, createdBy, ciState
}

func buildPacketRefs(opts PacketOptions) (packetRefs, error) {
	// buildPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	inputDir := filepath.Join(opts.OutDir, "inputs")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return packetRefs{}, err
	}
	return collectPacketRefs(inputDir, opts)
}

func collectPacketRefs(inputDir string, opts PacketOptions) (packetRefs, error) {
	// collectPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	diffRef, err := copyDiffRef(inputDir, opts.DiffPath)
	if err != nil {
		return packetRefs{}, err
	}
	return collectOptionalPacketRefs(inputDir, opts, diffRef)
}

func collectOptionalPacketRefs(inputDir string, opts PacketOptions, diffRef SafeRef) (packetRefs, error) {
	// collectOptionalPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	metadataRef, err := optionalMetadataRef(inputDir, opts.MetadataPath)
	if err != nil {
		return packetRefs{}, err
	}
	contextRefs, err := packetContextRefs(inputDir, opts.ContextPaths)
	if err != nil {
		return packetRefs{}, err
	}
	return packetRefsWithVerification(inputDir, opts, diffRef, metadataRef, contextRefs)
}

func copyDiffRef(inputDir, diffPath string) (SafeRef, error) {
	return copyInput(inputDir, "diff.patch", diffPath, RefKindDiff, ContentUnifiedDiff)
}

func packetRefsWithVerification(inputDir string, opts PacketOptions, diffRef SafeRef, metadataRef *SafeRef, contextRefs []SafeRef) (packetRefs, error) {
	// packetRefsWithVerification keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	verificationRefs, err := packetVerificationRefs(inputDir, opts.VerificationPaths)
	if err != nil {
		return packetRefs{}, err
	}
	return packetRefs{diff: diffRef, metadata: metadataRef, context: contextRefs, verification: verificationRefs}, nil
}

func optionalMetadataRef(inputDir, metadataPath string) (*SafeRef, error) {
	// optionalMetadataRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(metadataPath) == "" {
		return nil, nil
	}
	ref, err := copyInput(inputDir, "metadata.json", metadataPath, RefKindMetadata, contentType(metadataPath))
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func packetContextRefs(inputDir string, paths []string) ([]SafeRef, error) {
	// packetContextRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	refs, err := copyInputs(inputDir, "context", paths)
	if err != nil {
		return nil, err
	}
	for i := range refs {

		refs[i].Kind = contextKind(paths[i])
	}
	return refs, nil
}

func packetVerificationRefs(inputDir string, paths []string) ([]SafeRef, error) {
	// packetVerificationRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	refs, err := copyInputs(inputDir, "verification", paths)
	if err != nil {
		return nil, err
	}
	for i := range refs {

		refs[i].Kind = RefKindVerification
	}
	return refs, nil
}

func unavailablePacketFields(opts PacketOptions) []UnavailableField {
	// unavailablePacketFields keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	fields := []UnavailableField{}
	if strings.TrimSpace(opts.MetadataPath) == "" {

		fields = append(fields, UnavailableField{Field: "metadata_ref", State: StateNotAssessed, Reason: "metadata_input_not_provided"})
	}
	if len(opts.ContextPaths) == 0 {

		fields = append(fields, UnavailableField{Field: "context_refs", State: StateNotAssessed, Reason: "context_inputs_not_provided"})
	}
	if len(opts.VerificationPaths) == 0 {

		fields = append(fields, UnavailableField{Field: "verification_refs", State: StateNotAssessed, Reason: "verification_inputs_not_provided"})
	}
	return fields
}
func RunReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	// RunReview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validateProfile(profile); err != nil {
		return RunSet{}, nil, err
	}
	opts = normalizeRunOptions(opts)
	if opts.Preview {
		return RunSet{}, preview(packet, profile), nil
	}
	return runReview(packet, profile, opts)
}

func normalizeRunOptions(opts RunOptions) RunOptions {
	// normalizeRunOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if opts.Now.IsZero() {

		opts.Now = time.Now().UTC()
	}
	if opts.WorkDir == "" {

		opts.WorkDir = "."
	}
	return opts
}

func runReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	// runReview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	rawDir, err := prepareRunDirectories(opts.OutDir)
	if err != nil {
		return RunSet{}, nil, err
	}
	results, err := runReviewRoles(packet, profile.Roles, opts, rawDir)
	if err != nil {
		return RunSet{}, nil, err
	}
	runSet := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Results: results}

	if err := WriteJSON(filepath.Join(opts.OutDir, "results.json"), runSet); err != nil {
		return RunSet{}, nil, err
	}
	return runSet, nil, nil
}

func prepareRunDirectories(outDir string) (string, error) {
	// prepareRunDirectories keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if err := ensureNewDir(outDir); err != nil {
		return "", err
	}
	rawDir := filepath.Join(outDir, "raw")

	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return "", err
	}
	return rawDir, nil
}

func runReviewRoles(packet Packet, roles []ReviewRole, opts RunOptions, rawDir string) ([]ReviewerResult, error) {
	// runReviewRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	results := make([]ReviewerResult, 0, len(roles))
	for _, role := range roles {
		result, err := runRole(packet, role, opts, rawDir)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
func SynthesizeLedger(packet Packet, runs RunSet, existing *Ledger) Ledger {
	// SynthesizeLedger keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	findings := synthesizeLedgerFindings(runs, existingFindings(existing))
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packet.PacketDigest, Findings: findings}
}

func existingFindings(existing *Ledger) map[string]LedgerFinding {
	// existingFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	byFinding := map[string]LedgerFinding{}
	if existing != nil {
		for _, finding := range existing.Findings {
			byFinding[finding.ID] = finding
		}
	}
	return byFinding
}

func synthesizeLedgerFindings(runs RunSet, byFinding map[string]LedgerFinding) []LedgerFinding {
	// synthesizeLedgerFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	findings := []LedgerFinding{}
	for _, result := range runs.Results {
		for _, finding := range result.Findings {
			findings = append(findings, ledgerFindingFromReviewFinding(result, finding, byFinding))
		}
	}
	return findings
}

func ledgerFindingFromReviewFinding(result ReviewerResult, finding Finding, byFinding map[string]LedgerFinding) LedgerFinding {
	// ledgerFindingFromReviewFinding keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	id := ledgerFindingID(result, finding)
	return LedgerFinding{
		ID:           id,
		ReviewRunID:  result.ReviewRunID,
		Plane:        result.Plane,
		RoleID:       result.RoleID,
		Severity:     safeSeverity(finding.Severity),
		Summary:      safeText(finding.Summary),
		Citation:     finding.Citation,
		Disposition:  carriedLedgerDisposition(id, finding, byFinding),
		EvidenceRefs: finding.EvidenceRefs,
	}
}

func ledgerFindingID(result ReviewerResult, finding Finding) string {
	// ledgerFindingID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if finding.ID != "" {
		return finding.ID
	}

	return result.ReviewRunID + "-finding"
}

func carriedLedgerDisposition(id string, finding Finding, byFinding map[string]LedgerFinding) string {
	// carriedLedgerDisposition keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if prior, ok := byFinding[id]; ok && prior.Disposition != "" {

		return prior.Disposition
	}
	return defaultDisposition(finding.Severity)
}
func Validate(packet Packet, profile ReviewProfile, runs RunSet, ledger Ledger) Validation {
	// Validate keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	required := requiredPlaneSet(profile.RequiredPlanes)
	reasons := []string{}
	nextActions := []string{}
	cannotVerify := appendDigestValidation(packet, runs, ledger, &reasons, &nextActions)
	planeResults, usableCount, planesCannotVerify := validateRequiredPlanes(required, reviewRolesByID(profile.Roles), runs, &reasons, &nextActions)
	safeFindings, unresolved, findingsCannotVerify := validateLedgerFindings(packet, ledger, &reasons)
	state := reviewCoverageState(required, usableCount, cannotVerify || planesCannotVerify || findingsCannotVerify, unresolved)

	return validationResult(packet, planeResults, safeFindings, state, reasons, nextActions)
}

func reviewRolesByID(roles []ReviewRole) map[string]ReviewRole {
	// reviewRolesByID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	roleByID := map[string]ReviewRole{}
	for _, role := range roles {

		roleByID[role.RoleID] = role
	}
	return roleByID
}

func validationResult(packet Packet, planeResults []PlaneResult, findings []LedgerFinding, state string, reasons, nextActions []string) Validation {
	// validationResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return Validation{
		SchemaVersion:       SchemaVersionValidation,
		PacketDigest:        packet.PacketDigest,
		ReviewCoverageState: state,
		CIState:             packet.CIState,
		AuthorityScope:      AuthorityReviewRecordOnly,
		MergeDecision:       DecisionNotAuthorized,
		ReleaseDecision:     DecisionNotAuthorized,
		RiskAcceptance:      DecisionNotAuthorized,
		PlaneResults:        planeResults,
		Findings:            findings,
		Reasons:             uniqueStrings(reasons),
		NextActions:         uniqueStrings(nextActions),
	}
}

func requiredPlaneSet(planes []string) map[string]bool {
	// requiredPlaneSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	required := map[string]bool{}
	for _, plane := range planes {
		if plane != "" {

			required[plane] = true
		}
	}
	return required
}

func appendDigestValidation(packet Packet, runs RunSet, ledger Ledger, reasons, nextActions *[]string) bool {
	// appendDigestValidation keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	cannotVerify := false
	if packetDigestMismatch(packet, runs, ledger) {
		appendValidationAction(reasons, nextActions, "packet_digest_mismatch", "Create a new packet and rerun review for the current head.")
		cannotVerify = true
	}
	for _, result := range runs.Results {
		if result.PacketDigest != packet.PacketDigest {

			appendValidationAction(reasons, nextActions, "result_packet_digest_mismatch:"+safeID(result.ReviewRunID), "Discard stale reviewer results and rerun review for the current packet.")
			cannotVerify = true
		}
	}
	return cannotVerify
}

func packetDigestMismatch(packet Packet, runs RunSet, ledger Ledger) bool {
	return runs.PacketDigest != packet.PacketDigest || ledger.PacketDigest != packet.PacketDigest
}

func appendValidationAction(reasons, nextActions *[]string, reason, nextAction string) {
	// appendValidationAction keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	*reasons = append(*reasons, reason)
	*nextActions = append(*nextActions, nextAction)
}

func validateRequiredPlanes(required map[string]bool, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) ([]PlaneResult, int, bool) {
	// validateRequiredPlanes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	planeResults := make([]PlaneResult, 0, len(required))
	usableCount := 0
	cannotVerify := false
	for plane := range required {
		result := requiredPlaneResult(plane, roleByID, runs, reasons, nextActions)
		usableCount += boolCount(result.Usable)
		cannotVerify = cannotVerify || planeCannotVerify(result.Status)
		planeResults = append(planeResults, result)
	}
	sort.Slice(planeResults, func(i, j int) bool { return planeResults[i].Plane < planeResults[j].Plane })
	return planeResults, usableCount, cannotVerify
}

func requiredPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) PlaneResult {
	// requiredPlaneResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	result := bestPlaneResult(plane, roleByID, runs)
	appendPlaneValidationNotes(result, reasons, nextActions)
	return result
}

func boolCount(value bool) int {
	if value {
		return 1
	}
	return 0
}

func appendPlaneValidationNotes(result PlaneResult, reasons, nextActions *[]string) {
	// appendPlaneValidationNotes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if result.Reason != "" {

		*reasons = append(*reasons, fmt.Sprintf("%s:%s", result.Plane, result.Status))
	}
	if result.NextAction != "" && !result.Usable {

		*nextActions = append(*nextActions, result.NextAction)
	}
}

func bestPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet) PlaneResult {
	// bestPlaneResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	best := PlaneResult{Plane: plane, Status: StateNotAssessed, Usable: false, Reason: "required_plane_not_assessed", NextAction: "Run or import a reviewer result for this plane."}
	for _, result := range runs.Results {
		if result.Plane != plane {
			continue
		}
		candidate := planeResultWithModelCheck(roleByID[result.RoleID], result)
		if planeResultRank(candidate) > planeResultRank(best) {
			best = candidate
		}
	}
	return best
}

func planeResultRank(result PlaneResult) int {
	// Prefer retained, usable reviewer output over failed attempts. Findings
	// outrank no-findings so a later clean run cannot hide reported blockers.
	if result.Usable {
		return usablePlaneResultRank(result.Status)
	}
	return nonUsablePlaneResultRank(result.Status)
}

func usablePlaneResultRank(status string) int {
	// A usable findings result must dominate a usable no-findings result so a
	// later clean retry cannot erase blockers that remain in the ledger.
	if status == StatusFindingsReported {
		return 4
	}
	return 3
}

func nonUsablePlaneResultRank(status string) int {
	// Non-usable results still need ordering: verifier errors should be visible
	// ahead of lower-authority statuses, while not_assessed remains the floor.
	if planeCannotVerify(status) {
		return 2
	}
	if status != StateNotAssessed {
		return 1
	}
	return 0
}

func planeResultWithModelCheck(role ReviewRole, result ReviewerResult) PlaneResult {
	// planeResultWithModelCheck keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	best := planeResult(result)
	if best.Usable && modelMismatchWithoutFallback(role, result) {
		best.Usable = false
		best.Status = StatusCannotVerify
		best.Reason = "model_identity_mismatch"
		best.NextAction = "Rerun the reviewer or record fallback provenance for the observed model."
	}
	return best
}

func planeCannotVerify(status string) bool {
	// planeCannotVerify keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch status {

	case StatusCannotVerify, StatusTimedOut, StatusEmptyOutput, StatusOffTask, StatusParseFailed:
		return true
	default:
		return false
	}
}

func validateLedgerFindings(packet Packet, ledger Ledger, reasons *[]string) ([]LedgerFinding, bool, bool) {
	// validateLedgerFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	unresolved := false
	cannotVerify := false
	safeFindings := make([]LedgerFinding, 0, len(ledger.Findings))
	for _, finding := range ledger.Findings {
		finding.Summary = safeText(finding.Summary)
		unresolved = unresolved || ledgerFindingUnresolved(finding)
		cannotVerify = cannotVerify || appendCitationReasonIfUnresolvable(packet, finding, reasons)
		safeFindings = append(safeFindings, finding)
	}
	return safeFindings, unresolved, cannotVerify
}

func ledgerFindingUnresolved(finding LedgerFinding) bool {
	return (finding.Severity == SeverityCritical || finding.Severity == SeverityMajor) && finding.Disposition == DispositionUnresolvedReviewBlocker
}

func appendCitationReasonIfUnresolvable(packet Packet, finding LedgerFinding, reasons *[]string) bool {
	// appendCitationReasonIfUnresolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if citationResolvable(packet, finding.Citation) {
		return false
	}

	*reasons = append(*reasons, "finding_citation_cannot_verify")
	return true
}

func reviewCoverageState(required map[string]bool, usableCount int, cannotVerify, unresolved bool) string {
	// reviewCoverageState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if cannotVerify {
		return CoverageCannotVerify
	}
	if noReviewCoverage(required, usableCount) {
		return CoverageNotAssessed
	}
	return assessedReviewCoverageState(required, usableCount, unresolved)
}

func noReviewCoverage(required map[string]bool, usableCount int) bool {
	return len(required) == 0 || usableCount == 0
}

func assessedReviewCoverageState(required map[string]bool, usableCount int, unresolved bool) string {
	// assessedReviewCoverageState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if usableCount < len(required) {

		return CoveragePartial
	}
	if unresolved {

		return CoverageUnresolved
	}
	return CoverageSatisfied
}

func modelMismatchWithoutFallback(role ReviewRole, result ReviewerResult) bool {
	// modelMismatchWithoutFallback keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	requested := defaultString(role.RequestedModel, result.RequestedModel)
	observed := result.ObservedModel
	if modelIdentityMissing(requested) || modelIdentityMissing(observed) {
		return false
	}
	if requested == observed {
		return false
	}
	return fallbackMetadataMissing(result)
}

func modelIdentityMissing(model string) bool {
	return model == "" || model == StateNotAssessed
}

func fallbackMetadataMissing(result ReviewerResult) bool {
	return result.FallbackForModel == "" || result.FallbackReason == ""
}
func Summarize(validation Validation, ledger Ledger) string {
	// Summary is a rendered view over validation and ledger state; it does not
	// recompute coverage or mutate dispositions.
	var b strings.Builder
	writeSummaryHeader(&b, validation)
	writeSummaryPlanes(&b, validation.PlaneResults)
	writeSummaryFindings(&b, ledger.Findings)
	return b.String()
}

func writeSummaryHeader(b *strings.Builder, validation Validation) {
	// writeSummaryHeader keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	fmt.Fprintf(b, "Review coverage: %s\n", validation.ReviewCoverageState)
	fmt.Fprintf(b, "CI state: %s\n", validation.CIState)
	fmt.Fprintf(b, "Authority scope: %s\n", validation.AuthorityScope)
	fmt.Fprintf(b, "Merge decision: %s\n", validation.MergeDecision)
	fmt.Fprintf(b, "Release decision: %s\n", validation.ReleaseDecision)
	fmt.Fprintf(b, "Risk acceptance: %s\n", validation.RiskAcceptance)
	b.WriteString("This is review-record evidence only; merge, release, and risk decisions remain external.\n")
}

func writeSummaryPlanes(b *strings.Builder, planes []PlaneResult) {
	// writeSummaryPlanes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(planes) > 0 {
		b.WriteString("\nPlanes\n")
		for _, plane := range planes {

			fmt.Fprintf(b, "- %s: %s", plane.Plane, plane.Status)
			if plane.NextAction != "" {
				fmt.Fprintf(b, " next_action=%s", safeText(plane.NextAction))
			}
			b.WriteString("\n")
		}
	}
}

func writeSummaryFindings(b *strings.Builder, findings []LedgerFinding) {
	// writeSummaryFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(findings) > 0 {
		b.WriteString("\nFindings\n")
		for _, finding := range findings {

			fmt.Fprintf(b, "- %s [%s] %s (%s)\n", finding.ID, finding.Severity, safeText(finding.Summary), finding.Disposition)
		}
	}
}

func WriteJSON(path string, value any) error {
	// WriteJSON keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func ReadPacket(path string) (Packet, error) {
	// ReadPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var packet Packet
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {

		path = filepath.Join(path, "packet.json")
	}
	return packet, readJSON(path, &packet)
}

func ReadProfile(path string) (ReviewProfile, error) {
	// ReadProfile keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var profile ReviewProfile
	if err := readJSON(path, &profile); err != nil {
		return profile, err
	}

	return profile, validateProfile(profile)
}

func ReadRunSet(path string) (RunSet, error) {
	// ReadRunSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var runs RunSet

	path = runSetPath(path)
	if err := readJSON(path, &runs); err != nil {
		return runs, err
	}
	if err := validateRunSet(runs); err != nil {
		return runs, err
	}
	return runs, nil
}

func runSetPath(path string) string {
	// runSetPath keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {

		return filepath.Join(path, "results.json")
	}
	return path
}

func ReadLedger(path string) (Ledger, error) {
	var ledger Ledger
	return ledger, readJSON(path, &ledger)
}

func ReadValidation(path string) (Validation, error) {
	var validation Validation
	return validation, readJSON(path, &validation)
}

func readJSON(path string, value any) error {
	// readJSON keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func validatePacketOptions(opts PacketOptions) error {
	// validatePacketOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(opts.OutDir) == "" {

		return errors.New("pr_review_packet_requires_out")
	}
	if err := validatePacketIdentityOptions(opts); err != nil {
		return err
	}
	if err := validatePacketInputOptions(opts); err != nil {
		return err
	}
	return nil
}

func validatePacketIdentityOptions(opts PacketOptions) error {
	// validatePacketIdentityOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !repoIDPattern.MatchString(opts.RepoID) {
		return fmt.Errorf("unsafe_repo_id: repo_id must match %s", repoIDPattern.String())
	}
	if !changeRefPattern.MatchString(opts.ChangeRef) {
		return fmt.Errorf("unsafe_change_ref: change_ref must match %s", changeRefPattern.String())
	}
	if !validPacketCommits(opts) {
		return errors.New("invalid_commit_sha: base and head must be 40 lowercase hex characters")
	}
	return nil
}

func validPacketCommits(opts PacketOptions) bool {
	return sha40Pattern.MatchString(opts.BaseCommit) && sha40Pattern.MatchString(opts.HeadCommit)
}

func validatePacketInputOptions(opts PacketOptions) error {
	// validatePacketInputOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(opts.DiffPath) == "" {

		return errors.New("pr_review_packet_requires_diff")
	}
	if invalidPacketCIState(opts) {

		return fmt.Errorf("invalid_ci_state: %s", opts.CIState)
	}
	return nil
}

func invalidPacketCIState(opts PacketOptions) bool {
	return opts.CIState != "" && !validCIState(opts.CIState)
}

func validateRunSet(runs RunSet) error {
	// validateRunSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	seen := map[string]bool{}
	for _, result := range runs.Results {
		if strings.TrimSpace(result.ReviewRunID) == "" {

			return errors.New("review_result_requires_review_run_id")
		}
		if seen[result.ReviewRunID] {

			return fmt.Errorf("duplicate_review_run_id: %s", result.ReviewRunID)
		}
		seen[result.ReviewRunID] = true
	}
	return nil
}

func validateProfile(profile ReviewProfile) error {
	// validateProfile keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validateProfileHeader(profile); err != nil {
		return err
	}
	rolePlanes, err := validateProfileRoles(profile.Roles)
	if err != nil {
		return err
	}
	return validateRequiredPlaneRoles(profile.RequiredPlanes, rolePlanes)
}

func validateProfileHeader(profile ReviewProfile) error {
	// validateProfileHeader keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if profile.SchemaVersion != "" && profile.SchemaVersion != SchemaVersionProfile {

		return fmt.Errorf("invalid_profile_schema_version: %s", profile.SchemaVersion)
	}
	return requireProfileFields(profile)
}

func requireProfileFields(profile ReviewProfile) error {
	// requireProfileFields keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(profile.ProfileID) == "" {

		return errors.New("profile_requires_profile_id")
	}
	if len(profile.RequiredPlanes) == 0 {

		return errors.New("profile_requires_required_planes")
	}
	if len(profile.Roles) == 0 {

		return errors.New("profile_requires_roles")
	}
	return nil
}

func validateProfileRoles(roles []ReviewRole) (map[string]bool, error) {
	// validateProfileRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	rolePlanes := map[string]bool{}
	for _, role := range roles {
		if err := validateProfileRole(role); err != nil {
			return nil, err
		}
		rolePlanes[role.Plane] = true
	}
	return rolePlanes, nil
}

func validateProfileRole(role ReviewRole) error {
	// validateProfileRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if profileRoleMissingRequiredField(role) {

		return errors.New("profile_role_requires_id_plane_runner")
	}
	if !validRunner(role.Runner) {

		return fmt.Errorf("profile_role_invalid_runner: %s", role.Runner)
	}
	return nil
}

func profileRoleMissingRequiredField(role ReviewRole) bool {

	return role.RoleID == "" || role.Plane == "" || role.Runner == ""
}

func validateRequiredPlaneRoles(requiredPlanes []string, rolePlanes map[string]bool) error {
	// validateRequiredPlaneRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	for _, plane := range requiredPlanes {
		if !rolePlanes[plane] {
			return fmt.Errorf("profile_required_plane_without_role: %s", plane)
		}
	}
	return nil
}
func runRole(packet Packet, role ReviewRole, opts RunOptions, rawDir string) (ReviewerResult, error) {
	// runRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	result := newReviewerResult(packet, role, opts.Now)
	result.CommandDigest = commandDigest(role.Command)
	baseline, ready, err := prepareRoleRunner(&result, role, opts)
	if err != nil || !ready {
		return result, err
	}
	output, timedOut, err := runRoleCommand(role, opts)
	result.EndedAt = time.Now().UTC().Format(time.RFC3339)
	result = completeRoleResult(result, role, packet, opts.WorkDir, baseline, output, timedOut, err)
	return writeRawResult(result, rawDir, output)
}

func completeRoleResult(result ReviewerResult, role ReviewRole, packet Packet, workDir string, baseline *workingTreeBaseline, output []byte, timedOut bool, runErr error) ReviewerResult {
	// completeRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if completed, ok := completeUnparsedRoleResult(result, output, timedOut, runErr); ok {
		return completed
	}
	parsed, err := parseReviewerOutput(result, role, packet, output)
	return completeParsedRoleResult(markRoleParseFailure(parsed, err), role, workDir, baseline)
}

func completeUnparsedRoleResult(result ReviewerResult, output []byte, timedOut bool, runErr error) (ReviewerResult, bool) {
	// completeUnparsedRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if timedOut {
		result.Status = StatusTimedOut
		result.Reason = "runner_timed_out"
		return result, true
	}
	return completeUnparsedRoleErrorOrEmpty(result, output, runErr)
}

func completeUnparsedRoleErrorOrEmpty(result ReviewerResult, output []byte, runErr error) (ReviewerResult, bool) {
	// completeUnparsedRoleErrorOrEmpty keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch {
	case applyRunnerError(&result, runErr) != nil:
		return result, true
	case emptyReviewerOutput(output):
		result.Status = StatusEmptyOutput
		result.Reason = "runner_empty_output"
		return result, true
	default:
		return result, false
	}
}

func markRoleParseFailure(result ReviewerResult, err error) ReviewerResult {
	// markRoleParseFailure keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err == nil {
		return result
	}
	result.Status = StatusParseFailed
	result.Reason = "runner_output_parse_failed"
	return result
}

func completeParsedRoleResult(parsed ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) ReviewerResult {
	// completeParsedRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	applyOpenCodeMutationCheck(&parsed, role, workDir, baseline)
	return parsed
}

func emptyReviewerOutput(output []byte) bool {
	return len(strings.TrimSpace(string(output))) == 0
}

func newReviewerResult(packet Packet, role ReviewRole, now time.Time) ReviewerResult {
	// newReviewerResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	started := now.Format(time.RFC3339)

	return ReviewerResult{
		ReviewRunID:    safeID("run-" + role.RoleID),
		PacketDigest:   packet.PacketDigest,
		Plane:          role.Plane,
		RoleID:         role.RoleID,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		ObservedModel:  StateNotAssessed,
		ModelFamily:    StateNotAssessed,
		ModelVersion:   StateNotAssessed,
		Status:         StatusNotAssessed,
		StartedAt:      started,
		EndedAt:        started,
	}
}

func prepareRoleRunner(result *ReviewerResult, role ReviewRole, opts RunOptions) (*workingTreeBaseline, bool, error) {
	// prepareRoleRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if role.Runner != RunnerManualExternal && !opts.AllowedRunners[role.Runner] {
		return nil, false, fmt.Errorf("runner_not_allowed: %s", role.Runner)
	}
	if role.Runner == RunnerOpenCode {
		return prepareOpenCodeBaseline(result, role, opts.WorkDir)
	}
	return prepareCommandRunner(result, role)
}

func prepareCommandRunner(result *ReviewerResult, role ReviewRole) (*workingTreeBaseline, bool, error) {
	// prepareCommandRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if err := attachPromptRef(result, role); err != nil {
		return nil, false, nil
	}

	return nil, commandConfigured(result, role), nil
}

func attachPromptRef(result *ReviewerResult, role ReviewRole) error {
	// attachPromptRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	promptRef, err := promptSafeRef(role)
	if err != nil {

		result.Status = StatusCannotVerify
		result.Reason = "prompt_ref_cannot_verify"
		return err
	}
	result.PromptRef = promptRef
	return nil
}

func commandConfigured(result *ReviewerResult, role ReviewRole) bool {
	// commandConfigured keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(role.Command) == 0 {

		result.Reason = "runner_command_not_configured"
		return false
	}
	return true
}

func prepareOpenCodeBaseline(result *ReviewerResult, role ReviewRole, workDir string) (*workingTreeBaseline, bool, error) {
	// prepareOpenCodeBaseline keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !openCodeReadOnlyReady(result, role) {
		return nil, false, nil
	}
	baseline, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "working_tree_baseline_cannot_verify"
		return nil, false, nil
	}
	return openCodeBaselineReady(result, role, baseline)
}

func openCodeBaselineReady(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) (*workingTreeBaseline, bool, error) {
	// openCodeBaselineReady keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if !openCodeBaselineClean(result, role, baseline) {
		return nil, false, nil
	}

	return baseline, commandConfigured(result, role), nil
}

func openCodeReadOnlyReady(result *ReviewerResult, role ReviewRole) bool {
	// openCodeReadOnlyReady keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if !role.ReadOnlyEnforced {
		markOpenCodeReadOnlyMissing(result)
		return false
	}
	if err := attachPromptRef(result, role); err != nil {

		return false
	}
	return true
}

func openCodeBaselineClean(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) bool {
	// openCodeBaselineClean keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	mode := defaultString(role.WorkingTreeMode, "clean_required")
	if mode != "clean_required" || baseline.Count == 0 {
		return true
	}
	result.Status = StatusNotAssessed
	result.Reason = "working_tree_dirty"
	return false
}

func markOpenCodeReadOnlyMissing(result *ReviewerResult) {

	result.Reason = "opencode_read_only_not_enforced"
}

func runRoleCommand(role ReviewRole, opts RunOptions) ([]byte, bool, error) {
	// runRoleCommand keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	timeout := time.Duration(role.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, role.Command[0], role.Command[1:]...)
	cmd.Dir = opts.WorkDir
	output, err := cmd.Output()
	return output, errors.Is(ctx.Err(), context.DeadlineExceeded), err
}

func applyRunnerError(result *ReviewerResult, err error) error {
	// applyRunnerError keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err == nil {
		return nil
	}
	if errors.Is(err, exec.ErrNotFound) || strings.Contains(err.Error(), "executable file not found") {

		result.Status = StatusNotAssessed
		result.Reason = "runner_unavailable"
	} else {

		result.Status = StatusFailed
		result.Reason = "runner_failed"
	}
	return err
}

func applyOpenCodeMutationCheck(result *ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) {
	// applyOpenCodeMutationCheck keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !needsOpenCodeMutationCheck(role, baseline) {
		return
	}
	after, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		markBaselineCannotVerify(result)
		return
	}
	if baselineChanged(after, baseline) {
		result.Status = StatusCannotVerify
		result.Reason = "mutation_detected"
	}
}

func needsOpenCodeMutationCheck(role ReviewRole, baseline *workingTreeBaseline) bool {
	return role.Runner == RunnerOpenCode && baseline != nil
}

func markBaselineCannotVerify(result *ReviewerResult) {
	result.Status = StatusCannotVerify
	result.Reason = "working_tree_baseline_cannot_verify"
}

func baselineChanged(after, before *workingTreeBaseline) bool {
	return after.Digest != before.Digest || after.Count != before.Count
}

type workingTreeBaseline struct {
	Count  int
	Digest string
}

func captureWorkingTreeBaseline(workDir string) (*workingTreeBaseline, error) {
	// captureWorkingTreeBaseline keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	if strings.TrimSpace(string(output)) != "" {

		count = len(lines)
	}
	sum := sha256.Sum256(output)
	return &workingTreeBaseline{Count: count, Digest: hex.EncodeToString(sum[:])}, nil
}

func parseReviewerOutput(base ReviewerResult, role ReviewRole, packet Packet, output []byte) (ReviewerResult, error) {
	// Parsed output must echo the packet, plane, and role. A structurally valid
	// answer for another packet is off-task, not reusable evidence.
	var parsed ReviewerResult
	if err := decodeReviewerOutput(output, &parsed); err != nil {
		return base, err
	}
	if reviewerOutputMismatched(parsed, role, packet) {
		base.Status = StatusOffTask
		base.Reason = "reviewer_output_wrong_packet_plane_or_role"
		return base, nil
	}
	return normalizeParsedReviewerOutput(parsed, base, role), nil
}

func reviewerOutputMismatched(parsed ReviewerResult, role ReviewRole, packet Packet) bool {

	return parsed.PacketDigest != packet.PacketDigest || parsed.Plane != role.Plane || parsed.RoleID != role.RoleID
}

func decodeReviewerOutput(output []byte, parsed *ReviewerResult) error {
	// decodeReviewerOutput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	decoder := json.NewDecoder(strings.NewReader(string(output)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(parsed)
}

func normalizeParsedReviewerOutput(parsed, base ReviewerResult, role ReviewRole) ReviewerResult {
	// normalizeParsedReviewerOutput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.ReviewRunID = defaultString(parsed.ReviewRunID, base.ReviewRunID)
	parsed.Runner = defaultString(parsed.Runner, role.Runner)
	normalizeParsedReviewerModels(&parsed, role)
	attachParsedReviewerExecution(&parsed, base)
	normalizeParsedReviewerStatus(&parsed)
	return parsed
}

func normalizeParsedReviewerModels(parsed *ReviewerResult, role ReviewRole) {
	// normalizeParsedReviewerModels keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.RequestedModel = defaultString(parsed.RequestedModel, defaultString(role.RequestedModel, StateNotAssessed))
	parsed.ObservedModel = defaultString(parsed.ObservedModel, StateNotAssessed)
	parsed.ModelFamily = defaultString(parsed.ModelFamily, StateNotAssessed)
	parsed.ModelVersion = defaultString(parsed.ModelVersion, StateNotAssessed)
}

func attachParsedReviewerExecution(parsed *ReviewerResult, base ReviewerResult) {
	// attachParsedReviewerExecution keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.StartedAt = base.StartedAt
	parsed.EndedAt = base.EndedAt
	parsed.CommandDigest = base.CommandDigest
	parsed.PromptRef = base.PromptRef
}

func normalizeParsedReviewerStatus(parsed *ReviewerResult) {
	// normalizeParsedReviewerStatus keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if parsed.Status == "" {

		parsed.Status = defaultReviewerStatus(parsed.Findings)
	}
}

func defaultReviewerStatus(findings []Finding) string {
	// defaultReviewerStatus keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(findings) > 0 {

		return StatusFindingsReported
	}
	return StatusNoFindings
}

func writeRawResult(result ReviewerResult, rawDir string, output []byte) (ReviewerResult, error) {
	// writeRawResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if output != nil {

		name := safeID(result.ReviewRunID) + ".out"
		path := filepath.Join(rawDir, name)
		if err := os.WriteFile(path, output, 0o600); err != nil {
			return result, err
		}
		digest := sha256.Sum256(output)
		result.RawOutputRef = &SafeRef{ID: "raw-" + safeID(result.ReviewRunID), Kind: RefKindRawOutput, Ref: filepath.ToSlash(filepath.Join("raw", name)), DigestSHA256: hex.EncodeToString(digest[:]), ContentType: ContentText, RedactionState: RedactionDigestOnly}
	}
	return result, nil
}
func preview(packet Packet, profile ReviewProfile) *RunPreview {
	// preview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	roles := make([]PreviewRole, 0, len(profile.Roles))
	for _, role := range profile.Roles {
		roles = append(roles, previewRole(role))
	}
	return &RunPreview{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Roles: roles}
}

func previewRole(role ReviewRole) PreviewRole {
	// previewRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return PreviewRole{
		RoleID:         role.RoleID,
		Plane:          role.Plane,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		TimeoutSeconds: role.TimeoutSeconds,
		CommandDigest:  commandDigest(role.Command),
		PromptRef:      role.PromptTemplateRef,
		PromptDigest:   promptDigestForRole(role),
	}
}

func promptDigestForRole(role ReviewRole) string {
	// promptDigestForRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	promptRef, _ := promptSafeRef(role)
	if promptRef == nil {
		return ""
	}

	return promptRef.DigestSHA256
}

func promptSafeRef(role ReviewRole) (*SafeRef, error) {
	// promptSafeRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return nil, nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {

		return nil, err
	}
	return promptDigestRef(role, data), nil
}

func promptDigestRef(role ReviewRole, data []byte) *SafeRef {
	// promptDigestRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	sum := sha256.Sum256(data)
	return &SafeRef{
		ID:             "prompt-" + safeID(role.RoleID),
		Kind:           RefKindPrompt,
		Ref:            "digest-only:" + safeID(filepath.Base(role.PromptTemplateRef)),
		DigestSHA256:   hex.EncodeToString(sum[:]),
		ContentType:    contentType(role.PromptTemplateRef),
		RedactionState: RedactionDigestOnly,
	}
}

func copyInputs(inputDir, prefix string, paths []string) ([]SafeRef, error) {
	// copyInputs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	refs := make([]SafeRef, 0, len(paths))
	for i, path := range paths {
		ref, err := copiedInputRef(inputDir, prefix, i, path)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func copiedInputRef(inputDir, prefix string, index int, path string) (SafeRef, error) {
	// copiedInputRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	position := index + 1
	name := fmt.Sprintf("%s-%d%s", prefix, position, normalizedExt(path))
	ref, err := copyInput(inputDir, name, path, RefKindDoc, contentType(path))
	if err != nil {
		return SafeRef{}, err
	}
	ref.ID = fmt.Sprintf("%s-%d", prefix, position)
	return ref, nil
}

func copyInput(inputDir, name, source, kind, contentType string) (SafeRef, error) {
	// copyInput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	data, err := os.ReadFile(source)
	if err != nil {
		return SafeRef{}, err
	}
	if err := writeCopiedInput(inputDir, name, data); err != nil {
		return SafeRef{}, err
	}
	return copiedInputSafeRef(name, kind, contentType, data), nil
}

func writeCopiedInput(inputDir, name string, data []byte) error {
	// writeCopiedInput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	dest := filepath.Join(inputDir, name)
	if err := os.WriteFile(dest, data, 0o644); err != nil {

		return err
	}
	return nil
}

func copiedInputSafeRef(name, kind, contentType string, data []byte) SafeRef {
	// copiedInputSafeRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	digest := sha256.Sum256(data)
	return SafeRef{
		ID:             strings.TrimSuffix(name, filepath.Ext(name)),
		Kind:           kind,
		Ref:            filepath.ToSlash(filepath.Join("inputs", name)),
		DigestSHA256:   hex.EncodeToString(digest[:]),
		ContentType:    contentType,
		RedactionState: RedactionNone,
	}
}

func digestJSON(value any) (string, error) {
	// digestJSON keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	data, err := json.Marshal(value)
	if err != nil {

		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func packetDigest(packet Packet) (string, error) {
	// packetDigest keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	canonical := packet
	canonical.PacketDigest = ""
	return digestJSON(canonical)
}
func ensureNewDir(path string) error {
	// ensureNewDir keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(path) == "" {
		return errors.New("missing_output_path")
	}
	entries, err := os.ReadDir(path)
	if dirHasEntries(entries, err) {
		return fmt.Errorf("output_exists: %s", filepath.Base(path))
	}
	if readDirFailed(err) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func dirHasEntries(entries []os.DirEntry, err error) bool {
	return err == nil && len(entries) > 0
}

func readDirFailed(err error) bool {
	return err != nil && !errors.Is(err, os.ErrNotExist)
}

func validCIState(state string) bool {
	// validCIState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch state {
	case StatePass, StateFail, StatePending, StateNotAssessed, StateCannotVerify:

		return true
	default:
		return false
	}
}

func validRunner(runner string) bool {
	// validRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch runner {
	case RunnerPI, RunnerOpenCode, RunnerManualExternal:
		return true
	default:
		return false
	}
}

func planeResult(result ReviewerResult) PlaneResult {
	// planeResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	pr := PlaneResult{Plane: result.Plane, Status: result.Status, RunID: result.ReviewRunID}
	if reviewerResultUsable(result) {
		pr.Usable = true
		return pr
	}
	if reviewerStatusUsable(result.Status) {
		pr.Status = StatusCannotVerify
		pr.Reason, pr.NextAction = reviewerStatusAction(result.Status)
		return pr
	}
	pr.Reason, pr.NextAction = reviewerStatusAction(result.Status)
	return pr
}

func reviewerStatusUsable(status string) bool {
	return status == StatusFindingsReported || status == StatusNoFindings
}

func reviewerResultUsable(result ReviewerResult) bool {
	// Usable review coverage requires both a positive reviewer status and retained
	// raw output evidence. A hand-authored status without a digest-bound output
	// reference stays unverifiable.
	return reviewerStatusUsable(result.Status) && result.RawOutputRef != nil
}

func reviewerStatusAction(status string) (string, string) {
	// reviewerStatusAction keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	actions := map[string][2]string{
		StatusNotAssessed: {"reviewer_not_assessed", "Run a configured reviewer or import a usable result for this plane."},
		StatusTimedOut:    {"reviewer_timed_out", "Increase timeout or replace the reviewer for this plane."},
		StatusEmptyOutput: {"reviewer_empty_output", "Retry with a shorter bounded prompt or replace the reviewer."},
		StatusOffTask:     {"reviewer_off_task", "Rerun with the frozen packet and required output schema."},
		StatusParseFailed: {"reviewer_parse_failed", "Rerun with JSON-only output matching the required schema."},
	}
	if action, ok := actions[status]; ok {
		return action[0], action[1]
	}
	if reviewerStatusUsable(status) {
		return "reviewer_output_not_retained", "Attach digest-bound reviewer output before counting this plane."
	}
	return "reviewer_cannot_verify", "Replace or rerun the reviewer."
}

func defaultDisposition(severity string) string {
	// defaultDisposition keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch safeSeverity(severity) {
	case SeverityCritical, SeverityMajor:
		return DispositionUnresolvedReviewBlocker
	default:
		return DispositionDeferredNotAssessed
	}
}

func safeSeverity(severity string) string {
	// safeSeverity keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch severity {
	case SeverityCritical, SeverityMajor, SeverityMinor, SeverityInformational:

		return severity
	default:
		return SeverityInformational
	}
}

func citationResolvable(packet Packet, citation Citation) bool {
	// citationResolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !citationHasAnchor(citation) {
		return false
	}
	if resolvable, ok := citationRefResolvable(packet, citation); ok {
		return resolvable
	}
	return citation.SourceDigest != ""
}

func citationHasAnchor(citation Citation) bool {
	return citation.ContextRefID != "" || citation.SourceDigest != ""
}

func citationRefResolvable(packet Packet, citation Citation) (bool, bool) {
	// citationRefResolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	for _, resolver := range citationResolvers {
		if resolver.matches(packet, citation) {
			return resolver.resolvable(citation), true
		}
	}
	return false, false
}

type citationResolver struct {
	matches    func(Packet, Citation) bool
	resolvable func(Citation) bool
}

var citationResolvers = []citationResolver{
	{matches: citationMatchesDiff, resolvable: citationHasDiffLocation},
	{matches: citationMatchesContext, resolvable: citationHasContextLocation},
	{matches: citationMatchesVerification, resolvable: citationHasVerificationLocation},
}

func citationMatchesDiff(packet Packet, citation Citation) bool {
	return citation.ContextRefID == packet.DiffRef.ID || citation.ContextRefID == "diff"
}

func citationMatchesContext(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.ContextRefs, citation.ContextRefID)
}

func citationMatchesVerification(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.VerificationRefs, citation.ContextRefID)
}

func safeRefIDExists(refs []SafeRef, id string) bool {
	// safeRefIDExists keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	for _, ref := range refs {
		if id == ref.ID {

			return true
		}
	}
	return false
}

func citationHasDiffLocation(citation Citation) bool {

	return citation.DiffHunkID != "" || citation.SourceDigest != ""
}

func citationHasContextLocation(citation Citation) bool {

	return citation.DiffHunkID != "" || citation.SourceDigest != "" || citation.LineStart > 0
}

func citationHasVerificationLocation(citation Citation) bool {

	return citation.SourceDigest != "" || citation.LineStart > 0
}

func safeText(text string) string {
	// safeText keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if text == "" {
		return ""
	}
	if containsUnsafeTextMarker(text) || containsUnsafeTextPattern(text) {
		return "[redacted unsafe reviewer text]"
	}
	return text
}

func containsUnsafeTextMarker(text string) bool {
	// containsUnsafeTextMarker keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	unsafeMarkers := []string{"SYNTHETIC_", "Bearer ", "access_token=", "BEGIN PRIVATE KEY", "PRIVATE_KEY", "cookie=", "session=", "/Users/", "/private/"}
	for _, marker := range unsafeMarkers {
		if strings.Contains(text, marker) {

			return true
		}
	}
	return false
}

func containsUnsafeTextPattern(text string) bool {
	return (strings.Contains(text, "://") && strings.Contains(text, "@")) || strings.Contains(text, "token=")
}

func uniqueStrings(values []string) []string {
	// uniqueStrings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = safeText(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func commandDigest(command []string) string {
	// commandDigest keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if len(command) == 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(strings.Join(command, "\x00")))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func contextKind(path string) string {
	// contextKind keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".md" && ext != ".markdown" {
		return contextKindByExtension(ext)
	}
	return markdownContextKind(path)
}

func contextKindByExtension(ext string) string {
	// contextKindByExtension keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch ext {
	case ".json":

		return RefKindSchema
	default:
		return RefKindSourceExcerpt
	}
}

func markdownContextKind(path string) string {
	// markdownContextKind keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.Contains(strings.ToLower(filepath.Base(path)), "task") {

		return RefKindTask
	}
	return RefKindDoc
}

func contentType(path string) string {
	// contentType keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":

		return ContentJSON
	case ".md", ".markdown":
		return ContentMarkdown
	default:
		return ContentText
	}
}

func normalizedExt(path string) string {
	// normalizedExt keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".md", ".txt", ".diff", ".patch":

		return ext
	default:
		return ".txt"
	}
}

func safeID(value string) string {
	// safeID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	value = strings.ToLower(strings.TrimSpace(value))
	out := strings.Trim(strings.Map(safeIDMapper, value), "-.")
	if out == "" {

		return "item"
	}
	return out
}

const safeIDAllowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_.-"

func safeIDMapper(r rune) rune {
	// safeIDMapper keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if r <= 127 && strings.IndexByte(safeIDAllowedChars, byte(r)) >= 0 {
		return r
	}

	return '-'
}

func defaultString(value, fallback string) string {
	// defaultString keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return fallback
	}
	return value
}

func Copy(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}
