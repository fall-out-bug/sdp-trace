package adaptercapture

func containsSecret(value string) bool {
	// containsSecret preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if value == "" {
		return false
	}

	for _, needle := range secretMarkers {
		if containsFold(value, needle) {

			return true
		}
	}
	return false
}

var secretMarkers = []string{
	"secret-token",
	"password=",
	"token=",
	"bearer ",
	"access_token=",
	"credential=",
	"oidc_token",
	"session_id=",
	"raw prompt",
	"raw_prompt",
	"raw response",
	"raw_response",
	"raw_review_body",
	"tool_input_body",
	"tool_output_body",
	"model_request_payload",
	"model_response_payload",
	"adapter_config_raw",
}
