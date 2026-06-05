package adaptercapture

type validEventSpec struct {
	id        string
	eventType string
	sequence  int
}

var validEventSpecs = []validEventSpec{
	{id: "evt-run", eventType: "run_started", sequence: 1},
	{id: "evt-task", eventType: "task_locked", sequence: 2},
	{id: "evt-tool", eventType: "tool_call", sequence: 3},
	{id: "evt-command", eventType: "command_started", sequence: 4},
	{id: "evt-file", eventType: "file_mutation", sequence: 5},
	{id: "evt-model", eventType: "model_call_observed", sequence: 6},
	{id: "evt-test", eventType: "test_observed", sequence: 7},
	{id: "evt-close", eventType: "run_closed", sequence: 8},
}
