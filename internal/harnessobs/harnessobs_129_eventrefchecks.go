package harnessobs

func eventRefChecks(event Event) []eventRefCheck {
	// eventRefChecks keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return []eventRefCheck{
		{safeRef(event.SourceRef), "unsafe source_ref"},
		{sha256Pattern.MatchString(event.SourceDigest), "invalid source_digest"},
		{event.TaskRef == "" || safeRef(event.TaskRef), "unsafe task_ref"},
		{event.OperationRef == "" || safeOperationRef(event.OperationRef), "unsafe operation_ref"},
		{event.ActorRef == "" || safeRef(event.ActorRef), "unsafe actor_ref"},
	}
}
