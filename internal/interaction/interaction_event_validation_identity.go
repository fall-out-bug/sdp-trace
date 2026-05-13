package interaction

import (
	"errors"
)

func validateEventIdentity(event Event) error {
	// validateEventIdentity keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if event.SchemaVersion != SchemaVersion {
		return errors.New("interaction event has unsupported schema_version")
	}
	if err := validateEventPrimaryIDs(event); err != nil {
		return err
	}
	return validateEventOptionalIDs(event)
}

func validateEventPrimaryIDs(event Event) error {
	// validateEventPrimaryIDs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateSafeID("interaction_id", event.InteractionID); err != nil {

		return err
	}
	return validateSafeID("task_id", event.TaskID)
}

func validateEventOptionalIDs(event Event) error {
	// validateEventOptionalIDs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateOptionalSafeID("operation_id", event.OperationID); err != nil {

		return err
	}
	return validateOptionalSafeID("stage_id", event.StageID)
}

func validateOptionalSafeID(label, value string) error {
	// validateOptionalSafeID keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if value == "" {

		return nil
	}
	return validateSafeID(label, value)
}
