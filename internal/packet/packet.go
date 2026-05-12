package packet

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	PacketSchemaVersion = "change-evidence-packet.v0"
	BundleSchemaVersion = "evidence-bundle-manifest.v0"

	StatePass         = "pass"
	StatePartial      = "partial"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"
	StateNotInScope   = "not_in_scope"

	AuthoringToolGenerated = "tool_generated"
	ProjectionCanonical    = "canonical_markdown_artifact"
)

var RequiredRows = []string{
	"PC-CHANGE",
	"PC-INITIATOR",
	"PC-AGENT-ROUTE",
	"PC-MUTATION",
	"PC-VERIFICATION",
	"PC-REVIEW",
	"PC-AUTHORITY",
	"PC-THEATER",
	"PC-ATTESTATION",
	"PC-DECISION",
	"PC-RESIDUAL-GAPS",
}

var states = map[string]bool{
	StatePass:         true,
	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}

var missingReasonStates = map[string]bool{
	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}

var packetStates = map[string]bool{
	"draft":        true,
	"review_ready": true,
	"reviewed":     true,
	"superseded":   true,
}

var authoringMethods = map[string]bool{
	AuthoringToolGenerated:         true,
	"hand_authored_before_tooling": true,
}

var retainedForms = map[string]bool{
	"raw":          true,
	"redacted":     true,
	"digest_only":  true,
	"external_ref": true,
	"not_retained": true,
}

var redactionStatuses = map[string]bool{
	"not_needed":      true,
	"redacted":        true,
	"digest_only":     true,
	"withheld":        true,
	StateCannotVerify: true,
}

var theaterReasonCodes = map[string]bool{
	"agent_claimed_verification": true,
	"unbound_intent":             true,
	"ci_theater":                 true,
	"scope_theater":              true,
	"prompt_contamination":       true,
}

var requiredDecisions = []string{"merge", "release", "risk_acceptance", "security_review"}

type Packet struct {
	PacketVersion   string           `json:"packet_version"`
	PacketID        string           `json:"packet_id"`
	SourceChange    SourceChange     `json:"source_change"`
	GeneratedAt     string           `json:"generated_at"`
	AuthoringMethod string           `json:"authoring_method"`
	SelectedProfile string           `json:"selected_profile"`
	RedactionPolicy string           `json:"redaction_policy"`
	BundleRef       string           `json:"bundle_ref"`
	PacketState     string           `json:"packet_state"`
	Projection      Projection       `json:"projection"`
	Rows            []Row            `json:"rows"`
	TheaterFindings []TheaterFinding `json:"theater_findings,omitempty"`
	ResidualGaps    []ResidualGap    `json:"residual_gaps,omitempty"`
	DecisionOwners  []DecisionOwner  `json:"decision_owners,omitempty"`
	NonApproval     string           `json:"non_approval"`
	Extensions      map[string]any   `json:"extensions,omitempty"`
}

type SourceChange struct {
	Repository  string `json:"repository,omitempty"`
	ChangeID    string `json:"change_id,omitempty"`
	URL         string `json:"url,omitempty"`
	BaseRef     string `json:"base_ref,omitempty"`
	HeadRef     string `json:"head_ref,omitempty"`
	CommitRange string `json:"commit_range,omitempty"`
	HeadSHA     string `json:"head_sha,omitempty"`
}

type Projection struct {
	Kind        string `json:"kind"`
	Canonical   bool   `json:"canonical"`
	ArtifactRef string `json:"artifact_ref,omitempty"`
}

type Row struct {
	ID           string   `json:"id"`
	State        string   `json:"state"`
	Summary      string   `json:"summary"`
	EvidenceRefs []string `json:"evidence_refs"`
	Reason       string   `json:"reason,omitempty"`
	Owner        string   `json:"owner"`
}

type TheaterFinding struct {
	ReasonCode              string   `json:"reason_code"`
	State                   string   `json:"state"`
	Severity                string   `json:"severity,omitempty"`
	Finding                 string   `json:"finding"`
	TriggerEvidenceRefs     []string `json:"trigger_evidence_refs"`
	RequiredClosureEvidence string   `json:"required_closure_evidence,omitempty"`
}

type ResidualGap struct {
	RowID           string   `json:"row_id"`
	State           string   `json:"state"`
	Reason          string   `json:"reason"`
	EvidenceRefs    []string `json:"evidence_refs,omitempty"`
	ClosureEvidence string   `json:"closure_evidence,omitempty"`
}

type DecisionOwner struct {
	Decision string `json:"decision"`
	Owner    string `json:"owner"`
	State    string `json:"state"`
	Reason   string `json:"reason,omitempty"`
}

type BundleManifest struct {
	SchemaVersion string          `json:"schema_version"`
	BundleID      string          `json:"bundle_id"`
	PacketDigest  string          `json:"packet_digest,omitempty"`
	Entries       []BundleEntry   `json:"entries"`
	Resolvers     []ResolverEntry `json:"resolvers,omitempty"`
}

type BundleEntry struct {
	Ref                string   `json:"ref"`
	SourceClass        string   `json:"source_class"`
	Digest             string   `json:"digest,omitempty"`
	RetainedForm       string   `json:"retained_form"`
	RedactionStatus    string   `json:"redaction_status"`
	Resolver           string   `json:"resolver,omitempty"`
	ExpiresAt          string   `json:"expires_at,omitempty"`
	ArtifactAccess     string   `json:"artifact_access,omitempty"`
	ProjectionRole     string   `json:"projection_role,omitempty"`
	EvidenceKind       string   `json:"evidence_kind,omitempty"`
	ObservedComponents []string `json:"observed_components,omitempty"`
	ContradictsRef     string   `json:"contradicts_ref,omitempty"`
	ContradictsRowID   string   `json:"contradicts_row_id,omitempty"`
	Actor              string   `json:"actor,omitempty"`
	WriteAuthority     string   `json:"write_authority,omitempty"`
	GeneratedBy        string   `json:"generated_by,omitempty"`
	SourceCommitState  string   `json:"source_commit_state,omitempty"`
	SourceRef          string   `json:"source_ref,omitempty"`
}

type ResolverEntry struct {
	Ref      string `json:"ref"`
	Resolver string `json:"resolver"`
}

type Bundle struct {
	Packet   Packet         `json:"packet"`
	Manifest BundleManifest `json:"manifest"`
}

type GitHubPREvidenceInput struct {
	SchemaVersion          string              `json:"schema_version"`
	PR                     GitHubPR            `json:"pr"`
	CommitRange            GitHubCommitRange   `json:"commit_range"`
	Checks                 []GitHubCheck       `json:"checks,omitempty"`
	Artifacts              []GitHubArtifact    `json:"artifacts,omitempty"`
	Reviews                []GitHubReview      `json:"reviews,omitempty"`
	WorkflowRunID          string              `json:"workflow_run_id,omitempty"`
	RequirePromptBoundary  bool                `json:"require_prompt_boundary,omitempty"`
	AgentRouteRefs         []string            `json:"agent_route_refs,omitempty"`
	AgentRouteComponents   []string            `json:"agent_route_components,omitempty"`
	AgentRouteDigest       string              `json:"agent_route_digest,omitempty"`
	AgentRouteEvidenceKind string              `json:"agent_route_evidence_kind,omitempty"`
	PromptBoundary         PromptBoundary      `json:"prompt_boundary,omitempty"`
	IntegrationActions     []IntegrationAction `json:"integration_actions,omitempty"`
}

type PromptBoundary struct {
	Text          string `json:"text,omitempty"`
	Digest        string `json:"digest,omitempty"`
	CaptureActor  string `json:"capture_actor,omitempty"`
	CapturedAt    string `json:"captured_at,omitempty"`
	CaptureMethod string `json:"capture_method,omitempty"`
}

type PromptBoundaryClassification struct {
	Verdict          string   `json:"verdict"`
	RouteProofEffect string   `json:"route_proof_effect"`
	Reasons          []string `json:"reasons,omitempty"`
}

type IntegrationAction struct {
	Kind     string `json:"kind"`
	Actor    string `json:"actor"`
	Resolver string `json:"resolver"`
}

type BuildPRResult struct {
	State      string   `json:"state"`
	BundlePath string   `json:"bundle_path,omitempty"`
	PacketPath string   `json:"packet_path,omitempty"`
	ResultPath string   `json:"result_path,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

type GitHubPR struct {
	Number  int    `json:"number"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	BodyRef string `json:"body_ref,omitempty"`
	Author  string `json:"author"`
	BaseRef string `json:"base_ref"`
	HeadRef string `json:"head_ref"`
	HeadSHA string `json:"head_sha"`
}

type GitHubCommitRange struct {
	Base            string `json:"base"`
	Head            string `json:"head"`
	ChangedFilesRef string `json:"changed_files_ref,omitempty"`
}

type GitHubCheck struct {
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Conclusion   string   `json:"conclusion"`
	ArtifactRefs []string `json:"artifact_refs,omitempty"`
}

type GitHubArtifact struct {
	Name         string `json:"name"`
	Resolver     string `json:"resolver"`
	RetainedForm string `json:"retained_form"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	Digest       string `json:"digest,omitempty"`
}

type GitHubReview struct {
	Reviewer string `json:"reviewer"`
	Resolver string `json:"resolver"`
	State    string `json:"state"`
}

type Validation struct {
	State  string   `json:"state"`
	Errors []string `json:"errors,omitempty"`
}

func LoadBundle(path string) (Bundle, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Bundle{}, err
	}
	var bundle Bundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return Bundle{}, err
	}
	return bundle, nil
}

func LoadGitHubInput(path string) (GitHubPREvidenceInput, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return GitHubPREvidenceInput{}, err
	}
	var input GitHubPREvidenceInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return GitHubPREvidenceInput{}, err
	}
	return input, nil
}

func BuildFromGitHubInput(input GitHubPREvidenceInput, generatedAt time.Time) Bundle {
	packetID := fmt.Sprintf("github-pr-%d-change-evidence-packet", input.PR.Number)
	bundleID := fmt.Sprintf("%s-bundle", packetID)
	entries := githubEntries(input)
	rows := githubRows(input)
	packet := Packet{
		PacketVersion:   PacketSchemaVersion,
		PacketID:        packetID,
		SourceChange:    githubSourceChange(input),
		GeneratedAt:     generatedAt.UTC().Format(time.RFC3339),
		AuthoringMethod: AuthoringToolGenerated,
		SelectedProfile: "change-host-rich-v0",
		RedactionPolicy: "not_assessed",
		BundleRef:       bundleID,
		PacketState:     "draft",
		Projection:      Projection{Kind: ProjectionCanonical, Canonical: true, ArtifactRef: "packet:markdown"},
		Rows:            rows,
		ResidualGaps:    residualGapsForRows(rows),
		DecisionOwners:  defaultDecisionOwners(),
		NonApproval:     "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.",
	}
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if classification.Verdict == "contaminated" {
		packet.TheaterFindings = append(packet.TheaterFindings, TheaterFinding{
			ReasonCode:          "prompt_contamination",
			State:               StateFail,
			Severity:            "P0",
			Finding:             strings.Join(classification.Reasons, "; "),
			TriggerEvidenceRefs: []string{"prompt:boundary"},
		})
	}
	if len(input.IntegrationActions) > 0 {
		if packet.Extensions == nil {
			packet.Extensions = map[string]any{}
		}
		packet.Extensions["integration_actions"] = input.IntegrationActions
	}
	return Bundle{
		Packet: packet,
		Manifest: BundleManifest{
			SchemaVersion: BundleSchemaVersion,
			BundleID:      bundleID,
			PacketDigest:  PacketDigest(packet),
			Entries:       entries,
		},
	}
}

func ClassifyPromptBoundary(boundary PromptBoundary) PromptBoundaryClassification {
	text := strings.TrimSpace(boundary.Text)
	if text != "" {
		lower := strings.ToLower(text)
		for _, phrase := range forbiddenRecorderDutyPhrases() {
			if strings.Contains(lower, phrase) {
				return PromptBoundaryClassification{
					Verdict:          "contaminated",
					RouteProofEffect: StateFail,
					Reasons:          []string{"developer prompt contains recorder-duty phrase: " + phrase},
				}
			}
		}
		return PromptBoundaryClassification{Verdict: "clean", RouteProofEffect: StatePass}
	}
	if strings.TrimSpace(boundary.Digest) == "" &&
		strings.TrimSpace(boundary.CaptureActor) == "" &&
		strings.TrimSpace(boundary.CapturedAt) == "" &&
		strings.TrimSpace(boundary.CaptureMethod) == "" {
		return PromptBoundaryClassification{Verdict: "missing", RouteProofEffect: StateCannotVerify, Reasons: []string{"prompt boundary evidence missing"}}
	}
	if strings.TrimSpace(boundary.Digest) != "" &&
		strings.TrimSpace(boundary.CaptureActor) != "" &&
		strings.TrimSpace(boundary.CapturedAt) != "" &&
		strings.TrimSpace(boundary.CaptureMethod) != "" {
		if _, err := time.Parse(time.RFC3339, boundary.CapturedAt); err == nil {
			return PromptBoundaryClassification{Verdict: "digest_only", RouteProofEffect: StatePartial, Reasons: []string{"prompt text unavailable; digest metadata retained"}}
		}
	}
	return PromptBoundaryClassification{Verdict: "malformed", RouteProofEffect: StateCannotVerify, Reasons: []string{"prompt boundary metadata malformed"}}
}

func forbiddenRecorderDutyPhrases() []string {
	return []string{
		"sdp-trace",
		".sdp-trace",
		".evidence",
		"write evidence",
		"update evidence",
		"maintain provenance",
		"update provenance",
		"update packet",
		"update bundle",
		"close gate",
		"claim verification",
	}
}

func PacketDigest(packet Packet) string {
	clone := packet
	raw, err := json.Marshal(clone)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}

func Validate(bundle Bundle, now time.Time) Validation {
	validator := bundleValidator{
		bundle:        bundle,
		now:           now,
		entryByRef:    map[string]BundleEntry{},
		resolverByRef: map[string]string{},
	}
	return validator.validate()
}

func CheckDemoFirstPacket(bundle Bundle, now time.Time) Validation {
	validation := Validate(bundle, now)
	check := demoFirstPacketChecker{
		bundle:     bundle,
		now:        now,
		rows:       map[string]Row{},
		entryByRef: map[string]BundleEntry{},
		errors:     append([]string(nil), validation.Errors...),
	}
	return check.validate()
}

type demoFirstPacketChecker struct {
	bundle     Bundle
	now        time.Time
	rows       map[string]Row
	entryByRef map[string]BundleEntry
	errors     []string
}

func (c *demoFirstPacketChecker) validate() Validation {
	c.index()
	c.requireToolGenerated()
	c.requirePassOrPartialRows(4)
	c.requireRowEvidence("PC-CHANGE")
	c.requireRowEvidence("PC-MUTATION")
	c.requireAgentRouteEvidence()
	c.requireVerificationOrReviewAssessed()
	c.requireCannotVerifyClosureCap()
	state := StatePass
	if len(c.errors) > 0 {
		state = StateFail
	}
	return Validation{State: state, Errors: c.errors}
}

func (c *demoFirstPacketChecker) requireToolGenerated() {
	if c.bundle.Packet.AuthoringMethod != AuthoringToolGenerated {
		c.add("demo first-packet gate requires tool_generated authoring_method, got %s", c.bundle.Packet.AuthoringMethod)
	}
}

func (c *demoFirstPacketChecker) index() {
	for _, row := range c.bundle.Packet.Rows {
		c.rows[row.ID] = row
	}
	for _, entry := range c.bundle.Manifest.Entries {
		c.entryByRef[entry.Ref] = entry
	}
}

func (c *demoFirstPacketChecker) requirePassOrPartialRows(minimum int) {
	count := 0
	for _, row := range c.rows {
		if row.State == StatePass || row.State == StatePartial {
			count++
		}
	}
	if count < minimum {
		c.add("demo first-packet gate requires at least %d pass or partial rows, got %d", minimum, count)
	}
}

func (c *demoFirstPacketChecker) requireRowEvidence(rowID string) {
	row := c.rows[rowID]
	if len(row.EvidenceRefs) == 0 {
		c.add("demo first-packet gate requires %s retained evidence refs", rowID)
		return
	}
	for _, ref := range row.EvidenceRefs {
		entry, ok := c.entryByRef[ref]
		if !ok {
			continue
		}
		if !demoUsableEntry(entry, c.now) {
			c.add("demo first-packet gate requires %s evidence ref %q to be retained and usable", rowID, ref)
		}
	}
}

func (c *demoFirstPacketChecker) requireAgentRouteEvidence() {
	row := c.rows["PC-AGENT-ROUTE"]
	if row.State != StatePass && row.State != StatePartial {
		c.add("demo first-packet gate requires PC-AGENT-ROUTE must be pass or partial, got %s", row.State)
		return
	}
	if len(row.EvidenceRefs) == 0 {
		c.add("demo first-packet gate requires PC-AGENT-ROUTE retained evidence refs")
		return
	}
	for _, ref := range row.EvidenceRefs {
		entry := c.entryByRef[ref]
		if entry.SourceClass == "harness" && demoUsableEntry(entry, c.now) && demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry) {
			return
		}
	}
	c.add("demo first-packet gate requires PC-AGENT-ROUTE evidence from retained structured OpenCode/GSD/MiniMax harness route observation")
}

func (c *demoFirstPacketChecker) requireVerificationOrReviewAssessed() {
	if rowAssessed(c.rows["PC-VERIFICATION"]) || rowAssessed(c.rows["PC-REVIEW"]) {
		return
	}
	c.add("demo first-packet gate requires PC-VERIFICATION or PC-REVIEW to be pass, partial, or fail")
}

func rowAssessed(row Row) bool {
	return row.State == StatePass || row.State == StatePartial || row.State == StateFail
}

func (c *demoFirstPacketChecker) requireCannotVerifyClosureCap() {
	unclosed := 0
	for _, row := range c.rows {
		if row.State != StateCannotVerify {
			continue
		}
		if !gapForRowWithClosure(c.bundle.Packet.ResidualGaps, row.ID) {
			unclosed++
		}
	}
	if unclosed > 1 {
		c.add("demo first-packet gate allows at most one cannot_verify row without closure path, got %d", unclosed)
	}
}

func gapForRowWithClosure(gaps []ResidualGap, rowID string) bool {
	for _, gap := range gaps {
		if gap.RowID == rowID && strings.TrimSpace(gap.ClosureEvidence) != "" {
			return true
		}
	}
	return false
}

func (c *demoFirstPacketChecker) add(format string, args ...any) {
	c.errors = append(c.errors, fmt.Sprintf(format, args...))
}

func demoUsableEntry(entry BundleEntry, now time.Time) bool {
	if strings.TrimSpace(entry.Resolver) == "" || strings.TrimSpace(entry.Digest) == "" {
		return false
	}
	if entryExpired(entry, now) {
		return false
	}
	if passRefUnverifiable(entry) {
		return false
	}
	return true
}

func demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry BundleEntry) bool {
	if entry.EvidenceKind != "harness_route_observation" {
		return false
	}
	if syntheticEntryDigest(entry) {
		return false
	}
	components := map[string]bool{}
	for _, component := range entry.ObservedComponents {
		components[strings.ToLower(strings.TrimSpace(component))] = true
	}
	return components["opencode"] && components["gsd"] &&
		(components["minimax"] || components["minimax-m2.5"] || components["minimax-m2"])
}

func syntheticEntryDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Digest) == "" || entry.Digest == digestPlaceholder(entry.Ref+entry.Resolver)
}

type bundleValidator struct {
	bundle        Bundle
	now           time.Time
	entryByRef    map[string]BundleEntry
	resolverByRef map[string]string
	rows          map[string]Row
	errors        []string
}

func (v *bundleValidator) validate() Validation {
	v.validateMetadata()
	v.indexManifest()
	v.validateRows()
	v.validateFindingsAndGaps()
	state := StatePass
	if len(v.errors) > 0 {
		state = StateFail
	}
	return Validation{State: state, Errors: v.errors}
}

func (v *bundleValidator) validateMetadata() {
	if v.bundle.Packet.PacketVersion != PacketSchemaVersion {
		v.add("packet.packet_version must be %q", PacketSchemaVersion)
	}
	if v.bundle.Manifest.SchemaVersion != BundleSchemaVersion {
		v.add("manifest.schema_version must be %q", BundleSchemaVersion)
	}
	if strings.TrimSpace(v.bundle.Packet.PacketID) == "" {
		v.add("packet.packet_id is required")
	}
	if strings.TrimSpace(v.bundle.Packet.BundleRef) == "" {
		v.add("packet.bundle_ref is required")
	}
	if strings.TrimSpace(v.bundle.Manifest.BundleID) == "" {
		v.add("manifest.bundle_id is required")
	}
	if v.bundle.Packet.BundleRef != "" && v.bundle.Manifest.BundleID != "" && v.bundle.Packet.BundleRef != v.bundle.Manifest.BundleID {
		v.add("packet.bundle_ref %q must match manifest.bundle_id %q", v.bundle.Packet.BundleRef, v.bundle.Manifest.BundleID)
	}
	if strings.TrimSpace(v.bundle.Manifest.PacketDigest) == "" {
		v.add("manifest.packet_digest is required")
	} else if digest := PacketDigest(v.bundle.Packet); digest != "" && v.bundle.Manifest.PacketDigest != digest {
		v.add("manifest.packet_digest does not match packet content")
	}
	if strings.TrimSpace(v.bundle.Packet.NonApproval) == "" {
		v.add("packet.non_approval is required")
	}
	if !packetStates[v.bundle.Packet.PacketState] {
		v.add("packet.packet_state has unknown value %q", v.bundle.Packet.PacketState)
	}
	if !authoringMethods[v.bundle.Packet.AuthoringMethod] {
		v.add("packet.authoring_method has unknown value %q", v.bundle.Packet.AuthoringMethod)
	}
	if v.bundle.Packet.Projection.Canonical && v.bundle.Packet.Projection.Kind != ProjectionCanonical {
		v.add("canonical projection must be %q", ProjectionCanonical)
	}
	if !v.bundle.Packet.Projection.Canonical && strings.TrimSpace(v.bundle.Packet.Projection.ArtifactRef) == "" {
		v.add("non-canonical packet projection requires artifact_ref for canonical uploaded packet")
	}
}

func (v *bundleValidator) indexManifest() {
	for _, entry := range v.bundle.Manifest.Entries {
		if strings.TrimSpace(entry.Ref) == "" {
			v.add("manifest entry has empty ref")
			continue
		}
		if !retainedForms[entry.RetainedForm] {
			v.add("manifest entry %q has unknown retained_form %q", entry.Ref, entry.RetainedForm)
		}
		if !redactionStatuses[entry.RedactionStatus] {
			v.add("manifest entry %q has unknown redaction_status %q", entry.Ref, entry.RedactionStatus)
		}
		v.entryByRef[entry.Ref] = entry
		if strings.TrimSpace(entry.Resolver) != "" {
			v.resolverByRef[entry.Ref] = entry.Resolver
		}
	}
	for _, resolver := range v.bundle.Manifest.Resolvers {
		if strings.TrimSpace(resolver.Ref) != "" {
			v.resolverByRef[resolver.Ref] = resolver.Resolver
		}
	}
}

func (v *bundleValidator) validateRows() {
	rows := map[string]Row{}
	for _, row := range v.bundle.Packet.Rows {
		if !requiredRow(row.ID) {
			v.add("unknown row id %q", row.ID)
			continue
		}
		if rows[row.ID].ID != "" {
			v.add("duplicate row id %q", row.ID)
		}
		rows[row.ID] = row
		v.validateRow(row)
	}
	for _, id := range RequiredRows {
		if rows[id].ID == "" {
			v.add("missing required row %q", id)
		}
	}
	v.rows = rows
	v.validateContradictions(rows)
	v.validateResidualCoverage(rows)
}

func (v *bundleValidator) validateRow(row Row) {
	if !states[row.State] {
		v.add("%s has unknown state %q", row.ID, row.State)
	}
	if strings.TrimSpace(row.Summary) == "" {
		v.add("%s requires summary", row.ID)
	}
	if strings.TrimSpace(row.Owner) == "" {
		v.add("%s requires owner", row.ID)
	}
	if missingReasonStates[row.State] && strings.TrimSpace(row.Reason) == "" {
		v.add("%s state %s requires reason", row.ID, row.State)
	}
	if row.State == StatePass && len(row.EvidenceRefs) == 0 {
		v.add("%s pass requires retained evidence refs", row.ID)
	}
	for _, ref := range row.EvidenceRefs {
		v.validateEvidenceRef(row.ID, row.State, ref)
	}
}

func (v *bundleValidator) validateEvidenceRef(rowID, state, ref string) {
	entry, ok := v.entryByRef[ref]
	if !ok {
		v.add("%s evidence ref %q is absent from manifest", rowID, ref)
		return
	}
	if strings.TrimSpace(v.resolverByRef[ref]) == "" {
		v.add("%s evidence ref %q has no resolver entry", rowID, ref)
	}
	if state == StatePass && entryExpired(entry, v.now) {
		v.add("%s pass cites expired artifact ref %q", rowID, ref)
	}
	if state == StatePass && passRefUnverifiable(entry) {
		v.add("%s pass cites unverifiable artifact ref %q", rowID, ref)
	}
}

func (v *bundleValidator) validateContradictions(rows map[string]Row) {
	for _, entry := range v.entryByRef {
		rowID := entry.ContradictsRowID
		if rowID == "" {
			rowID = rowIDForRef(rows, entry.Ref)
		}
		if entry.ContradictsRef == "" || rowID == "" {
			continue
		}
		row := rows[rowID]
		if row.State != StatePartial {
			v.add("%s has contradictory evidence but state is %s, want partial", rowID, row.State)
		}
		if !gapForRow(v.bundle.Packet.ResidualGaps, rowID) {
			v.add("%s contradictory evidence requires residual gap explanation", rowID)
		}
	}
}

func (v *bundleValidator) validateFindingsAndGaps() {
	v.validateTheaterState()
	v.validateDecisionOwners()
	for _, finding := range v.bundle.Packet.TheaterFindings {
		if strings.TrimSpace(finding.ReasonCode) == "" {
			v.add("theater finding requires reason_code")
		} else if !theaterReasonCodes[finding.ReasonCode] {
			v.add("theater finding has unknown reason_code %q", finding.ReasonCode)
		}
		for _, ref := range finding.TriggerEvidenceRefs {
			v.validateEvidenceRef("theater finding "+finding.ReasonCode, StatePartial, ref)
		}
	}
	for _, gap := range v.bundle.Packet.ResidualGaps {
		if !requiredRow(gap.RowID) {
			v.add("residual gap has unknown row id %q", gap.RowID)
		}
		if strings.TrimSpace(gap.Reason) == "" {
			v.add("residual gap for %s requires reason", gap.RowID)
		}
	}
}

func (v *bundleValidator) validateResidualCoverage(rows map[string]Row) {
	for _, row := range rows {
		if row.ID == "PC-RESIDUAL-GAPS" || row.State == StatePass {
			continue
		}
		if !gapForRow(v.bundle.Packet.ResidualGaps, row.ID) {
			v.add("%s non-pass row requires residual gap explanation", row.ID)
		}
	}
}

func (v *bundleValidator) validateTheaterState() {
	row := v.rows["PC-THEATER"]
	if len(v.bundle.Packet.TheaterFindings) == 0 {
		return
	}
	if row.State == StatePass {
		v.add("PC-THEATER cannot be pass when theater findings are present")
	}
	if row.State != StatePartial && row.State != StateFail && row.State != StateCannotVerify {
		v.add("PC-THEATER with theater findings must be partial, fail, or cannot_verify")
	}
}

func (v *bundleValidator) validateDecisionOwners() {
	owners := map[string]DecisionOwner{}
	for _, owner := range v.bundle.Packet.DecisionOwners {
		decision := strings.TrimSpace(owner.Decision)
		if decision == "" {
			v.add("decision owner requires decision")
			continue
		}
		owners[decision] = owner
		if strings.TrimSpace(owner.Owner) == "" {
			v.add("decision %s requires owner", decision)
		}
		if !states[owner.State] {
			v.add("decision %s has unknown state %q", decision, owner.State)
		}
		if missingReasonStates[owner.State] && strings.TrimSpace(owner.Reason) == "" {
			v.add("decision %s state %s requires reason", decision, owner.State)
		}
	}
	for _, decision := range requiredDecisions {
		if owners[decision].Decision == "" {
			v.add("missing decision owner %q", decision)
		}
	}
}

func (v *bundleValidator) add(format string, args ...any) {
	v.errors = append(v.errors, fmt.Sprintf(format, args...))
}

func requiredRow(id string) bool {
	for _, required := range RequiredRows {
		if id == required {
			return true
		}
	}
	return false
}

func rowIDForRef(rows map[string]Row, ref string) string {
	for id, row := range rows {
		for _, rowRef := range row.EvidenceRefs {
			if rowRef == ref {
				return id
			}
		}
	}
	return ""
}

func gapForRow(gaps []ResidualGap, rowID string) bool {
	for _, gap := range gaps {
		if gap.RowID == rowID && strings.TrimSpace(gap.Reason) != "" {
			return true
		}
	}
	return false
}

func entryExpired(entry BundleEntry, now time.Time) bool {
	if strings.TrimSpace(entry.ExpiresAt) == "" {
		return false
	}
	expiresAt, err := time.Parse(time.RFC3339, entry.ExpiresAt)
	if err != nil {
		return true
	}
	return !expiresAt.After(now)
}

func passRefUnverifiable(entry BundleEntry) bool {
	if entry.RedactionStatus == StateCannotVerify || entry.RetainedForm == "not_retained" {
		return true
	}
	switch entry.ArtifactAccess {
	case "", "present":
		return false
	case "expired", "inaccessible", "malformed", "not_assessed", StateCannotVerify:
		return true
	default:
		return false
	}
}

func githubSourceChange(input GitHubPREvidenceInput) SourceChange {
	return SourceChange{
		Repository:  input.PR.URL,
		ChangeID:    fmt.Sprintf("PR-%d", input.PR.Number),
		URL:         input.PR.URL,
		BaseRef:     input.PR.BaseRef,
		HeadRef:     input.PR.HeadRef,
		CommitRange: input.CommitRange.Base + ".." + input.CommitRange.Head,
		HeadSHA:     input.PR.HeadSHA,
	}
}

func githubRows(input GitHubPREvidenceInput) []Row {
	change := githubChangeRow(input)
	mutation := githubMutationRow(input)
	return []Row{
		change,
		githubInitiatorRow(input),
		githubAgentRouteRow(input),
		mutation,
		githubVerificationRow(input),
		githubReviewRow(input),
		githubRow("PC-AUTHORITY", StateNotAssessed, "Authority was not assessed for this generated GitHub input.", nil, "authority profile was not provided"),
		githubRow("PC-THEATER", StatePass, "No P0 theater finding triggered by the minimal GitHub input builder.", []string{"theater:builder"}, ""),
		githubRow("PC-ATTESTATION", StateNotAssessed, "Signed or external attestation was not assessed.", nil, "signed trust inputs were not provided"),
		githubRow("PC-DECISION", StateNotAssessed, "Default decision owner placeholders are recorded.", []string{"decision:owners"}, "decision owners are placeholders, not bound approval or ownership evidence"),
		githubRow("PC-RESIDUAL-GAPS", StatePartial, "Non-pass rows remain explicit in residual gaps.", []string{"gap:generated"}, "generated packet contains explicit non-pass rows"),
	}
}

func githubChangeRow(input GitHubPREvidenceInput) Row {
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {
		return githubRow("PC-CHANGE", StateCannotVerify, "Change-host metadata is retained but commit range is incomplete.", []string{"github:pr"}, "missing commit range base or head")
	}
	return githubRow("PC-CHANGE", StatePass, "Change-host metadata and commit range are retained.", []string{"github:pr", "git:commit-range"}, "")
}

func githubMutationRow(input GitHubPREvidenceInput) Row {
	if strings.TrimSpace(input.CommitRange.Base) == "" || strings.TrimSpace(input.CommitRange.Head) == "" {
		return githubRow("PC-MUTATION", StateCannotVerify, "Commit range is incomplete.", nil, "missing commit range base or head")
	}
	return githubRow("PC-MUTATION", StatePass, "Commit range and changed files are retained.", []string{"git:commit-range"}, "")
}

func githubInitiatorRow(input GitHubPREvidenceInput) Row {
	if input.PR.BodyRef != "" {
		return githubRow("PC-INITIATOR", StatePartial, "PR body task source is retained.", []string{"github:pr-body"}, "PR body is weaker than a dedicated issue binding")
	}
	return githubRow("PC-INITIATOR", StateNotAssessed, "No task or initiator evidence was provided.", nil, "missing PR body, issue, or retained task artifact")
}

func githubAgentRouteRow(input GitHubPREvidenceInput) Row {
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if input.RequirePromptBoundary && classification.RouteProofEffect == StateFail {
		return githubRow("PC-AGENT-ROUTE", StateFail, "Developer prompt contains recorder duties.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
	}
	if input.RequirePromptBoundary && classification.RouteProofEffect == StateCannotVerify {
		return githubRow("PC-AGENT-ROUTE", StateCannotVerify, "Prompt boundary evidence cannot verify developer-route independence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
	}
	if len(input.AgentRouteRefs) > 0 {
		if input.RequirePromptBoundary && classification.RouteProofEffect == StatePartial {
			return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs and digest-only prompt boundary are retained.", []string{"agent:route", "prompt:boundary"}, "prompt text is unavailable; digest-only boundary supports partial route proof")
		}
		return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs are retained.", []string{"agent:route"}, "route refs are input refs, not a complete observed delegation chain")
	}
	return githubRow("PC-AGENT-ROUTE", StateNotAssessed, "Agent route evidence was not provided.", nil, "missing OpenCode/GSD observation ref")
}

func githubVerificationRow(input GitHubPREvidenceInput) Row {
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if input.RequirePromptBoundary && (classification.RouteProofEffect == StateFail || classification.RouteProofEffect == StateCannotVerify) {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "Verification cannot pass without clean or partially retained prompt-boundary evidence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
	}
	if len(input.Checks) == 0 {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "No GitHub check evidence was provided.", nil, "missing GitHub check or workflow run evidence")
	}
	if input.RequirePromptBoundary && strings.TrimSpace(input.WorkflowRunID) == "" {
		return githubRow("PC-VERIFICATION", StateCannotVerify, "No current workflow run id was provided.", []string{"github:check"}, "missing workflow run id for CI-owned packet generation")
	}
	if !checksHaveRetainedArtifactRefs(input) {
		return githubRow("PC-VERIFICATION", StatePartial, "GitHub check evidence is retained without retained artifact binding.", []string{"github:check"}, "GitHub CI green is not verification pass without retained artifact evidence")
	}
	for _, check := range input.Checks {
		if check.Conclusion != "success" {
			return githubRow("PC-VERIFICATION", StatePartial, "GitHub checks include non-success conclusions.", []string{"github:check"}, "not all retained checks concluded success")
		}
	}
	refs := append([]string{"github:check"}, artifactEvidenceRefs(input)...)
	if strings.TrimSpace(input.WorkflowRunID) != "" {
		return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained for workflow run "+input.WorkflowRunID+".", refs, "")
	}
	return githubRow("PC-VERIFICATION", StatePass, "GitHub check and retained artifact evidence are retained.", refs, "")
}

func artifactEvidenceRefs(input GitHubPREvidenceInput) []string {
	refs := []string{}
	seen := map[string]bool{}
	for _, check := range input.Checks {
		for _, ref := range check.ArtifactRefs {
			artifactRef := "artifact:" + ref
			if !seen[artifactRef] {
				refs = append(refs, artifactRef)
				seen[artifactRef] = true
			}
		}
	}
	return refs
}

func checksHaveRetainedArtifactRefs(input GitHubPREvidenceInput) bool {
	artifacts := map[string]bool{}
	for _, artifact := range input.Artifacts {
		if strings.TrimSpace(artifact.Name) != "" && artifact.RetainedForm != "not_retained" {
			artifacts[artifact.Name] = true
		}
	}
	for _, check := range input.Checks {
		if len(check.ArtifactRefs) == 0 {
			return false
		}
		for _, ref := range check.ArtifactRefs {
			if !artifacts[ref] {
				return false
			}
		}
	}
	return true
}

func githubReviewRow(input GitHubPREvidenceInput) Row {
	if len(input.Reviews) == 0 {
		return githubRow("PC-REVIEW", StateNotAssessed, "Review evidence was not provided.", nil, "missing GitHub review or retained external review")
	}
	for _, review := range input.Reviews {
		if review.State != StatePass {
			return githubRow("PC-REVIEW", StatePartial, "Review evidence is retained with non-pass state.", []string{"github:review"}, "review evidence did not fully pass")
		}
	}
	return githubRow("PC-REVIEW", StatePass, "Review evidence is retained.", []string{"github:review"}, "")
}

func githubRow(id, state, summary string, refs []string, reason string) Row {
	return Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Reason: reason, Owner: "maintainer"}
}

func githubEntries(input GitHubPREvidenceInput) []BundleEntry {
	entries := []BundleEntry{
		authorityEntry(bundleEntry("github:pr", "change_host", input.PR.URL, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("git:commit-range", "git", input.CommitRange.Base+".."+input.CommitRange.Head, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("theater:builder", "witness", "sdp-trace packet build-pr", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("decision:owners", "manual", "default generated decision owners", "raw"), "operator", "operator_authored", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
		authorityEntry(bundleEntry("gap:generated", "manual", "generated residual gaps", "raw"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID),
	}
	if input.RequirePromptBoundary || strings.TrimSpace(input.PromptBoundary.Text) != "" || strings.TrimSpace(input.PromptBoundary.Digest) != "" {
		entries = append(entries, authorityEntry(bundleEntry("prompt:boundary", "harness", promptBoundaryResolver(input.PromptBoundary), promptBoundaryRetainedForm(input.PromptBoundary)), "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.PromptBoundary.Digest))
	}
	if input.PR.BodyRef != "" {
		entries = append(entries, authorityEntry(bundleEntry("github:pr-body", "change_host", input.PR.BodyRef, "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID))
	}
	if len(input.AgentRouteRefs) > 0 {
		entry := bundleEntry("agent:route", "harness", strings.Join(input.AgentRouteRefs, ", "), "external_ref")
		if strings.TrimSpace(input.AgentRouteDigest) != "" {
			entry.Digest = input.AgentRouteDigest
		}
		entry.EvidenceKind = input.AgentRouteEvidenceKind
		entry.ObservedComponents = input.AgentRouteComponents
		entry = authorityEntry(entry, "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.AgentRouteDigest)
		entries = append(entries, entry)
	}
	if len(input.Checks) > 0 {
		entries = append(entries, authorityEntry(bundleEntry("github:check", "ci", checkResolvers(input.Checks), "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID))
	}
	if len(input.Reviews) > 0 {
		entries = append(entries, authorityEntry(bundleEntry("github:review", "review", reviewResolvers(input.Reviews), "external_ref"), "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID))
	}
	for _, artifact := range input.Artifacts {
		entry := bundleEntry("artifact:"+artifact.Name, "ci", artifact.Resolver, artifact.RetainedForm)
		entry.ExpiresAt = artifact.ExpiresAt
		entry.Digest = artifact.Digest
		entry = authorityEntry(entry, "ci_packet_builder", "ci_generated", "sdp-trace packet build-pr", "github_workflow_run", input.WorkflowRunID)
		entries = append(entries, entry)
	}
	for _, action := range input.IntegrationActions {
		entry := bundleEntry("integration:"+action.Kind, "manual", action.Resolver, "external_ref")
		entry = authorityEntry(entry, "integration", "integration_authored", action.Actor, "github_workflow_run", input.WorkflowRunID)
		entries = append(entries, entry)
	}
	return entries
}

func authorityEntry(entry BundleEntry, actor, writeAuthority, generatedBy, sourceCommitState, sourceRef string) BundleEntry {
	entry.Actor = actor
	entry.WriteAuthority = writeAuthority
	entry.GeneratedBy = generatedBy
	entry.SourceCommitState = sourceCommitState
	entry.SourceRef = strings.TrimSpace(sourceRef)
	return entry
}

func promptBoundaryResolver(boundary PromptBoundary) string {
	if strings.TrimSpace(boundary.Text) != "" {
		return "prompt:text-retained"
	}
	if strings.TrimSpace(boundary.Digest) != "" {
		return "prompt:digest:" + boundary.Digest
	}
	return "prompt:missing"
}

func promptBoundaryRetainedForm(boundary PromptBoundary) string {
	if strings.TrimSpace(boundary.Text) != "" {
		return "redacted"
	}
	if strings.TrimSpace(boundary.Digest) != "" {
		return "digest_only"
	}
	return "not_retained"
}

func bundleEntry(ref, sourceClass, resolver, retainedForm string) BundleEntry {
	resolver = redactSecretLike(resolver)
	return BundleEntry{
		Ref:             ref,
		SourceClass:     sourceClass,
		Digest:          digestPlaceholder(ref + resolver),
		RetainedForm:    retainedForm,
		RedactionStatus: "not_needed",
		Resolver:        resolver,
		ArtifactAccess:  "present",
	}
}

func redactSecretLike(value string) string {
	redacted := value
	for _, marker := range []string{"SECRET", "TOKEN", "Authorization:"} {
		if strings.Contains(strings.ToUpper(redacted), strings.ToUpper(marker)) {
			return "[redacted-secret]"
		}
	}
	return redacted
}

func residualGapsForRows(rows []Row) []ResidualGap {
	gaps := []ResidualGap{}
	for _, row := range rows {
		if row.State == StatePass || row.ID == "PC-RESIDUAL-GAPS" {
			continue
		}
		gaps = append(gaps, ResidualGap{RowID: row.ID, State: row.State, Reason: row.Reason, ClosureEvidence: "provide retained evidence for " + row.ID})
	}
	return gaps
}

func defaultDecisionOwners() []DecisionOwner {
	return []DecisionOwner{
		{Decision: "merge", Owner: "maintainer", State: StateNotAssessed, Reason: "packet is not approval"},
		{Decision: "release", Owner: "release owner", State: StateNotAssessed, Reason: "packet is not release approval"},
		{Decision: "risk_acceptance", Owner: "risk owner", State: StateNotAssessed, Reason: "packet is not risk acceptance"},
		{Decision: "security_review", Owner: "security owner", State: StateNotAssessed, Reason: "packet is not security review"},
	}
}

func checkResolvers(checks []GitHubCheck) string {
	values := []string{}
	for _, check := range checks {
		values = append(values, check.Name+"="+check.URL)
	}
	return strings.Join(values, ", ")
}

func reviewResolvers(reviews []GitHubReview) string {
	values := []string{}
	for _, review := range reviews {
		values = append(values, review.Reviewer+"="+review.Resolver)
	}
	return strings.Join(values, ", ")
}

func digestPlaceholder(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}

func RenderMarkdown(bundle Bundle) (string, error) {
	validation := Validate(bundle, time.Now().UTC())
	if validation.State != StatePass {
		return "", errors.New(strings.Join(validation.Errors, "; "))
	}
	packet := bundle.Packet
	var out bytes.Buffer
	fmt.Fprintf(&out, "# Change Evidence Packet v0\n\n")
	fmt.Fprintf(&out, "This packet is evidence organization, not merge, release, compliance, production trust, or quality approval.\n\n")
	renderExecutiveSummary(&out, packet)
	renderMetadata(&out, packet)
	renderRows(&out, packet.Rows)
	renderTheater(&out, packet)
	renderDecisions(&out, packet.DecisionOwners)
	renderEvidence(&out, bundle.Manifest)
	renderResidualGaps(&out, packet.ResidualGaps)
	renderNonProof(&out, packet)
	return out.String(), nil
}

func renderExecutiveSummary(out *bytes.Buffer, packet Packet) {
	fmt.Fprintf(out, "## Executive Summary\n\n")
	fmt.Fprintf(out, "- Source change: %s %s.\n", packet.SourceChange.Repository, packet.SourceChange.ChangeID)
	fmt.Fprintf(out, "- Packet state: %s.\n", packet.PacketState)
	fmt.Fprintf(out, "- Selected evidence profile: %s.\n", packet.SelectedProfile)
	fmt.Fprintf(out, "- Required rows preserve pass, partial, fail, cannot_verify, not_assessed, and not_in_scope states without a score.\n")
	fmt.Fprintf(out, "- Next decision ownership is recorded separately from approval.\n\n")
}

func renderMetadata(out *bytes.Buffer, packet Packet) {
	fmt.Fprintf(out, "## Packet Metadata\n\n")
	fmt.Fprintf(out, "| field | value |\n| --- | --- |\n")
	fields := [][2]string{
		{"packet_id", packet.PacketID},
		{"schema", packet.PacketVersion},
		{"generated_from", packet.SourceChange.URL},
		{"generated_at", packet.GeneratedAt},
		{"authoring_method", packet.AuthoringMethod},
		{"selected_profile", packet.SelectedProfile},
		{"redaction_policy", packet.RedactionPolicy},
		{"bundle_ref", packet.BundleRef},
		{"packet_state", packet.PacketState},
	}
	for _, field := range fields {
		fmt.Fprintf(out, "| %s | %s |\n", field[0], md(field[1]))
	}
	fmt.Fprintln(out)
}

func renderRows(out *bytes.Buffer, rows []Row) {
	fmt.Fprintf(out, "## Required Rows\n\n")
	fmt.Fprintf(out, "| row id | state | answer | evidence refs | gap / next evidence | owner |\n| --- | --- | --- | --- | --- | --- |\n")
	ordered := append([]Row(nil), rows...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return requiredRowIndex(ordered[i].ID) < requiredRowIndex(ordered[j].ID)
	})
	for _, row := range ordered {
		gap := row.Reason
		if gap == "" {
			gap = "none"
		}
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s |\n", row.ID, row.State, md(row.Summary), md(strings.Join(row.EvidenceRefs, ", ")), md(gap), md(row.Owner))
	}
	fmt.Fprintln(out)
}

func renderTheater(out *bytes.Buffer, packet Packet) {
	fmt.Fprintf(out, "## Theater Findings\n\n")
	fmt.Fprintf(out, "| reason code | state | severity | finding | trigger evidence | required closure evidence |\n| --- | --- | --- | --- | --- | --- |\n")
	theater := rowByID(packet.Rows, "PC-THEATER")
	if len(packet.TheaterFindings) == 0 {
		fmt.Fprintf(out, "| none | %s | none | %s | PC-THEATER row | %s |\n\n", theater.State, md(theater.Summary), md(theater.Reason))
		return
	}
	for _, finding := range packet.TheaterFindings {
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s |\n", finding.ReasonCode, finding.State, md(finding.Severity), md(finding.Finding), md(strings.Join(finding.TriggerEvidenceRefs, ", ")), md(finding.RequiredClosureEvidence))
	}
	fmt.Fprintln(out)
}

func rowByID(rows []Row, id string) Row {
	for _, row := range rows {
		if row.ID == id {
			return row
		}
	}
	return Row{ID: id, State: StateCannotVerify, Summary: "row missing", Reason: "row missing"}
}

func renderDecisions(out *bytes.Buffer, owners []DecisionOwner) {
	fmt.Fprintf(out, "## Decision Ownership\n\n")
	fmt.Fprintf(out, "| decision | owner | state | reason |\n| --- | --- | --- | --- |\n")
	for _, owner := range owners {
		fmt.Fprintf(out, "| %s | %s | %s | %s |\n", md(owner.Decision), md(owner.Owner), owner.State, md(owner.Reason))
	}
	fmt.Fprintln(out)
}

func renderEvidence(out *bytes.Buffer, manifest BundleManifest) {
	fmt.Fprintf(out, "## Evidence Bundle\n\n")
	fmt.Fprintf(out, "Manifest: `%s`\n\n", md(manifest.BundleID))
	fmt.Fprintf(out, "| ref | source class | retained form | redaction status | resolver |\n| --- | --- | --- | --- | --- |\n")
	for _, entry := range manifest.Entries {
		resolver := entry.Resolver
		if resolver == "" {
			resolver = resolverFromList(manifest.Resolvers, entry.Ref)
		}
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s |\n", md(entry.Ref), md(entry.SourceClass), md(entry.RetainedForm), md(entry.RedactionStatus), md(resolver))
	}
	fmt.Fprintln(out)
}

func renderResidualGaps(out *bytes.Buffer, gaps []ResidualGap) {
	fmt.Fprintf(out, "## Residual Gaps\n\n")
	if len(gaps) == 0 {
		fmt.Fprintf(out, "No residual gaps recorded beyond row states.\n\n")
		return
	}
	fmt.Fprintf(out, "| row id | state | reason | closure evidence |\n| --- | --- | --- | --- |\n")
	for _, gap := range gaps {
		fmt.Fprintf(out, "| %s | %s | %s | %s |\n", gap.RowID, gap.State, md(gap.Reason), md(gap.ClosureEvidence))
	}
	fmt.Fprintln(out)
}

func renderNonProof(out *bytes.Buffer, packet Packet) {
	fmt.Fprintf(out, "## What This Packet Does Not Prove\n\n")
	if strings.TrimSpace(packet.NonApproval) != "" {
		fmt.Fprintf(out, "%s\n\n", packet.NonApproval)
		return
	}
	fmt.Fprintf(out, "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.\n\n")
}

func requiredRowIndex(id string) int {
	for i, required := range RequiredRows {
		if id == required {
			return i
		}
	}
	return len(RequiredRows)
}

func resolverFromList(resolvers []ResolverEntry, ref string) string {
	for _, resolver := range resolvers {
		if resolver.Ref == ref {
			return resolver.Resolver
		}
	}
	return ""
}

func md(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}
