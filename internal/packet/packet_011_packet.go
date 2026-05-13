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
