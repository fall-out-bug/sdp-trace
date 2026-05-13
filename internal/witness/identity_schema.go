package witness

// SourceIdentity binds a witness to the repository revision that produced or
// verified the evidence.
type SourceIdentity struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	CommitSHA  string `json:"commit_sha"`
}

// CIIdentity captures the provider run identity before profile logic decides
// whether the fields are independently witnessed.
type CIIdentity struct {
	Provider   string `json:"provider"`
	ServerURL  string `json:"server_url"`
	Workflow   string `json:"workflow"`
	Job        string `json:"job"`
	RunID      string `json:"run_id"`
	RunAttempt string `json:"run_attempt"`
	Actor      string `json:"actor"`
}
