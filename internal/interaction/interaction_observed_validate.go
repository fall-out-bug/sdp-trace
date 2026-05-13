package interaction

import (
	"errors"
	"fmt"
)

func NewObservedEvent(opts RelayOptions, body []byte, sequence int) (Event, error) {
	// NewObservedEvent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateObservedEventOptions(opts); err != nil {
		return Event{}, err
	}
	if unsafeCount(body) > 0 {
		return Event{}, errors.New("interaction content contains unsafe material and cannot be retained")
	}
	return observedEvent(opts, body, sequence), nil
}

func validateObservedEventOptions(opts RelayOptions) error {
	// validateObservedEventOptions keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateObservedEventIDs(opts); err != nil {
		return err
	}

	return validateObservedEventCatalog(opts)
}

func validateObservedEventIDs(opts RelayOptions) error {
	// validateObservedEventIDs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return err
	}
	if err := validateSafeID("target", opts.Target); err != nil {
		return err
	}
	return validateObservedActorID(opts)
}
func validateObservedActorID(opts RelayOptions) error {
	// validateObservedActorID keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if opts.ActorID == "" {

		opts.ActorID = opts.ActorType
	}
	return validateSafeID("actor_id", opts.ActorID)
}
func validateObservedEventCatalog(opts RelayOptions) error {
	// validateObservedEventCatalog keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if !validActorType(opts.ActorType) {
		return fmt.Errorf("unsupported actor_type %q", opts.ActorType)
	}
	if !validEventType(opts.EventType) {
		return fmt.Errorf("unsupported event_type %q", opts.EventType)
	}
	return nil
}
