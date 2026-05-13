package packet

type ResidualGap struct {
	RowID           string   `json:"row_id"`
	State           string   `json:"state"`
	Reason          string   `json:"reason"`
	EvidenceRefs    []string `json:"evidence_refs,omitempty"`
	ClosureEvidence string   `json:"closure_evidence,omitempty"`
}
