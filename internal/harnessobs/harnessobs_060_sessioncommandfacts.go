package harnessobs

import (
	"encoding/json"
)

func sessionCommandFacts(session SessionRun) []Event {
	// sessionCommandFacts keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !hasSessionCommandModel(session) {

		return nil
	}
	event := sessionCommandModelEvent(session)
	data, err := json.Marshal(event)
	if err != nil {
		return nil
	}

	event.SourceDigest = digestLine(data)

	return []Event{event}
}
