package repoobserver

// Gap explains why an observer surface has not reached its required state.
type Gap struct {
	SurfaceID  string `json:"surface_id"`
	ReasonCode string `json:"reason_code"`
	Detail     string `json:"detail"`
}

// NextAction gives concrete remediation without treating the action text as
// proof.
type NextAction struct {
	SurfaceID  string `json:"surface_id"`
	ActionText string `json:"action_text"`
	Blocking   bool   `json:"blocking"`
}
