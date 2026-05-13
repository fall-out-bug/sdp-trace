package trace

import (
	"bytes"
	"sort"
)

// This file owns canonical object rendering for trace digest bytes.

func writeCanonicalMap(buf *bytes.Buffer, value map[string]any) error {
	// writeCanonicalMap preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// Map keys are sorted at this boundary because Go map iteration order is
	// explicitly not evidence.
	keys := sortedMapKeys(value)
	buf.WriteByte('{')
	// The loop emits one complete key/value pair per sorted key.
	for i, key := range keys {
		if i > 0 {
			// Separators are explicit so the encoder cannot add whitespace.
			buf.WriteByte(',')
		}
		writeJSONString(buf, key)
		buf.WriteByte(':')
		if err := writeCanonicalJSON(buf, value[key]); err != nil {
			// Nested rendering errors identify an unhashable child value.
			return err
		}
	}
	buf.WriteByte('}')
	// Object close is written by this helper so nested callers cannot omit it.
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
	// Sorting is lexical and stable across platforms.
	return keys
}
