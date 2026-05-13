package trace

import (
	"bytes"
)

// This file owns scalar canonicalization for trace digests. Every helper keeps
// JSON-compatible spelling deterministic before bytes are hashed.

func writeCanonicalScalar(buf *bytes.Buffer, value any) error {
	// Keep the public scalar dispatch separate from nil/non-nil handling.
	return writeCanonicalScalarValue(buf, value)
}
func writeCanonicalScalarValue(buf *bytes.Buffer, value any) error {
	// writeCanonicalScalarValue preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if value == nil {
		writeJSONNull(buf)
		return nil
	}
	return writeNonNilCanonicalScalar(buf, value)
}
func writeNonNilCanonicalScalar(buf *bytes.Buffer, value any) error {
	// writeNonNilCanonicalScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch typed := value.(type) {
	case string:
		writeJSONString(buf, typed)
		return nil
	case bool:
		writeJSONBool(buf, typed)
		return nil
	default:
		return writeCanonicalFallbackScalar(buf, typed)
	}
}
func writeJSONNull(buf *bytes.Buffer) {
	// Null is an explicit payload value, not an omitted digest field.
	buf.WriteString("null")
}
func writeCanonicalFallbackScalar(buf *bytes.Buffer, value any) error {
	// writeCanonicalFallbackScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if writeNumericScalar(buf, value) {
		return nil
	}
	return writeJSONFallback(buf, value)
}
