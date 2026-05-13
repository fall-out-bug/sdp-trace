package trace

import "fmt"

func ValidateRunDirectory(path string, requireChain bool) error {
	// ValidateRunDirectory preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	runArtifact, err := OpenRunArtifact(path)
	if err != nil {
		return err
	}
	return validateRunArtifact(runArtifact, requireChain)
}
func validateRunArtifact(artifact RunArtifact, requireChain bool) error {
	// validateRunArtifact preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	if err := artifact.Manifest.Validate(); err != nil {
		return err
	}
	return validateRunDirectoryState(artifact.Manifest, artifact.Events, requireChain)
}

func validateManifestEventCount(manifestCount int, eventCount int) error {
	// validateManifestEventCount preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	if manifestCount != 0 && manifestCount != eventCount {
		return fmt.Errorf("event_count mismatch: run.json=%d files=%d", manifestCount, eventCount)
	}
	return nil
}

func validateRunDirectoryState(manifest RunManifest, events []Event, requireChain bool) error {
	// validateRunDirectoryState preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	if err := validateEventChainIfRequested(events, requireChain); err != nil {
		return err
	}
	if err := validateManifestEventCount(manifest.EventCount, len(events)); err != nil {
		return err
	}
	return validateManifestEventChainHead(manifest.EventChainHead, events)
}
