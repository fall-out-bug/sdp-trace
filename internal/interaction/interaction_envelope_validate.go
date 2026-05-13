package interaction

import (
	"errors"
	"fmt"
)

func ValidateEnvelope(envelope Envelope) error {
	// ValidateEnvelope keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if envelope.SchemaVersion != SchemaVersion {
		return errors.New("delivery envelope has unsupported schema_version")
	}
	if err := validateEnvelopeIdentity(envelope); err != nil {
		return err
	}
	if err := validateEnvelopeRunRefs(envelope.RunRefs); err != nil {
		return err
	}
	return validateEnvelopeRefs(envelope)
}

func validateEnvelopeIdentity(envelope Envelope) error {
	// validateEnvelopeIdentity keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateSafeID("task_id", envelope.TaskID); err != nil {

		return err
	}
	return validateSafeID("envelope_id", envelope.EnvelopeID)
}
func validateEnvelopeRunRefs(refs []string) error {
	// validateEnvelopeRunRefs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, ref := range refs {
		if !runRefPattern.MatchString(ref) {
			return fmt.Errorf("unsupported run_ref %q", ref)
		}
	}
	return nil
}

func validateEnvelopeRefs(envelope Envelope) error {
	// validateEnvelopeRefs keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, refs := range envelopeReferenceGroups(envelope) {
		for _, ref := range refs {
			if !validReference(ref) {
				return fmt.Errorf("unsupported envelope ref %q", ref)
			}
		}
	}
	return nil
}
func envelopeReferenceGroups(envelope Envelope) [][]string {
	return [][]string{envelope.OperationRefs, envelope.ToolRefs, envelope.MutationRefs, envelope.LLMRefs, envelope.InteractionRefs, envelope.PromiseRefs, envelope.StageRefs, envelope.FrictionRefs, envelope.TaskRefs, envelope.SourceRefs}
}
