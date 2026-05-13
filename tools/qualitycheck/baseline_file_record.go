package main

// fileMIBaselineRecord stores one below-threshold file by normalized path key.
type fileMIBaselineRecord struct {
	// File duplicates the key as display text for generated JSON review.
	Key                  string  `json:"key"`
	File                 string  `json:"file"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
}
