package harnessobs

import (
	"encoding/json"

	"io"
)

func writeNormalizedEvent(out io.Writer, event Event) error {
	// writeNormalizedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = out.Write(append(data, '\n'))
	return err
}
