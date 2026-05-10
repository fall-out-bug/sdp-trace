package forensic

import "sort"

const (
	SchemaVersion = "block18-forensic-retention-assessment-v1"

	ProfileForensicRetention = "forensic_retention"
	TrustScopeForensic       = "forensic_retention_observed"
	TrustScopeLocalObserved  = "local_observed"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StateNotAssessed  = "not_assessed"

	RetentionModeDigestOnly          = "digest_only"
	RetentionModeSanitizedExcerpt    = "sanitized_excerpt"
	RetentionModeEncryptedRawRef     = "encrypted_raw_ref"
	RetentionModeExternalArtifactRef = "external_artifact_ref"
	RetentionModeNotAssessed         = "not_assessed"

	RedactionActionApplyRule       = "apply_rule"
	RedactionActionWithhold        = "withhold"
	RedactionActionMarkUnavailable = "mark_unavailable"

	AuthorityVerified    = "verified"
	AuthoritySelfClaimed = "self_claimed"
	AuthorityNotAssessed = "not_assessed"

	AccessStateVerifiedAvailable = "verified_available"
	AccessStateRestricted        = "restricted"
	AccessStateUnavailable       = "unavailable"
	AccessStateRevoked           = "revoked"
	AccessStateNotAssessed       = "not_assessed"

	KeyCustodyNotApplicable = "not_applicable"
	KeyCustodyHolderKnown   = "holder_known"
	KeyCustodyEscrowed      = "escrowed"
	KeyCustodyDestroyed     = "destroyed"
	KeyCustodyCompromised   = "compromised"
	KeyCustodyUnknown       = "unknown"
	KeyCustodyNotAssessed   = "not_assessed"

	RetentionLifecycleActive      = "active"
	RetentionLifecycleExpired     = "expired"
	RetentionLifecycleRevoked     = "revoked"
	RetentionLifecycleDeleted     = "deleted"
	RetentionLifecycleRotated     = "rotated"
	RetentionLifecycleNotAssessed = "not_assessed"

	UnavailableReasonNotAssessed        = "not_assessed"
	UnavailableReasonMissingReference   = "missing_reference"
	UnavailableReasonAccessDenied       = "access_denied"
	UnavailableReasonExpired            = "expired"
	UnavailableReasonKeyUnavailable     = "key_unavailable"
	UnavailableReasonStoreUnreachable   = "store_unreachable"
	UnavailableReasonDigestUnverifiable = "digest_unverifiable"
)

var forensicConditionIDs = []string{
	"redaction_policy_bound",
	"redaction_prewrite_applied",
	"redaction_unresolved_visible",
	"retention_mode_declared",
	"critical_evidence_reconstructable",
	"raw_reference_bound",
	"forensic_profile_not_overclaimed",
	"profile_selection_accountable",
}

type Input struct {
	Policy Policy
	Run    RunEvidence
}

type Policy struct {
	PolicyID                      string                 `json:"policy_id"`
	SchemaVersion                 string                 `json:"schema_version"`
	PolicyDigest                  string                 `json:"policy_digest"`
	PolicyProvenance              Provenance             `json:"policy_provenance"`
	AllowedRetentionModes         []string               `json:"allowed_retention_modes"`
	RedactionActions              []string               `json:"redaction_actions"`
	ForbiddenPersistenceClasses   []string               `json:"forbidden_committed_persistence_classes"`
	CriticalEventFamilies         []string               `json:"critical_event_families"`
	NonCriticalEventFamilyReasons []CriticalityDowngrade `json:"non_critical_event_family_reasons,omitempty"`
	WithholdingRequiresAuthority  bool                   `json:"withholding_requires_authority,omitempty"`
	Authority                     AuthorityRef           `json:"authority"`
	ProfileMappings               []ProfileMapping       `json:"profile_mappings"`
	UnresolvedRedactionImpact     string                 `json:"unresolved_redaction_impact"`
	Rules                         []Rule                 `json:"rules,omitempty"`
}

type Provenance struct {
	Source string `json:"source"`
	Digest string `json:"digest"`
}

type CriticalityDowngrade struct {
	EventType   string `json:"event_type"`
	Reason      string `json:"reason"`
	AuthorityID string `json:"authority_id"`
}

type Rule struct {
	RuleID         string `json:"rule_id"`
	DetectorFamily string `json:"detector_family"`
	RuleVersion    string `json:"rule_version"`
	Action         string `json:"action"`
	RetentionMode  string `json:"retention_mode,omitempty"`
}

type ProfileMapping struct {
	EventFamily            string       `json:"event_family"`
	RequiredRetentionModes []string     `json:"required_retention_modes"`
	Critical               bool         `json:"critical,omitempty"`
	DowngradeReason        string       `json:"downgrade_reason,omitempty"`
	Authority              AuthorityRef `json:"authority,omitempty"`
}

type RunEvidence struct {
	RunID                 string           `json:"run_id"`
	SelectedProfile       string           `json:"selected_profile"`
	RedactionPolicyDigest string           `json:"redaction_policy_digest"`
	ProfileSelection      ProfileSelection `json:"profile_selection"`
	Events                []EventRetention `json:"events"`
}

type ProfileSelection struct {
	ActorID                 string `json:"actor_id"`
	SelectedProfile         string `json:"selected_profile"`
	RedactionPolicyDigest   string `json:"redaction_policy_digest"`
	Justification           string `json:"justification"`
	AuthorityVerified       bool   `json:"authority_verified"`
	SelectionEvidenceDigest string `json:"selection_evidence_digest"`
}

type EventRetention struct {
	EventType              string        `json:"event_type"`
	RetentionMode          string        `json:"retention_mode"`
	ForensicImportance     string        `json:"forensic_importance"`
	RedactionPolicyDigest  string        `json:"redaction_policy_digest"`
	RedactionInputDigest   string        `json:"redaction_input_digest"`
	RedactedPayloadDigest  string        `json:"redacted_payload_digest"`
	RedactionAction        string        `json:"redaction_action"`
	RedactionRuleRefs      []string      `json:"redaction_rule_refs,omitempty"`
	RedactionUnresolved    bool          `json:"redaction_unresolved,omitempty"`
	SecretLikeValuePresent bool          `json:"secret_like_value_present,omitempty"`
	RedactionAuthority     AuthorityRef  `json:"redaction_authority"`
	Withholding            *Withholding  `json:"withholding,omitempty"`
	RawReference           *RawReference `json:"raw_reference,omitempty"`
}

type AuthorityRef struct {
	ActorID           string `json:"actor_id"`
	VerificationState string `json:"verification_state"`
}

type Withholding struct {
	Authority     AuthorityRef `json:"authority"`
	Requestor     AuthorityRef `json:"requestor,omitempty"`
	ReasonCode    string       `json:"reason_code"`
	Justification string       `json:"justification"`
}

type RawReference struct {
	ReferenceType           string             `json:"reference_type"`
	ReferenceURI            string             `json:"reference_uri"`
	Digest                  Digest             `json:"digest"`
	AccessState             string             `json:"access_state"`
	AccessStateLastVerified string             `json:"access_state_last_verified"`
	KeyCustodyState         string             `json:"key_custody_state"`
	RetentionLifecycle      RetentionLifecycle `json:"retention_lifecycle"`
	UnavailableReason       string             `json:"unavailable_reason"`
}

type Digest struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

type RetentionLifecycle struct {
	State     string `json:"state"`
	PolicyRef string `json:"policy_ref"`
	ExpiresAt string `json:"expires_at"`
}

type AssessmentResult struct {
	SchemaVersion               string      `json:"schema_version"`
	SelectedProfile             string      `json:"selected_profile"`
	ForensicRetentionAssessment string      `json:"forensic_retention_assessment"`
	TrustScope                  string      `json:"trust_scope"`
	ForensicConditions          []Condition `json:"forensic_conditions"`
	Reasons                     []string    `json:"reasons"`
	NextActions                 []string    `json:"next_actions"`
}

type Condition struct {
	ID                    string `json:"id"`
	State                 string `json:"state"`
	ReasonCode            string `json:"reason_code"`
	Reason                string `json:"reason"`
	NextAction            string `json:"next_action,omitempty"`
	CappedToRetentionMode string `json:"capped_to_retention_mode,omitempty"`
}

func Evaluate(input Input) AssessmentResult {
	conditions := []Condition{
		policyCondition(input),
		prewriteCondition(input),
		unresolvedCondition(input),
		retentionModeCondition(input),
		criticalEvidenceCondition(input),
		rawReferenceCondition(input),
		overclaimCondition(input),
		profileSelectionCondition(input),
	}
	result := AssessmentResult{
		SchemaVersion:               SchemaVersion,
		SelectedProfile:             ProfileForensicRetention,
		ForensicRetentionAssessment: topLevel(conditions),
		TrustScope:                  TrustScopeForensic,
		ForensicConditions:          conditions,
	}
	if result.ForensicRetentionAssessment != StatePass {
		result.TrustScope = TrustScopeLocalObserved
	}
	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}

func policyCondition(input Input) Condition {
	if input.Policy.PolicyID == "" || input.Policy.PolicyDigest == "" {
		return cannotVerify("redaction_policy_bound", "missing_redaction_policy", "redaction policy is required", "Supply a redaction policy with stable id, version, digest, and provenance.")
	}
	if len(input.Policy.RedactionActions) == 0 || len(input.Policy.ForbiddenPersistenceClasses) == 0 || input.Policy.Authority.ActorID == "" || len(input.Policy.ProfileMappings) == 0 || input.Policy.UnresolvedRedactionImpact == "" {
		return cannotVerify("redaction_policy_bound", "redaction_policy_contract_incomplete", "redaction policy contract is incomplete", "Supply redaction actions, forbidden persistence classes, authority, profile mappings, and unresolved-redaction impact.")
	}
	if input.Policy.Authority.VerificationState == AuthoritySelfClaimed {
		return cannotVerify("redaction_policy_bound", "authority_self_claimed", "redaction policy authority is self-claimed", "Use a provenance or accountability-bound redaction policy authority.")
	}
	if input.Run.RedactionPolicyDigest != input.Policy.PolicyDigest {
		return fail("redaction_policy_bound", "redaction_policy_mismatch", "run evidence is not bound to the selected redaction policy", "Regenerate or select evidence bound to the redaction policy.")
	}
	for _, event := range input.Run.Events {
		if event.RedactionPolicyDigest != "" && event.RedactionPolicyDigest != input.Policy.PolicyDigest {
			return fail("redaction_policy_bound", "redaction_policy_mismatch", "event evidence contradicts the selected redaction policy digest", "Use a run whose event redaction policy digests match.")
		}
		if event.RedactionAuthority.VerificationState == AuthoritySelfClaimed {
			return cannotVerify("redaction_policy_bound", "authority_self_claimed", "redaction authority is self-claimed", "Use a provenance or accountability-bound redaction authority.")
		}
	}
	return pass("redaction_policy_bound", "redaction_policy_bound", "redaction policy digest and authority evidence are bound")
}

func prewriteCondition(input Input) Condition {
	if input.Policy.PolicyID == "" || input.Policy.PolicyDigest == "" {
		return cannotVerify("redaction_prewrite_applied", "redaction_policy_missing", "redaction rule coverage cannot be checked without the selected policy", "Supply the selected redaction policy before assessing rule coverage.")
	}
	rules := policyRules(input.Policy)
	for _, event := range input.Run.Events {
		if event.SecretLikeValuePresent {
			return fail("redaction_prewrite_applied", "secret_like_value_persisted", "secret-like value is marked as persisted in retained metadata", "Apply pre-write redaction and retain only digests or safe references.")
		}
		if event.RedactionInputDigest == "" || event.RedactedPayloadDigest == "" {
			return cannotVerify("redaction_prewrite_applied", "redaction_digest_missing", "pre-write redaction digests are missing", "Record pre-redaction and redacted payload digests.")
		}
		if event.RedactionAction == RedactionActionApplyRule && len(event.RedactionRuleRefs) == 0 {
			return cannotVerify("redaction_prewrite_applied", "redaction_rule_refs_missing", "redaction rule references are missing", "Record the redaction rule ids applied before persistence.")
		}
		for _, ruleRef := range event.RedactionRuleRefs {
			rule, ok := rules[ruleRef]
			if !ok {
				return fail("redaction_prewrite_applied", "redaction_rule_unknown", "event references a redaction rule that is absent from the selected policy", "Use event redaction_rule_refs from the selected redaction policy.")
			}
			if rule.Action != "" && rule.Action != event.RedactionAction {
				return fail("redaction_prewrite_applied", "redaction_rule_action_mismatch", "event redaction action contradicts the selected policy rule", "Align event redaction action with the selected policy rule.")
			}
		}
	}
	return pass("redaction_prewrite_applied", "redaction_prewrite_applied", "pre-write redaction metadata is verifier-readable")
}

func unresolvedCondition(input Input) Condition {
	for _, event := range input.Run.Events {
		if event.RedactionUnresolved {
			return fail("redaction_unresolved_visible", "redaction_unresolved", "unresolved redaction is visible and blocks forensic retention", "Resolve redaction or lower the forensic claim.")
		}
		if event.RedactionAction == RedactionActionWithhold {
			if event.Withholding == nil || event.Withholding.Authority.ActorID == "" || event.Withholding.ReasonCode == "" || event.Withholding.Justification == "" {
				return cannotVerify("redaction_unresolved_visible", "withholding_audit_missing", "withholding lacks required audit evidence", "Record withholding authority, requestor when different, reason, and justification.")
			}
			if event.Withholding.Authority.VerificationState != AuthorityVerified {
				return cannotVerify("redaction_unresolved_visible", "withholding_authority_unverifiable", "withholding authority is not provenance or accountability verified", "Record verified withholding authority evidence.")
			}
		}
	}
	return pass("redaction_unresolved_visible", "redaction_resolved", "redaction states are resolved or explicitly non-blocking")
}

func retentionModeCondition(input Input) Condition {
	allowed := allowedRetentionModes(input.Policy)
	for _, event := range input.Run.Events {
		if !validRetentionMode(event.RetentionMode) {
			return fail("retention_mode_declared", "invalid_retention_mode", "event declares a non-FR-054 retention mode", "Use digest_only, sanitized_excerpt, encrypted_raw_ref, external_artifact_ref, or not_assessed.")
		}
		if len(allowed) > 0 && !allowed[event.RetentionMode] {
			return fail("retention_mode_declared", "retention_mode_not_policy_allowed", "event retention mode is not allowed by the selected redaction policy", "Use a retention mode allowed by the selected policy.")
		}
	}
	return pass("retention_mode_declared", "retention_mode_declared", "events declare FR-054 retention modes")
}

func criticalEvidenceCondition(input Input) Condition {
	critical := criticalEvents(input)
	for _, event := range input.Run.Events {
		if !critical[event.EventType] && event.ForensicImportance != "critical" {
			continue
		}
		switch event.RetentionMode {
		case RetentionModeSanitizedExcerpt:
			continue
		case RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef:
			if event.RawReference == nil {
				return cannotVerify("critical_evidence_reconstructable", "raw_reference_missing", "critical raw reference evidence is missing", "Bind critical evidence to encrypted or external raw reference metadata.")
			}
		case RetentionModeDigestOnly:
			return failWithCap("critical_evidence_reconstructable", "critical_evidence_digest_only", "critical evidence is digest-only and not reconstructable", RetentionModeDigestOnly, "Retain sanitized excerpts, encrypted raw references, or external artifact references for critical event families.")
		case RetentionModeNotAssessed:
			return failWithCap("critical_evidence_reconstructable", "critical_evidence_not_assessed", "critical evidence retention is not assessed", RetentionModeNotAssessed, "Capture critical evidence or keep forensic retention open.")
		}
	}
	return pass("critical_evidence_reconstructable", "critical_evidence_reconstructable", "critical event families have reconstructable retention")
}

func rawReferenceCondition(input Input) Condition {
	for _, event := range input.Run.Events {
		ref := event.RawReference
		if ref == nil {
			continue
		}
		if condition, ok := validateRawReference(ref); !ok {
			return condition
		}
	}
	return pass("raw_reference_bound", "raw_reference_bound", "raw references are digest-bound and access-verifiable")
}

func validateRawReference(ref *RawReference) (Condition, bool) {
	if ref.Digest.Algorithm != "sha256" || len(ref.Digest.Value) != 64 {
		return fail("raw_reference_bound", "weak_digest", "raw reference digest is weak, unknown, or malformed", "Use SHA-256 or stronger digest binding for raw references."), false
	}
	if ref.ReferenceType != RetentionModeEncryptedRawRef && ref.ReferenceType != RetentionModeExternalArtifactRef {
		return fail("raw_reference_bound", "raw_reference_type_invalid", "raw reference type is not an accepted FR-054 raw reference mode", "Use encrypted_raw_ref or external_artifact_ref."), false
	}
	if ref.ReferenceURI == "" {
		return cannotVerify("raw_reference_bound", "missing_reference", "raw reference URI is missing", "Provide a stable encrypted or external raw reference."), false
	}
	if rawReferenceAccessUnverifiable(ref) {
		return cannotVerify("raw_reference_bound", "access_unverifiable", "raw reference access state is not verifiably available", "Record current access verification state and time."), false
	}
	if ref.AccessStateLastVerified == "" {
		return cannotVerify("raw_reference_bound", "access_unverifiable", "raw reference access verification time is missing", "Record access_state_last_verified for the assessment."), false
	}
	if encryptedKeyCustodyUnverifiable(ref) {
		return cannotVerify("raw_reference_bound", "key_custody_unverifiable", "encrypted raw reference key custody is not verifiable", "Record holder_known or escrowed key custody state."), false
	}
	if retentionLifecycleUnverifiable(ref.RetentionLifecycle.State) {
		return cannotVerify("raw_reference_bound", "retention_lifecycle_unverifiable", "raw reference retention lifecycle is not active", "Record active retention lifecycle evidence."), false
	}
	return Condition{}, true
}

func rawReferenceAccessUnverifiable(ref *RawReference) bool {
	switch ref.AccessState {
	case "", AccessStateNotAssessed, AccessStateUnavailable, AccessStateRevoked:
		return true
	default:
		return false
	}
}

func encryptedKeyCustodyUnverifiable(ref *RawReference) bool {
	if ref.ReferenceType != RetentionModeEncryptedRawRef {
		return false
	}
	switch ref.KeyCustodyState {
	case "", KeyCustodyUnknown, KeyCustodyNotAssessed, KeyCustodyCompromised, KeyCustodyDestroyed:
		return true
	default:
		return false
	}
}

func retentionLifecycleUnverifiable(state string) bool {
	switch state {
	case "", RetentionLifecycleNotAssessed, RetentionLifecycleExpired, RetentionLifecycleRevoked, RetentionLifecycleDeleted:
		return true
	default:
		return false
	}
}

func overclaimCondition(input Input) Condition {
	critical := criticalEvents(input)
	for _, event := range input.Run.Events {
		if !critical[event.EventType] && event.ForensicImportance != "critical" {
			continue
		}
		if event.RetentionMode == RetentionModeDigestOnly || event.RetentionMode == RetentionModeNotAssessed {
			return fail("forensic_profile_not_overclaimed", "forensic_profile_capped", "forensic retention output is capped by insufficient critical evidence", "Do not claim forensic reconstruction for digest-only or not-assessed critical evidence.")
		}
	}
	return pass("forensic_profile_not_overclaimed", "forensic_profile_not_overclaimed", "forensic output does not exceed retained evidence")
}

func profileSelectionCondition(input Input) Condition {
	selection := input.Run.ProfileSelection
	if selection.SelectedProfile == "" {
		return Condition{ID: "profile_selection_accountable", State: StateNotAssessed, ReasonCode: "profile_selection_not_assessed", Reason: "profile selection accountability is not recorded", NextAction: "Record actor, profile, policy digest, and justification when policy requires it."}
	}
	if selection.SelectedProfile != ProfileForensicRetention || selection.RedactionPolicyDigest != input.Policy.PolicyDigest || selection.ActorID == "" || selection.Justification == "" || !selection.AuthorityVerified {
		return cannotVerify("profile_selection_accountable", "profile_selection_unverifiable", "forensic profile selection accountability cannot be verified", "Record accountable forensic profile selection evidence.")
	}
	return pass("profile_selection_accountable", "profile_selection_accountable", "forensic profile selection is accountable")
}

func criticalEvents(input Input) map[string]bool {
	out := map[string]bool{}
	defaults := []string{
		"command_started",
		"command_finished",
		"test_output_observed",
		"file_mutation_observed",
		"artifact_captured",
		"model_identity_observed",
		"harness_identity_observed",
		"requirement_superseded",
		"redaction_applied",
		"run_closed",
	}
	for _, eventType := range defaults {
		out[eventType] = true
	}
	for _, eventType := range input.Policy.CriticalEventFamilies {
		out[eventType] = true
	}
	for _, downgrade := range input.Policy.NonCriticalEventFamilyReasons {
		if downgrade.EventType != "" && downgrade.Reason != "" && downgrade.AuthorityID != "" {
			delete(out, downgrade.EventType)
		}
	}
	return out
}

func validRetentionMode(mode string) bool {
	switch mode {
	case RetentionModeDigestOnly, RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef, RetentionModeNotAssessed:
		return true
	default:
		return false
	}
}

func policyRules(policy Policy) map[string]Rule {
	out := map[string]Rule{}
	for _, rule := range policy.Rules {
		if rule.RuleID != "" {
			out[rule.RuleID] = rule
		}
	}
	return out
}

func allowedRetentionModes(policy Policy) map[string]bool {
	out := map[string]bool{}
	for _, mode := range policy.AllowedRetentionModes {
		out[mode] = true
	}
	return out
}

func topLevel(conditions []Condition) string {
	highest := StatePass
	for _, condition := range conditions {
		if condition.State == StateFail {
			return StateFail
		}
		if condition.State == StateCannotVerify || condition.State == StateNotAssessed {
			highest = StateCannotVerify
		}
	}
	return highest
}

func reasons(conditions []Condition) []string {
	out := []string{}
	for _, condition := range conditions {
		if condition.State != StatePass {
			out = append(out, condition.ReasonCode+": "+condition.Reason)
		}
	}
	sort.Strings(out)
	return out
}

func nextActions(conditions []Condition) []string {
	set := map[string]bool{}
	for _, condition := range conditions {
		if condition.State != StatePass && condition.NextAction != "" {
			set[condition.NextAction] = true
		}
	}
	out := []string{}
	for action := range set {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

func pass(id, code, reason string) Condition {
	return Condition{ID: id, State: StatePass, ReasonCode: code, Reason: reason}
}

func fail(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateFail, ReasonCode: code, Reason: reason, NextAction: action}
}

func failWithCap(id, code, reason, cap, action string) Condition {
	return Condition{ID: id, State: StateFail, ReasonCode: code, Reason: reason, CappedToRetentionMode: cap, NextAction: action}
}

func cannotVerify(id, code, reason, action string) Condition {
	return Condition{ID: id, State: StateCannotVerify, ReasonCode: code, Reason: reason, NextAction: action}
}

// ValidTestInput is used by CLI tests to write representative fixtures without
// duplicating Block 18 policy/run semantics outside this package.
func ValidTestInput() Input {
	return validTestInput()
}

func validTestInput() Input {
	policyDigest := "1111111111111111111111111111111111111111111111111111111111111111"
	return Input{
		Policy: Policy{
			PolicyID:              "customer-forensic-policy-v1",
			SchemaVersion:         "1.0.0",
			PolicyDigest:          policyDigest,
			PolicyProvenance:      Provenance{Source: "vcs", Digest: policyDigest},
			AllowedRetentionModes: []string{RetentionModeDigestOnly, RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef, RetentionModeNotAssessed},
			RedactionActions:      []string{RedactionActionApplyRule, RedactionActionWithhold, RedactionActionMarkUnavailable},
			ForbiddenPersistenceClasses: []string{
				"credentials",
				"tokens",
				"raw_prompts",
				"raw_model_responses",
				"source_snippets",
				"stdout_stderr_bodies",
				"oidc_tokens",
				"adapter_secrets",
				"gateway_tokens",
				"checkpoint_key_material",
			},
			CriticalEventFamilies: []string{"command_finished", "test_output_observed"},
			Authority:             AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
			ProfileMappings: []ProfileMapping{
				{
					EventFamily:            "command_finished",
					RequiredRetentionModes: []string{RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef},
					Critical:               true,
					Authority:              AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				},
				{
					EventFamily:            "test_output_observed",
					RequiredRetentionModes: []string{RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef},
					Critical:               true,
					Authority:              AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				},
			},
			UnresolvedRedactionImpact: "fail_forensic_retention",
			Rules: []Rule{
				{RuleID: "secret-token-v1", DetectorFamily: "secret", RuleVersion: "1.0.0", Action: RedactionActionApplyRule, RetentionMode: RetentionModeSanitizedExcerpt},
				{RuleID: "withhold-privacy-v1", DetectorFamily: "privacy", RuleVersion: "1.0.0", Action: RedactionActionWithhold, RetentionMode: RetentionModeNotAssessed},
			},
		},
		Run: RunEvidence{
			RunID:                 "forensic-run-1",
			SelectedProfile:       ProfileForensicRetention,
			RedactionPolicyDigest: policyDigest,
			ProfileSelection: ProfileSelection{
				ActorID:                 "human:security-owner",
				SelectedProfile:         ProfileForensicRetention,
				RedactionPolicyDigest:   policyDigest,
				Justification:           "incident review",
				AuthorityVerified:       true,
				SelectionEvidenceDigest: "2222222222222222222222222222222222222222222222222222222222222222",
			},
			Events: []EventRetention{
				{
					EventType:              "command_finished",
					RetentionMode:          RetentionModeSanitizedExcerpt,
					ForensicImportance:     "critical",
					RedactionPolicyDigest:  policyDigest,
					RedactionInputDigest:   "3333333333333333333333333333333333333333333333333333333333333333",
					RedactedPayloadDigest:  "4444444444444444444444444444444444444444444444444444444444444444",
					RedactionAction:        RedactionActionApplyRule,
					RedactionRuleRefs:      []string{"secret-token-v1"},
					SecretLikeValuePresent: false,
					RedactionAuthority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				},
				{
					EventType:              "test_output_observed",
					RetentionMode:          RetentionModeExternalArtifactRef,
					ForensicImportance:     "critical",
					RedactionPolicyDigest:  policyDigest,
					RedactionInputDigest:   "5555555555555555555555555555555555555555555555555555555555555555",
					RedactedPayloadDigest:  "6666666666666666666666666666666666666666666666666666666666666666",
					RedactionAction:        RedactionActionApplyRule,
					RedactionRuleRefs:      []string{"secret-token-v1"},
					SecretLikeValuePresent: false,
					RedactionAuthority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
					RawReference: &RawReference{
						ReferenceType:           RetentionModeExternalArtifactRef,
						ReferenceURI:            "artifact://ci/run-1/test-output",
						Digest:                  Digest{Algorithm: "sha256", Value: "7777777777777777777777777777777777777777777777777777777777777777"},
						AccessState:             AccessStateVerifiedAvailable,
						AccessStateLastVerified: "2026-05-07T10:00:00Z",
						KeyCustodyState:         KeyCustodyNotApplicable,
						RetentionLifecycle:      RetentionLifecycle{State: RetentionLifecycleActive, PolicyRef: "policy:retain-30d", ExpiresAt: "2026-06-06T10:00:00Z"},
						UnavailableReason:       UnavailableReasonNotAssessed,
					},
				},
			},
		},
	}
}
