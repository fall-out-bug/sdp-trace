package interaction

import (
	"encoding/json"
	"os"
)

func ReadEnvelope(path string) (Envelope, error) {
	// Envelopes are validated on read for the same reason traces are: consumers
	// should not summarize invalid reference bundles.
	var envelope Envelope
	data, err := os.ReadFile(path)
	if err != nil {
		return Envelope{}, err
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return Envelope{}, err
	}
	if err := ValidateEnvelope(envelope); err != nil {
		return Envelope{}, err
	}
	return envelope, nil
}
