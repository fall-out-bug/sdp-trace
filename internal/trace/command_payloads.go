package trace

// RecordedCommandPayload stores common command lifecycle metadata.
type RecordedCommandPayload struct {
	CommandID string `json:"command_id"`
	Command   string `json:"command"`
}

// CommandStartedPayload captures launch metadata.
type CommandStartedPayload struct {
	RecordedCommandPayload
	CmdPID         int    `json:"cmd_pid"`
	WrapperName    string `json:"wrapper_name"`
	Cwd            string `json:"cwd"`
	ContractDigest string `json:"contract_digest"`
}

// CommandFinishedPayload captures close metadata.
type CommandFinishedPayload struct {
	RecordedCommandPayload
	ExitCode       int    `json:"exit_code"`
	Signal         string `json:"signal,omitempty"`
	StdoutDigest   string `json:"stdout_digest"`
	StderrDigest   string `json:"stderr_digest"`
	ClockMonotonic int64  `json:"clock_monotonic,omitempty"`
}

// RecorderAttachedPayload captures launch metadata for the recorder.
type RecorderAttachedPayload struct {
	RecorderVersion string `json:"recorder_version"`
	RunNonce        string `json:"run_nonce"`
	Pid             int    `json:"pid"`
}

// RunStartedPayload captures command + contract context.
type RunStartedPayload struct {
	TaskRef        string `json:"task_ref"`
	Command        string `json:"command"`
	ContractID     string `json:"contract_id"`
	ContractDigest string `json:"contract_digest"`
	ContractPath   string `json:"contract_path,omitempty"`
}

// RunClosedPayload summarizes terminal state.
type RunClosedPayload struct {
	ChainHead     string `json:"chain_head"`
	ClosureState  string `json:"closure_state"`
	CommandID     string `json:"command_id"`
	MissingCount  int    `json:"missing_evidence_count"`
	ObservedCount int    `json:"observed_event_count"`
}
