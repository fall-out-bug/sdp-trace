package trace

import (
	"bytes"
	"encoding/json"
	"errors"
)

// This file owns payload digest construction for trace event payload bytes.

func CanonicalEventPayloadDigest(payload json.RawMessage) (string, error) {
	// CanonicalEventPayloadDigest preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	if len(payload) == 0 {
		// Missing payload bytes cannot produce replayable digest evidence.
		return "", errors.New("payload is required")
	}
	var decoded any
	// UseNumber preserves numeric payload tokens until canonical rendering.
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	if err := dec.Decode(&decoded); err != nil {
		// Payload must decode before map-order normalization can happen.
		// Invalid payload JSON cannot produce a trustworthy digest.
		return "", err
	}
	canonical, err := CanonicalJSON(decoded)
	if err != nil {
		// Canonical payload rendering is the authority for the digest bytes.
		return "", err
	}
	// Payload digest uses the same hash encoding as event hashes.
	return eventHashHex(canonical), nil
}

func ComputeEventPayloadDigest(event Event) (string, error) {
	// Event payload bytes are already synchronized by Event helpers.
	// This wrapper keeps the Event method surface small and delegates policy.
	return CanonicalEventPayloadDigest(event.Payload)
}
