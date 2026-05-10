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

type SafeRef struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Ref            string `json:"ref"`
	DigestSHA256   string `json:"digest_sha256"`
	ContentType    string `json:"content_type"`
	RedactionState string `json:"redaction_state"`
}

type UnavailableField struct {
	Field  string `json:"field"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}

type ReviewProfile struct {
	SchemaVersion  string       `json:"schema_version"`
	ProfileID      string       `json:"profile_id"`
	RequiredPlanes []string     `json:"required_planes"`
	Roles          []ReviewRole `json:"roles"`
}

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

type RunOptions struct {
	OutDir         string
	AllowedRunners map[string]bool
	Preview        bool
	Now            time.Time
	WorkDir        string
}

type RunPreview struct {
	SchemaVersion string        `json:"schema_version"`
	PacketDigest  string        `json:"packet_digest"`
	Roles         []PreviewRole `json:"roles"`
}

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

type RunSet struct {
	SchemaVersion string           `json:"schema_version"`
	PacketDigest  string           `json:"packet_digest"`
	Results       []ReviewerResult `json:"results"`
}

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

type Citation struct {
	ContextRefID string `json:"context_ref_id,omitempty"`
	DiffHunkID   string `json:"diff_hunk_id,omitempty"`
	SourceDigest string `json:"source_digest,omitempty"`
	LineStart    int    `json:"line_start,omitempty"`
	LineEnd      int    `json:"line_end,omitempty"`
}

type Ledger struct {
	SchemaVersion string          `json:"schema_version"`
	PacketDigest  string          `json:"packet_digest"`
	Findings      []LedgerFinding `json:"findings"`
}

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

type PlaneResult struct {
	Plane      string `json:"plane"`
	Status     string `json:"status"`
	Usable     bool   `json:"usable"`
	RunID      string `json:"review_run_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
	NextAction string `json:"next_action,omitempty"`
}

func BuildPacket(opts PacketOptions) (Packet, error) {
	if err := validatePacketOptions(opts); err != nil {
		return Packet{}, err
	}
	if err := ensureNewDir(opts.OutDir); err != nil {
		return Packet{}, err
	}
	now, createdBy, ciState := packetDefaults(opts)
	refs, err := buildPacketRefs(opts)
	if err != nil {
		return Packet{}, err
	}
	packet := Packet{
		SchemaVersion:     SchemaVersionPacket,
		PacketID:          fmt.Sprintf("%s-%s-%s", opts.RepoID, opts.ChangeRef, opts.HeadCommit[:12]),
		RepoID:            opts.RepoID,
		ChangeRef:         opts.ChangeRef,
		BaseCommit:        opts.BaseCommit,
		HeadCommit:        opts.HeadCommit,
		DiffRef:           refs.diff,
		MetadataRef:       refs.metadata,
		ContextRefs:       refs.context,
		VerificationRefs:  refs.verification,
		CIState:           ciState,
		CreatedAt:         now.Format(time.RFC3339),
		CreatedBy:         createdBy,
		RedactionState:    RedactionNone,
		UnavailableFields: unavailablePacketFields(opts),
	}
	digest, err := packetDigest(packet)
	if err != nil {
		return Packet{}, err
	}
	packet.PacketDigest = "sha256:" + digest
	if err := WriteJSON(filepath.Join(opts.OutDir, "packet.json"), packet); err != nil {
		return Packet{}, err
	}
	return packet, nil
}

type packetRefs struct {
	diff         SafeRef
	metadata     *SafeRef
	context      []SafeRef
	verification []SafeRef
}

func packetDefaults(opts PacketOptions) (time.Time, string, string) {
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
	inputDir := filepath.Join(opts.OutDir, "inputs")
	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return packetRefs{}, err
	}
	diffRef, err := copyInput(inputDir, "diff.patch", opts.DiffPath, RefKindDiff, ContentUnifiedDiff)
	if err != nil {
		return packetRefs{}, err
	}
	metadataRef, err := optionalMetadataRef(inputDir, opts.MetadataPath)
	if err != nil {
		return packetRefs{}, err
	}
	contextRefs, err := packetContextRefs(inputDir, opts.ContextPaths)
	if err != nil {
		return packetRefs{}, err
	}
	verificationRefs, err := packetVerificationRefs(inputDir, opts.VerificationPaths)
	if err != nil {
		return packetRefs{}, err
	}
	return packetRefs{diff: diffRef, metadata: metadataRef, context: contextRefs, verification: verificationRefs}, nil
}

func optionalMetadataRef(inputDir, metadataPath string) (*SafeRef, error) {
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
	if err := validateProfile(profile); err != nil {
		return RunSet{}, nil, err
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if opts.WorkDir == "" {
		opts.WorkDir = "."
	}
	if opts.Preview {
		return RunSet{}, preview(packet, profile), nil
	}
	if err := ensureNewDir(opts.OutDir); err != nil {
		return RunSet{}, nil, err
	}
	rawDir := filepath.Join(opts.OutDir, "raw")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return RunSet{}, nil, err
	}
	results := make([]ReviewerResult, 0, len(profile.Roles))
	for _, role := range profile.Roles {
		result, err := runRole(packet, role, opts, rawDir)
		if err != nil {
			return RunSet{}, nil, err
		}
		results = append(results, result)
	}
	runSet := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Results: results}
	if err := WriteJSON(filepath.Join(opts.OutDir, "results.json"), runSet); err != nil {
		return RunSet{}, nil, err
	}
	return runSet, nil, nil
}

func SynthesizeLedger(packet Packet, runs RunSet, existing *Ledger) Ledger {
	byFinding := map[string]LedgerFinding{}
	if existing != nil {
		for _, finding := range existing.Findings {
			byFinding[finding.ID] = finding
		}
	}
	findings := []LedgerFinding{}
	for _, result := range runs.Results {
		for _, finding := range result.Findings {
			id := finding.ID
			if id == "" {
				id = result.ReviewRunID + "-finding"
			}
			disposition := defaultDisposition(finding.Severity)
			if prior, ok := byFinding[id]; ok && prior.Disposition != "" {
				disposition = prior.Disposition
			}
			findings = append(findings, LedgerFinding{
				ID:           id,
				ReviewRunID:  result.ReviewRunID,
				Plane:        result.Plane,
				RoleID:       result.RoleID,
				Severity:     safeSeverity(finding.Severity),
				Summary:      safeText(finding.Summary),
				Citation:     finding.Citation,
				Disposition:  disposition,
				EvidenceRefs: finding.EvidenceRefs,
			})
		}
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packet.PacketDigest, Findings: findings}
}

func Validate(packet Packet, profile ReviewProfile, runs RunSet, ledger Ledger) Validation {
	roleByID := map[string]ReviewRole{}
	for _, role := range profile.Roles {
		roleByID[role.RoleID] = role
	}

	required := requiredPlaneSet(profile.RequiredPlanes)
	reasons := []string{}
	nextActions := []string{}
	cannotVerify := appendDigestValidation(packet, runs, ledger, &reasons, &nextActions)
	planeResults, usableCount, planesCannotVerify := validateRequiredPlanes(required, roleByID, runs, &reasons, &nextActions)
	safeFindings, unresolved, findingsCannotVerify := validateLedgerFindings(packet, ledger, &reasons)
	state := reviewCoverageState(required, usableCount, cannotVerify || planesCannotVerify || findingsCannotVerify, unresolved)

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
		Findings:            safeFindings,
		Reasons:             uniqueStrings(reasons),
		NextActions:         uniqueStrings(nextActions),
	}
}

func requiredPlaneSet(planes []string) map[string]bool {
	required := map[string]bool{}
	for _, plane := range planes {
		if plane != "" {
			required[plane] = true
		}
	}
	return required
}

func appendDigestValidation(packet Packet, runs RunSet, ledger Ledger, reasons, nextActions *[]string) bool {
	cannotVerify := false
	if runs.PacketDigest != packet.PacketDigest || ledger.PacketDigest != packet.PacketDigest {
		cannotVerify = true
		*reasons = append(*reasons, "packet_digest_mismatch")
		*nextActions = append(*nextActions, "Create a new packet and rerun review for the current head.")
	}
	for _, result := range runs.Results {
		if result.PacketDigest != packet.PacketDigest {
			cannotVerify = true
			*reasons = append(*reasons, "result_packet_digest_mismatch:"+safeID(result.ReviewRunID))
			*nextActions = append(*nextActions, "Discard stale reviewer results and rerun review for the current packet.")
		}
	}
	return cannotVerify
}

func validateRequiredPlanes(required map[string]bool, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) ([]PlaneResult, int, bool) {
	planeResults := make([]PlaneResult, 0, len(required))
	usableCount := 0
	cannotVerify := false
	for plane := range required {
		best := bestPlaneResult(plane, roleByID, runs)
		if best.Usable {
			usableCount++
		}
		if planeCannotVerify(best.Status) {
			cannotVerify = true
		}
		if best.Reason != "" {
			*reasons = append(*reasons, fmt.Sprintf("%s:%s", plane, best.Status))
		}
		if best.NextAction != "" && !best.Usable {
			*nextActions = append(*nextActions, best.NextAction)
		}
		planeResults = append(planeResults, best)
	}
	sort.Slice(planeResults, func(i, j int) bool { return planeResults[i].Plane < planeResults[j].Plane })
	return planeResults, usableCount, cannotVerify
}

func bestPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet) PlaneResult {
	best := PlaneResult{Plane: plane, Status: StateNotAssessed, Usable: false, Reason: "required_plane_not_assessed", NextAction: "Run or import a reviewer result for this plane."}
	for _, result := range runs.Results {
		if result.Plane != plane {
			continue
		}
		best = planeResult(result)
		if best.Usable && modelMismatchWithoutFallback(roleByID[result.RoleID], result) {
			best.Usable = false
			best.Status = StatusCannotVerify
			best.Reason = "model_identity_mismatch"
			best.NextAction = "Rerun the reviewer or record fallback provenance for the observed model."
		}
		break
	}
	return best
}

func planeCannotVerify(status string) bool {
	switch status {
	// StatusCannotVerify and CoverageCannotVerify currently share the same
	// wire value. Plane results carry reviewer statuses, so use the status
	// constant here and keep coverage-state selection in reviewCoverageState.
	case StatusCannotVerify, StatusTimedOut, StatusEmptyOutput, StatusOffTask, StatusParseFailed:
		return true
	default:
		return false
	}
}

func validateLedgerFindings(packet Packet, ledger Ledger, reasons *[]string) ([]LedgerFinding, bool, bool) {
	unresolved := false
	cannotVerify := false
	safeFindings := make([]LedgerFinding, 0, len(ledger.Findings))
	for _, finding := range ledger.Findings {
		finding.Summary = safeText(finding.Summary)
		if (finding.Severity == SeverityCritical || finding.Severity == SeverityMajor) && finding.Disposition == DispositionUnresolvedReviewBlocker {
			unresolved = true
		}
		if !citationResolvable(packet, finding.Citation) {
			cannotVerify = true
			*reasons = append(*reasons, "finding_citation_cannot_verify")
		}
		safeFindings = append(safeFindings, finding)
	}
	return safeFindings, unresolved, cannotVerify
}

func reviewCoverageState(required map[string]bool, usableCount int, cannotVerify, unresolved bool) string {
	switch {
	case cannotVerify:
		return CoverageCannotVerify
	case len(required) == 0 || usableCount == 0:
		return CoverageNotAssessed
	case usableCount < len(required):
		return CoveragePartial
	case unresolved:
		return CoverageUnresolved
	default:
		return CoverageSatisfied
	}
}

func modelMismatchWithoutFallback(role ReviewRole, result ReviewerResult) bool {
	requested := defaultString(role.RequestedModel, result.RequestedModel)
	observed := result.ObservedModel
	if requested == "" || requested == StateNotAssessed || observed == "" || observed == StateNotAssessed {
		return false
	}
	if requested == observed {
		return false
	}
	return result.FallbackForModel == "" || result.FallbackReason == ""
}

func Summarize(validation Validation, ledger Ledger) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Review coverage: %s\n", validation.ReviewCoverageState)
	fmt.Fprintf(&b, "CI state: %s\n", validation.CIState)
	fmt.Fprintf(&b, "Authority scope: %s\n", validation.AuthorityScope)
	fmt.Fprintf(&b, "Merge decision: %s\n", validation.MergeDecision)
	fmt.Fprintf(&b, "Release decision: %s\n", validation.ReleaseDecision)
	fmt.Fprintf(&b, "Risk acceptance: %s\n", validation.RiskAcceptance)
	b.WriteString("This is review-record evidence only; merge, release, and risk decisions remain external.\n")
	if len(validation.PlaneResults) > 0 {
		b.WriteString("\nPlanes\n")
		for _, plane := range validation.PlaneResults {
			fmt.Fprintf(&b, "- %s: %s", plane.Plane, plane.Status)
			if plane.NextAction != "" {
				fmt.Fprintf(&b, " next_action=%s", safeText(plane.NextAction))
			}
			b.WriteString("\n")
		}
	}
	if len(ledger.Findings) > 0 {
		b.WriteString("\nFindings\n")
		for _, finding := range ledger.Findings {
			fmt.Fprintf(&b, "- %s [%s] %s (%s)\n", finding.ID, finding.Severity, safeText(finding.Summary), finding.Disposition)
		}
	}
	return b.String()
}

func WriteJSON(path string, value any) error {
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
	var packet Packet
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		path = filepath.Join(path, "packet.json")
	}
	return packet, readJSON(path, &packet)
}

func ReadProfile(path string) (ReviewProfile, error) {
	var profile ReviewProfile
	if err := readJSON(path, &profile); err != nil {
		return profile, err
	}
	return profile, validateProfile(profile)
}

func ReadRunSet(path string) (RunSet, error) {
	var runs RunSet
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		path = filepath.Join(path, "results.json")
	}
	if err := readJSON(path, &runs); err != nil {
		return runs, err
	}
	if err := validateRunSet(runs); err != nil {
		return runs, err
	}
	return runs, nil
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
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func validatePacketOptions(opts PacketOptions) error {
	if strings.TrimSpace(opts.OutDir) == "" {
		return errors.New("pr_review_packet_requires_out")
	}
	if !repoIDPattern.MatchString(opts.RepoID) {
		return fmt.Errorf("unsafe_repo_id: repo_id must match %s", repoIDPattern.String())
	}
	if !changeRefPattern.MatchString(opts.ChangeRef) {
		return fmt.Errorf("unsafe_change_ref: change_ref must match %s", changeRefPattern.String())
	}
	if !sha40Pattern.MatchString(opts.BaseCommit) || !sha40Pattern.MatchString(opts.HeadCommit) {
		return errors.New("invalid_commit_sha: base and head must be 40 lowercase hex characters")
	}
	if strings.TrimSpace(opts.DiffPath) == "" {
		return errors.New("pr_review_packet_requires_diff")
	}
	if opts.CIState != "" && !validCIState(opts.CIState) {
		return fmt.Errorf("invalid_ci_state: %s", opts.CIState)
	}
	return nil
}

func validateRunSet(runs RunSet) error {
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
	if profile.SchemaVersion != "" && profile.SchemaVersion != SchemaVersionProfile {
		return fmt.Errorf("invalid_profile_schema_version: %s", profile.SchemaVersion)
	}
	return requireProfileFields(profile)
}

func requireProfileFields(profile ReviewProfile) error {
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
	for _, plane := range requiredPlanes {
		if !rolePlanes[plane] {
			return fmt.Errorf("profile_required_plane_without_role: %s", plane)
		}
	}
	return nil
}

func runRole(packet Packet, role ReviewRole, opts RunOptions, rawDir string) (ReviewerResult, error) {
	result := newReviewerResult(packet, role, opts.Now)
	result.CommandDigest = commandDigest(role.Command)
	promptRef, err := promptSafeRef(role)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "prompt_ref_cannot_verify"
		return result, nil
	}
	result.PromptRef = promptRef
	baseline, ready, err := prepareRoleRunner(&result, role, opts)
	if err != nil || !ready {
		return result, err
	}
	output, timedOut, err := runRoleCommand(role, opts)
	ended := time.Now().UTC().Format(time.RFC3339)
	result.EndedAt = ended
	if timedOut {
		result.Status = StatusTimedOut
		result.Reason = "runner_timed_out"
		return writeRawResult(result, rawDir, output)
	}
	if err := applyRunnerError(&result, err); err != nil {
		return writeRawResult(result, rawDir, output)
	}
	if len(strings.TrimSpace(string(output))) == 0 {
		result.Status = StatusEmptyOutput
		result.Reason = "runner_empty_output"
		return writeRawResult(result, rawDir, output)
	}
	result, err = parseReviewerOutput(result, role, packet, output)
	if err != nil {
		result.Status = StatusParseFailed
		result.Reason = "runner_output_parse_failed"
	}
	applyOpenCodeMutationCheck(&result, role, opts.WorkDir, baseline)
	return writeRawResult(result, rawDir, output)
}

func newReviewerResult(packet Packet, role ReviewRole, now time.Time) ReviewerResult {
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
	if role.Runner != RunnerManualExternal && !opts.AllowedRunners[role.Runner] {
		return nil, false, fmt.Errorf("runner_not_allowed: %s", role.Runner)
	}
	if role.Runner == RunnerOpenCode && !role.ReadOnlyEnforced {
		// Keep the default not_assessed status: this is a safety preflight
		// refusal, not a runner execution failure.
		result.Reason = "opencode_read_only_not_enforced"
		return nil, false, nil
	}
	if role.Runner == RunnerOpenCode {
		return prepareOpenCodeBaseline(result, role, opts.WorkDir)
	}
	if len(role.Command) == 0 {
		result.Reason = "runner_command_not_configured"
		return nil, false, nil
	}
	return nil, true, nil
}

func prepareOpenCodeBaseline(result *ReviewerResult, role ReviewRole, workDir string) (*workingTreeBaseline, bool, error) {
	mode := defaultString(role.WorkingTreeMode, "clean_required")
	baseline, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "working_tree_baseline_cannot_verify"
		return nil, false, nil
	}
	if mode == "clean_required" && baseline.Count > 0 {
		result.Status = StatusNotAssessed
		result.Reason = "working_tree_dirty"
		return nil, false, nil
	}
	if len(role.Command) == 0 {
		result.Reason = "runner_command_not_configured"
		return nil, false, nil
	}
	return baseline, true, nil
}

func runRoleCommand(role ReviewRole, opts RunOptions) ([]byte, bool, error) {
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
	if role.Runner != RunnerOpenCode || baseline == nil {
		return
	}
	after, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "working_tree_baseline_cannot_verify"
		return
	}
	if after.Digest != baseline.Digest || after.Count != baseline.Count {
		result.Status = StatusCannotVerify
		result.Reason = "mutation_detected"
	}
}

type workingTreeBaseline struct {
	Count  int
	Digest string
}

func captureWorkingTreeBaseline(workDir string) (*workingTreeBaseline, error) {
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
	var parsed ReviewerResult
	// RequiredOutputSchema identifies the declared schema contract; this parser
	// enforces the concrete Go contract with unknown-field rejection.
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&parsed); err != nil {
		return base, err
	}
	if parsed.PacketDigest != packet.PacketDigest || parsed.Plane != role.Plane || parsed.RoleID != role.RoleID {
		base.Status = StatusOffTask
		base.Reason = "reviewer_output_wrong_packet_plane_or_role"
		return base, nil
	}
	parsed.ReviewRunID = defaultString(parsed.ReviewRunID, base.ReviewRunID)
	parsed.Runner = defaultString(parsed.Runner, role.Runner)
	parsed.RequestedModel = defaultString(parsed.RequestedModel, defaultString(role.RequestedModel, StateNotAssessed))
	parsed.ObservedModel = defaultString(parsed.ObservedModel, StateNotAssessed)
	parsed.ModelFamily = defaultString(parsed.ModelFamily, StateNotAssessed)
	parsed.ModelVersion = defaultString(parsed.ModelVersion, StateNotAssessed)
	parsed.StartedAt = base.StartedAt
	parsed.EndedAt = base.EndedAt
	parsed.CommandDigest = base.CommandDigest
	parsed.PromptRef = base.PromptRef
	if parsed.Status == "" {
		if len(parsed.Findings) > 0 {
			parsed.Status = StatusFindingsReported
		} else {
			parsed.Status = StatusNoFindings
		}
	}
	return parsed, nil
}

func writeRawResult(result ReviewerResult, rawDir string, output []byte) (ReviewerResult, error) {
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
	roles := make([]PreviewRole, 0, len(profile.Roles))
	for _, role := range profile.Roles {
		promptRef, _ := promptSafeRef(role)
		promptDigest := ""
		if promptRef != nil {
			promptDigest = promptRef.DigestSHA256
		}
		roles = append(roles, PreviewRole{
			RoleID:         role.RoleID,
			Plane:          role.Plane,
			Runner:         role.Runner,
			RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
			TimeoutSeconds: role.TimeoutSeconds,
			CommandDigest:  commandDigest(role.Command),
			PromptRef:      role.PromptTemplateRef,
			PromptDigest:   promptDigest,
		})
	}
	return &RunPreview{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Roles: roles}
}

func promptSafeRef(role ReviewRole) (*SafeRef, error) {
	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return nil, nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(data)
	return &SafeRef{
		ID:             "prompt-" + safeID(role.RoleID),
		Kind:           RefKindPrompt,
		Ref:            "digest-only:" + safeID(filepath.Base(role.PromptTemplateRef)),
		DigestSHA256:   hex.EncodeToString(sum[:]),
		ContentType:    contentType(role.PromptTemplateRef),
		RedactionState: RedactionDigestOnly,
	}, nil
}

func copyInputs(inputDir, prefix string, paths []string) ([]SafeRef, error) {
	refs := make([]SafeRef, 0, len(paths))
	for i, path := range paths {
		name := fmt.Sprintf("%s-%d%s", prefix, i+1, normalizedExt(path))
		ref, err := copyInput(inputDir, name, path, RefKindDoc, contentType(path))
		if err != nil {
			return nil, err
		}
		ref.ID = fmt.Sprintf("%s-%d", prefix, i+1)
		refs = append(refs, ref)
	}
	return refs, nil
}

func copyInput(inputDir, name, source, kind, contentType string) (SafeRef, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return SafeRef{}, err
	}
	dest := filepath.Join(inputDir, name)
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return SafeRef{}, err
	}
	digest := sha256.Sum256(data)
	return SafeRef{
		ID:             strings.TrimSuffix(name, filepath.Ext(name)),
		Kind:           kind,
		Ref:            filepath.ToSlash(filepath.Join("inputs", name)),
		DigestSHA256:   hex.EncodeToString(digest[:]),
		ContentType:    contentType,
		RedactionState: RedactionNone,
	}, nil
}

func digestJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func packetDigest(packet Packet) (string, error) {
	canonical := packet
	canonical.PacketDigest = ""
	return digestJSON(canonical)
}

func ensureNewDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("missing_output_path")
	}
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		return fmt.Errorf("output_exists: %s", filepath.Base(path))
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func validCIState(state string) bool {
	switch state {
	case StatePass, StateFail, StatePending, StateNotAssessed, StateCannotVerify:
		return true
	default:
		return false
	}
}

func validRunner(runner string) bool {
	switch runner {
	case RunnerPI, RunnerOpenCode, RunnerManualExternal:
		return true
	default:
		return false
	}
}

func planeResult(result ReviewerResult) PlaneResult {
	pr := PlaneResult{Plane: result.Plane, Status: result.Status, RunID: result.ReviewRunID}
	switch result.Status {
	case StatusFindingsReported, StatusNoFindings:
		pr.Usable = true
	case StatusNotAssessed:
		pr.Reason = "reviewer_not_assessed"
		pr.NextAction = "Run a configured reviewer or import a usable result for this plane."
	case StatusTimedOut:
		pr.Status = StatusTimedOut
		pr.Reason = "reviewer_timed_out"
		pr.NextAction = "Increase timeout or replace the reviewer for this plane."
	case StatusEmptyOutput:
		pr.Reason = "reviewer_empty_output"
		pr.NextAction = "Retry with a shorter bounded prompt or replace the reviewer."
	case StatusOffTask:
		pr.Reason = "reviewer_off_task"
		pr.NextAction = "Rerun with the frozen packet and required output schema."
	case StatusParseFailed:
		pr.Reason = "reviewer_parse_failed"
		pr.NextAction = "Rerun with JSON-only output matching the required schema."
	default:
		pr.Reason = "reviewer_cannot_verify"
		pr.NextAction = "Replace or rerun the reviewer."
	}
	return pr
}

func defaultDisposition(severity string) string {
	switch safeSeverity(severity) {
	case SeverityCritical, SeverityMajor:
		return DispositionUnresolvedReviewBlocker
	default:
		return DispositionDeferredNotAssessed
	}
}

func safeSeverity(severity string) string {
	switch severity {
	case SeverityCritical, SeverityMajor, SeverityMinor, SeverityInformational:
		return severity
	default:
		return SeverityInformational
	}
}

func citationResolvable(packet Packet, citation Citation) bool {
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
	if text == "" {
		return ""
	}
	unsafeMarkers := []string{"SYNTHETIC_", "Bearer ", "access_token=", "BEGIN PRIVATE KEY", "PRIVATE_KEY", "cookie=", "session=", "/Users/", "/private/"}
	for _, marker := range unsafeMarkers {
		if strings.Contains(text, marker) {
			return "[redacted unsafe reviewer text]"
		}
	}
	if (strings.Contains(text, "://") && strings.Contains(text, "@")) || strings.Contains(text, "token=") {
		return "[redacted unsafe reviewer text]"
	}
	return text
}

func uniqueStrings(values []string) []string {
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
	if len(command) == 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(strings.Join(command, "\x00")))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func contextKind(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return RefKindSchema
	case ".md", ".markdown":
		if strings.Contains(strings.ToLower(filepath.Base(path)), "task") {
			return RefKindTask
		}
		return RefKindDoc
	default:
		return RefKindSourceExcerpt
	}
}

func contentType(path string) string {
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
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".md", ".txt", ".diff", ".patch":
		return ext
	default:
		return ".txt"
	}
}

func safeID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	out := strings.Trim(strings.Map(safeIDMapper, value), "-.")
	if out == "" {
		return "item"
	}
	return out
}

const safeIDAllowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_.-"

func safeIDMapper(r rune) rune {
	if r <= 127 && strings.IndexByte(safeIDAllowedChars, byte(r)) >= 0 {
		return r
	}
	return '-'
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func Copy(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}
