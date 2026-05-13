package prreview

type ReviewRole struct {
	RoleID               string   `json:"role_id"`
	Plane                string   `json:"plane"`
	Runner               string   `json:"runner"`
	RequestedModel       string   `json:"requested_model"`
	Command              []string `json:"command,omitempty"`
	TimeoutSeconds       int      `json:"timeout_seconds,omitempty"`
	PromptTemplateRef    string   `json:"prompt_template_ref,omitempty"`
	RequiredOutputSchema string   `json:"required_output_schema,omitempty"`
	RawOutputRetention   string   `json:"raw_output_retention,omitempty"`
	ReadOnlyEnforced     bool     `json:"read_only_enforced,omitempty"`
	WorkingTreeMode      string   `json:"working_tree_mode,omitempty"`
}

// RunOptions contains local execution policy; it is intentionally separate
// from the portable packet so packet evidence stays harness-neutral.
