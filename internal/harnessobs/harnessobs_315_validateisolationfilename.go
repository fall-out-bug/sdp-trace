package harnessobs

import (
	"errors"

	"strings"
)

func validateIsolationFilename(base string) error {
	// validateIsolationFilename keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(base) == "" || strings.ContainsAny(base, `/\`) {

		return errors.New("unsafe isolation filename")
	}
	return nil
}
