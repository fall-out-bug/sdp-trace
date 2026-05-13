package trace

import (
	"bytes"
	"encoding/json"
)

// This file owns object/list canonicalization. Scalar spelling and digest
// hashing live in separate files so replay logic stays reviewable by boundary.
// The functions here deliberately avoid encoding/json for final output because
// digest evidence requires byte-for-byte stable ordering and separators.

// CanonicalJSON encodes a value in deterministic key order.
func CanonicalJSON(value any) ([]byte, error) {
	// This is the public canonicalization boundary used by event hashes and
	// payload digests. Keep normalization, recursive rendering, and returned
	// bytes together so callers cannot accidentally hash Go map iteration order,
	// raw numeric spelling, or encoder-specific whitespace.
	normalized, err := normalizeJSONValue(value)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	// Rendering writes into a fresh buffer so callers cannot accidentally reuse
	// prior digest material.
	if err := writeCanonicalJSON(&buffer, normalized); err != nil {
		// Rendering failure means the normalized JSON tree still contains an
		// unsupported value shape.
		return nil, err
	}
	// Returned bytes are the exact digest input for callers.
	return buffer.Bytes(), nil
}
func normalizeJSONValue(value any) (any, error) {
	// normalizeJSONValue preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// The json.Number decoder prevents float conversion before canonical
	// rendering decides the final spelling.

	raw, err := json.Marshal(value)
	if err != nil {
		// Marshal errors mean the caller supplied a value that cannot become
		// replayable JSON.
		return nil, err
	}
	var normalized any
	// Decode through json.Decoder rather than json.Unmarshal so UseNumber is
	// available for later numeric spelling decisions.
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&normalized); err != nil {
		// Decode errors keep invalid JSON-shaped values out of hash material.
		return nil, err
	}
	// The normalized tree contains map[string]any, []any, json.Number, and JSON
	// scalar leaves only.
	return normalized, nil
}
func writeCanonicalJSON(buf *bytes.Buffer, value any) error {
	// writeCanonicalJSON preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch typed := value.(type) {
	case map[string]any:
		// Objects need key sorting before recursive rendering.
		return writeCanonicalMap(buf, typed)
	case []any:
		// Arrays keep source order but recursively canonicalize each item.
		return writeCanonicalList(buf, typed)
	default:
		// Scalars own their spelling rules in event_scalar.go.
		return writeCanonicalScalar(buf, typed)
	}
}
