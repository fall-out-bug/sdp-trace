package main

const (
	defaultRunRoot   = ".sdp-trace-runs"
	defaultReportDir = ".sdp-trace-report"
)

type doctorReport struct {
	Command            string        `json:"command"`
	Result             string        `json:"result"`
	Environment        []doctorCheck `json:"environment"`
	ControlPoints      []doctorCheck `json:"control_points"`
	SafeRetentionModes []string      `json:"safe_retention_modes"`
}

type doctorCheck struct {
	ID        string   `json:"id"`
	State     string   `json:"state"`
	Reason    string   `json:"reason"`
	Contract  string   `json:"contract_id,omitempty"`
	Missing   []string `json:"missing,omitempty"`
	Reference string   `json:"reference,omitempty"`
}

type doctorOptions struct {
	ContractPath string
	OutputDir    string
	ReportDir    string
	Env          map[string]string
}

type previewBoundary struct {
	Boundary string `json:"boundary"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
}

type previewOfflineImplication struct {
	Requirement string `json:"requirement"`
	State       string `json:"state"`
	Reason      string `json:"reason"`
}
