package main

// functionMIBaselineRecord stores one below-threshold function by stable key.
type functionMIBaselineRecord struct {
	// Key is the lookup identity used by the ratchet; the remaining fields are
	// review context.
	Key                  string  `json:"key"`
	File                 string  `json:"file"`
	Line                 int     `json:"line"`
	Name                 string  `json:"name"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
}
