package harnessobs

import (
	"encoding/json"
)

func addNormalizedSourceDigests(events []Event) ([]Event, error) {
	// addNormalizedSourceDigests keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i := range events {
		data, err := json.Marshal(events[i])
		if err != nil {
			return nil, err
		}

		events[i].SourceDigest = digestLine(data)
	}
	return events, nil
}
