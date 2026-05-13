package query

type assessmentEnvelope struct {
	SchemaVersion            string                `json:"schema_version"`
	ForensicConditions       []assessmentCondition `json:"forensic_conditions,omitempty"`
	AdapterCaptureConditions []assessmentCondition `json:"adapter_capture_conditions,omitempty"`
	ForensicAssessment       string                `json:"forensic_retention_assessment,omitempty"`
	AdapterCaptureAssessment string                `json:"adapter_capture_assessment,omitempty"`
}

type assessmentCondition struct {
	ID                    string `json:"id"`
	State                 string `json:"state"`
	ReasonCode            string `json:"reason_code"`
	Reason                string `json:"reason"`
	CappedToRetentionMode string `json:"capped_to_retention_mode,omitempty"`
}
