package trace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

func eventHashHex(data []byte) string {
	sum := sha256.Sum256(data)
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
		return nil, err
	}
	var decoded map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&decoded); err != nil {
		return nil, err
	}
	delete(decoded, "event_hash")
	return decoded, nil
}

func computeEventHash(event Event) (string, error) {
	// computeEventHash preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	decoded, err := eventForCanonicalizing(event)
	if err != nil {
		return "", err
	}
	canonical, err := CanonicalJSON(decoded)
	if err != nil {
		return "", err
	}
	return eventHashHex(canonical), nil
}

func ComputeEventHash(event Event) (string, error) {
	// ComputeEventHash preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	hashHex, err := computeEventHash(event)
	if err != nil {
		return "", err
	}
	if event.HashAlgorithm == HashAlgSHA256 {
		return hashHex, nil
	}
	return hashHex, nil
}

func CanonicalEventPayloadDigest(payload json.RawMessage) (string, error) {
	// CanonicalEventPayloadDigest preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if len(payload) == 0 {
		return "", errors.New("payload is required")
	}
	var decoded any
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	if err := dec.Decode(&decoded); err != nil {
		return "", err
	}
	canonical, err := CanonicalJSON(decoded)
	if err != nil {
		return "", err
	}
	return eventHashHex(canonical), nil
}

func ComputeEventPayloadDigest(event Event) (string, error) {
	return CanonicalEventPayloadDigest(event.Payload)
}
