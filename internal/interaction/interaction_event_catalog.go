package interaction

import "fmt"

var eventTypes = []string{
	"clarification_answer",
	"clarification_request",
	"boundary_violation",
	"corrective_feedback",
	"evidence_correction",
	"pause_requested",
	"plan_approved",
	"plan_proposed",
	"plan_rejected",
	"resume_approved",
	"task_assignment",
	"tool_or_model_drift",
}

func EventTypes() []string {
	return append([]string(nil), eventTypes...)
}

func validateEventTypeAndFriction(event Event) error {
	// validateEventTypeAndFriction keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if !validEventType(event.EventType) {
		return fmt.Errorf("unsupported event_type %q", event.EventType)
	}
	if event.FrictionClass != frictionClass(event.EventType) {
		return fmt.Errorf("event friction_class %q does not match event_type %q", event.FrictionClass, event.EventType)
	}
	return nil
}
