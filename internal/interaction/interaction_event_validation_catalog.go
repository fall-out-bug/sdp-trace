package interaction

func validateEventCatalog(event Event) error {
	// validateEventCatalog keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateSafeID("actor.id", event.Actor.ID); err != nil {
		return err
	}
	if err := validateEventTypeAndFriction(event); err != nil {
		return err
	}
	if err := validateEventActorAndState(event); err != nil {
		return err
	}
	return validateEventRetentionStates(event)
}
