package prreview

type UnavailableField struct {
	Field  string `json:"field"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}

// ReviewProfile declares which planes must be reviewed and which runners may
// produce evidence for those planes.
