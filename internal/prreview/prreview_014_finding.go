package prreview

type Finding struct {
	ID           string   `json:"id"`
	Severity     string   `json:"severity"`
	Citation     Citation `json:"citation"`
	Summary      string   `json:"summary"`
	Rationale    string   `json:"rationale,omitempty"`
	SuggestedFix string   `json:"suggested_fix,omitempty"`
	Question     string   `json:"question,omitempty"`
	EvidenceRefs []string `json:"evidence_refs,omitempty"`
}

// Citation points a finding back to packet evidence; validation rejects claims
// that cannot resolve to a packet ref or digest.
