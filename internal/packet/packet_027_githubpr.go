package packet

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
