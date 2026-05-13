package harnessobs

func writeNormalizedEvents(outPath string, events []Event) error {
	// writeNormalizedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	out, err := createNormalizedEventsFile(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, event := range events {

		if err := writeNormalizedEvent(out, event); err != nil {
			return err
		}
	}
	return nil
}
