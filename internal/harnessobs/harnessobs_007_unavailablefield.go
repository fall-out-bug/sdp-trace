package harnessobs

type UnavailableField struct {
	Field      string `json:"field"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
}
