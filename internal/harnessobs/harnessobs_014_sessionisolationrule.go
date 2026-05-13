package harnessobs

type SessionIsolationRule struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	Required   bool   `json:"required"`
}
