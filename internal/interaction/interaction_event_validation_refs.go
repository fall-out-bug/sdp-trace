package interaction

import (
	"errors"
	"fmt"
	"time"
)

func validateEventTiming(event Event) error {
	// validateEventTiming keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if _, err := time.Parse(time.RFC3339, event.ObservedAt); err != nil {
		return errors.New("observed_at must be RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, event.CreatedAt); err != nil {
		return errors.New("created_at must be RFC3339")
	}
	if event.SourceSequence < 0 {
		return errors.New("source_sequence must be non-negative")
	}
	return nil
}

func validateEventRefs(event Event) error {
	// validateEventRefs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateReferenceRefs(event.ReferenceRefs); err != nil {

		return err
	}
	return validateLLMRefs(event.LLMRefs)
}

func validateReferenceRefs(refs []string) error {
	// validateReferenceRefs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, ref := range refs {
		if !validReference(ref) {
			return fmt.Errorf("unsupported reference_ref %q", ref)
		}
	}
	return nil
}

func validateLLMRefs(refs []LLMRef) error {
	// validateLLMRefs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, ref := range refs {
		if !validLLMLinkageState(ref.LinkageState) {
			return fmt.Errorf("unsupported llm linkage_state %q", ref.LinkageState)
		}
	}
	return nil
}

func validLLMLinkageState(state string) bool {
	return state == StateNotAssessed || state == StateCannotVerify || state == StateReferenced
}
