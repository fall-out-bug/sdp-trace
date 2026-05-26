package harnessobs

type Limits struct {
	MaxLineBytes int `json:"max_line_bytes,omitempty"`
	MaxEvents    int `json:"max_events,omitempty"`
}
