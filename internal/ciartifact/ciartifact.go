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
	source, sourceSafe := sanitizeSource(manifest.SelectedSource)
	run, runSafe := sanitizeRun(manifest.SelectedRun)
	reqs := requirements(manifest.RequiredFamilies)
	families := evaluateFamilies(reqs, manifest.ArtifactFamilies)
	index := evaluateIndex(manifest.ArtifactIndex)
	safety := evaluateSafety(manifest.OutputSafety)
	identityCannotVerify := !sourceSafe || !runSafe
	state := topLevel(families, index, safety, len(reqs), identityCannotVerify)
	result := ObservationResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileCIArtifactObservation,
		AuthorityScope:           safeAuthorityScope(manifest.AuthorityScope),
		ArtifactObservationState: state,
		SelectedSource:           source,
		SelectedRun:              run,
		ProducerScope:            aggregateProducerScope(families),
		ArtifactAccessState:      aggregateAccessState(families),
		RequiredFamilies:         orderedRequirements(reqs),
		ArtifactFamilies:         families,
		Bindings:                 bindingSummary(families),
		ArtifactIndex:            index,
		OutputSafety:             safety,
		SafetyRuleset:            defaultSafetyRuleset(manifest.SafetyRuleset),
	}
	result.Reasons = reasons(families, index, safety, identityCannotVerify)
	result.NextActions = nextActions(families, index, safety, identityCannotVerify)
	return result
}

func requirements(input []FamilyRequirement) map[string]FamilyRequirement {
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
	seen := map[string]bool{}
	var out []FamilyObservation
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, evaluateFamily(req, observed[family], true))
			seen[family] = true
		}
	}
	var extra []string
	for family := range observed {
		if !seen[family] {
			extra = append(extra, family)
		}
	}
	sort.Strings(extra)
	for _, family := range extra {
		out = append(out, evaluateFamily(FamilyRequirement{Family: family}, observed[family], false))
	}
	return out
}

func evaluateFamily(req FamilyRequirement, input FamilyInput, required bool) FamilyObservation {
	family := req.Family
	if family == "" {
		family = input.Family
	}
	producer := safeProducerScope(input.ProducerScope)
	access := AccessAbsent
	if input.ArtifactAccessState != "" {
		access = safeAccessState(input.ArtifactAccessState)
	}
	binding := BindingAbsent
	if input.BindingState != "" {
		binding = safeBindingState(input.BindingState)
	}
	result := FamilyObservation{
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
	if !required {
		return result
	}
	result.FamilyState = StatePass
	result.ReasonCode = "family_observed"
	result.Reason = "required artifact family was observed with selected proof level"
	result.NextAction = ""
	switch access {
	case AccessUnsafe:
		result.FamilyState = StateFail
		result.ReasonCode = "unsafe_artifact_output"
		result.Reason = "artifact family matched a forbidden output-safety class"
		result.NextAction = "Remove unsafe artifact content and regenerate the observation."
	case AccessAbsent:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "family_absent_in_ci_bundle"
		result.Reason = "required artifact family is absent from the selected CI bundle"
		result.NextAction = "Upload the required artifact family or mark it outside profile scope."
	case AccessPartial:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "family_partial_in_ci_bundle"
		result.Reason = "required artifact family is only partially present"
		result.NextAction = "Upload every required artifact for the selected family."
	case AccessExpired:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "artifact_expired_before_inspection"
		result.Reason = "artifact family expired before inspection"
		result.NextAction = "Regenerate CI artifacts or preserve them in an accepted external store."
	case AccessInaccessible:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "artifact_inaccessible"
		result.Reason = "artifact family could not be accessed under the selected profile"
		result.NextAction = "Provide accessible artifact evidence or mark access not assessed."
	case AccessMalformed:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "artifact_malformed"
		result.Reason = "artifact family metadata is malformed"
		result.NextAction = "Fix artifact metadata and rerun observation."
	case AccessCannotVerify:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "artifact_access_cannot_verify"
		result.Reason = "artifact family access could not be verified"
		result.NextAction = "Provide verifier-readable artifact access metadata."
	}
	if result.FamilyState != StatePass {
		return result
	}
	if req.RequiredProducerScope == ProducerCIUploaded && producer != ProducerCIUploaded {
		result.FamilyState = StateCannotVerify
		result.ReasonCode = lowerAuthorityReason(producer)
		result.Reason = "artifact family was observed below the selected CI-uploaded proof level"
		result.NextAction = "Provide CI-uploaded artifact evidence for the selected family."
		return result
	}
	switch binding {
	case BindingMismatch:
		result.FamilyState = StateFail
		result.ReasonCode = "source_run_binding_mismatch"
		result.Reason = "artifact family binding contradicts the selected source or run"
		result.NextAction = "Regenerate artifact evidence for the selected source and run."
	case BindingAbsent:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "source_run_binding_missing"
		result.Reason = "artifact family lacks selected source or run binding"
		result.NextAction = "Record source and run binding for the selected artifact family."
	case BindingUnverifiable:
		result.FamilyState = StateCannotVerify
		result.ReasonCode = "external_artifact_ref_unverifiable"
		result.Reason = "artifact family binding is unverifiable under the selected profile"
		result.NextAction = "Provide digest-checkable artifact binding evidence."
	}
	return result
}

func evaluateIndex(input ArtifactIndexInput) ArtifactIndexResult {
	state := defaultString(input.State, IndexNotAssessed)
	result := ArtifactIndexResult{State: state, Result: StateNotAssessed, ReasonCode: "artifact_index_not_assessed", Reason: "artifact index was outside selected profile scope"}
	switch state {
	case IndexValid:
		result.Result = StatePass
		result.ReasonCode = "artifact_index_valid"
		result.Reason = "artifact index is present and valid"
	case IndexSelfReference:
		result.Result = StateFail
		result.ReasonCode = "artifact_index_self_reference"
		result.Reason = "artifact index includes itself as an indexed entry"
	case IndexDigestMismatch:
		result.Result = StateFail
		result.ReasonCode = "artifact_digest_mismatch"
		result.Reason = "artifact digest contradicts selected artifact metadata"
	case IndexMissing:
		result.Result = StateCannotVerify
		result.ReasonCode = "artifact_index_missing"
		result.Reason = "required artifact index is missing"
	case IndexUnverifiable:
		result.Result = StateCannotVerify
		result.ReasonCode = "artifact_index_unverifiable"
		result.Reason = "artifact index could not be verified"
	case IndexNotAssessed:
	default:
		result.State = IndexUnverifiable
		result.Result = StateCannotVerify
		result.ReasonCode = "artifact_index_unverifiable"
		result.Reason = "artifact index state is unrecognized under selected profile"
	}
	return result
}

func evaluateSafety(input OutputSafetyInput) OutputSafetyResult {
	state := defaultString(input.State, StateNotAssessed)
	out := OutputSafetyResult{State: state, UnsafeClasses: safeClasses(input.UnsafeClasses), ReasonCode: "output_safety_not_assessed", Reason: "output safety was outside selected profile scope"}
	switch state {
	case StatePass:
		out.ReasonCode = "output_safety_pass"
		out.Reason = "observation output safety classes are absent"
	case StateFail:
		out.ReasonCode = "unsafe_artifact_output"
		out.Reason = "observation detected forbidden output-safety classes"
	case StateCannotVerify:
		out.ReasonCode = "output_safety_cannot_verify"
		out.Reason = "observation output safety could not be verified"
	case StateNotAssessed:
	default:
		out.State = StateCannotVerify
		out.ReasonCode = "output_safety_cannot_verify"
		out.Reason = "observation output safety state is unrecognized under selected profile"
	}
	return out
}

func topLevel(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, requiredCount int, identityCannotVerify bool) string {
	if anyFamilyState(families, StateFail) || index.Result == StateFail || safety.State == StateFail {
		return StateFail
	}
	if identityCannotVerify || anyFamilyState(families, StateCannotVerify) || index.Result == StateCannotVerify || safety.State == StateCannotVerify {
		return StateCannotVerify
	}
	if requiredCount == 0 {
		return StateNotAssessed
	}
	return StatePass
}

func anyFamilyState(families []FamilyObservation, state string) bool {
	for _, family := range families {
		if family.FamilyState == state {
			return true
		}
	}
	return false
}

func reasons(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	set := map[string]bool{}
	for _, family := range families {
		if family.FamilyState != StatePass && family.FamilyState != StateNotAssessed {
			set[family.ReasonCode+": "+family.Reason] = true
		}
	}
	if index.Result != StatePass && index.Result != StateNotAssessed {
		set[index.ReasonCode+": "+index.Reason] = true
	}
	if safety.State != StatePass && safety.State != StateNotAssessed {
		set[safety.ReasonCode+": "+safety.Reason] = true
	}
	if identityCannotVerify {
		set["unsafe_identity_metadata: selected source or run identity contained unsafe or unsupported metadata"] = true
	}
	return sortedKeys(set)
}

func nextActions(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	set := map[string]bool{}
	for _, family := range families {
		if family.NextAction != "" {
			set[family.NextAction] = true
		}
	}
	switch index.Result {
	case StateFail, StateCannotVerify:
		set["Regenerate or supply a verifier-readable artifact index."] = true
	}
	if safety.State == StateFail || safety.State == StateCannotVerify {
		set["Use the recorded safety ruleset id to remove unsafe artifact output before rerun."] = true
	}
	if identityCannotVerify {
		set["Provide safe source and run identity metadata before using this observation as CI-backed proof."] = true
	}
	return sortedKeys(set)
}

func bindingSummary(families []FamilyObservation) BindingSummary {
	sourceRun := BindingNotAssessed
	producer := BindingNotAssessed
	for _, family := range families {
		if !family.Required {
			continue
		}
		sourceRun = worseBinding(sourceRun, family.BindingState)
		if family.RequiredProducer == ProducerCIUploaded && family.ProducerScope != ProducerCIUploaded {
			producer = worseBinding(producer, BindingMismatch)
		} else {
			producer = worseBinding(producer, BindingMatched)
		}
	}
	return BindingSummary{SourceBindingState: sourceRun, RunBindingState: sourceRun, ProducerBindingState: producer}
}

func worseBinding(current, candidate string) string {
	if bindingRank(candidate) > bindingRank(current) {
		return candidate
	}
	return current
}

func bindingRank(state string) int {
	switch state {
	case BindingMismatch:
		return 5
	case BindingAbsent:
		return 4
	case BindingUnverifiable:
		return 3
	case BindingMatched:
		return 2
	default:
		return 1
	}
}

func aggregateProducerScope(families []FamilyObservation) string {
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {
			set[family.ProducerScope] = true
		}
	}
	return aggregate(set, ProducerNotAssessed)
}

func aggregateAccessState(families []FamilyObservation) string {
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {
			set[family.ArtifactAccessState] = true
		}
	}
	return aggregate(set, AccessNotAssessed)
}

func aggregate(set map[string]bool, empty string) string {
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
	var out []FamilyRequirement
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, req)
		}
	}
	return out
}

func lowerAuthorityReason(producer string) string {
	switch producer {
	case ProducerCheckedIn:
		return "checked_in_claim_contradicts_ci_artifacts"
	case ProducerAgentReported:
		return "agent_reported_claim_without_observed_family"
	case ProducerExternalArtifactRef:
		return "external_artifact_ref_unverifiable"
	default:
		return "lower_authority_producer_scope"
	}
}

func canonicalFamily(family string) string {
	family = strings.TrimSpace(family)
	if family == "pr_ci" {
		return "change_ci"
	}
	return family
}

func validFamily(family string) bool {
	for _, allowed := range familyOrder {
		if family == allowed {
			return true
		}
	}
	return false
}

func safeProducerScope(value string) string {
	switch value {
	case ProducerCIUploaded, ProducerCheckedIn, ProducerLocalGenerated, ProducerAgentReported, ProducerHarnessObserved, ProducerExternalArtifactRef, ProducerNotAssessed:
		return value
	default:
		return ProducerNotAssessed
	}
}

func safeRequiredProducerScope(value string) string {
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
	if safeToken(value) {
		return defaultString(value, AuthorityScopeObservation)
	}
	return AuthorityScopeObservation
}

func safeToken(value string) bool {
	return len(value) > 0 && len(value) <= 128 && safeIdentityToken(value, "_.:-")
}

func safeAccessState(value string) string {
	switch value {
	case AccessPresent, AccessAbsent, AccessPartial, AccessExpired, AccessInaccessible, AccessMalformed, AccessUnsafe, AccessNotAssessed, AccessCannotVerify:
		return value
	default:
		return AccessCannotVerify
	}
}

func safeBindingState(value string) string {
	switch value {
	case BindingMatched, BindingMismatch, BindingAbsent, BindingUnverifiable, BindingNotAssessed:
		return value
	default:
		return BindingUnverifiable
	}
}

func sanitizeSource(input SourceIdentity) (SourceIdentity, bool) {
	out := SourceIdentity{}
	ok := true
	if safeIdentityToken(input.Repository, "/._-") {
		out.Repository = input.Repository
	} else if input.Repository != "" {
		ok = false
	}
	if safeRef(input.Ref) {
		out.Ref = input.Ref
	} else if input.Ref != "" {
		ok = false
	}
	if safeHex(input.CommitSHA, 40) || safeHex(input.CommitSHA, 64) {
		out.CommitSHA = input.CommitSHA
	} else if input.CommitSHA != "" {
		ok = false
	}
	return out, ok
}

func sanitizeRun(input RunIdentity) (RunIdentity, bool) {
	out := RunIdentity{}
	ok := true
	if safeIdentityToken(input.Provider, "._-") {
		out.Provider = input.Provider
	} else if input.Provider != "" {
		ok = false
	}
	if safeIdentityToken(input.RunID, "._:-") {
		out.RunID = input.RunID
	} else if input.RunID != "" {
		ok = false
	}
	if safeIdentityToken(input.RunAttempt, "._:-") {
		out.RunAttempt = input.RunAttempt
	} else if input.RunAttempt != "" {
		ok = false
	}
	if safeIdentityToken(input.WorkflowID, "._:-") {
		out.WorkflowID = input.WorkflowID
	} else if input.WorkflowID != "" {
		ok = false
	}
	if safeIdentityToken(input.JobID, "._:-") {
		out.JobID = input.JobID
	} else if input.JobID != "" {
		ok = false
	}
	return out, ok
}

func safeRef(value string) bool {
	if value == "" {
		return true
	}
	if !(strings.HasPrefix(value, "refs/heads/") || strings.HasPrefix(value, "refs/tags/") || strings.HasPrefix(value, "refs/pull/") || strings.HasPrefix(value, "refs/merge-requests/")) {
		return false
	}
	return safeIdentityToken(value, "/._-")
}

func safeIdentityToken(value, extra string) bool {
	if value == "" || len(value) > 256 {
		return value == ""
	}
	if unsafeIdentityValue(value) {
		return false
	}
	for _, r := range value {
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') {
			continue
		}
		if strings.ContainsRune(extra, r) {
			continue
		}
		return false
	}
	return true
}

func unsafeIdentityValue(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{
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
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return strings.HasPrefix(value, "/") || strings.HasPrefix(value, "~")
}

func safeHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F') {
			continue
		}
		return false
	}
	return true
}

func safeClasses(input []string) []string {
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
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
