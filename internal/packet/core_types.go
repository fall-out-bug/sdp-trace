package packet

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
