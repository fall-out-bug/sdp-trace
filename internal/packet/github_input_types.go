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

type PromptBoundary struct {
	Text          string `json:"text,omitempty"`
	Digest        string `json:"digest,omitempty"`
	CaptureActor  string `json:"capture_actor,omitempty"`
	CapturedAt    string `json:"captured_at,omitempty"`
	CaptureMethod string `json:"capture_method,omitempty"`
}

type PromptBoundaryClassification struct {
	Verdict          string   `json:"verdict"`
	RouteProofEffect string   `json:"route_proof_effect"`
	Reasons          []string `json:"reasons,omitempty"`
}

type IntegrationAction struct {
	Kind     string `json:"kind"`
	Actor    string `json:"actor"`
	Resolver string `json:"resolver"`
}

type BuildPRResult struct {
	State      string   `json:"state"`
	BundlePath string   `json:"bundle_path,omitempty"`
	PacketPath string   `json:"packet_path,omitempty"`
	ResultPath string   `json:"result_path,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

type GitHubPR struct {
	Number  int    `json:"number"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	BodyRef string `json:"body_ref,omitempty"`
	Author  string `json:"author"`
	BaseRef string `json:"base_ref"`
	HeadRef string `json:"head_ref"`
	HeadSHA string `json:"head_sha"`
}

type GitHubCommitRange struct {
	Base            string `json:"base"`
	Head            string `json:"head"`
	ChangedFilesRef string `json:"changed_files_ref,omitempty"`
}

type GitHubCheck struct {
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Conclusion   string   `json:"conclusion"`
	ArtifactRefs []string `json:"artifact_refs,omitempty"`
}

type GitHubArtifact struct {
	Name         string `json:"name"`
	Resolver     string `json:"resolver"`
	RetainedForm string `json:"retained_form"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	Digest       string `json:"digest,omitempty"`
}

type GitHubReview struct {
	Reviewer string `json:"reviewer"`
	Resolver string `json:"resolver"`
	State    string `json:"state"`
}
