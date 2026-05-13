package interaction

import (
	"bufio"
	"os"
)

func readJSONLEvents(path string) ([]Event, error) {
	// readJSONLEvents keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return scanJSONLEvents(file)
}
func scanJSONLEvents(file *os.File) ([]Event, error) {
	// scanJSONLEvents keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), MaxBodyBytes*4)
	events := make([]Event, 0)
	for scanner.Scan() {
		var err error
		events, err = appendJSONLEventLine(events, scanner.Text())
		if err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
