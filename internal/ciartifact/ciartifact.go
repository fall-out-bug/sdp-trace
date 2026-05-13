package ciartifact

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

const (
	SchemaVersion = "block26-ci-artifact-observation-v1"

	ProfileCIArtifactObservation = "ci_artifact_observation"
	AuthorityScopeObservation    = "ci_artifact_observation"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"

	ProducerCIUploaded          = "ci_uploaded"
	ProducerCheckedIn           = "checked_in"
	ProducerLocalGenerated      = "local_generated"
	ProducerAgentReported       = "agent_reported"
	ProducerHarnessObserved     = "harness_observed"
	ProducerExternalArtifactRef = "external_artifact_ref"
	ProducerNotAssessed         = "not_assessed"

	AccessPresent      = "present"
	AccessAbsent       = "absent"
	AccessPartial      = "partial"
	AccessExpired      = "expired"
	AccessInaccessible = "inaccessible"
	AccessMalformed    = "malformed"
	AccessUnsafe       = "unsafe"
	AccessNotAssessed  = "not_assessed"
	AccessCannotVerify = "cannot_verify"

	BindingMatched       = "matched"
	BindingMismatch      = "mismatch"
	BindingAbsent        = "absent"
	BindingUnverifiable  = "unverifiable"
	BindingNotAssessed   = "not_assessed"
	IndexValid           = "valid"
	IndexSelfReference   = "self_reference"
	IndexDigestMismatch  = "digest_mismatch"
	IndexMissing         = "missing"
	IndexUnverifiable    = "unverifiable"
	IndexNotAssessed     = "not_assessed"
	SafetyRulesetDefault = "block26-default-output-safety-v1"
)

var familyOrder = []string{
	"run",
	"report",
	"witness",
	"provenance",
	"evidence",
	"trace",
	"artifact_index",
	"redaction_scan",
	"review",
	"change_ci",
}

var validFamilies = map[string]bool{
	"run":            true,
	"report":         true,
	"witness":        true,
	"provenance":     true,
	"evidence":       true,
	"trace":          true,
	"artifact_index": true,
	"redaction_scan": true,
	"review":         true,
	"change_ci":      true,
}

type Manifest struct {
	SchemaVersion    string              `json:"schema_version"`
	SelectedProfile  string              `json:"selected_profile"`
	AuthorityScope   string              `json:"authority_scope"`
	SelectedSource   SourceIdentity      `json:"selected_source"`
	SelectedRun      RunIdentity         `json:"selected_run"`
	RequiredFamilies []FamilyRequirement `json:"required_families"`
	ArtifactFamilies []FamilyInput       `json:"artifact_families"`
	ArtifactIndex    ArtifactIndexInput  `json:"artifact_index"`
	OutputSafety     OutputSafetyInput   `json:"output_safety"`
	SafetyRuleset    SafetyRuleset       `json:"safety_ruleset"`
}

type SourceIdentity struct {
	Repository string `json:"repository,omitempty"`
	Ref        string `json:"ref,omitempty"`
	CommitSHA  string `json:"commit_sha,omitempty"`
}

type RunIdentity struct {
	Provider   string `json:"provider,omitempty"`
	RunID      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	WorkflowID string `json:"workflow_id,omitempty"`
	JobID      string `json:"job_id,omitempty"`
}

type FamilyRequirement struct {
	Family                string `json:"family"`
	RequiredProducerScope string `json:"required_producer_scope"`
}

type FamilyInput struct {
	Family              string `json:"family"`
	ProducerScope       string `json:"producer_scope"`
	ArtifactAccessState string `json:"artifact_access_state"`
	BindingState        string `json:"binding_state"`
	ClaimSource         string `json:"claim_source,omitempty"`
}

type ArtifactIndexInput struct {
	State string `json:"state"`
}

type OutputSafetyInput struct {
	State         string   `json:"state"`
	UnsafeClasses []string `json:"unsafe_classes,omitempty"`
}

type SafetyRuleset struct {
	ID     string `json:"id"`
	SHA256 string `json:"sha256"`
}

type ObservationResult struct {
	SchemaVersion            string              `json:"schema_version"`
	SelectedProfile          string              `json:"selected_profile"`
	AuthorityScope           string              `json:"authority_scope"`
	ArtifactObservationState string              `json:"artifact_observation_state"`
	SelectedSource           SourceIdentity      `json:"selected_source"`
	SelectedRun              RunIdentity         `json:"selected_run"`
	ProducerScope            string              `json:"producer_scope"`
	ArtifactAccessState      string              `json:"artifact_access_state"`
	RequiredFamilies         []FamilyRequirement `json:"required_families"`
	ArtifactFamilies         []FamilyObservation `json:"artifact_families"`
	Bindings                 BindingSummary      `json:"bindings"`
	ArtifactIndex            ArtifactIndexResult `json:"artifact_index"`
	OutputSafety             OutputSafetyResult  `json:"output_safety"`
	SafetyRuleset            SafetyRuleset       `json:"safety_ruleset"`
	Reasons                  []string            `json:"reasons"`
	NextActions              []string            `json:"next_actions"`
}

type FamilyObservation struct {
	Family              string `json:"family"`
	Required            bool   `json:"required"`
	RequiredProducer    string `json:"required_producer_scope,omitempty"`
	ProducerScope       string `json:"producer_scope"`
	ArtifactAccessState string `json:"artifact_access_state"`
	BindingState        string `json:"binding_state"`
	FamilyState         string `json:"family_state"`
	ReasonCode          string `json:"reason_code"`
	Reason              string `json:"reason"`
	NextAction          string `json:"next_action,omitempty"`
}

type BindingSummary struct {
	SourceBindingState   string `json:"source_binding_state"`
	RunBindingState      string `json:"run_binding_state"`
	ProducerBindingState string `json:"producer_binding_state"`
}

type ArtifactIndexResult struct {
	State      string `json:"state"`
	Result     string `json:"result"`
	ReasonCode string `json:"reason_code"`
	Reason     string `json:"reason"`
}

type OutputSafetyResult struct {
	State         string   `json:"state"`
	UnsafeClasses []string `json:"unsafe_classes,omitempty"`
	ReasonCode    string   `json:"reason_code"`
	Reason        string   `json:"reason"`
}

func Evaluate(manifest Manifest) ObservationResult {
	// Evaluation starts from the manifest inputs and produces evidence states,
	// not a scalar health score or an implicit pass claim.
	// Each later helper keeps missing, unsafe, and contradictory evidence separate.

	source, sourceSafe := sanitizeSource(manifest.SelectedSource)
	run, runSafe := sanitizeRun(manifest.SelectedRun)
	inputs := evaluatedManifestInputs(manifest)
	identityCannotVerify := !sourceSafe || !runSafe
	state := topLevel(inputs.families, inputs.index, inputs.safety, len(inputs.reqs), identityCannotVerify)
	return observationResult(manifest, source, run, inputs, state, identityCannotVerify)
}
func observationResult(manifest Manifest, source SourceIdentity, run RunIdentity, inputs evaluatedInputs, state string, identityCannotVerify bool) ObservationResult {
	// Result assembly keeps assessment state, reasons, and next actions aligned.
	// The output is a replayable CI-artifact observation, not external CI proof.

	result := baseObservationResult(manifest, source, run, inputs, state)
	addObservationGaps(&result, inputs, identityCannotVerify)
	return result
}

func baseObservationResult(manifest Manifest, source SourceIdentity, run RunIdentity, inputs evaluatedInputs, state string) ObservationResult {
	// The base result begins in cannot_verify until concrete manifest evidence
	// raises or fails individual artifact-family conditions.
	result := ObservationResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileCIArtifactObservation,
		AuthorityScope:           safeAuthorityScope(manifest.AuthorityScope),
		ArtifactObservationState: state,

		SelectedSource: source,
		SelectedRun:    run,

		ProducerScope:       aggregateProducerScope(inputs.families),
		ArtifactAccessState: aggregateAccessState(inputs.families),
		RequiredFamilies:    orderedRequirements(inputs.reqs),
		ArtifactFamilies:    inputs.families,

		Bindings:      bindingSummary(inputs.families),
		ArtifactIndex: inputs.index,
		OutputSafety:  inputs.safety,

		SafetyRuleset: defaultSafetyRuleset(manifest.SafetyRuleset),
	}
	return result
}

func addObservationGaps(result *ObservationResult, inputs evaluatedInputs, identityCannotVerify bool) {
	// Gap rows are attached only after family and safety checks have run.
	// This avoids presenting skipped checks as successful artifact coverage.

	result.Reasons = reasons(inputs.families, inputs.index, inputs.safety, identityCannotVerify)
	result.NextActions = nextActions(inputs.families, inputs.index, inputs.safety, identityCannotVerify)
}

type evaluatedInputs struct {
	reqs     map[string]FamilyRequirement
	families []FamilyObservation
	index    ArtifactIndexResult
	safety   OutputSafetyResult
}

func evaluatedManifestInputs(manifest Manifest) evaluatedInputs {
	// Manifest input evaluation is the boundary where raw manifest fields become
	// normalized artifact-family observations.

	reqs := requirements(manifest.RequiredFamilies)
	return evaluatedInputs{
		reqs:     reqs,
		families: evaluateFamilies(reqs, manifest.ArtifactFamilies),
		index:    evaluateIndex(manifest.ArtifactIndex),
		safety:   evaluateSafety(manifest.OutputSafety),
	}
}

func requirements(input []FamilyRequirement) map[string]FamilyRequirement {
	// Requirements are derived from declared producer and access expectations,
	// keeping absent declarations distinct from failed declarations.
	reqs := map[string]FamilyRequirement{}
	for _, req := range input {

		family := canonicalFamily(req.Family)
		if family == "" {
			continue
		}
		req.Family = family
		req.RequiredProducerScope = safeRequiredProducerScope(req.RequiredProducerScope)
		reqs[family] = req
	}
	return reqs
}

func evaluateFamilies(reqs map[string]FamilyRequirement, inputs []FamilyInput) []FamilyObservation {
	// Family evaluation compares required and observed artifact groups without
	// collapsing individual family verdicts into an opaque aggregate.

	observed := observedFamilies(inputs)
	out, seen := requiredFamilyObservations(reqs, observed)
	for _, family := range extraFamilies(observed, seen) {
		out = append(out, evaluateFamily(FamilyRequirement{Family: family}, observed[family], false))
	}
	return out
}

func observedFamilies(inputs []FamilyInput) map[string]FamilyInput {
	// Observed families are collected from sanitized manifest data so unsafe labels
	// cannot become source references in the final assessment.
	observed := map[string]FamilyInput{}
	for _, input := range inputs {

		input.Family = canonicalFamily(input.Family)
		if !validFamily(input.Family) {
			continue
		}
		input.ProducerScope = safeProducerScope(input.ProducerScope)
		input.ArtifactAccessState = safeAccessState(input.ArtifactAccessState)
		input.BindingState = safeBindingState(input.BindingState)

		observed[input.Family] = input
	}
	return observed
}

func requiredFamilyObservations(reqs map[string]FamilyRequirement, observed map[string]FamilyInput) ([]FamilyObservation, map[string]bool) {
	// Required family observations preserve the requested producer scope and access
	// state before observed artifacts are matched.

	seen := map[string]bool{}
	out := make([]FamilyObservation, 0, len(reqs)+len(observed))
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, evaluateFamily(req, observed[family], true))
			seen[family] = true
		}
	}
	return out, seen
}

func extraFamilies(observed map[string]FamilyInput, seen map[string]bool) []string {
	// Extra family detection reports unexpected artifact groups without treating
	// them as proof that required families were covered.
	var extra []string
	for family := range observed {
		if !seen[family] {

			extra = append(extra, family)
		}
	}
	sort.Strings(extra)
	return extra
}

func evaluateFamily(req FamilyRequirement, input FamilyInput, required bool) FamilyObservation {
	// A single family verdict is built from access, producer, and binding evidence.
	// The helper keeps those dimensions independently reviewable.

	state := familyInputState(input)
	result := initialFamilyObservation(req, input, required, state.producer, state.access, state.binding)
	if !required {

		return result
	}
	markRequiredFamilyObserved(&result)
	if applyAccessResult(&result, state.access) {

		return result
	}
	if applyRequiredProducerResult(&result, req.RequiredProducerScope, state.producer) {

		return result
	}
	applyBindingResult(&result, state.binding)
	return result
}

type familyInput struct {
	producer string
	access   string
	binding  string
}

func familyInputState(input FamilyInput) familyInput {
	// familyInputState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	return familyInput{
		producer: safeProducerScope(input.ProducerScope),
		access:   familyAccessState(input),
		binding:  familyBindingState(input),
	}
}

func markRequiredFamilyObserved(result *FamilyObservation) {
	// markRequiredFamilyObserved keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	result.FamilyState = StatePass
	result.ReasonCode = "family_observed"
	result.Reason = "required artifact family was observed with selected proof level"
	result.NextAction = ""
}

func familyAccessState(input FamilyInput) string {
	// familyAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if input.ArtifactAccessState == "" {

		return AccessAbsent
	}
	return safeAccessState(input.ArtifactAccessState)
}

func familyBindingState(input FamilyInput) string {
	// familyBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if input.BindingState == "" {

		return BindingAbsent
	}
	return safeBindingState(input.BindingState)
}

func initialFamilyObservation(req FamilyRequirement, input FamilyInput, required bool, producer, access, binding string) FamilyObservation {
	// initialFamilyObservation keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	family := req.Family
	if family == "" {

		family = input.Family
	}

	return FamilyObservation{
		Family:              family,
		Required:            required,
		RequiredProducer:    req.RequiredProducerScope,
		ProducerScope:       producer,
		ArtifactAccessState: access,
		BindingState:        binding,
		FamilyState:         StateNotAssessed,
		ReasonCode:          "family_not_selected",
		Reason:              "artifact family was outside the selected profile scope",
	}
}

type familyOutcome struct {
	state  string
	code   string
	reason string
	action string
}

var accessResults = map[string]familyOutcome{
	AccessUnsafe:       {StateFail, "unsafe_artifact_output", "artifact family matched a forbidden output-safety class", "Remove unsafe artifact content and regenerate the observation."},
	AccessAbsent:       {StateCannotVerify, "family_absent_in_ci_bundle", "required artifact family is absent from the selected CI bundle", "Upload the required artifact family or mark it outside profile scope."},
	AccessPartial:      {StateCannotVerify, "family_partial_in_ci_bundle", "required artifact family is only partially present", "Upload every required artifact for the selected family."},
	AccessExpired:      {StateCannotVerify, "artifact_expired_before_inspection", "artifact family expired before inspection", "Regenerate CI artifacts or preserve them in an accepted external store."},
	AccessInaccessible: {StateCannotVerify, "artifact_inaccessible", "artifact family could not be accessed under the selected profile", "Provide accessible artifact evidence or mark access not assessed."},
	AccessMalformed:    {StateCannotVerify, "artifact_malformed", "artifact family metadata is malformed", "Fix artifact metadata and rerun observation."},
	AccessCannotVerify: {StateCannotVerify, "artifact_access_cannot_verify", "artifact family access could not be verified", "Provide verifier-readable artifact access metadata."},
}

var bindingResults = map[string]familyOutcome{
	BindingMismatch:     {StateFail, "source_run_binding_mismatch", "artifact family binding contradicts the selected source or run", "Regenerate artifact evidence for the selected source and run."},
	BindingAbsent:       {StateCannotVerify, "source_run_binding_missing", "artifact family lacks selected source or run binding", "Record source and run binding for the selected artifact family."},
	BindingUnverifiable: {StateCannotVerify, "external_artifact_ref_unverifiable", "artifact family binding is unverifiable under the selected profile", "Provide digest-checkable artifact binding evidence."},
}

func applyAccessResult(result *FamilyObservation, access string) bool {
	// applyAccessResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if outcome, ok := accessResults[access]; ok {

		setFamilyResult(result, outcome)
	}
	return result.FamilyState != StatePass
}

func applyRequiredProducerResult(result *FamilyObservation, requiredProducer, producer string) bool {
	// applyRequiredProducerResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	requiresCIUploaded := requiredProducer == ProducerCIUploaded
	producerIsCIUploaded := producer == ProducerCIUploaded
	if !requiresCIUploaded || producerIsCIUploaded {

		return false
	}

	setFamilyResult(result, familyOutcome{StateCannotVerify, lowerAuthorityReason(producer), "artifact family was observed below the selected CI-uploaded proof level", "Provide CI-uploaded artifact evidence for the selected family."})
	return true
}

func applyBindingResult(result *FamilyObservation, binding string) {
	// applyBindingResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if outcome, ok := bindingResults[binding]; ok {

		setFamilyResult(result, outcome)
	}
}

func setFamilyResult(result *FamilyObservation, outcome familyOutcome) {
	// setFamilyResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	result.FamilyState = outcome.state
	result.ReasonCode = outcome.code
	result.Reason = outcome.reason
	result.NextAction = outcome.action
}

func evaluateIndex(input ArtifactIndexInput) ArtifactIndexResult {
	// Index evaluation checks whether the manifest index itself is present and safe
	// before any artifact-family evidence is trusted.

	state := defaultString(input.State, IndexNotAssessed)
	if outcome, ok := indexOutcomes[state]; ok {
		return ArtifactIndexResult{State: state, Result: outcome.state, ReasonCode: outcome.code, Reason: outcome.reason}
	}
	outcome := indexOutcomes[IndexUnverifiable]

	return ArtifactIndexResult{State: IndexUnverifiable, Result: outcome.state, ReasonCode: outcome.code, Reason: "artifact index state is unrecognized under selected profile"}
}

var indexOutcomes = map[string]familyOutcome{
	IndexValid:          {StatePass, "artifact_index_valid", "artifact index is present and valid", ""},
	IndexSelfReference:  {StateFail, "artifact_index_self_reference", "artifact index includes itself as an indexed entry", ""},
	IndexDigestMismatch: {StateFail, "artifact_digest_mismatch", "artifact digest contradicts selected artifact metadata", ""},
	IndexMissing:        {StateCannotVerify, "artifact_index_missing", "required artifact index is missing", ""},
	IndexUnverifiable:   {StateCannotVerify, "artifact_index_unverifiable", "artifact index could not be verified", ""},
	IndexNotAssessed:    {StateNotAssessed, "artifact_index_not_assessed", "artifact index was outside selected profile scope", ""},
}

func evaluateSafety(input OutputSafetyInput) OutputSafetyResult {
	// Safety evaluation scans source and run identity fields before they are copied
	// into human-facing reasons or machine-readable refs.

	state := defaultString(input.State, StateNotAssessed)
	outcome, ok := safetyOutcomes[state]
	if !ok {

		state = StateCannotVerify
		outcome = familyOutcome{StateCannotVerify, "output_safety_cannot_verify", "observation output safety state is unrecognized under selected profile", ""}
	}
	return OutputSafetyResult{State: state, UnsafeClasses: safeClasses(input.UnsafeClasses), ReasonCode: outcome.code, Reason: outcome.reason}
}

var safetyOutcomes = map[string]familyOutcome{
	StatePass:         {StatePass, "output_safety_pass", "observation output safety classes are absent", ""},
	StateFail:         {StateFail, "unsafe_artifact_output", "observation detected forbidden output-safety classes", ""},
	StateCannotVerify: {StateCannotVerify, "output_safety_cannot_verify", "observation output safety could not be verified", ""},
	StateNotAssessed:  {StateNotAssessed, "output_safety_not_assessed", "output safety was outside selected profile scope", ""},
}

func topLevel(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, requiredCount int, identityCannotVerify bool) string {
	// Top-level aggregation reports the worst live evidence state without hiding
	// lower-level family failures.
	if artifactAssessmentHasState(families, index, safety, StateFail) {

		return StateFail
	}
	if identityOrArtifactCannotVerify(identityCannotVerify, families, index, safety) {

		return StateCannotVerify
	}
	if requiredCount == 0 {

		return StateNotAssessed
	}
	return StatePass
}

func identityOrArtifactCannotVerify(identityCannotVerify bool, families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult) bool {

	return identityCannotVerify || artifactAssessmentHasState(families, index, safety, StateCannotVerify)
}

func artifactAssessmentHasState(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, state string) bool {

	return anyFamilyState(families, state) || index.Result == state || safety.State == state
}

func anyFamilyState(families []FamilyObservation, state string) bool {
	// anyFamilyState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {
		if family.FamilyState == state {

			return true
		}
	}
	return false
}

func reasons(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	// Reasons are derived from the recorded family states so prose follows machine
	// evidence instead of becoming independent authority.
	set := map[string]bool{}

	addFamilyReasons(set, families)
	addVisibleReason(set, index.Result, index.ReasonCode, index.Reason)
	addVisibleReason(set, safety.State, safety.ReasonCode, safety.Reason)
	addIdentityReason(set, identityCannotVerify)
	return sortedKeys(set)
}

func addFamilyReasons(set map[string]bool, families []FamilyObservation) {
	// addFamilyReasons keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {

		addVisibleReason(set, family.FamilyState, family.ReasonCode, family.Reason)
	}
}

func addVisibleReason(set map[string]bool, state, code, reason string) {
	// addVisibleReason keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if familyReasonVisible(state) {

		set[code+": "+reason] = true
	}
}

func addIdentityReason(set map[string]bool, identityCannotVerify bool) {
	if identityCannotVerify {
		set["unsafe_identity_metadata: selected source or run identity contained unsafe or unsupported metadata"] = true
	}
}

func familyReasonVisible(state string) bool {
	return state != StatePass && state != StateNotAssessed
}

func nextActions(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	// Next actions name the smallest missing or unsafe evidence boundary needed to
	// move the CI-artifact observation forward.
	set := map[string]bool{}

	addFamilyActions(set, families)
	addConditionalAction(set, resultNeedsAction(index.Result), "Regenerate or supply a verifier-readable artifact index.")
	addConditionalAction(set, resultNeedsAction(safety.State), "Use the recorded safety ruleset id to remove unsafe artifact output before rerun.")
	addConditionalAction(set, identityCannotVerify, "Provide safe source and run identity metadata before using this observation as CI-backed proof.")
	return sortedKeys(set)
}

func addFamilyActions(set map[string]bool, families []FamilyObservation) {
	// addFamilyActions keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {

		addConditionalAction(set, family.NextAction != "", family.NextAction)
	}
}

func addConditionalAction(set map[string]bool, include bool, action string) {
	if include {
		set[action] = true
	}
}

func resultNeedsAction(state string) bool {
	return state == StateFail || state == StateCannotVerify
}

func bindingSummary(families []FamilyObservation) BindingSummary {
	// Binding summaries keep producer scope and access state separate because each
	// can fail independently.

	sourceRun := BindingNotAssessed
	producer := BindingNotAssessed
	for _, family := range families {
		if !family.Required {

			continue
		}
		sourceRun = worseBinding(sourceRun, family.BindingState)
		producer = worseBinding(producer, producerBindingState(family))
	}
	return BindingSummary{SourceBindingState: sourceRun, RunBindingState: sourceRun, ProducerBindingState: producer}
}

func producerBindingState(family FamilyObservation) string {
	// producerBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if family.RequiredProducer == ProducerCIUploaded && family.ProducerScope != ProducerCIUploaded {

		return BindingMismatch
	}
	return BindingMatched
}

func worseBinding(current, candidate string) string {
	// worseBinding keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if bindingRank(candidate) > bindingRank(current) {

		return candidate
	}
	return current
}

func bindingRank(state string) int {
	// bindingRank keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if rank, ok := bindingRanks[state]; ok {
		return rank
	}

	return 1
}

var bindingRanks = map[string]int{
	BindingMismatch:     5,
	BindingAbsent:       4,
	BindingUnverifiable: 3,
	BindingMatched:      2,
}

func aggregateProducerScope(families []FamilyObservation) string {
	// aggregateProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {

			set[family.ProducerScope] = true
		}
	}
	return aggregate(set, ProducerNotAssessed)
}

func aggregateAccessState(families []FamilyObservation) string {
	// aggregateAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {

			set[family.ArtifactAccessState] = true
		}
	}
	return aggregate(set, AccessNotAssessed)
}

func aggregate(set map[string]bool, empty string) string {
	// Aggregation preserves cannot_verify and fail precedence rather than averaging
	// artifact-family outcomes.
	if len(set) == 0 {

		return empty
	}
	if len(set) == 1 {
		for value := range set {

			return value
		}
	}

	return "mixed"
}

func orderedRequirements(reqs map[string]FamilyRequirement) []FamilyRequirement {
	// orderedRequirements keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	out := make([]FamilyRequirement, 0, len(reqs))
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, req)
		}
	}
	return out
}

func lowerAuthorityReason(producer string) string {
	// lowerAuthorityReason keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if reason, ok := lowerAuthorityReasons[producer]; ok {
		return reason
	}

	return "lower_authority_producer_scope"
}

var lowerAuthorityReasons = map[string]string{
	ProducerCheckedIn:           "checked_in_claim_contradicts_ci_artifacts",
	ProducerAgentReported:       "agent_reported_claim_without_observed_family",
	ProducerExternalArtifactRef: "external_artifact_ref_unverifiable",
}

func canonicalFamily(family string) string {
	// canonicalFamily keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	family = strings.TrimSpace(family)
	if family == "pr_ci" {

		return "change_ci"
	}
	return family
}

func validFamily(family string) bool {
	return validFamilies[family]
}

func safeProducerScope(value string) string {
	// safeProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validProducerScopes[value] {
		return value
	}

	return ProducerNotAssessed
}

var validProducerScopes = map[string]bool{
	ProducerCIUploaded:          true,
	ProducerCheckedIn:           true,
	ProducerLocalGenerated:      true,
	ProducerAgentReported:       true,
	ProducerHarnessObserved:     true,
	ProducerExternalArtifactRef: true,
	ProducerNotAssessed:         true,
}

func safeRequiredProducerScope(value string) string {
	// safeRequiredProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if strings.TrimSpace(value) == "" {

		return ProducerCIUploaded
	}
	scope := safeProducerScope(value)
	if scope == ProducerNotAssessed && value != ProducerNotAssessed {

		return ProducerCIUploaded
	}
	return scope
}

func safeAuthorityScope(value string) string {
	// safeAuthorityScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if safeToken(value) {
		return defaultString(value, AuthorityScopeObservation)
	}

	return AuthorityScopeObservation
}

func safeToken(value string) bool {
	return len(value) > 0 && len(value) <= 128 && safeIdentityToken(value, "_.:-")
}

func safeAccessState(value string) string {
	// safeAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validAccessStates[value] {
		return value
	}

	return AccessCannotVerify
}

var validAccessStates = map[string]bool{
	AccessPresent:      true,
	AccessAbsent:       true,
	AccessPartial:      true,
	AccessExpired:      true,
	AccessInaccessible: true,
	AccessMalformed:    true,
	AccessUnsafe:       true,
	AccessNotAssessed:  true,
	AccessCannotVerify: true,
}

func safeBindingState(value string) string {
	// safeBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validBindingStates[value] {
		return value
	}

	return BindingUnverifiable
}

var validBindingStates = map[string]bool{
	BindingMatched:      true,
	BindingMismatch:     true,
	BindingAbsent:       true,
	BindingUnverifiable: true,
	BindingNotAssessed:  true,
}

func sanitizeSource(input SourceIdentity) (SourceIdentity, bool) {
	// sanitizeSource keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	repository, repositoryOK := sanitizeSourceField(input.Repository, func(value string) bool { return safeIdentityToken(value, "/._-") })
	ref, refOK := sanitizeSourceField(input.Ref, safeRef)
	commitSHA, commitOK := sanitizeSourceField(input.CommitSHA, safeCommitSHA)
	return SourceIdentity{Repository: repository, Ref: ref, CommitSHA: commitSHA}, allTrue(repositoryOK, refOK, commitOK)
}

func sanitizeSourceField(value string, valid func(string) bool) (string, bool) {
	// sanitizeSourceField keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" || valid(value) {
		return value, true
	}

	return "", false
}

func safeCommitSHA(value string) bool {
	return safeHex(value, 40) || safeHex(value, 64)
}

func sanitizeRun(input RunIdentity) (RunIdentity, bool) {
	// Run sanitization strips unsafe identity fields before run metadata can appear
	// in report surfaces.

	out := RunIdentity{}
	var okProvider, okRunID, okRunAttempt, okWorkflowID, okJobID bool
	out.Provider, okProvider = sanitizeRunField(input.Provider, "._-")
	out.RunID, okRunID = sanitizeRunField(input.RunID, "._:-")
	out.RunAttempt, okRunAttempt = sanitizeRunField(input.RunAttempt, "._:-")
	out.WorkflowID, okWorkflowID = sanitizeRunField(input.WorkflowID, "._:-")
	out.JobID, okJobID = sanitizeRunField(input.JobID, "._:-")
	return out, allTrue(okProvider, okRunID, okRunAttempt, okWorkflowID, okJobID)
}

func sanitizeRunField(value, extra string) (string, bool) {
	// sanitizeRunField keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" || safeIdentityToken(value, extra) {
		return value, true
	}

	return "", false
}

func allTrue(values ...bool) bool {
	// allTrue keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, value := range values {
		if !value {

			return false
		}
	}
	return true
}

func safeRef(value string) bool {
	// safeRef keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" {
		return true
	}
	if !safeRefPrefix(value) {

		return false
	}
	return safeIdentityToken(value, "/._-")
}

func safeRefPrefix(value string) bool {
	// safeRefPrefix keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	return strings.HasPrefix(value, "refs/heads/") ||
		strings.HasPrefix(value, "refs/tags/") ||
		strings.HasPrefix(value, "refs/pull/") ||
		strings.HasPrefix(value, "refs/merge-requests/")
}

func safeIdentityToken(value, extra string) bool {
	// Identity tokens are allow-listed because artifact source labels can be echoed
	// into refs and reasons.
	if value == "" {
		return true
	}
	if !safeIdentityTokenLength(value) {

		return false
	}
	if unsafeIdentityValue(value) {

		return false
	}
	return safeIdentityTokenCharacters(value, extra)
}

func safeIdentityTokenLength(value string) bool {
	return len(value) <= 256
}

func safeIdentityTokenCharacters(value, extra string) bool {
	// safeIdentityTokenCharacters keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, r := range value {
		if !safeIdentityTokenRune(r, extra) {

			return false
		}
	}
	return true
}

func safeIdentityTokenRune(r rune, extra string) bool {
	return safeIdentityTokenAlnum(r) || strings.ContainsRune(extra, r)
}

func safeIdentityTokenAlnum(r rune) bool {
	return strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", r)
}

func unsafeIdentityValue(value string) bool {
	// unsafeIdentityValue keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	lower := strings.ToLower(value)

	if containsUnsafeIdentityMarker(lower) {
		return true
	}

	return strings.HasPrefix(value, "/") || strings.HasPrefix(value, "~")
}

func containsUnsafeIdentityMarker(lower string) bool {
	// containsUnsafeIdentityMarker keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, marker := range unsafeIdentityMarkers {
		if strings.Contains(lower, marker) {

			return true
		}
	}
	return false
}

var unsafeIdentityMarkers = []string{
	"token=",
	"access_token",
	"oidc_token",
	"bearer",
	"private_artifact_url",
	"raw prompt",
	"raw response",
	"raw_job_log",
	"begin_private_key",
	"begin-private-key",
	"eyj",
}

func safeHex(value string, length int) bool {
	// safeHex keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if len(value) != length {

		return false
	}
	for _, r := range value {
		if !isHexRune(r) {

			return false
		}
	}
	return true
}

func isHexRune(r rune) bool {
	return strings.ContainsRune("0123456789abcdefABCDEF", r)
}

func safeClasses(input []string) []string {
	// safeClasses keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	allowed := map[string]bool{
		"token_like": true, "jwt_token": true, "private_key": true,
		"cloud_credential": true, "provider_token": true,
		"private_artifact_url": true, "private_filesystem_path": true,
		"prompt_or_model_payload": true, "raw_job_log": true,
		"high_entropy_secret": true,
	}
	var out []string
	for _, value := range input {
		if allowed[value] {

			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func defaultSafetyRuleset(input SafetyRuleset) SafetyRuleset {
	// The default safety ruleset is explicit product behavior, not hidden policy.
	// Callers can compare its digest to understand which checks ran.
	if !safeToken(input.ID) {

		input.ID = SafetyRulesetDefault
	}
	if !safeHex(input.SHA256, 64) {

		sum := sha256.Sum256([]byte(defaultSafetyRulesetContent()))
		input.SHA256 = hex.EncodeToString(sum[:])
	}
	return input
}

func defaultSafetyRulesetContent() string {
	// Ruleset content stays deterministic so its digest is stable evidence for the
	// safety rules applied to this observation.

	return strings.Join([]string{
		SafetyRulesetDefault,
		"token_like",
		"jwt_token",
		"private_key",
		"cloud_credential",
		"provider_token",
		"private_artifact_url",
		"private_filesystem_path",
		"prompt_or_model_payload",
		"raw_job_log",
		"high_entropy_secret",
	}, "\n")
}

func defaultString(value, fallback string) string {
	// defaultString keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if strings.TrimSpace(value) == "" {

		return fallback
	}
	return value
}

func sortedKeys(set map[string]bool) []string {
	// sortedKeys keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
