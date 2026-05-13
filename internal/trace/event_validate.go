package trace

import (
	"fmt"
)

// Validate checks local invariants before persistence and replay.
func (event Event) Validate() error {
	// Event validation keeps required identity/hash shape separate from payload
	// representation decoding.
	if err := event.validateRequiredFields(); err != nil {
		return err
	}
	if err := event.validateHashAlgorithm(); err != nil {
		return err
	}
	_, err := event.syncPayloadRepresentation()
	return err
}

func (event Event) validateRequiredFields() error {
	// Required fields are the minimum replay keys for event-chain verification.
	return firstValidationError(
		requiredString(event.SchemaVersion, "missing schema_version"),
		requiredString(event.RunID, "missing run_id"),
		requiredString(event.EventID, "missing event_id"),
		requiredString(event.Timestamp, "missing timestamp"),
		requiredString(event.EventHash, "missing event_hash"),
		validSequence(event.Sequence),
	)
}

func (event Event) validateHashAlgorithm() error {
	if event.HashAlgorithm != "" && event.HashAlgorithm != HashAlgSHA256 {
		// Empty preserves older events; explicit algorithms must match the only
		// hash supported by the replay contract.
		return fmt.Errorf("unsupported hash_algorithm %s", event.HashAlgorithm)
	}
	return nil
}
