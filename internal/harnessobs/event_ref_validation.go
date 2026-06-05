package harnessobs

import "errors"

type eventRefCheck struct {
	ok  bool
	err string
}

// Event reference validation keeps all external references on the same safe
// identifier policy before the event can contribute evidence.
func validateEventRefs(event Event) error {
	for _, check := range eventRefChecks(event) {
		if !check.ok {
			return errors.New(check.err)
		}
	}
	return nil
}

// eventRefChecks keeps the validation messages stable for event validation
// callers while sharing one policy list with validateEventRefs.
func eventRefChecks(event Event) []eventRefCheck {
	return []eventRefCheck{
		{safeRef(event.SourceRef), "unsafe source_ref"},
		{sha256Pattern.MatchString(event.SourceDigest), "invalid source_digest"},
		{event.TaskRef == "" || safeRef(event.TaskRef), "unsafe task_ref"},
		{event.OperationRef == "" || safeOperationRef(event.OperationRef), "unsafe operation_ref"},
		{event.ActorRef == "" || safeRef(event.ActorRef), "unsafe actor_ref"},
	}
}
