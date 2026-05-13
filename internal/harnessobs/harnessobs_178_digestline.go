package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func digestLine(line []byte) string {
	// digestLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err == nil {
		raw["source_digest"] = ""
		canonical, err := json.Marshal(raw)
		if err == nil {

			sum := sha256.Sum256(canonical)
			return hex.EncodeToString(sum[:])
		}
	}

	sum := sha256.Sum256(line)
	return hex.EncodeToString(sum[:])
}
