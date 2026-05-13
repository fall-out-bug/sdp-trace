package packet

type Row struct {
	ID           string   `json:"id"`
	State        string   `json:"state"`
	Summary      string   `json:"summary"`
	EvidenceRefs []string `json:"evidence_refs"`
	Reason       string   `json:"reason,omitempty"`
	Owner        string   `json:"owner"`
}
