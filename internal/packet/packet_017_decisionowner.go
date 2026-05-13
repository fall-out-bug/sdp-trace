package packet

type DecisionOwner struct {
	Decision string `json:"decision"`
	Owner    string `json:"owner"`
	State    string `json:"state"`
	Reason   string `json:"reason,omitempty"`
}
