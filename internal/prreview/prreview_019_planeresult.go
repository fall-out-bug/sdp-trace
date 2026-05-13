package prreview

type PlaneResult struct {
	Plane      string `json:"plane"`
	Status     string `json:"status"`
	Usable     bool   `json:"usable"`
	RunID      string `json:"review_run_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
	NextAction string `json:"next_action,omitempty"`
}
