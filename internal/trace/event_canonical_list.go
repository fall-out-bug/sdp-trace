package trace

import "bytes"

// This file owns canonical array rendering for trace digest bytes.

func writeCanonicalList(buf *bytes.Buffer, value []any) error {
	// writeCanonicalList preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	buf.WriteByte('[')
	// The loop preserves array order because array position is evidence.
	for i, item := range value {
		if i > 0 {
			// Array separators are explicit for the same byte-for-byte digest.
			buf.WriteByte(',')
		}
		if err := writeCanonicalJSON(buf, item); err != nil {
			// Fail the whole payload when any child cannot be canonicalized.
			return err
		}
	}
	buf.WriteByte(']')
	// Array close is written here so recursive callers receive complete JSON.
	return nil
}
