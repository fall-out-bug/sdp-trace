package trace

import "fmt"

// This file owns contiguous event-chain validation.

// ValidateEventChain checks that hashes and prev-event links are consistent.
func ValidateEventChain(events []Event) error {
	// Chain replay starts at the explicit null sentinel.
	prevEventHash := NullEventHash
	for i, event := range events {
		// Each verified event becomes the expected predecessor for the next row.
		if err := validateChainEvent(i, event, prevEventHash); err != nil {
			return err
		}
		prevEventHash = event.EventHash
	}
	// Reaching the end means every event linked to its predecessor.
	return nil
}

func validateChainEvent(index int, event Event, prevEventHash string) error {
	// validateChainEvent preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.
	computed, err := event.WithComputedEventHash()
	if err != nil {
		// Hash recomputation failure blocks all later chain checks.
		return fmt.Errorf("event %d (%s) hash generation failed: %w", index, event.EventID, err)
	}
	if err := validateEventHash(index, event, computed.EventHash); err != nil {
		// Event hash mismatch is the strongest local tamper signal.
		return err
	}
	if err := event.VerifyPayloadDigest(); err != nil {
		// Payload digest mismatch means the event body is not replayable.
		return fmt.Errorf("event %d (%s) payload_digest invalid: %w", index, event.EventID, err)
	}
	// Position validation binds the event to the previous verified hash.
	return validateEventPosition(index, event, prevEventHash)
}

func validateEventHash(index int, event Event, computedHash string) error {
	// validateEventHash preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.
	if event.EventHash != computedHash {
		return fmt.Errorf("event %d (%s) event_hash mismatch: expected %s got %s", index, event.EventID, computedHash, event.EventHash)
	}
	return nil
}

func validateEventPosition(index int, event Event, prevEventHash string) error {
	// validateEventPosition preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.
	if event.Sequence != index {
		return fmt.Errorf("event %d has non-zero-based sequence %d", index, event.Sequence)
	}
	if event.PrevEventHash != prevEventHash {
		return fmt.Errorf("event %d (%s) prev_event_hash expected %s", index, event.EventID, prevEventHash)
	}
	return nil
}
