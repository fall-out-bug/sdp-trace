package adaptercapture

func digestOnlyValidEvent(eventType string) bool {
	return digestOnlyValidEvents[eventType]
}

var digestOnlyValidEvents = map[string]bool{
	"run_started":   true,
	"task_locked":   true,
	"run_closed":    true,
	"file_mutation": true,
}
