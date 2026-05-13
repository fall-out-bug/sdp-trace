package packet

type PromptBoundary struct {
	Text          string `json:"text,omitempty"`
	Digest        string `json:"digest,omitempty"`
	CaptureActor  string `json:"capture_actor,omitempty"`
	CapturedAt    string `json:"captured_at,omitempty"`
	CaptureMethod string `json:"capture_method,omitempty"`
}
