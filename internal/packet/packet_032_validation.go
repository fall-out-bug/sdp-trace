package packet

type Validation struct {
	State  string   `json:"state"`
	Errors []string `json:"errors,omitempty"`
}
