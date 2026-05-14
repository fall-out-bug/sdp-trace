package main

var interactionRelayStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"task-id", ""},
	{"actor-type", "human_user"},
	{"actor-id", ""},
	{"target", "agent"},
	{"event-type", "corrective_feedback"},
	{"operation-id", ""},
	{"stage-id", ""},
	{"out", ""},
}
