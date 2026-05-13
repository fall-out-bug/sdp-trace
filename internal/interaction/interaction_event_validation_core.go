package interaction

func ValidateEvent(event Event) error {
	// ValidateEvent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	return firstValidationError(
		func() error { return validateEventIdentity(event) },
		func() error { return validateEventCatalog(event) },
		func() error { return validateEventSource(event) },
		func() error { return validateEventContent(event) },
		func() error { return validateEventTiming(event) },
		func() error { return validateEventRefs(event) },
	)
}

func firstValidationError(checks ...func() error) error {
	// firstValidationError keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}
