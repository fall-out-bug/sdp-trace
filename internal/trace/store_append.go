package trace

// AppendRunEvent appends one local event to an existing run artifact and updates the run manifest chain head.
func AppendRunEvent(runDir string, eventType EventType, payload map[string]any, observedBy string) (Event, error) {
	// Appends always start from persisted run state so sequence and previous hash
	// reflect retained evidence rather than caller memory.
	artifact, err := OpenRunArtifact(runDir)
	if err != nil {
		return Event{}, err
	}
	return appendEventToArtifact(artifact, eventType, payload, observedBy)
}

func appendEventToArtifact(artifact RunArtifact, eventType EventType, payload map[string]any, observedBy string) (Event, error) {
	// The event is fully hashed before any filesystem write is attempted.
	event, err := newAppendedRunEvent(artifact, eventType, payload, observedBy)
	if err != nil {
		return Event{}, err
	}
	return persistAppendedEvent(artifact, event)
}

func persistAppendedEvent(artifact RunArtifact, event Event) (Event, error) {
	// Event write precedes manifest advancement so a failed write cannot move
	// chain heads past retained event files.
	if err := artifact.Layout.WriteEvent(event); err != nil {
		return Event{}, err
	}
	return appendRunEventManifest(artifact, event)
}

func appendRunEventManifest(artifact RunArtifact, event Event) (Event, error) {
	// Manifest counters and chain heads advance together to one latest event.
	artifact.Manifest.EventCount = event.Sequence + 1
	artifact.Manifest.EventChainHead = event.EventHash
	artifact.Manifest.FinalChainHead = event.EventHash
	if err := artifact.Layout.WriteRun(artifact.Manifest); err != nil {
		return Event{}, err
	}
	return event, nil
}
