package prreview

type PreviewRole struct {
	RoleID         string `json:"role_id"`
	Plane          string `json:"plane"`
	Runner         string `json:"runner"`
	RequestedModel string `json:"requested_model"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	CommandDigest  string `json:"command_digest"`
	PromptRef      string `json:"prompt_template_ref,omitempty"`
	PromptDigest   string `json:"prompt_digest,omitempty"`
}

// RunSet records reviewer outputs for a single packet digest.
