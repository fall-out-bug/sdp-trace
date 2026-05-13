package query

type QueryPackRow struct {
	ID                   string   `json:"id"`
	Query                string   `json:"query"`
	EvidenceState        string   `json:"evidence_state"`
	EvidenceFamily       string   `json:"evidence_family"`
	Reconstructable      *bool    `json:"reconstructable,omitempty"`
	SourceRef            string   `json:"source_ref"`
	SourceConditionID    string   `json:"source_condition_id,omitempty"`
	SourceConditionState string   `json:"source_condition_state,omitempty"`
	ReasonCode           string   `json:"reason_code,omitempty"`
	EvidenceGap          string   `json:"evidence_gap,omitempty"`
	RelatedRows          []string `json:"related_rows,omitempty"`
}
