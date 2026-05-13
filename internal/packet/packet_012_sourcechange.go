package packet

type SourceChange struct {
	Repository  string `json:"repository,omitempty"`
	ChangeID    string `json:"change_id,omitempty"`
	URL         string `json:"url,omitempty"`
	BaseRef     string `json:"base_ref,omitempty"`
	HeadRef     string `json:"head_ref,omitempty"`
	CommitRange string `json:"commit_range,omitempty"`
	HeadSHA     string `json:"head_sha,omitempty"`
}
