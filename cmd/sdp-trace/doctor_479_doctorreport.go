package main

type doctorReport struct {
	Command            string        `json:"command"`
	Result             string        `json:"result"`
	Environment        []doctorCheck `json:"environment"`
	ControlPoints      []doctorCheck `json:"control_points"`
	SafeRetentionModes []string      `json:"safe_retention_modes"`
}
