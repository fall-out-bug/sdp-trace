package interaction

import (
	"fmt"
)

func validateEventActorAndState(event Event) error {
	// validateEventActorAndState keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if !validActorType(event.Actor.ActorType) {
		return fmt.Errorf("unsupported actor_type %q", event.Actor.ActorType)
	}
	if !validEventState(event.State) {
		return fmt.Errorf("unsupported state %q", event.State)
	}
	return nil
}

func validateEventRetentionStates(event Event) error {
	// validateEventRetentionStates keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if !validRetention(event.Retention) {
		return fmt.Errorf("unsupported retention %q", event.Retention)
	}
	if !validCompleteness(event.CompletenessState) {
		return fmt.Errorf("unsupported completeness_state %q", event.CompletenessState)
	}
	if !validChannelExclusivity(event.ChannelExclusivity) {
		return fmt.Errorf("unsupported channel_exclusivity_state %q", event.ChannelExclusivity)
	}
	return nil
}
