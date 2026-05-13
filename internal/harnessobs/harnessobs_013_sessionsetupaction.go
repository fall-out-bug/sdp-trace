package harnessobs

type SessionSetupAction struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
}
