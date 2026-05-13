package trace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// This file owns digest construction. Canonical rendering is delegated so hash
// callers share one replayable byte representation.
// Event and payload hashes intentionally share the same byte and hex policy.

func eventHashHex(data []byte) string {
	// Hex output is lowercase and stable for manifest and event-chain fields.
	// SHA-256 is computed over canonical bytes supplied by the caller.
	sum := sha256.Sum256(data)
	// Encode the raw digest bytes, never the formatted input string.
	// All callers receive the same lowercase hex representation.
	return hex.EncodeToString(sum[:])
}

// eventForCanonicalizing rebuilds events into a map excluding event_hash.
func eventForCanonicalizing(event Event) (map[string]any, error) {
	// eventForCanonicalizing preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// The event hash excludes the stored event_hash field itself while retaining
	// all other encoded fields exactly as the JSON boundary sees them.

	raw, err := json.Marshal(event)
	if err != nil {
		// Event structs must be JSON-shapeable before they can be hashed.
		return nil, err
	}
	var decoded map[string]any
	// UseNumber keeps numeric tokens under the same canonical scalar policy as
	// payload digests.
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&decoded); err != nil {
		// A failed decode means canonical hash input could not be reconstructed.
		return nil, err
	}
	// The digest field is not allowed to prove itself.
	delete(decoded, "event_hash")
	// The remaining map is ready for deterministic key sorting.
	// No caller may re-add event_hash before canonical rendering.
	return decoded, nil
}

func computeEventHash(event Event) (string, error) {
	// computeEventHash preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// Event hashes are computed over the event object with event_hash removed.

	decoded, err := eventForCanonicalizing(event)
	if err != nil {
		// Decoding failure means the event cannot be reduced to canonical input.
		return "", err
	}
	canonical, err := CanonicalJSON(decoded)
	if err != nil {
		// Canonicalization failure blocks hash generation.
		return "", err
	}
	// The canonical bytes are the only event-hash input.
	return eventHashHex(canonical), nil
}

func ComputeEventHash(event Event) (string, error) {
	// ComputeEventHash preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// The public API does not inspect payload trust; it only hashes the event
	// shape supplied by the caller.

	hashHex, err := computeEventHash(event)
	if err != nil {
		// Preserve canonicalization errors for verifier diagnostics.
		return "", err
	}
	// Algorithm validation is separate from byte construction.
	if event.HashAlgorithm == HashAlgSHA256 {
		// SHA-256 is currently the only accepted algorithm.
		// Returning here keeps the accepted path explicit for future algorithms.
		return hashHex, nil
	}
	// Unknown algorithm handling remains non-upgrading for compatibility.
	// Callers validate algorithm acceptance elsewhere before trusting the event.
	return hashHex, nil
}
