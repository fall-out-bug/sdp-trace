package harnessobs

import (
	"bytes"
	"encoding/json"

	"errors"

	"io"
)

func decodeStrictJSON(data []byte, target any) error {
	// decodeStrictJSON keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	// A second decode distinguishes trailing JSON from ordinary EOF.
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("strict JSON input contains trailing data")
	}
	return nil
}
