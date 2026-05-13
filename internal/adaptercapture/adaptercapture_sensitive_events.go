package adaptercapture

func sensitiveEvent(eventType string) bool {
	return sensitiveEventTypes[eventType]
}

var sensitiveEventTypes = map[string]bool{
	"tool_call":           true,
	"command_started":     true,
	"model_call_observed": true,
	"test_observed":       true,
}
