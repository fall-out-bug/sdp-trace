package main

type doctorCheck struct {
	ID        string   `json:"id"`
	State     string   `json:"state"`
	Reason    string   `json:"reason"`
	Contract  string   `json:"contract_id,omitempty"`
	Missing   []string `json:"missing,omitempty"`
	Reference string   `json:"reference,omitempty"`
}
