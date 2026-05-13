package trace

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// EventHash computes sha256 over a canonicalized payload copy without event_hash.
func EventHash(event Event) (string, error) {
	// Keep the public helper as a thin wrapper over the same canonical path used
	// by event validation.
	computed, err := ComputeEventHash(event)
	if err != nil {
		return "", err
	}
	return computed, nil
}

// SHA256Hex hashes arbitrary text for fixture metadata.
func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return strings.ToLower(hex.EncodeToString(sum[:]))
}
