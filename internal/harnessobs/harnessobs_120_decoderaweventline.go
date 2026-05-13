package harnessobs

import (
	"encoding/json"

	"fmt"
)

func decodeRawEventLine(line []byte, lineNo int) (map[string]any, error) {
	// decodeRawEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {

		return nil, fmt.Errorf("source line %d: malformed_jsonl", lineNo)
	}
	return raw, nil
}
