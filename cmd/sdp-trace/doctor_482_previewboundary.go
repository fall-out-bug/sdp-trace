package main

type previewBoundary struct {
	Boundary string `json:"boundary"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
}
