package harnessobs

import (
	"fmt"
)

func validateLoadedRun(run Run) error {
	// validateLoadedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if run.SchemaVersion != RunSchemaVersion {

		return fmt.Errorf("unsupported run schema_version: %s", run.SchemaVersion)
	}
	return nil
}
