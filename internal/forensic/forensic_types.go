package forensic

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

type policyContractCheck struct {
	failed    bool
	condition Condition
}

type conditionFailure struct {
	matched   bool
	condition Condition
}
