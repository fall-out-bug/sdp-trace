package prreview

type ReviewerResult struct {
	ReviewRunID      string    `json:"review_run_id"`
	PacketDigest     string    `json:"packet_digest"`
	Plane            string    `json:"plane"`
	RoleID           string    `json:"role_id"`
	Runner           string    `json:"runner"`
	RequestedModel   string    `json:"requested_model"`
	ObservedModel    string    `json:"observed_model"`
	ModelFamily      string    `json:"model_family"`
	ModelVersion     string    `json:"model_version"`
	FallbackForModel string    `json:"fallback_for_model,omitempty"`
	FallbackReason   string    `json:"fallback_reason,omitempty"`
	Status           string    `json:"status"`
	Findings         []Finding `json:"findings"`
	CommandDigest    string    `json:"command_digest,omitempty"`
	RawOutputRef     *SafeRef  `json:"raw_output_ref,omitempty"`
	PromptRef        *SafeRef  `json:"prompt_ref,omitempty"`
	ContextRefs      []string  `json:"context_refs,omitempty"`
	StartedAt        string    `json:"started_at,omitempty"`
	EndedAt          string    `json:"ended_at,omitempty"`
	Reason           string    `json:"reason,omitempty"`
}

// Finding is a reviewer claim that must be backed by a citation before
// validation can treat it as replayable.
