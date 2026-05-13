package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func validationDigest(validation Validation) string {
	// validationDigest keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	copy := validation
	copy.ValidationDigest = ""

	data, _ := json.Marshal(copy)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
