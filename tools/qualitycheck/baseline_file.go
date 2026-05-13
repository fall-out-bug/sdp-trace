package main

// fileMIBaseline is the persisted ratchet for existing file-MI debt.
type fileMIBaseline struct {
	// SchemaVersion keeps file baselines distinct from function baselines.
	SchemaVersion string                 `json:"schema_version"`
	Threshold     float64                `json:"threshold"`
	Files         []fileMIBaselineRecord `json:"files"`
}
