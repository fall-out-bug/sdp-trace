package packet

type GitHubArtifact struct {
	Name         string `json:"name"`
	Resolver     string `json:"resolver"`
	RetainedForm string `json:"retained_form"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	Digest       string `json:"digest,omitempty"`
}
