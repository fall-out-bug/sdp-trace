package harnessobs

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func writeNormalizedEvents(outPath string, events []Event) error {
	// Create the destination only after normalization has produced a complete
	// event list; callers should not observe partially normalized raw input as
	// trusted evidence.
	out, err := createNormalizedEventsFile(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, event := range events {
		// Each event is written as one compact JSONL record so source digests can
		// be replayed line-by-line by the normal event scanner.
		if err := writeNormalizedEvent(out, event); err != nil {
			return err
		}
	}
	return nil
}

func createNormalizedEventsFile(outPath string) (*os.File, error) {
	// Raw normalization can materialize nested event_source_path values after
	// profile-relative output-file safety has already accepted the path.
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return nil, err
	}
	return os.Create(outPath)
}

func writeNormalizedEvent(out io.Writer, event Event) error {
	// Marshal before appending the newline so digest calculation and JSONL
	// formatting use the same compact representation.
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = out.Write(append(data, '\n'))
	return err
}
