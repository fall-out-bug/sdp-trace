package managed

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
