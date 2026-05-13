package trace

import "fmt"

func validateManifestEventChainHead(manifestHead string, events []Event) error {
	// validateManifestEventChainHead preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	if manifestHead == "" || len(events) == 0 {
		return nil
	}
	if events[len(events)-1].EventHash != manifestHead {
		return fmt.Errorf("run manifest event_chain_head does not match last event hash")
	}
	return nil
}
