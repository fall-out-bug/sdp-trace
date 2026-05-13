package trace

import (
	"bytes"
	"encoding/json"
	"sort"
)

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
	if err := writeCanonicalJSON(&buffer, normalized); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func normalizeJSONValue(value any) (any, error) {
	// normalizeJSONValue preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var normalized any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}
func writeCanonicalJSON(buf *bytes.Buffer, value any) error {
	// writeCanonicalJSON preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch typed := value.(type) {
	case map[string]any:
		return writeCanonicalMap(buf, typed)
	case []any:
		return writeCanonicalList(buf, typed)
	default:
		return writeCanonicalScalar(buf, typed)
	}
}
func writeCanonicalMap(buf *bytes.Buffer, value map[string]any) error {
	// writeCanonicalMap preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	keys := sortedMapKeys(value)
	buf.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeJSONString(buf, key)
		buf.WriteByte(':')
		if err := writeCanonicalJSON(buf, value[key]); err != nil {
			return err
		}
	}
	buf.WriteByte('}')
	return nil
}
func sortedMapKeys(value map[string]any) []string {
	// sortedMapKeys preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
func writeCanonicalList(buf *bytes.Buffer, value []any) error {
	// writeCanonicalList preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	buf.WriteByte('[')
	for i, item := range value {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := writeCanonicalJSON(buf, item); err != nil {
			return err
		}
	}
	buf.WriteByte(']')
	return nil
}
