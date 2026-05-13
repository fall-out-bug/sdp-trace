package packet

type GitHubPREvidenceInput struct {
	// GitHub input is source evidence for generation; it is not accepted as a
	// verdict until rows and manifest refs are built and validated.
	SchemaVersion          string              `json:"schema_version"`
	PR                     GitHubPR            `json:"pr"`
	CommitRange            GitHubCommitRange   `json:"commit_range"`
	Checks                 []GitHubCheck       `json:"checks,omitempty"`
	Artifacts              []GitHubArtifact    `json:"artifacts,omitempty"`
	Reviews                []GitHubReview      `json:"reviews,omitempty"`
	WorkflowRunID          string              `json:"workflow_run_id,omitempty"`
	RequirePromptBoundary  bool                `json:"require_prompt_boundary,omitempty"`
	AgentRouteRefs         []string            `json:"agent_route_refs,omitempty"`
	AgentRouteComponents   []string            `json:"agent_route_components,omitempty"`
	AgentRouteDigest       string              `json:"agent_route_digest,omitempty"`
	AgentRouteEvidenceKind string              `json:"agent_route_evidence_kind,omitempty"`
	PromptBoundary         PromptBoundary      `json:"prompt_boundary,omitempty"`
	IntegrationActions     []IntegrationAction `json:"integration_actions,omitempty"`
}
