package harnessobs

import "path/filepath"

func writeObservationEvents(outDir string, events []Event) error {
	for _, event := range events {
		path := filepath.Join(outDir, "events", event.EventID+".json")
		if err := writeJSON(path, event); err != nil {
			return err
		}
	}
	return nil
}
