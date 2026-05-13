package interaction

import (
	"errors"
)

func normalizeTranscriptEvents(events []Event, opts ImportOptions) error {
	// normalizeTranscriptEvents keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if len(events) == 0 {
		return errors.New("interaction import-transcript requires at least one event")
	}
	for i := range events {
		if err := normalizeTranscriptEvent(&events[i], opts); err != nil {
			return err
		}
	}
	return validateOrdering(events)
}
