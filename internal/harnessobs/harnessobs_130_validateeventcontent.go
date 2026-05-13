package harnessobs

import (
	"errors"
)

func validateEventContent(event Event) error {
	// validateEventContent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !validContentState(event.ContentState) {

		return errors.New("invalid content_state")
	}
	return nil
}
