package interaction

import (
	"errors"
	"fmt"
)

func validateEventSource(event Event) error {
	// validateEventSource keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if !validSourceType(event.Source.SourceType) {
		return fmt.Errorf("unsupported source_type %q", event.Source.SourceType)
	}
	if err := validateSafeID("source_id", event.SourceID); err != nil {
		return err
	}
	return nil
}

func validateEventContent(event Event) error {
	// validateEventContent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if err := validateEventDigest(event); err != nil {
		return err
	}

	return validateEventContentRef(event)
}

func validateEventDigest(event Event) error {
	// validateEventDigest keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if event.DigestAlgorithm != DigestAlgorithmSHA256 || !sha256Pattern.MatchString(event.ContentDigest) {

		return errors.New("interaction event requires sha256 content digest")
	}
	return nil
}

func validateEventContentRef(event Event) error {
	// validateEventContentRef keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateContentRefFormat(event.ContentRef); err != nil {
		return err
	}
	if event.ContentRef == "" && event.NotRetainedReason == "" {
		return errors.New("interaction event without content_ref requires not_retained_reason")
	}
	return nil
}
func validateContentRefFormat(contentRef string) error {
	// validateContentRefFormat keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if contentRef != "" && !contentRefPattern.MatchString(contentRef) {

		return fmt.Errorf("unsupported content_ref %q", contentRef)
	}
	return nil
}
