package harnessobs

import (
	"errors"

	"os"

	"strings"
)

func readOptionalLines(path string) ([]string, error) {
	// readOptionalLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {

		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}
