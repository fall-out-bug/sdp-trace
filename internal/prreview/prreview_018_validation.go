package prreview

type Validation struct {
	SchemaVersion       string          `json:"schema_version"`
	PacketDigest        string          `json:"packet_digest"`
	ReviewCoverageState string          `json:"review_coverage_state"`
	CIState             string          `json:"ci_state"`
	AuthorityScope      string          `json:"authority_scope"`
	MergeDecision       string          `json:"merge_decision"`
	ReleaseDecision     string          `json:"release_decision"`
	RiskAcceptance      string          `json:"risk_acceptance"`
	PlaneResults        []PlaneResult   `json:"plane_results"`
	Findings            []LedgerFinding `json:"findings"`
	Reasons             []string        `json:"reasons"`
	NextActions         []string        `json:"next_actions"`
}

// PlaneResult states whether one required review plane has usable evidence and
// what action remains when it does not.
