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
	actions := sortedObservedActions(pkg.ObservedActions)
	bindings := evaluateBindings(pkg.EvidenceBindings, actions)
	evaluations := evaluateActions(pkg, env, envState, envReason, actions, bindings)
	return authorityResult(pkg, actions, bindings, evaluations, envState)
}

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

func sortedObservedActions(actions []ObservedAction) []ObservedAction {

	out := append([]ObservedAction(nil), actions...)
	sort.Slice(out, func(i, j int) bool { return out[i].EventID < out[j].EventID })
	return out
}

func evaluateActions(pkg Package, env AuthorityEnvelope, envState, envReason string, actions []ObservedAction, bindings []EvidenceBinding) []AuthorityEvaluation {

	bindingByEvent := bindingStatesByEvent(bindings)
	resolution := evidenceResolutionIndex(pkg.EvidenceResolution)
	evaluations := make([]AuthorityEvaluation, 0, len(actions))
	for i, action := range actions {

		evaluationID := fmt.Sprintf("authority-evaluation-%03d", i+1)
		evaluations = append(evaluations, evaluateAction(evaluationID, pkg.SelectedPolicyID, env, envState, envReason, action, bindingByEvent[action.EventID], resolution))
	}
	return evaluations
}

func authorityResult(pkg Package, actions []ObservedAction, bindings []EvidenceBinding, evaluations []AuthorityEvaluation, envState string) Result {

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
	return selectMatchingEnvelope(matchingEnvelopes(pkg.AuthorityEnvelopes, selected))
}

func selectMatchingEnvelope(matches []AuthorityEnvelope) (AuthorityEnvelope, string, string) {

	switch len(matches) {
	case 0:
		return AuthorityEnvelope{}, StateNotAssessed, "selected_policy_not_found"
	case 1:
		return selectedEnvelope(matches[0])
	default:
		return matches[0], StateCannotVerify, "selected_policy_ambiguous"
	}
}

func matchingEnvelopes(envelopes []AuthorityEnvelope, selected string) []AuthorityEnvelope {
	var matches []AuthorityEnvelope
	for _, env := range envelopes {

		if env.PolicyID == selected {
			matches = append(matches, env)
		}
	}
	return matches
}

func selectedEnvelope(env AuthorityEnvelope) (AuthorityEnvelope, string, string) {
	if reason := validateEnvelope(env); reason != "" {

		return env, StateCannotVerify, reason
	}
	return env, "", ""
}

func validateEnvelope(env AuthorityEnvelope) string {
	if env.PolicyID == "" || env.ActorRef == "" || env.TaskID == "" {

		return "authority_envelope_malformed"
	}
	return firstReason(
		validateEventSet(env.AllowedEvents, env.DeniedEvents),
		validateTargetRules(env),
		validateTargetRuleOverlap(env.TargetRules),
	)
}

func firstReason(reasons ...string) string {
	for _, reason := range reasons {
		if reason != "" {

			return reason
		}
	}
	return ""
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
	if targetRuleMalformed(rule) {

		return "target_rule_malformed"
	}
	if validateEventSet(rule.AllowedEvents, rule.DeniedEvents) != "" {

		return "target_rule_conflict"
	}
	if targetRuleConflictsWithTopLevel(env, rule) {
		return "target_rule_conflicts_with_top_level"
	}
	return ""
}

func targetRuleMalformed(rule TargetRule) bool {
	return rule.RuleID == "" || rule.TargetPattern == ""
}

func targetRuleConflictsWithTopLevel(env AuthorityEnvelope, rule TargetRule) bool {
	return eventSetIntersects(rule.AllowedEvents, env.DeniedEvents) ||
		eventSetIntersects(rule.DeniedEvents, env.AllowedEvents)
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

	return firstReason(
		unsupportedEventReason(allowed),
		unsupportedEventReason(denied),
		allowDenyConflictReason(allowed, denied),
	)
}

func unsupportedEventReason(events []string) string {
	for _, event := range events {
		if !validEventType(event) {

			return "unsupported_event_type"
		}
	}
	return ""
}

func allowDenyConflictReason(allowed, denied []string) string {
	if eventSetIntersects(allowed, denied) {

		return "allow_deny_event_conflict"
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
	return eventSetsConflict(a.AllowedEvents, a.DeniedEvents, b.AllowedEvents, b.DeniedEvents)
}

func eventSetsConflict(aAllowed, aDenied, bAllowed, bDenied []string) bool {
	return eventSetIntersects(aAllowed, bDenied) || eventSetIntersects(aDenied, bAllowed)
}

func eventSetIntersects(left, right []string) bool {
	for _, event := range left {
		if contains(right, event) {

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

	eval := baseAuthorityEvaluation(evaluationID, selectedPolicyID, action)
	applyBindingAttribution(&eval, action, eventBindings)
	return eval
}

func baseAuthorityEvaluation(evaluationID, selectedPolicyID string, action ObservedAction) AuthorityEvaluation {

	return AuthorityEvaluation{
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
}

func applyBindingAttribution(eval *AuthorityEvaluation, action ObservedAction, eventBindings []EvidenceBinding) {
	if action.SourceType == "harness_log" && action.OperationID != "" {

		eval.ToolAttribution = AttributionVerified
	}
	if hasVerifiedGatewayBinding(action, eventBindings) {

		eval.ModelAttribution = AttributionVerified
	}
}

func applyPreDecisionBlockers(eval *AuthorityEvaluation, env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) bool {

	state, reason := preDecisionBlocker(env, envState, envReason, action, eventBindings, resolution)
	if reason == "" {
		return false
	}
	eval.State = state
	eval.ReasonCode = reason
	return true
}

func preDecisionBlocker(env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) (string, string) {
	if envState != "" {

		return envState, envReason
	}
	if !validEventType(action.EventType) {

		return StateCannotVerify, "unsupported_event_type"
	}
	if state, reason := preDecisionReason(env, action, eventBindings, resolution); reason != "" {
		return state, reason
	}
	return "", ""
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

		eval.State = approvalFailureState(reason)
		eval.ReasonCode = reason
		return
	}
	eval.State = decision.state
	eval.ReasonCode = decision.reasonCode
}

func approvalFailureState(reason string) string {

	if reason == "approval_evidence_missing" {
		return StateOutsideAuthority
	}
	return StateCannotVerify
}

type matchResult struct {
	state      string
	reasonCode string
	ruleRef    string
}

func matchDecision(env AuthorityEnvelope, action ObservedAction) matchResult {

	result := topLevelDecision(env, action)
	matchedTargetState := ""
	for _, rule := range env.TargetRules {
		next, targetState, ok := targetRuleDecision(rule, action)
		if !ok {
			continue
		}
		if targetStatesConflict(matchedTargetState, targetState) {

			return matchResult{state: StateCannotVerify, reasonCode: "overlapping_target_rules_conflict", ruleRef: result.ruleRef + "," + rule.RuleID}
		}
		matchedTargetState = targetState
		result = next
	}
	return result
}

func topLevelDecision(env AuthorityEnvelope, action ObservedAction) matchResult {
	if contains(env.DeniedEvents, action.EventType) {

		return matchResult{state: StateOutsideAuthority, reasonCode: "event_denied", ruleRef: "denied_events"}
	}
	if contains(env.AllowedEvents, action.EventType) {

		return matchResult{state: StateWithinAuthority, reasonCode: "event_allowed", ruleRef: "allowed_events"}
	}
	return matchResult{state: StateNotAssessed, reasonCode: "no_applicable_authority_rule"}
}

func targetRuleDecision(rule TargetRule, action ObservedAction) (matchResult, string, bool) {
	if !targetMatches(rule.TargetPattern, action.Target) {

		return matchResult{}, "", false
	}
	if contains(rule.DeniedEvents, action.EventType) {

		return matchResult{state: StateOutsideAuthority, reasonCode: "target_event_denied", ruleRef: rule.RuleID}, StateOutsideAuthority, true
	}
	if contains(rule.AllowedEvents, action.EventType) {

		return matchResult{state: StateWithinAuthority, reasonCode: "target_event_allowed", ruleRef: rule.RuleID}, StateWithinAuthority, true
	}
	return matchResult{}, "", false
}

func targetStatesConflict(previous, next string) bool {
	return (previous == StateOutsideAuthority && next == StateWithinAuthority) ||
		(previous == StateWithinAuthority && next == StateOutsideAuthority)
}

func targetMatches(pattern, target string) bool {
	if pattern == "" || target == "" {

		return false
	}
	if !strings.Contains(pattern, "**") {

		return pathMatches(pattern, target)
	}
	return recursivePathMatches(pattern, target)
}

func pathMatches(pattern, target string) bool {
	ok, err := path.Match(pattern, target)
	return err == nil && ok
}

func recursivePathMatches(pattern, target string) bool {

	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*\*`, `.*`)
	re = strings.ReplaceAll(re, `\*`, `[^/]*`)
	ok, err := regexp.MatchString("^"+re+"$", target)
	return err == nil && ok
}

func approvalReason(env AuthorityEnvelope, action ObservedAction, ruleRef string, resolution map[string]string) string {
	for _, req := range env.ApprovalRequirements {

		if reason := approvalRequirementReason(req, action, ruleRef, resolution); reason != "" {
			return reason
		}
	}
	return ""
}

func approvalRequirementReason(req ApprovalRequirement, action ObservedAction, ruleRef string, resolution map[string]string) string {
	if !approvalRequirementApplies(req, action, ruleRef) {
		return ""
	}
	if strings.TrimSpace(req.ApprovalEvidenceRef) == "" {

		return "approval_evidence_missing"
	}
	return evidenceRefsReason([]string{req.ApprovalEvidenceRef}, resolution)
}

func approvalRequirementApplies(req ApprovalRequirement, action ObservedAction, ruleRef string) bool {
	return (req.EventType == "" || req.EventType == action.EventType) &&
		(req.TargetRuleRef == "" || req.TargetRuleRef == ruleRef)
}

func evaluateBindings(inputs []EvidenceBindingInput, actions []ObservedAction) []EvidenceBinding {
	actionIDs := map[string]bool{}
	for _, action := range actions {

		actionIDs[action.EventID] = true
	}
	out := make([]EvidenceBinding, 0, len(inputs))
	for _, input := range inputs {
		out = append(out, evaluateBinding(input, actionIDs))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BindingID < out[j].BindingID })
	return out
}

func evaluateBinding(input EvidenceBindingInput, actionIDs map[string]bool) EvidenceBinding {

	state, reason := bindingStateAndReason(input, actionIDs)
	return EvidenceBinding{
		BindingID:     input.BindingID,
		LeftEventID:   input.LeftEventID,
		RightEventID:  input.RightEventID,
		BindingType:   input.BindingType,
		BindingState:  state,
		MatchedFields: input.MatchedFields,
		EvidenceRef:   input.EvidenceRef,
		ReasonCode:    reason,
	}
}

func bindingStateAndReason(input EvidenceBindingInput, actionIDs map[string]bool) (string, string) {
	if !actionIDs[input.LeftEventID] || !actionIDs[input.RightEventID] {

		return BindingNotAssessed, "binding_source_event_absent"
	}
	return knownBindingStateAndReason(input.BindingState)
}

func knownBindingStateAndReason(state string) (string, string) {

	switch state {
	case BindingVerified:
		return BindingVerified, "binding_verified"
	case BindingNotAssessed:
		return BindingNotAssessed, "binding_not_assessed"
	default:
		return BindingCannotVerify, "binding_cannot_verify"
	}
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
	addEvidenceResolution(out, input.ResolvedExternalRefs, "resolved")
	addEvidenceResolution(out, input.InaccessibleRefs, "inaccessible")
	addEvidenceResolution(out, input.MalformedRefs, "malformed")
	addEvidenceResolution(out, input.StaleRefs, "stale")
	return out
}

func addEvidenceResolution(out map[string]string, refs []string, state string) {
	for _, ref := range refs {

		out[ref] = state
	}
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

		addReasonCode(reasons, eval.ReasonCode)
	}
	for _, binding := range result.BindingEvaluations {
		addReasonCode(reasons, binding.ReasonCode)
	}
	return mapKeys(reasons)
}

func addReasonCode(reasons map[string]bool, code string) {
	if code != "" {
		reasons[code] = true
	}
}

func nextActions(result Result) []string {
	if action, ok := authorityStateNextActions[result.AuthorityEvaluationState]; ok {
		return []string{action}
	}

	return []string{"Retain evidence references if downstream consumers need replay."}
}

var authorityStateNextActions = map[string]string{
	StateCannotVerify:     "Fix malformed, stale, inaccessible, or conflicting authority evidence before using these facts.",
	StateNotAssessed:      "Supply a selected policy_id, authority envelope, applicable rule, or required evidence before claiming authority compliance.",
	StateOutsideAuthority: "External policy consumers decide whether outside_authority blocks, contaminates, or requires escalation.",
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
