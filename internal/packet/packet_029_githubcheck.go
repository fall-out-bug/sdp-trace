package packet

type GitHubCheck struct {
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Conclusion   string   `json:"conclusion"`
	ArtifactRefs []string `json:"artifact_refs,omitempty"`
}
