package harnessobs

type SessionIsolationResult struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	TargetPath string `json:"target_path"`
	Pattern    string `json:"pattern"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	SHA256     string `json:"sha256,omitempty"`
}
