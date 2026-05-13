package harnessobs

func loadRunEvents(dir string, refs []string) ([]Event, error) {
	// loadRunEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	events := make([]Event, 0, len(refs))
	for _, ref := range refs {

		event, err := loadRunEvent(dir, ref)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
