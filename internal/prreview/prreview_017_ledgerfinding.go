package prreview

type LedgerFinding struct {
	ID                  string   `json:"id"`
	ReviewRunID         string   `json:"review_run_id"`
	Plane               string   `json:"plane"`
	RoleID              string   `json:"role_id"`
	Severity            string   `json:"severity"`
	Summary             string   `json:"summary"`
	Citation            Citation `json:"citation"`
	Disposition         string   `json:"disposition"`
	EvidenceRefs        []string `json:"evidence_refs,omitempty"`
	DispositionEvidence string   `json:"disposition_evidence,omitempty"`
}

// Validation is the portable coverage verdict for review evidence only. It does
// not authorize merge, release, or risk acceptance.
