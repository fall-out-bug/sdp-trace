package packet

type GitHubReview struct {
	Reviewer string `json:"reviewer"`
	Resolver string `json:"resolver"`
	State    string `json:"state"`
}
