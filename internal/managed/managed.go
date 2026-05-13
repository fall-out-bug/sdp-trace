package managed

import "sort"

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
	// Evaluation turns managed-adapter evidence into explicit condition states,
	// never into an opaque health score or implicit managed-mode approval.
	conditions := managedConditions(input)
	result := managedAssessmentResult(conditions)
	if result.ManagedHarnessAssessment != StatePass {

		result.TrustScope = "local_observed"
	}
	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}

func managedConditions(input Input) []Condition {
	// Conditions are grouped by authority, observation, and closure so missing
	// evidence stays separate from failed evidence.

	conditions := managedAuthorityConditions(input)
	conditions = append(conditions, managedObservationConditions(input)...)
	conditions = append(conditions, managedClosureConditions(input)...)
	return conditions
}

func managedAuthorityConditions(input Input) []Condition {
	// Authority conditions decide whether the selected adapter was permitted before
	// observation evidence can be trusted.

	return []Condition{
		pass("managed_profile_explicitly_selected", "managed_profile_selected", "managed harness profile was explicitly selected"),
		policyCondition(input.Policy),
		registryCondition(input.Registry),
		boundaryCondition(input),
		adapterIdentityCondition(input),
		capabilityCondition(input),
	}
}

func managedObservationConditions(input Input) []Condition {
	// Observation conditions check adapter activity, event coverage, suppression,
	// bypasses, and witness data without upgrading missing observations.

	return []Condition{
		adapterActivationCondition(input),
		adapterConnectionCondition(input),
		eventGroupCondition(input, "required_harness_events_observed", "harness"),
		eventGroupCondition(input, "required_tool_events_observed", "tool"),
		eventGroupCondition(input, "required_file_mutations_observed", "file"),
		testProvenanceCondition(input),
		suppressionCondition(input),
		bypassCondition(input),
	}
}

func managedClosureConditions(input Input) []Condition {
	// Closure conditions bind override and witness evidence after observation checks
	// have identified the adapter and event surface.

	return []Condition{
		witnessCondition(input),
		overrideCondition(input),
	}
}

func managedAssessmentResult(conditions []Condition) AssessmentResult {
	// Result assembly keeps condition states, reasons, and next actions tied to the
	// machine evidence emitted by each managed gate.

	return AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileManagedHarness,
		ManagedHarnessAssessment: topLevel(conditions),
		TrustScope:               TrustScopeManaged,
		ManagedConditions:        conditions,
	}
}
func policyCondition(policy Policy) Condition {
	// Policy evidence is required before adapter authorization can be considered
	// managed rather than local observation.
	if policy.PolicyID == "" {

		return cannotVerify("managed_policy_loaded", "missing_managed_policy", "managed policy is required", "Supply a managed policy anchored before the run.")
	}
	if !preRunProvenance(policy.PolicyProvenance.Source) || policy.PolicyProvenance.Digest == "" {

		return fail("managed_policy_loaded", "post_hoc_policy", "managed policy provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent policy.")
	}
	return pass("managed_policy_loaded", "managed_policy_loaded", "managed policy is readable and anchored before the run")
}

func registryCondition(registry Registry) Condition {
	// Registry evidence names whether the adapter exists in the declared registry
	// instead of treating any selected adapter as trusted.
	if registry.RegistryID == "" {

		return cannotVerify("adapter_registry_loaded", "missing_adapter_registry", "adapter registry is required", "Supply an adapter registry anchored before the run.")
	}
	if !preRunProvenance(registry.Provenance.Source) || registry.Provenance.Digest == "" {

		return fail("adapter_registry_loaded", "post_hoc_registry", "adapter registry provenance is not anchored before the run", "Use a VCS, CI, human-signed, or customer policy equivalent registry.")
	}
	return pass("adapter_registry_loaded", "adapter_registry_loaded", "adapter registry is readable and anchored before the run")
}

func boundaryCondition(input Input) Condition {
	// Boundary evidence keeps selected adapter identity and task boundary checks
	// explicit before capability or event coverage is evaluated.
	boundary := input.Run.ManagedBoundaryEnrolled
	if boundary == nil {

		return fail("managed_boundary_enrolled_before_run", "run_not_managed", "selected run has no managed boundary enrollment event", "Run through the managed wrapper or lower the claim.")
	}
	if boundaryNotBeforeLaunch(*boundary, input.Run.ChildLaunch) {

		return fail("managed_boundary_enrolled_before_run", "late_enrollment", "managed boundary enrollment is not before child launch", "Enroll the managed boundary before child harness execution.")
	}
	if boundaryBindingMismatch(*boundary, input) {

		return fail("managed_boundary_enrolled_before_run", "managed_boundary_not_in_chain", "managed boundary event does not bind the selected policy, registry, or run nonce", "Regenerate the run under the selected managed policy and registry.")
	}
	return pass("managed_boundary_enrolled_before_run", "managed_boundary_enrolled", "managed boundary enrollment is in chain before child launch")
}

func boundaryNotBeforeLaunch(boundary ManagedBoundaryEnrolled, launch LaunchEvent) bool {
	return boundary.Sequence >= launch.Sequence || launch.Sequence == 0
}

func boundaryBindingMismatch(boundary ManagedBoundaryEnrolled, input Input) bool {
	return boundary.ManagedPolicyDigest != input.Policy.PolicyProvenance.Digest || boundary.AdapterRegistryDigest != input.Registry.Provenance.Digest || boundary.RunNonce != input.Run.RunNonce
}

func adapterIdentityCondition(input Input) Condition {
	// adapterIdentityCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	adapter, ok := selectedAdapter(input)
	if !ok {

		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "selected adapter is not present in the registry", "Register an authorized adapter before the run.")
	}
	if !adapterIdentityAuthorized(input, adapter) {
		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "adapter identity is not verified and authorized by policy", "Use a verified adapter identity authorized by managed policy.")
	}
	return pass("adapter_identity_authorized", "adapter_identity_authorized", "adapter identity is verified and authorized by policy")
}

func adapterIdentityAuthorized(input Input, adapter Adapter) bool {
	// adapterIdentityAuthorized preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if adapter.IdentityState != IdentityVerified {

		return false
	}
	_, ok := selectedAuthorizedAdapter(input, adapter)
	return ok
}

func capabilityCondition(input Input) Condition {
	// Capability checks compare policy requirements to selected adapter claims
	// without assuming adapter self-description is sufficient proof.
	adapter, ok := selectedAdapter(input)
	if !ok {

		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "selected adapter capabilities cannot be verified", "Supply an authorized adapter with declared capabilities.")
	}
	return selectedAdapterCapabilityCondition(input, adapter)
}

func selectedAdapterCapabilityCondition(input Input, adapter Adapter) Condition {
	// selectedAdapterCapabilityCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
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
	// managedCapabilityPolicy preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	authorized, ok := selectedAuthorizedAdapter(input, adapter)
	if !ok || len(authorized.CapabilityIDs) == 0 {

		return AuthorizedAdapter{}, cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "managed policy does not name required adapter capabilities", "Supply managed policy capability requirements for the selected adapter."), false
	}
	return authorized, Condition{}, true
}

func adapterSatisfiesPolicyCapabilities(adapter Adapter, authorized AuthorizedAdapter) bool {
	// adapterSatisfiesPolicyCapabilities preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
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
	// adapterCapabilityIDs preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
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
	// adapterCapabilitiesCoverEvents preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	capEvents := authorizedCapabilityEvents(adapter, authorized)
	for _, eventType := range requiredEventTypes(input) {
		if !capEvents[eventType] {

			return false
		}
	}
	return true
}

func authorizedCapabilityEvents(adapter Adapter, authorized AuthorizedAdapter) map[string]bool {
	// authorizedCapabilityEvents preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
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
	// adapterActivationCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if input.Run.ManagedBoundaryEnrolled == nil || input.Run.ManagedBoundaryEnrolled.AdapterID == "" {

		return cannotVerify("adapter_activation_observed", "adapter_activation_missing", "adapter activation cannot be verified", "Record adapter activation before child launch.")
	}
	return pass("adapter_activation_observed", "adapter_activation_observed", "adapter activation is bound to managed enrollment")
}

func adapterConnectionCondition(input Input) Condition {
	// adapterConnectionCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if input.Run.AdapterDisconnectObserved {
		return fail("adapter_connection_continuous", "adapter_disconnect_observed", "adapter disconnected during required managed observation window", "Rerun with continuous managed adapter connection.")
	}

	return pass("adapter_connection_continuous", "adapter_connection_continuous", "adapter connection has no observed disconnect during required window")
}

func eventGroupCondition(input Input, id, group string) Condition {
	// Event group checks preserve missing, suppressed, and observed event states as
	// separate managed-mode evidence outcomes.
	required := eventTypesForGroup(input, group)
	if len(required) == 0 {

		return pass(id, "condition_pass", "no required events for group")
	}
	scopes := acceptableScopesForGroup(input, group)
	if !allEventsObserved(input.Run.ObservedEvents, required, scopes) {
		return missingEventGroupCondition(input, id, group)
	}
	return pass(id, "condition_pass", "required "+group+" events are observed")
}

func allEventsObserved(events []EvidenceEvent, required, scopes []string) bool {
	// allEventsObserved preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, eventType := range required {
		if !eventObserved(events, eventType, scopes) {

			return false
		}
	}
	return true
}

func missingEventGroupCondition(input Input, id, group string) Condition {
	// missingEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	reasonPrefix := group
	if group == "file" {

		reasonPrefix = "file_mutation"
	}
	if condition, ok := suppressedEventGroupCondition(input, id, group, reasonPrefix); ok {
		return condition
	}
	return Condition{ID: id, State: StateMissingTelemetry, ReasonCode: reasonPrefix + "_event_missing", Reason: "required " + group + " event is missing", NextAction: "Run through a managed boundary that emits required " + group + " events."}
}

func suppressedEventGroupCondition(input Input, id, group, reasonPrefix string) (Condition, bool) {
	// suppressedEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	suppressed, valid, satisfies := suppressionForGroup(input, group)
	if !suppressed {

		return Condition{}, false
	}
	return validSuppressedEventGroupCondition(id, group, reasonPrefix, valid, satisfies)
}

func validSuppressedEventGroupCondition(id, group, reasonPrefix string, valid, satisfies bool) (Condition, bool) {
	// validSuppressedEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if valid && satisfies {

		return pass(id, reasonPrefix+"_event_suppressed_by_policy", "required "+group+" event is suppressed by policy for this profile"), true
	}
	if valid {

		return Condition{ID: id, State: StateSuppressed, ReasonCode: reasonPrefix + "_event_suppressed", Reason: "required " + group + " event is suppressed but does not satisfy the managed profile", NextAction: "Capture the required " + group + " event or authorize satisfying suppression in pre-run policy."}, true
	}
	return Condition{}, false
}

func testProvenanceCondition(input Input) Condition {
	// Test provenance evidence is evaluated independently because observed events
	// do not prove tests actually ran.
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"local_observed", "ci_witnessed"}) {

		return pass("test_provenance_not_agent_reported", "test_provenance_not_agent_reported", "test evidence is wrapper or CI observed")
	}
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"agent_reported"}) {

		return fail("test_provenance_not_agent_reported", "agent_reported_test_not_executed", "test evidence is only agent-reported", "Record test execution through the managed wrapper or CI.")
	}
	return Condition{ID: "test_provenance_not_agent_reported", State: StateMissingTelemetry, ReasonCode: "test_provenance_missing", Reason: "test execution evidence is missing", NextAction: "Record test execution through the managed wrapper or CI."}
}

func suppressionCondition(input Input) Condition {
	// suppressionCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, suppressed := range input.Run.SuppressedEventGroups {
		if !suppressionVerified(input.Policy, suppressed) {

			return fail("suppression_policy_valid", "suppression_unverified", "suppression is not authorized by pre-run policy provenance", "Use pre-run policy authority for suppression or capture the event.")
		}
	}
	return pass("suppression_policy_valid", "suppression_policy_valid", "suppression policy is valid or no suppression is present")
}

func suppressionVerified(policy Policy, suppressed SuppressedEventGroup) bool {
	// suppressionVerified preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	rule, ok := suppressionRuleForGroup(policy, suppressed.EventGroup)

	return ok &&
		suppressed.AuthorizedByPolicy &&
		preRunProvenance(suppressed.PolicyProvenanceSource) &&
		preRunProvenance(rule.PolicyProvenanceSource)
}
func bypassCondition(input Input) Condition {
	// Bypass evidence remains explicit so an intentional bypass cannot look like
	// a passing managed-adapter observation.
	if input.Run.BypassObserved {
		return fail("bypass_not_observed", "bypass_observed", "managed boundary bypass was observed", "Rerun without bypass or lower the claim.")
	}

	return pass("bypass_not_observed", "bypass_not_observed", "no managed boundary bypass is observed")
}
func witnessCondition(input Input) Condition {
	// Witness evidence binds managed output back to run/report artifacts rather
	// than trusting a checked-in witness record by itself.
	witness := input.Witness
	if condition, ok := missingManagedWitnessCondition(input.Run, witness); ok {
		return condition
	}
	if managedWitnessMismatches(input) {

		return fail("managed_witness_bound", "managed_witness_mismatch", "managed witness does not bind the selected run, policy, registry, chain, or artifacts", "Regenerate managed witness evidence for the selected run.")
	}
	return pass("managed_witness_bound", "managed_witness_bound", "managed witness binds source, run, policy, registry, chain, and artifacts")
}

func missingManagedWitnessCondition(run RunEvidence, witness Witness) (Condition, bool) {
	// missingManagedWitnessCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if witness.WitnessID == "" {

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness evidence is required", "Supply managed witness evidence bound to the run."), true
	}
	return invalidManagedWitnessCondition(run, witness)
}

func invalidManagedWitnessCondition(run RunEvidence, witness Witness) (Condition, bool) {
	// invalidManagedWitnessCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch {
	case missingManagedWitnessPassState(witness):

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness is missing pass/freshness state", "Supply fresh managed witness evidence."), true
	case missingWitnessArtifacts(run, witness):

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness artifact binding is required", "Supply managed witness evidence with output artifact digests."), true
	default:
		return Condition{}, false
	}
}

func missingManagedWitnessPassState(witness Witness) bool {
	return witness.Status != StatePass || witness.FreshnessState != StatePass
}

func missingWitnessArtifacts(run RunEvidence, witness Witness) bool {
	return len(run.OutputArtifacts) == 0 || len(witness.ArtifactDigests) == 0
}

func managedWitnessMismatches(input Input) bool {
	// managedWitnessMismatches preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	boundary := input.Run.ManagedBoundaryEnrolled
	if boundary == nil {

		return true
	}
	return managedWitnessBindingMismatch(input, *boundary)
}

func managedWitnessBindingMismatch(input Input, boundary ManagedBoundaryEnrolled) bool {
	// managedWitnessBindingMismatch preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return !witnessMatchesRun(input.Witness, input.Run) ||
		!witnessMatchesAuthority(input.Witness, input.Policy, input.Registry) ||
		!witnessMatchesEvents(input.Witness, input.Run, boundary) ||
		!artifactsMatch(input.Run.OutputArtifacts, input.Witness.ArtifactDigests)
}

func witnessMatchesRun(witness Witness, run RunEvidence) bool {
	return witnessRunIdentityMatches(witness, run) &&
		witnessRunTraceMatches(witness, run)
}

func witnessRunIdentityMatches(witness Witness, run RunEvidence) bool {
	return witness.RunID == run.RunID && witness.RunNonce == run.RunNonce && witness.SourceCommit == run.SourceCommit
}

func witnessRunTraceMatches(witness Witness, run RunEvidence) bool {
	return witness.ChainHead == run.ChainHead &&
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
	// Override evidence is terminal closure context; it does not erase earlier
	// adapter, event, or witness failures.
	if input.Run.OverrideAttemptsTrustUpgrade {

		return fail("override_does_not_upgrade_managed_profile", "override_upgrade_rejected", "override artifact attempts to upgrade managed profile state", "Record override as non-upgrading evidence only.")
	}
	if input.Run.OverridePresent {
		return pass("override_does_not_upgrade_managed_profile", "override_present_non_upgrading", "override request is visible and non-upgrading")
	}
	return pass("override_does_not_upgrade_managed_profile", "override_absent_non_upgrading", "no override request is available to upgrade managed profile")
}

func selectedAdapter(input Input) (Adapter, bool) {
	// selectedAdapter preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
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
	// selectedAuthorizedAdapter preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, allowed := range input.Policy.AuthorizedAdapters {
		if authorizedAdapterMatches(allowed, adapter) {

			return allowed, true
		}
	}
	return AuthorizedAdapter{}, false
}

func authorizedAdapterMatches(allowed AuthorizedAdapter, adapter Adapter) bool {
	// authorizedAdapterMatches preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return allowed.AdapterID == adapter.AdapterID &&
		allowed.HarnessID == adapter.HarnessID &&
		allowed.AuthorityRef == adapter.AuthorityRef &&
		allowed.DeploymentRef == adapter.DeploymentRef
}
func requiredEventTypes(input Input) []string {
	// requiredEventTypes preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if len(input.Contract.RequiredEventTypes) > 0 {

		return input.Contract.RequiredEventTypes
	}
	return policyRequiredEventTypes(input.Policy)
}

func policyRequiredEventTypes(policy Policy) []string {
	// policyRequiredEventTypes preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	var out []string
	for _, group := range policy.RequiredEventGroups {

		out = append(out, group.EventTypes...)
	}
	return out
}

func stringSet(values []string) map[string]bool {
	// stringSet preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	out := map[string]bool{}
	for _, value := range values {

		out[value] = true
	}
	return out
}

func eventTypesForGroup(input Input, groupID string) []string {
	// eventTypesForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {

			return group.EventTypes
		}
	}
	return nil
}

func acceptableScopesForGroup(input Input, groupID string) []string {
	// acceptableScopesForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, group := range input.Policy.RequiredEventGroups {
		if group.ID == groupID {

			return group.AcceptableProvenanceScopes
		}
	}
	return nil
}
func suppressionForGroup(input Input, groupID string) (bool, bool, bool) {
	// suppressionForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	suppressed, ok := selectedSuppressedEventGroup(input.Run.SuppressedEventGroups, groupID)
	if !ok {
		return false, false, false
	}
	rule, ok := verifiedSuppressionRule(input.Policy, suppressed)
	if !ok {

		return true, false, false
	}

	return true, true, rule.SuppressionMaySatisfyProfile && suppressed.AuthorizedByPolicy
}

func selectedSuppressedEventGroup(groups []SuppressedEventGroup, groupID string) (SuppressedEventGroup, bool) {
	// selectedSuppressedEventGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, suppressed := range groups {
		if suppressed.EventGroup == groupID {

			return suppressed, true
		}
	}
	return SuppressedEventGroup{}, false
}

func verifiedSuppressionRule(policy Policy, suppressed SuppressedEventGroup) (SuppressionRule, bool) {
	// verifiedSuppressionRule preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	rule, ok := suppressionRuleForGroup(policy, suppressed.EventGroup)
	if !ok || !suppressionVerified(policy, suppressed) {

		return SuppressionRule{}, false
	}
	return rule, true
}

func suppressionRuleForGroup(policy Policy, groupID string) (SuppressionRule, bool) {
	// suppressionRuleForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, rule := range policy.SuppressionRules {
		if rule.EventGroup == groupID {

			return rule, true
		}
	}
	return SuppressionRule{}, false
}

func eventObserved(events []EvidenceEvent, eventType string, scopes []string) bool {
	// eventObserved preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	scopeSet := scopeSet(scopes)
	for _, event := range events {
		if eventObservedInScope(event, eventType, scopeSet) {

			return true
		}
	}
	return false
}

func scopeSet(scopes []string) map[string]bool {
	// scopeSet preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	scopeSet := map[string]bool{}
	for _, scope := range scopes {

		scopeSet[scope] = true
	}
	return scopeSet
}

func eventObservedInScope(event EvidenceEvent, eventType string, scopes map[string]bool) bool {
	return event.EventType == eventType && (len(scopes) == 0 || scopes[event.ProvenanceScope])
}

func preRunProvenance(source string) bool {
	// preRunProvenance preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch source {
	case "vcs", "ci_config", "human_signed", "customer_policy_equivalent":

		return true
	default:
		return false
	}
}

func artifactsMatch(expected, observed []ArtifactDigest) bool {
	// artifactsMatch preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if len(expected) == 0 || len(observed) == 0 {

		return false
	}
	want := artifactDigestsByPath(expected)
	if !consumeMatchingArtifacts(want, observed) {
		return false
	}
	return len(want) == 0
}
func consumeMatchingArtifacts(want map[string]string, observed []ArtifactDigest) bool {
	// consumeMatchingArtifacts preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, artifact := range observed {
		if want[artifact.Path] != artifact.SHA256 {

			return false
		}
		delete(want, artifact.Path)
	}
	return true
}

func artifactDigestsByPath(artifacts []ArtifactDigest) map[string]string {
	// artifactDigestsByPath preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	want := map[string]string{}
	for _, artifact := range artifacts {

		want[artifact.Path] = artifact.SHA256
	}
	return want
}

func topLevel(conditions []Condition) string {
	// Top-level state reports the strongest non-pass condition without averaging
	// managed adapter evidence.
	state := StatePass
	for _, condition := range conditions {

		state = worse(state, condition.State)
	}
	return state
}

func worse(current, next string) string {
	// worse preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if severity(next) > severity(current) {

		return topLevelState(next)
	}
	return topLevelState(current)
}

func topLevelState(state string) string {
	// topLevelState preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch state {
	case StateMissingTelemetry, StateNotIntegrated, StateUnsupported, StateSuppressed:

		return StateCannotVerify
	default:
		return state
	}
}

func severity(state string) int {
	return managedSeverityByState(topLevelState(state))
}

func managedSeverityByState(state string) int {
	// managedSeverityByState preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return map[string]int{
		StateFail:         4,
		StateCannotVerify: 3,
		StateNotAssessed:  2,
		StatePass:         1,
	}[state]
}

func reasons(conditions []Condition) []string {
	// Reasons are derived from condition IDs and states so prose does not become
	// independent authority.
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
	// Next actions point at concrete missing evidence boundaries needed for a
	// replayable managed-mode verdict.

	ordered := orderConditions(conditions)
	seen := map[string]bool{}
	return collectNextActions(ordered, seen)
}
func collectNextActions(ordered []Condition, seen map[string]bool) []string {
	// collectNextActions preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	out := []string{}
	for _, condition := range ordered {
		if skipNextAction(condition, seen) {
			continue
		}

		seen[condition.NextAction] = true
		out = append(out, condition.NextAction)
	}
	return out
}

func skipNextAction(condition Condition, seen map[string]bool) bool {
	return condition.State == StatePass || condition.NextAction == "" || seen[condition.NextAction]
}

func orderConditions(conditions []Condition) []Condition {
	// orderConditions preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	index := conditionOrderIndex()
	ordered := append([]Condition(nil), conditions...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return conditionLess(ordered, index, i, j)
	})
	return ordered
}

func conditionOrderIndex() map[string]int {
	// conditionOrderIndex preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	index := map[string]int{}
	for i, id := range managedConditionIDs {

		index[id] = i
	}
	return index
}

func conditionLess(ordered []Condition, index map[string]int, i, j int) bool {
	// conditionLess preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if severity(ordered[i].State) != severity(ordered[j].State) {

		return severity(ordered[i].State) > severity(ordered[j].State)
	}
	return index[ordered[i].ID] < index[ordered[j].ID]
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
