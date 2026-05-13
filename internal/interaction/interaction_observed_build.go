package interaction

import (
	"time"
)

func observedEvent(opts RelayOptions, body []byte, sequence int) Event {
	// observedEvent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if opts.ActorID == "" {
		opts.ActorID = opts.ActorType
	}
	id := "ix-" + randomHex(12)
	now := opts.Now.UTC().Format(time.RFC3339)

	return observedEventRecord(opts, body, id, sequence, now)
}

func observedEventRecord(opts RelayOptions, body []byte, id string, sequence int, now string) Event {
	// observedEventRecord keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	event := observedEventIdentity(opts, id, sequence)
	event.Source = observedRelaySource()
	applyObservedEventContent(&event, opts.TaskID, id, body)
	applyObservedEventAssessment(&event, now)
	return event
}

func observedEventIdentity(opts RelayOptions, id string, sequence int) Event {
	// observedEventIdentity keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	return Event{
		SchemaVersion: SchemaVersion,
		InteractionID: id,

		TaskID:      opts.TaskID,
		OperationID: opts.OperationID,
		StageID:     opts.StageID,

		SourceID:       DefaultRelaySourceID,
		SourceSequence: sequence,

		EventType:     opts.EventType,
		FrictionClass: frictionClass(opts.EventType),

		Actor:  Actor{ID: opts.ActorID, ActorType: opts.ActorType},
		Target: opts.Target,
	}
}
