package trace

import (
	"bytes"
	"reflect"
)

// This file owns reflected integer spelling for canonical trace rendering.

func writeReflectedSignedInteger(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedSignedInteger preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	switch reflected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// All signed integer widths share one canonical base-10 spelling.
		// Width is type metadata and never part of trace digest material.
		writeJSONInteger(buf, reflected.Int())
		return true
	default:
		// Non-signed kinds are checked by the unsigned helper or fallback path.
		// Returning false preserves that ordering.
		return false
	}
}

func writeReflectedUnsignedInteger(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedUnsignedInteger preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	switch reflected.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// All unsigned integer widths share one canonical base-10 spelling.
		writeJSONUnsignedInteger(buf, reflected.Uint())
		return true
	default:
		// Non-unsigned kinds are not numeric values handled here.
		// Returning false allows generic JSON fallback where appropriate.
		return false
	}
}
