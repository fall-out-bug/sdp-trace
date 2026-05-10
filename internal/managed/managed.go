package managed

import (
	"sort"
)

const (
	SchemaVersion = "block17-managed-assessment-v1"

	ProfileManagedHarness = "managed_harness"
	TrustScopeManaged     = "managed_harness_observed"

	StatePass             = "pass"
	StateFail             = "fail"
	StateCannotVerify     = "cannot_verify"
	StateNotAssessed      = "not_assessed"
	StateMissingTelemetry = "missing_telemetry"
	StateNotIntegrated    = "not_integrated"
	StateUnsupported      = "unsupported"
	StateSuppressed       = "suppressed"

	IdentityVerified     = "verified"
	IdentitySelfClaimed  = "self_claimed"
	IdentityUnauthorized = "unauthorized"
)

var managedConditionIDs = []string{
	"managed_profile_explicitly_selected",
	"managed_policy_loaded",
	"adapter_registry_loaded",
	"managed_boundary_enrolled_before_run",
	"adapter_identity_authorized",
	"adapter_capabilities_satisfy_contract",
	"adapter_activation_observed",
	"adapter_connection_continuous",
	"required_harness_events_observed",
	"required_tool_events_observed",
	"required_file_mutations_observed",
	"test_provenance_not_agent_reported",
	"suppression_policy_valid",
	"bypass_not_observed",
	"managed_witness_bound",
	"override_does_not_upgrade_managed_profile",
}

type Input struct {
	Contract Contract
	Policy   Policy
	Registry Registry
	Run      RunEvidence
	Witness  Witness
}

type Contract struct {
	RequiredEventTypes []string `json:"required_event_types"`
}

type Policy struct {
	PolicyID            string               `json:"policy_id"`
	PolicyProvenance    Provenance           `json:"policy_provenance"`
	AuthorizedAdapters  []AuthorizedAdapter  `json:"authorized_adapters"`
	RequiredEventGroups []RequiredEventGroup `json:"required_event_groups"`
	SuppressionRules    []SuppressionRule    `json:"suppression_rules,omitempty"`
}

type Provenance struct {
	Source string `json:"source"`
	Digest string `json:"digest"`
}

type AuthorizedAdapter struct {
	AdapterID       string   `json:"adapter_id"`
	HarnessID       string   `json:"harness_id"`
	AuthorityRef    string   `json:"authority_ref"`
	DeploymentRef   string   `json:"deployment_ref"`
	VersionRequired string   `json:"version_required"`
	CapabilityIDs   []string `json:"capability_ids"`
}

type RequiredEventGroup struct {
	ID                           string   `json:"id"`
	EventTypes                   []string `json:"event_types"`
	AcceptableProvenanceScopes   []string `json:"acceptable_provenance_scopes"`
	SuppressionMaySatisfyProfile bool     `json:"suppression_may_satisfy_profile,omitempty"`
}

type SuppressionRule struct {
	EventGroup                   string `json:"event_group"`
	AuthorityRef                 string `json:"authority_ref"`
	PolicyProvenanceSource       string `json:"policy_provenance_source"`
	SuppressionMaySatisfyProfile bool   `json:"suppression_may_satisfy_profile"`
}

type Registry struct {
	RegistryID string     `json:"registry_id"`
	Provenance Provenance `json:"provenance"`
	Adapters   []Adapter  `json:"adapters"`
}

type Adapter struct {
	AdapterID      string       `json:"adapter_id"`
	HarnessID      string       `json:"harness_id"`
	Version        string       `json:"version"`
	DeploymentRef  string       `json:"deployment_ref"`
	IdentityState  string       `json:"identity_state"`
	AuthorityRef   string       `json:"authority_ref"`
	AllowedEvents  []string     `json:"allowed_events"`
	CapabilityRefs []string     `json:"capability_refs"`
	Capabilities   []Capability `json:"capabilities"`
}

type Capability struct {
	ID              string   `json:"id"`
	Version         string   `json:"version,omitempty"`
	EventTypes      []string `json:"event_types"`
	ProvenanceScope string   `json:"provenance_scope"`
}

type RunEvidence struct {
	RunID                        string                   `json:"run_id"`
	RunNonce                     string                   `json:"run_nonce"`
	SourceCommit                 string                   `json:"source_commit"`
	ChainHead                    string                   `json:"chain_head"`
	EventCount                   int                      `json:"event_count"`
	ManagedBoundaryEnrolled      *ManagedBoundaryEnrolled `json:"managed_boundary_enrolled,omitempty"`
	ChildLaunch                  LaunchEvent              `json:"child_launch"`
	ObservedEvents               []EvidenceEvent          `json:"observed_events"`
	TestEvidence                 []EvidenceEvent          `json:"test_evidence"`
	SuppressedEventGroups        []SuppressedEventGroup   `json:"suppressed_event_groups,omitempty"`
	AdapterDisconnectObserved    bool                     `json:"adapter_disconnect_observed,omitempty"`
	BypassObserved               bool                     `json:"bypass_observed,omitempty"`
	OutputArtifacts              []ArtifactDigest         `json:"output_artifacts"`
	OverrideAttemptsTrustUpgrade bool                     `json:"override_attempts_trust_upgrade,omitempty"`
	OverridePresent              bool                     `json:"override_present,omitempty"`
}

type ManagedBoundaryEnrolled struct {
	Sequence              int    `json:"sequence"`
	ManagedPolicyDigest   string `json:"managed_policy_digest"`
	AdapterRegistryDigest string `json:"adapter_registry_digest"`
	AdapterID             string `json:"adapter_id"`
	EnrollmentSource      string `json:"enrollment_source"`
	ManagedProfileID      string `json:"managed_profile_id"`
	ParentRunID           string `json:"parent_run_id"`
	RunNonce              string `json:"run_nonce"`
	EventDigest           string `json:"event_digest"`
}

type LaunchEvent struct {
	Sequence    int    `json:"sequence"`
	EventDigest string `json:"event_digest"`
}

type EvidenceEvent struct {
	EventType       string `json:"event_type"`
	ProvenanceScope string `json:"provenance_scope"`
}

type SuppressedEventGroup struct {
	EventGroup             string `json:"event_group"`
	AuthorizedByPolicy     bool   `json:"authorized_by_policy"`
	PolicyProvenanceSource string `json:"policy_provenance_source"`
}

type ArtifactDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type Witness struct {
	WitnessID             string           `json:"witness_id"`
	Status                string           `json:"status"`
	RunID                 string           `json:"run_id"`
	RunNonce              string           `json:"run_nonce"`
	SourceCommit          string           `json:"source_commit"`
	ManagedPolicyDigest   string           `json:"managed_policy_digest"`
	AdapterRegistryDigest string           `json:"adapter_registry_digest"`
	AdapterIdentityDigest string           `json:"adapter_identity_digest"`
	EnrollmentEventDigest string           `json:"enrollment_event_digest"`
	LaunchEventDigest     string           `json:"launch_event_digest"`
	ChainHead             string           `json:"chain_head"`
	EventCount            int              `json:"event_count"`
	FreshnessState        string           `json:"freshness_state"`
	ArtifactDigests       []ArtifactDigest `json:"artifact_digests"`
}

type AssessmentResult struct {
	SchemaVersion            string      `json:"schema_version"`
	SelectedProfile          string      `json:"selected_profile"`
	ManagedHarnessAssessment string      `json:"managed_harness_assessment"`
	TrustScope               string      `json:"trust_scope"`
	ManagedConditions        []Condition `json:"managed_conditions"`
	Reasons                  []string    `json:"reasons"`
	NextActions              []string    `json:"next_actions"`
}

type Condition struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	Reason     string `json:"reason"`
	NextAction string `json:"next_action,omitempty"`
}

func Evaluate(input Input) AssessmentResult {
	conditions := []Condition{
		pass("managed_profile_explicitly_selected", "managed_profile_selected", "managed harness profile was explicitly selected"),
		policyCondition(input.Policy),
		registryCondition(input.Registry),
		boundaryCondition(input),
		adapterIdentityCondition(input),
		capabilityCondition(input),
		adapterActivationCondition(input),
		adapterConnectionCondition(input),
		eventGroupCondition(input, "required_harness_events_observed", "harness"),
		eventGroupCondition(input, "required_tool_events_observed", "tool"),
		eventGroupCondition(input, "required_file_mutations_observed", "file"),
		testProvenanceCondition(input),
		suppressionCondition(input),
		bypassCondition(input),
		witnessCondition(input),
		overrideCondition(input),
	}
	result := AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileManagedHarness,
		ManagedHarnessAssessment: topLevel(conditions),
		TrustScope:               TrustScopeManaged,
		ManagedConditions:        conditions,
	}
	if result.ManagedHarnessAssessment != StatePass {
		result.TrustScope = "local_observed"
	}
	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}

func policyCondition(policy Policy) Condition {
	if policy.PolicyID == "" {
		return cannotVerify("managed_policy_loaded", "missing_managed_policy", "managed policy is required", "Supply a managed policy anchored before the run.")
	}
	if !preRunProvenance(policy.PolicyProvenance.Source) || policy.PolicyProvenance.Digest == "" {
		return fail("managed_policy_loaded", "post_hoc_policy", "managed policy provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent policy.")
	}
	return pass("managed_policy_loaded", "managed_policy_loaded", "managed policy is readable and anchored before the run")
}

func registryCondition(registry Registry) Condition {
	if registry.RegistryID == "" {
		return cannotVerify("adapter_registry_loaded", "missing_adapter_registry", "adapter registry is required", "Supply an adapter registry anchored before the run.")
	}
	if !preRunProvenance(registry.Provenance.Source) || registry.Provenance.Digest == "" {
		return fail("adapter_registry_loaded", "post_hoc_registry", "adapter registry provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent registry.")
	}
	return pass("adapter_registry_loaded", "adapter_registry_loaded", "adapter registry is readable and anchored before the run")
}

func boundaryCondition(input Input) Condition {
	boundary := input.Run.ManagedBoundaryEnrolled
	if boundary == nil {
		return fail("managed_boundary_enrolled_before_run", "run_not_managed", "selected run has no managed boundary enrollment event", "Run through the managed wrapper or lower the claim.")
	}
	if boundary.Sequence >= input.Run.ChildLaunch.Sequence || input.Run.ChildLaunch.Sequence == 0 {
		return fail("managed_boundary_enrolled_before_run", "late_enrollment", "managed boundary enrollment is not before child launch", "Enroll the managed boundary before child harness execution.")
	}
	if boundary.ManagedPolicyDigest != input.Policy.PolicyProvenance.Digest || boundary.AdapterRegistryDigest != input.Registry.Provenance.Digest || boundary.RunNonce != input.Run.RunNonce {
		return fail("managed_boundary_enrolled_before_run", "managed_boundary_not_in_chain", "managed boundary event does not bind the selected policy, registry, or run nonce", "Regenerate the run under the selected managed policy and registry.")
	}
	return pass("managed_boundary_enrolled_before_run", "managed_boundary_enrolled", "managed boundary enrollment is in chain before child launch")
}

func adapterIdentityCondition(input Input) Condition {
	adapter, ok := selectedAdapter(input)
	if !ok {
		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "selected adapter is not present in the registry", "Register an authorized adapter before the run.")
	}
	authorized := false
	for _, allowed := range input.Policy.AuthorizedAdapters {
		if allowed.AdapterID == adapter.AdapterID && allowed.HarnessID == adapter.HarnessID && allowed.AuthorityRef == adapter.AuthorityRef && allowed.DeploymentRef == adapter.DeploymentRef {
			authorized = true
			break
		}
	}
	if !authorized || adapter.IdentityState != IdentityVerified {
		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "adapter identity is not verified and authorized by policy", "Use a verified adapter identity authorized by managed policy.")
	}
	return pass("adapter_identity_authorized", "adapter_identity_authorized", "adapter identity is verified and authorized by policy")
}

func capabilityCondition(input Input) Condition {
	adapter, ok := selectedAdapter(input)
	if !ok {
		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "selected adapter capabilities cannot be verified", "Supply an authorized adapter with declared capabilities.")
	}
	return selectedAdapterCapabilityCondition(input, adapter)
}

func selectedAdapterCapabilityCondition(input Input, adapter Adapter) Condition {
	authorized, condition, ok := managedCapabilityPolicy(input, adapter)
	if !ok {
		return condition
	}
	if !adapterSatisfiesPolicyCapabilities(adapter, authorized) {
		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "adapter capability references do not satisfy managed policy", "Use an adapter whose declared capabilities match managed policy requirements.")
	}
	if !adapterCapabilitiesCoverEvents(input, adapter, authorized) {
		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "adapter capability set does not cover a required event type", "Use an adapter whose capabilities cover the managed contract.")
	}
	return pass("adapter_capabilities_satisfy_contract", "adapter_capabilities_satisfy_contract", "adapter capabilities cover required event types")
}

func managedCapabilityPolicy(input Input, adapter Adapter) (AuthorizedAdapter, Condition, bool) {
	authorized, ok := selectedAuthorizedAdapter(input, adapter)
	if !ok || len(authorized.CapabilityIDs) == 0 {
		return AuthorizedAdapter{}, cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "managed policy does not name required adapter capabilities", "Supply managed policy capability requirements for the selected adapter."), false
	}
	return authorized, Condition{}, true
}

func adapterSatisfiesPolicyCapabilities(adapter Adapter, authorized AuthorizedAdapter) bool {
	capabilityRefs := stringSet(adapter.CapabilityRefs)
	capabilityIDs := adapterCapabilityIDs(adapter)
	for _, capabilityID := range authorized.CapabilityIDs {
		if !capabilityDeclared(capabilityID, capabilityRefs, capabilityIDs) {
			return false
		}
	}
	return true
}

func adapterCapabilityIDs(adapter Adapter) map[string]bool {
	capabilityIDs := map[string]bool{}
	for _, capability := range adapter.Capabilities {
		capabilityIDs[capability.ID] = true
	}
	return capabilityIDs
}

func capabilityDeclared(capabilityID string, refs, ids map[string]bool) bool {
	return refs[capabilityID] && ids[capabilityID]
}

func adapterCapabilitiesCoverEvents(input Input, adapter Adapter, authorized AuthorizedAdapter) bool {
	capEvents := authorizedCapabilityEvents(adapter, authorized)
	for _, eventType := range requiredEventTypes(input) {
		if !capEvents[eventType] {
			return false
		}
	}
	return true
}

func authorizedCapabilityEvents(adapter Adapter, authorized AuthorizedAdapter) map[string]bool {
	authorizedCapabilityIDs := stringSet(authorized.CapabilityIDs)
	capEvents := map[string]bool{}
	for _, capability := range adapter.Capabilities {
		if !authorizedCapabilityIDs[capability.ID] {
			continue
		}
		for _, eventType := range capability.EventTypes {
			capEvents[eventType] = true
		}
	}
	return capEvents
}

func adapterActivationCondition(input Input) Condition {
	if input.Run.ManagedBoundaryEnrolled == nil || input.Run.ManagedBoundaryEnrolled.AdapterID == "" {
		return cannotVerify("adapter_activation_observed", "adapter_activation_missing", "adapter activation cannot be verified", "Record adapter activation before child launch.")
	}
	return pass("adapter_activation_observed", "adapter_activation_observed", "adapter activation is bound to managed enrollment")
}

func adapterConnectionCondition(input Input) Condition {
	if input.Run.AdapterDisconnectObserved {
		return fail("adapter_connection_continuous", "adapter_disconnect_observed", "adapter disconnected during required managed observation window", "Rerun with continuous managed adapter connection.")
	}
	return pass("adapter_connection_continuous", "adapter_connection_continuous", "adapter connection has no observed disconnect during required window")
}

func eventGroupCondition(input Input, id, group string) Condition {
	required := eventTypesForGroup(input, group)
	if len(required) == 0 {
		return pass(id, "condition_pass", "no required events for group")
	}
	scopes := acceptableScopesForGroup(input, group)
	for _, eventType := range required {
		if eventObserved(input.Run.ObservedEvents, eventType, scopes) {
			continue
		}
		return missingEventGroupCondition(input, id, group)
	}
	return pass(id, "condition_pass", "required "+group+" events are observed")
}

func missingEventGroupCondition(input Input, id, group string) Condition {
	reasonPrefix := groupReasonPrefix(group)
	if condition, ok := suppressedEventGroupCondition(input, id, group, reasonPrefix); ok {
		return condition
	}
	return Condition{ID: id, State: StateMissingTelemetry, ReasonCode: reasonPrefix + "_event_missing", Reason: "required " + group + " event is missing", NextAction: "Run through a managed boundary that emits required " + group + " events."}
}

func suppressedEventGroupCondition(input Input, id, group, reasonPrefix string) (Condition, bool) {
	suppressed, valid, satisfies := suppressionForGroup(input, group)
	if !suppressed {
		return Condition{}, false
	}
	if valid && satisfies {
		return pass(id, reasonPrefix+"_event_suppressed_by_policy", "required "+group+" event is suppressed by policy for this profile"), true
	}
	if valid {
		return Condition{ID: id, State: StateSuppressed, ReasonCode: reasonPrefix + "_event_suppressed", Reason: "required " + group + " event is suppressed but does not satisfy the managed profile", NextAction: "Capture the required " + group + " event or authorize satisfying suppression in pre-run policy."}, true
	}
	return Condition{}, false
}

func testProvenanceCondition(input Input) Condition {
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"local_observed", "ci_witnessed"}) {
		return pass("test_provenance_not_agent_reported", "test_provenance_not_agent_reported", "test evidence is wrapper or CI observed")
	}
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"agent_reported"}) {
		return fail("test_provenance_not_agent_reported", "agent_reported_test_not_executed", "test evidence is only agent-reported", "Record test execution through the managed wrapper or CI.")
	}
	return Condition{ID: "test_provenance_not_agent_reported", State: StateMissingTelemetry, ReasonCode: "test_provenance_missing", Reason: "test execution evidence is missing", NextAction: "Record test execution through the managed wrapper or CI."}
}

func suppressionCondition(input Input) Condition {
	for _, suppressed := range input.Run.SuppressedEventGroups {
		rule, ok := suppressionRuleForGroup(input.Policy, suppressed.EventGroup)
		if !ok || !suppressed.AuthorizedByPolicy || !preRunProvenance(suppressed.PolicyProvenanceSource) || !preRunProvenance(rule.PolicyProvenanceSource) {
			return fail("suppression_policy_valid", "suppression_unverified", "suppression is not authorized by pre-run policy provenance", "Use pre-run policy authority for suppression or capture the event.")
		}
	}
	return pass("suppression_policy_valid", "suppression_policy_valid", "suppression policy is valid or no suppression is present")
}

func bypassCondition(input Input) Condition {
	if input.Run.BypassObserved {
		return fail("bypass_not_observed", "bypass_observed", "managed boundary bypass was observed", "Rerun without bypass or lower the claim.")
	}
	return pass("bypass_not_observed", "bypass_not_observed", "no managed boundary bypass is observed")
}

func witnessCondition(input Input) Condition {
	witness := input.Witness
	if witness.WitnessID == "" {
		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness evidence is required", "Supply managed witness evidence bound to the run.")
	}
	if witness.Status != StatePass || witness.FreshnessState != StatePass {
		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness is missing pass/freshness state", "Supply fresh managed witness evidence.")
	}
	if missingWitnessArtifacts(input.Run, witness) {
		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness artifact binding is required", "Supply managed witness evidence with output artifact digests.")
	}
	if managedWitnessMismatches(input) {
		return fail("managed_witness_bound", "managed_witness_mismatch", "managed witness does not bind the selected run, policy, registry, chain, or artifacts", "Regenerate managed witness evidence for the selected run.")
	}
	return pass("managed_witness_bound", "managed_witness_bound", "managed witness binds source, run, policy, registry, chain, and artifacts")
}

func missingWitnessArtifacts(run RunEvidence, witness Witness) bool {
	return len(run.OutputArtifacts) == 0 || len(witness.ArtifactDigests) == 0
}

func managedWitnessMismatches(input Input) bool {
	boundary := input.Run.ManagedBoundaryEnrolled
	return boundary == nil ||
		!witnessMatchesRun(input.Witness, input.Run) ||
		!witnessMatchesAuthority(input.Witness, input.Policy, input.Registry) ||
		!witnessMatchesEvents(input.Witness, input.Run, *boundary) ||
		!artifactsMatch(input.Run.OutputArtifacts, input.Witness.ArtifactDigests)
}

func witnessMatchesRun(witness Witness, run RunEvidence) bool {
	return witness.RunID == run.RunID &&
		witness.RunNonce == run.RunNonce &&
		witness.SourceCommit == run.SourceCommit &&
		witness.ChainHead == run.ChainHead &&
		witness.EventCount == run.EventCount
}

func witnessMatchesAuthority(witness Witness, policy Policy, registry Registry) bool {
	return witness.ManagedPolicyDigest == policy.PolicyProvenance.Digest &&
		witness.AdapterRegistryDigest == registry.Provenance.Digest
}

func witnessMatchesEvents(witness Witness, run RunEvidence, boundary ManagedBoundaryEnrolled) bool {
	return witness.EnrollmentEventDigest == boundary.EventDigest &&
		witness.LaunchEventDigest == run.ChildLaunch.EventDigest
}

func overrideCondition(input Input) Condition {
	if input.Run.OverrideAttemptsTrustUpgrade {
		return fail("override_does_not_upgrade_managed_profile", "override_upgrade_rejected", "override artifact attempts to upgrade managed profile state", "Record override as non-upgrading evidence only.")
	}
	if input.Run.OverridePresent {
		return pass("override_does_not_upgrade_managed_profile", "override_present_non_upgrading", "override request is visible and non-upgrading")
	}
	return pass("override_does_not_upgrade_managed_profile", "override_absent_non_upgrading", "no override request is available to upgrade managed profile")
}

func selectedAdapter(input Input) (Adapter, bool) {
	if input.Run.ManagedBoundaryEnrolled == nil {
		return Adapter{}, false
	}
	for _, adapter := range input.Registry.Adapters {
		if adapter.AdapterID == input.Run.ManagedBoundaryEnrolled.AdapterID {
			return adapter, true
		}
	}
	return Adapter{}, false
}

func selectedAuthorizedAdapter(input Input, adapter Adapter) (AuthorizedAdapter, bool) {
	for _, allowed := range input.Policy.AuthorizedAdapters {
		if allowed.AdapterID == adapter.AdapterID && allowed.HarnessID == adapter.HarnessID && allowed.AuthorityRef == adapter.AuthorityRef && allowed.DeploymentRef == adapter.DeploymentRef {
			return allowed, true
		}
	}
	return AuthorizedAdapter{}, false
}

func requiredEventTypes(input Input) []string {
	if len(input.Contract.RequiredEventTypes) > 0 {
		return input.Contract.RequiredEventTypes
	}
	var out []string
	for _, group := range input.Policy.RequiredEventGroups {
		out = append(out, group.EventTypes...)
	}
	return out
}

func stringSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		out[value] = true
	}
	return out
}

func eventTypesForGroup(input Input, groupID string) []string {
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {
			return group.EventTypes
		}
	}
	return nil
}

func acceptableScopesForGroup(input Input, groupID string) []string {
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {
			return group.AcceptableProvenanceScopes
		}
	}
	return nil
}

func suppressionForGroup(input Input, groupID string) (bool, bool, bool) {
	for _, suppressed := range input.Run.SuppressedEventGroups {
		if suppressed.EventGroup != groupID {
			continue
		}
		rule, ok := suppressionRuleForGroup(input.Policy, groupID)
		if !ok || !suppressed.AuthorizedByPolicy || !preRunProvenance(suppressed.PolicyProvenanceSource) || !preRunProvenance(rule.PolicyProvenanceSource) {
			return true, false, false
		}
		return true, true, rule.SuppressionMaySatisfyProfile && suppressed.AuthorizedByPolicy
	}
	return false, false, false
}

func suppressionRuleForGroup(policy Policy, groupID string) (SuppressionRule, bool) {
	for _, rule := range policy.SuppressionRules {
		if rule.EventGroup == groupID {
			return rule, true
		}
	}
	return SuppressionRule{}, false
}

func groupReasonPrefix(group string) string {
	if group == "file" {
		return "file_mutation"
	}
	return group
}

func eventObserved(events []EvidenceEvent, eventType string, scopes []string) bool {
	scopeSet := map[string]bool{}
	for _, scope := range scopes {
		scopeSet[scope] = true
	}
	for _, event := range events {
		if event.EventType == eventType && (len(scopeSet) == 0 || scopeSet[event.ProvenanceScope]) {
			return true
		}
	}
	return false
}

func preRunProvenance(source string) bool {
	switch source {
	case "vcs", "ci_config", "human_signed", "customer_policy_equivalent":
		return true
	default:
		return false
	}
}

func artifactsMatch(expected, observed []ArtifactDigest) bool {
	if len(expected) == 0 || len(observed) == 0 {
		return false
	}
	want := map[string]string{}
	for _, artifact := range expected {
		want[artifact.Path] = artifact.SHA256
	}
	for _, artifact := range observed {
		if want[artifact.Path] != artifact.SHA256 {
			return false
		}
		delete(want, artifact.Path)
	}
	return len(want) == 0
}

func topLevel(conditions []Condition) string {
	state := StatePass
	for _, condition := range conditions {
		state = worse(state, condition.State)
	}
	return state
}

func worse(current, next string) string {
	if severity(next) > severity(current) {
		return topLevelState(next)
	}
	return topLevelState(current)
}

func topLevelState(state string) string {
	switch state {
	case StateMissingTelemetry, StateNotIntegrated, StateUnsupported, StateSuppressed:
		return StateCannotVerify
	default:
		return state
	}
}

func severity(state string) int {
	switch topLevelState(state) {
	case StateFail:
		return 4
	case StateCannotVerify:
		return 3
	case StateNotAssessed:
		return 2
	case StatePass:
		return 1
	default:
		return 0
	}
}

func reasons(conditions []Condition) []string {
	ordered := orderConditions(conditions)
	out := []string{}
	for _, condition := range ordered {
		if condition.State != StatePass {
			out = append(out, condition.ReasonCode+": "+condition.Reason)
		}
	}
	return out
}

func nextActions(conditions []Condition) []string {
	ordered := orderConditions(conditions)
	out := []string{}
	seen := map[string]bool{}
	for _, condition := range ordered {
		if condition.State == StatePass || condition.NextAction == "" || seen[condition.NextAction] {
			continue
		}
		seen[condition.NextAction] = true
		out = append(out, condition.NextAction)
	}
	return out
}

func orderConditions(conditions []Condition) []Condition {
	index := map[string]int{}
	for i, id := range managedConditionIDs {
		index[id] = i
	}
	ordered := append([]Condition(nil), conditions...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if severity(ordered[i].State) != severity(ordered[j].State) {
			return severity(ordered[i].State) > severity(ordered[j].State)
		}
		return index[ordered[i].ID] < index[ordered[j].ID]
	})
	return ordered
}

func pass(id, code, reason string) Condition {
	return Condition{ID: id, State: StatePass, ReasonCode: code, Reason: reason}
}

func fail(id, code, reason, next string) Condition {
	return Condition{ID: id, State: StateFail, ReasonCode: code, Reason: reason, NextAction: next}
}

func cannotVerify(id, code, reason, next string) Condition {
	return Condition{ID: id, State: StateCannotVerify, ReasonCode: code, Reason: reason, NextAction: next}
}
