package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func digestCommand(command []string) string {
	// digestCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	data, _ := json.Marshal(command)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
