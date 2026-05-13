package packet

type TheaterFinding struct {
	ReasonCode              string   `json:"reason_code"`
	State                   string   `json:"state"`
	Severity                string   `json:"severity,omitempty"`
	Finding                 string   `json:"finding"`
	TriggerEvidenceRefs     []string `json:"trigger_evidence_refs"`
	RequiredClosureEvidence string   `json:"required_closure_evidence,omitempty"`
}
