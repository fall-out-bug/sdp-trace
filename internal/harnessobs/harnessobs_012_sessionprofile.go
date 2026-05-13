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
