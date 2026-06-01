package harnessobs

type SessionProfile struct {
	SchemaVersion      string                 `json:"schema_version"`
	ProfileID          string                 `json:"profile_id"`
	HarnessProfilePath string                 `json:"harness_profile_path"`
	EventSourcePath    string                 `json:"event_source_path"`
	RawEventSourcePath string                 `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string                 `json:"raw_event_format,omitempty"`
	SetupActions       []SessionSetupAction   `json:"setup_actions,omitempty"`
	IsolationRules     []SessionIsolationRule `json:"isolation_rules,omitempty"`
	StreamCapture      string                 `json:"stream_capture"`
}

type SessionSetupAction struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
}

type SessionIsolationRule struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	Required   bool   `json:"required"`
}

type SessionRun struct {
	SchemaVersion      string                   `json:"schema_version"`
	ProfileID          string                   `json:"profile_id"`
	HarnessProfilePath string                   `json:"harness_profile_path"`
	EventSourcePath    string                   `json:"event_source_path"`
	RawEventSourcePath string                   `json:"raw_event_source_path,omitempty"`
	RawEventFormat     string                   `json:"raw_event_format,omitempty"`
	SetupActionIDs     []string                 `json:"setup_action_ids,omitempty"`
	IsolationResults   []SessionIsolationResult `json:"isolation_results,omitempty"`
	CommandDigest      string                   `json:"command_digest,omitempty"`
	CommandDigestState string                   `json:"command_digest_state,omitempty"`
	CommandModel       string                   `json:"command_model,omitempty"`
	CommandModelState  string                   `json:"command_model_state,omitempty"`
	ProcessID          int                      `json:"process_id,omitempty"`
	ProcessIDState     string                   `json:"process_id_state,omitempty"`
	StartTime          string                   `json:"start_time,omitempty"`
	EndTime            string                   `json:"end_time,omitempty"`
	SourceCommit       string                   `json:"source_commit,omitempty"`
	SourceCommitState  string                   `json:"source_commit_state,omitempty"`
	ObservedRunDir     string                   `json:"observed_run_dir,omitempty"`
	OutputDigest       string                   `json:"output_digest,omitempty"`
	NormalizedDigest   string                   `json:"normalized_digest,omitempty"`
	CollectionState    string                   `json:"collection_state,omitempty"`
	CollectionReason   string                   `json:"collection_reason,omitempty"`
	CreatedAt          string                   `json:"created_at"`
}

// SessionIsolationResult records the live readback result of a setup rule; it
// is evidence about the local setup artifact, not proof that a harness obeyed it.
type SessionIsolationResult struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	SHA256     string `json:"sha256,omitempty"`
}
