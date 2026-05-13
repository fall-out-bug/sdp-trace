package harnessobs

import (
	"strings"
)

func unretainedRawToolInputField(path, key string, value any) bool {
	// unretainedRawToolInputField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if key != "prompt" {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}

	segments := strings.Split(path, ".")
	if len(segments) < 3 {
		return false
	}
	return path == "part.state.input.prompt"
}
