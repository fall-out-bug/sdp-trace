package authority

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const (
	PackageSchemaVersion = "authority-observation-package-v1"
	ResultSchemaVersion  = "authority-evaluation-result-v1"
	Profile              = "authority-envelope"

	StateWithinAuthority  = "within_authority"
	StateOutsideAuthority = "outside_authority"
	StateNotAssessed      = "not_assessed"
	StateCannotVerify     = "cannot_verify"

	AttributionVerified     = "verified"
	AttributionNotAssessed  = "not_assessed"
	AttributionCannotVerify = "cannot_verify"

	BindingVerified     = "verified"
	BindingNotAssessed  = "not_assessed"
	BindingCannotVerify = "cannot_verify"
)

var (
	standardEventTypes = map[string]bool{
		"observe":           true,
		"review":            true,
		"feedback":          true,
		"direct_mutation":   true,
		"commit":            true,
		"merge":             true,
		"ci_run":            true,
		"harness_tool_call": true,
		"gateway_request":   true,
	}
	evidenceRefPattern           = regexp.MustCompile(`^(file:[A-Za-z0-9_./:-]+|artifact:[A-Za-z0-9_.:-]+#[A-Za-z0-9_./#:-]+|git:[A-Fa-f0-9]{40,64}#[A-Za-z0-9_./:-]+|external:[A-Za-z0-9_.:-]+)$`)
	unsafeRefPattern             = regexp.MustCompile(`(?i)(bearer |access_token=|oidc_token|BEGIN [A-Z ]*PRIVATE KEY|raw prompt|raw response|raw_job_log|private_artifact_url)`)
	evidenceRefResolutionReasons = map[string]string{
		"inaccessible": "evidence_ref_inaccessible",
		"malformed":    "evidence_ref_malformed",
		"stale":        "evidence_ref_stale",
	}
)

var aggregateStatePriority = map[string]int{
	StateCannotVerify:     3,
	StateOutsideAuthority: 2,
	StateWithinAuthority:  1,
	StateNotAssessed:      0,
}

var aggregateStateByRank = []string{
	StateNotAssessed,
	StateWithinAuthority,
	StateOutsideAuthority,
	StateCannotVerify,
}

type Package struct {
	SchemaVersion      string                 `json:"schema_version"`
	SelectedPolicyID   string                 `json:"selected_policy_id"`
	Actors             []ActorDeclaration     `json:"actors"`
	AuthorityEnvelopes []AuthorityEnvelope    `json:"authority_envelopes"`
	ObservedActions    []ObservedAction       `json:"observed_actions"`
	EvidenceBindings   []EvidenceBindingInput `json:"evidence_bindings,omitempty"`
	EvidenceResolution EvidenceResolution     `json:"evidence_resolution,omitempty"`
}

type AuthorityEnvelope struct {
	SchemaVersion        string                `json:"schema_version"`
	TaskID               string                `json:"task_id"`
	PolicyID             string                `json:"policy_id"`
	AuthorityScope       string                `json:"authority_scope"`
	ActorRef             string                `json:"actor_ref"`
	AllowedEvents        []string              `json:"allowed_events"`
	DeniedEvents         []string              `json:"denied_events"`
	TargetRules          []TargetRule          `json:"target_rules"`
	ApprovalRequirements []ApprovalRequirement `json:"approval_requirements,omitempty"`
	EffectiveFromEventID string                `json:"effective_from_event_id,omitempty"`
	SupersedesPolicyID   string                `json:"supersedes_policy_id,omitempty"`
}

type TargetRule struct {
	RuleID        string   `json:"rule_id"`
	TargetPattern string   `json:"target_pattern"`
	AllowedEvents []string `json:"allowed_events"`
	DeniedEvents  []string `json:"denied_events"`
}

type ApprovalRequirement struct {
	RequirementID       string `json:"requirement_id"`
	EventType           string `json:"event_type,omitempty"`
	TargetRuleRef       string `json:"target_rule_ref,omitempty"`
	ApprovalEvidenceRef string `json:"approval_evidence_ref"`
}

type ActorDeclaration struct {
	ActorID                string `json:"actor_id"`
	ActorType              string `json:"actor_type"`
	DeclaredRole           string `json:"declared_role"`
	Harness                string `json:"harness,omitempty"`
	Model                  string `json:"model,omitempty"`
	ModelAttributionSource string `json:"model_attribution_source,omitempty"`
	OperationID            string `json:"operation_id,omitempty"`
}

type ObservedAction struct {
	EventID               string   `json:"event_id"`
	TaskID                string   `json:"task_id,omitempty"`
	EventType             string   `json:"event_type"`
	Target                string   `json:"target,omitempty"`
	SourceType            string   `json:"source_type"`
	EvidenceRefs          []string `json:"evidence_refs"`
	ActorID               string   `json:"actor_id,omitempty"`
	OperationID           string   `json:"operation_id,omitempty"`
	ObservedAt            string   `json:"observed_at"`
	ObservationConfidence string   `json:"observation_confidence"`
}

type EvidenceBindingInput struct {
	BindingID     string   `json:"binding_id"`
	LeftEventID   string   `json:"left_event_id"`
	RightEventID  string   `json:"right_event_id"`
	BindingType   string   `json:"binding_type"`
	BindingState  string   `json:"binding_state"`
	MatchedFields []string `json:"matched_fields"`
	EvidenceRef   string   `json:"evidence_ref"`
}

type EvidenceResolution struct {
	ResolvedExternalRefs []string `json:"resolved_external_refs,omitempty"`
	InaccessibleRefs     []string `json:"inaccessible_refs,omitempty"`
	MalformedRefs        []string `json:"malformed_refs,omitempty"`
	StaleRefs            []string `json:"stale_refs,omitempty"`
}

type Result struct {
	SchemaVersion            string                `json:"schema_version"`
	SelectedProfile          string                `json:"selected_profile"`
	SelectedPolicyID         string                `json:"selected_policy_id"`
	AuthorityEvaluationState string                `json:"authority_evaluation_state"`
	Evaluations              []AuthorityEvaluation `json:"evaluations"`
	BindingEvaluations       []EvidenceBinding     `json:"binding_evaluations"`
	SourceCoverage           []string              `json:"source_coverage"`
	Reasons                  []string              `json:"reasons"`
	NextActions              []string              `json:"next_actions"`
}

type AuthorityEvaluation struct {
	EvaluationID      string   `json:"evaluation_id"`
	EventID           string   `json:"event_id"`
	PolicyID          string   `json:"policy_id"`
	State             string   `json:"state"`
	ReasonCode        string   `json:"reason_code"`
	MatchedRuleRef    string   `json:"matched_rule_ref,omitempty"`
	ActorAttribution  string   `json:"actor_attribution"`
	ToolAttribution   string   `json:"tool_attribution"`
	ModelAttribution  string   `json:"model_attribution"`
	SourceCoverage    []string `json:"source_coverage"`
	EvidenceRefs      []string `json:"evidence_refs"`
	ActorID           string   `json:"actor_id,omitempty"`
	OperationID       string   `json:"operation_id,omitempty"`
	MissingAttributes []string `json:"missing_attributes,omitempty"`
}

type EvidenceBinding struct {
	BindingID     string   `json:"binding_id"`
	LeftEventID   string   `json:"left_event_id"`
	RightEventID  string   `json:"right_event_id"`
	BindingType   string   `json:"binding_type"`
	BindingState  string   `json:"binding_state"`
	MatchedFields []string `json:"matched_fields"`
	EvidenceRef   string   `json:"evidence_ref"`
	ReasonCode    string   `json:"reason_code"`
}

func ReadPackage(path string) (Package, error) {
	var pkg Package
	raw, err := os.ReadFile(path)
	if err != nil {
		return pkg, err
	}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return pkg, err
	}
	return pkg, nil
}

func Write(path string, result Result) error {
	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

func Evaluate(pkg Package) Result {
	env, envState, envReason := selectEnvelope(pkg)
	actions := append([]ObservedAction(nil), pkg.ObservedActions...)
	sort.Slice(actions, func(i, j int) bool { return actions[i].EventID < actions[j].EventID })
	bindings := evaluateBindings(pkg.EvidenceBindings, actions)
	bindingByEvent := bindingStatesByEvent(bindings)
	resolution := evidenceResolutionIndex(pkg.EvidenceResolution)
	evaluations := make([]AuthorityEvaluation, 0, len(actions))
	for i, action := range actions {
		evaluationID := fmt.Sprintf("authority-evaluation-%03d", i+1)
		evaluations = append(evaluations, evaluateAction(evaluationID, pkg.SelectedPolicyID, env, envState, envReason, action, bindingByEvent[action.EventID], resolution))
	}
	result := Result{
		SchemaVersion:            ResultSchemaVersion,
		SelectedProfile:          Profile,
		SelectedPolicyID:         strings.TrimSpace(pkg.SelectedPolicyID),
		AuthorityEvaluationState: aggregateState(evaluations, envState),
		Evaluations:              evaluations,
		BindingEvaluations:       bindings,
		SourceCoverage:           sourceCoverage(actions),
	}
	result.Reasons = resultReasons(result)
	result.NextActions = nextActions(result)
	return result
}

func selectEnvelope(pkg Package) (AuthorityEnvelope, string, string) {
	selected := strings.TrimSpace(pkg.SelectedPolicyID)
	if selected == "" {
		return AuthorityEnvelope{}, StateNotAssessed, "policy_not_selected"
	}
	var matches []AuthorityEnvelope
	for _, env := range pkg.AuthorityEnvelopes {
		if env.PolicyID == selected {
			matches = append(matches, env)
		}
	}
	if len(matches) == 0 {
		return AuthorityEnvelope{}, StateNotAssessed, "selected_policy_not_found"
	}
	if len(matches) > 1 {
		return matches[0], StateCannotVerify, "selected_policy_ambiguous"
	}
	if reason := validateEnvelope(matches[0]); reason != "" {
		return matches[0], StateCannotVerify, reason
	}
	return matches[0], "", ""
}

func validateEnvelope(env AuthorityEnvelope) string {
	if env.PolicyID == "" || env.ActorRef == "" || env.TaskID == "" {
		return "authority_envelope_malformed"
	}
	if reason := validateEventSet(env.AllowedEvents, env.DeniedEvents); reason != "" {
		return reason
	}
	if reason := validateTargetRules(env); reason != "" {
		return reason
	}
	return validateTargetRuleOverlap(env.TargetRules)
}

func validateTargetRules(env AuthorityEnvelope) string {
	for _, rule := range env.TargetRules {
		if reason := validateTargetRule(env, rule); reason != "" {
			return reason
		}
	}
	return ""
}

func validateTargetRule(env AuthorityEnvelope, rule TargetRule) string {
	if rule.RuleID == "" || rule.TargetPattern == "" {
		return "target_rule_malformed"
	}
	if reason := validateEventSet(rule.AllowedEvents, rule.DeniedEvents); reason != "" {
		return "target_rule_conflict"
	}
	if targetRuleConflictsWithTopLevel(env, rule) {
		return "target_rule_conflicts_with_top_level"
	}
	return ""
}

func targetRuleConflictsWithTopLevel(env AuthorityEnvelope, rule TargetRule) bool {
	for _, event := range rule.AllowedEvents {
		if contains(env.DeniedEvents, event) {
			return true
		}
	}
	for _, event := range rule.DeniedEvents {
		if contains(env.AllowedEvents, event) {
			return true
		}
	}
	return false
}

func validateTargetRuleOverlap(rules []TargetRule) string {
	for i := range rules {
		if targetRuleConflictsWithAny(rules[i], rules[i+1:]) {
			return "overlapping_target_rules_conflict"
		}
	}
	return ""
}

func targetRuleConflictsWithAny(rule TargetRule, others []TargetRule) bool {
	for _, other := range others {
		if targetRulesConflict(rule, other) {
			return true
		}
	}
	return false
}

func validateEventSet(allowed, denied []string) string {
	for _, event := range append(append([]string{}, allowed...), denied...) {
		if !validEventType(event) {
			return "unsupported_event_type"
		}
	}
	for _, event := range allowed {
		if contains(denied, event) {
			return "allow_deny_event_conflict"
		}
	}
	return ""
}

func validEventType(event string) bool {
	return standardEventTypes[event] || strings.HasPrefix(event, "custom:")
}

func targetRulesConflict(a, b TargetRule) bool {
	if a.TargetPattern != b.TargetPattern {
		return false
	}
	for _, event := range a.AllowedEvents {
		if contains(b.DeniedEvents, event) {
			return true
		}
	}
	for _, event := range a.DeniedEvents {
		if contains(b.AllowedEvents, event) {
			return true
		}
	}
	return false
}

func evaluateAction(evaluationID, selectedPolicyID string, env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) AuthorityEvaluation {
	eval := newAuthorityEvaluation(evaluationID, selectedPolicyID, action, eventBindings)
	eval.MissingAttributes = missingAttributes(eval)
	if applyPreDecisionBlockers(&eval, env, envState, envReason, action, eventBindings, resolution) {
		return eval
	}
	decision := matchDecision(env, action)
	eval.MatchedRuleRef = decision.ruleRef
	applyDecision(&eval, env, action, decision, resolution)
	return eval
}

func newAuthorityEvaluation(evaluationID, selectedPolicyID string, action ObservedAction, eventBindings []EvidenceBinding) AuthorityEvaluation {
	eval := AuthorityEvaluation{
		EvaluationID:     evaluationID,
		EventID:          action.EventID,
		PolicyID:         selectedPolicyID,
		ActorAttribution: actorAttributionState(action),
		ToolAttribution:  AttributionNotAssessed,
		ModelAttribution: AttributionNotAssessed,
		SourceCoverage:   uniqueStrings([]string{action.SourceType}),
		EvidenceRefs:     safeRefs(action.EvidenceRefs),
		ActorID:          action.ActorID,
		OperationID:      action.OperationID,
	}
	if action.SourceType == "harness_log" && action.OperationID != "" {
		eval.ToolAttribution = AttributionVerified
	}
	if hasVerifiedGatewayBinding(action, eventBindings) {
		eval.ModelAttribution = AttributionVerified
	}
	return eval
}

func applyPreDecisionBlockers(eval *AuthorityEvaluation, env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) bool {
	if envState != "" {
		eval.State = envState
		eval.ReasonCode = envReason
		return true
	}
	if !validEventType(action.EventType) {
		eval.State = StateCannotVerify
		eval.ReasonCode = "unsupported_event_type"
		return true
	}
	if state, reason := preDecisionReason(env, action, eventBindings, resolution); reason != "" {
		eval.State = state
		eval.ReasonCode = reason
		return true
	}
	return false
}

func preDecisionReason(env AuthorityEnvelope, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) (string, string) {
	if state, reason := taskScopeReason(env, action); reason != "" {
		return state, reason
	}
	if reason := evidenceRefsReason(action.EvidenceRefs, resolution); reason != "" {
		return StateCannotVerify, reason
	}
	if bindingCannotVerify(eventBindings) {
		return StateCannotVerify, "evidence_binding_cannot_verify"
	}
	return "", ""
}

func taskScopeReason(env AuthorityEnvelope, action ObservedAction) (string, string) {
	if env.AuthorityScope == "repository" {
		return "", ""
	}
	if action.TaskID == "" {
		return StateNotAssessed, "task_not_assessed"
	}
	if action.TaskID != env.TaskID {
		return StateNotAssessed, "task_outside_selected_envelope"
	}
	return "", ""
}

func applyDecision(eval *AuthorityEvaluation, env AuthorityEnvelope, action ObservedAction, decision matchResult, resolution map[string]string) {
	if decision.state == StateCannotVerify {
		eval.State, eval.ReasonCode = StateCannotVerify, decision.reasonCode
		return
	}
	if decision.state == StateNotAssessed {
		eval.State, eval.ReasonCode = StateNotAssessed, "no_applicable_authority_rule"
		return
	}
	if reason := approvalReason(env, action, decision.ruleRef, resolution); reason != "" {
		if reason == "approval_evidence_missing" {
			eval.State = StateOutsideAuthority
		} else {
			eval.State = StateCannotVerify
		}
		eval.ReasonCode = reason
		return
	}
	eval.State = decision.state
	eval.ReasonCode = decision.reasonCode
}

type matchResult struct {
	state      string
	reasonCode string
	ruleRef    string
}

func matchDecision(env AuthorityEnvelope, action ObservedAction) matchResult {
	state := StateNotAssessed
	reason := "no_applicable_authority_rule"
	ruleRef := ""
	var matchedTargetState string
	if contains(env.AllowedEvents, action.EventType) {
		state = StateWithinAuthority
		reason = "event_allowed"
		ruleRef = "allowed_events"
	}
	if contains(env.DeniedEvents, action.EventType) {
		state = StateOutsideAuthority
		reason = "event_denied"
		ruleRef = "denied_events"
	}
	for _, rule := range env.TargetRules {
		if !targetMatches(rule.TargetPattern, action.Target) {
			continue
		}
		if contains(rule.AllowedEvents, action.EventType) {
			if matchedTargetState == StateOutsideAuthority {
				return matchResult{state: StateCannotVerify, reasonCode: "overlapping_target_rules_conflict", ruleRef: ruleRef + "," + rule.RuleID}
			}
			matchedTargetState = StateWithinAuthority
			state = StateWithinAuthority
			reason = "target_event_allowed"
			ruleRef = rule.RuleID
		}
		if contains(rule.DeniedEvents, action.EventType) {
			if matchedTargetState == StateWithinAuthority {
				return matchResult{state: StateCannotVerify, reasonCode: "overlapping_target_rules_conflict", ruleRef: ruleRef + "," + rule.RuleID}
			}
			matchedTargetState = StateOutsideAuthority
			state = StateOutsideAuthority
			reason = "target_event_denied"
			ruleRef = rule.RuleID
		}
	}
	return matchResult{state: state, reasonCode: reason, ruleRef: ruleRef}
}

func targetMatches(pattern, target string) bool {
	if pattern == "" || target == "" {
		return false
	}
	if !strings.Contains(pattern, "**") {
		ok, err := path.Match(pattern, target)
		return err == nil && ok
	}
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*\*`, `.*`)
	re = strings.ReplaceAll(re, `\*`, `[^/]*`)
	ok, err := regexp.MatchString("^"+re+"$", target)
	return err == nil && ok
}

func approvalReason(env AuthorityEnvelope, action ObservedAction, ruleRef string, resolution map[string]string) string {
	for _, req := range env.ApprovalRequirements {
		if req.EventType != "" && req.EventType != action.EventType {
			continue
		}
		if req.TargetRuleRef != "" && req.TargetRuleRef != ruleRef {
			continue
		}
		if strings.TrimSpace(req.ApprovalEvidenceRef) == "" {
			return "approval_evidence_missing"
		}
		if reason := evidenceRefsReason([]string{req.ApprovalEvidenceRef}, resolution); reason != "" {
			return reason
		}
	}
	return ""
}

func evaluateBindings(inputs []EvidenceBindingInput, actions []ObservedAction) []EvidenceBinding {
	actionIDs := map[string]bool{}
	for _, action := range actions {
		actionIDs[action.EventID] = true
	}
	out := make([]EvidenceBinding, 0, len(inputs))
	for _, input := range inputs {
		state := input.BindingState
		reason := ""
		switch {
		case !actionIDs[input.LeftEventID] || !actionIDs[input.RightEventID]:
			state = BindingNotAssessed
			reason = "binding_source_event_absent"
		case state == BindingVerified:
			reason = "binding_verified"
		case state == BindingNotAssessed:
			reason = "binding_not_assessed"
		default:
			state = BindingCannotVerify
			reason = "binding_cannot_verify"
		}
		out = append(out, EvidenceBinding{
			BindingID:     input.BindingID,
			LeftEventID:   input.LeftEventID,
			RightEventID:  input.RightEventID,
			BindingType:   input.BindingType,
			BindingState:  state,
			MatchedFields: input.MatchedFields,
			EvidenceRef:   input.EvidenceRef,
			ReasonCode:    reason,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BindingID < out[j].BindingID })
	return out
}

func bindingStatesByEvent(bindings []EvidenceBinding) map[string][]EvidenceBinding {
	out := map[string][]EvidenceBinding{}
	for _, binding := range bindings {
		out[binding.LeftEventID] = append(out[binding.LeftEventID], binding)
		out[binding.RightEventID] = append(out[binding.RightEventID], binding)
	}
	return out
}

func hasVerifiedGatewayBinding(action ObservedAction, bindings []EvidenceBinding) bool {
	for _, binding := range bindings {
		if binding.BindingState == BindingVerified && binding.BindingType == "same_gateway_request" {
			return true
		}
	}
	return false
}

func bindingCannotVerify(bindings []EvidenceBinding) bool {
	for _, binding := range bindings {
		if binding.BindingState == BindingCannotVerify {
			return true
		}
	}
	return false
}

func evidenceResolutionIndex(input EvidenceResolution) map[string]string {
	out := map[string]string{}
	for _, ref := range input.ResolvedExternalRefs {
		out[ref] = "resolved"
	}
	for _, ref := range input.InaccessibleRefs {
		out[ref] = "inaccessible"
	}
	for _, ref := range input.MalformedRefs {
		out[ref] = "malformed"
	}
	for _, ref := range input.StaleRefs {
		out[ref] = "stale"
	}
	return out
}

func evidenceRefsReason(refs []string, resolution map[string]string) string {
	if len(refs) == 0 {
		return "evidence_ref_missing"
	}
	for _, ref := range refs {
		if reason := evidenceRefReason(ref, resolution[ref]); reason != "" {
			return reason
		}
	}
	return ""
}

func evidenceRefReason(ref string, resolution string) string {
	if malformedEvidenceRef(ref) {
		return "evidence_ref_malformed"
	}
	if reason, ok := evidenceRefResolutionReasons[resolution]; ok {
		return reason
	}
	if unresolvedExternalEvidenceRef(ref, resolution) {
		return "external_evidence_unresolved"
	}
	return ""
}

func malformedEvidenceRef(ref string) bool {
	return unsafeRefPattern.MatchString(ref) || !evidenceRefPattern.MatchString(ref)
}

func unresolvedExternalEvidenceRef(ref string, resolution string) bool {
	return strings.HasPrefix(ref, "external:") && resolution != "resolved"
}

func aggregateState(evaluations []AuthorityEvaluation, envState string) string {
	if envState == StateCannotVerify {
		return StateCannotVerify
	}
	if len(evaluations) == 0 {
		return StateNotAssessed
	}
	rank := highestEvaluationStateRank(evaluations)
	if rank < 0 {
		return StateCannotVerify
	}
	return aggregateStateByRank[rank]
}

func highestEvaluationStateRank(evaluations []AuthorityEvaluation) int {
	highestRank := -1
	for _, eval := range evaluations {
		rank, ok := aggregateStatePriority[eval.State]
		if ok && rank > highestRank {
			highestRank = rank
		}
	}
	return highestRank
}

func resultReasons(result Result) []string {
	reasons := map[string]bool{}
	for _, eval := range result.Evaluations {
		if eval.ReasonCode != "" {
			reasons[eval.ReasonCode] = true
		}
	}
	for _, binding := range result.BindingEvaluations {
		if binding.ReasonCode != "" {
			reasons[binding.ReasonCode] = true
		}
	}
	return mapKeys(reasons)
}

func nextActions(result Result) []string {
	switch result.AuthorityEvaluationState {
	case StateCannotVerify:
		return []string{"Fix malformed, stale, inaccessible, or conflicting authority evidence before using these facts."}
	case StateNotAssessed:
		return []string{"Supply a selected policy_id, authority envelope, applicable rule, or required evidence before claiming authority compliance."}
	case StateOutsideAuthority:
		return []string{"External policy consumers decide whether outside_authority blocks, contaminates, or requires escalation."}
	default:
		return []string{"Retain evidence references if downstream consumers need replay."}
	}
}

func actorAttributionState(action ObservedAction) string {
	if strings.TrimSpace(action.ActorID) == "" {
		return AttributionNotAssessed
	}
	switch action.SourceType {
	case "harness_log", "manual_import", "pr_api", "ci_artifact":
		return AttributionVerified
	default:
		return AttributionNotAssessed
	}
}

func missingAttributes(eval AuthorityEvaluation) []string {
	var missing []string
	if eval.ActorAttribution == AttributionNotAssessed {
		missing = append(missing, "actor")
	}
	if eval.ToolAttribution == AttributionNotAssessed {
		missing = append(missing, "tool")
	}
	if eval.ModelAttribution == AttributionNotAssessed {
		missing = append(missing, "model")
	}
	return missing
}

func sourceCoverage(actions []ObservedAction) []string {
	var sources []string
	for _, action := range actions {
		sources = append(sources, action.SourceType)
	}
	return uniqueStrings(sources)
}

func safeRefs(refs []string) []string {
	out := make([]string, 0, len(refs))
	for _, ref := range refs {
		if evidenceRefPattern.MatchString(ref) && !unsafeRefPattern.MatchString(ref) {
			out = append(out, ref)
		}
	}
	return out
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func mapKeys(values map[string]bool) []string {
	out := make([]string, 0, len(values))
	for key := range values {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
