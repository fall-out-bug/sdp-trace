package trace

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
)

// This file owns numeric scalar detection for canonical trace rendering.
// Numeric values are detected before generic JSON fallback so equivalent Go
// numeric types produce stable digest bytes.
// Final spelling remains delegated to event_json_write.go.

func invalidFloat(value float64) bool {
	// JSON has no NaN or infinity spelling, so those values collapse later.
	// The caller decides the replacement spelling.
	return math.IsNaN(value) || math.IsInf(value, 0)
}

func writeNumericScalar(buf *bytes.Buffer, value any) bool {
	// writeNumericScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	if number, ok := value.(json.Number); ok {
		// Decoded JSON numbers retain their source token until normalized here.
		// json.Number is the most source-faithful numeric representation.
		writeJSONNumber(buf, number)
		return true
	}
	reflected := reflect.ValueOf(value)
	// Reflection covers typed Go numeric values that bypass JSON decoding.
	if !reflected.IsValid() {
		// Invalid reflection values are not numeric scalars.
		return false
	}
	// A valid reflected value is tested by numeric kind, not by concrete type.
	return writeReflectedNumericScalar(buf, reflected)
}

func writeReflectedNumericScalar(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedNumericScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	if writeReflectedFloat(buf, reflected) {
		// Float kinds are terminal once matched.
		return true
	}
	if writeReflectedSignedInteger(buf, reflected) {
		// Signed integer kinds are terminal once matched.
		return true
	}
	if writeReflectedUnsignedInteger(buf, reflected) {
		// Unsigned integer kinds are terminal once matched.
		return true
	}
	// Non-numeric reflected values fall back to generic JSON rendering.
	return false
}

func writeReflectedFloat(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedFloat preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	switch reflected.Kind() {
	case reflect.Float32:
		// Float32 widens before spelling so both float widths share one policy.
		// The reflected conversion is explicit to avoid architecture surprises.
		writeJSONNumber(buf, trimFloatToJSON(reflected.Convert(reflect.TypeOf(float64(0))).Float()))
		return true
	case reflect.Float64:
		// Float64 is already in the canonical width used by trimFloatToJSON.
		// It still goes through trimFloatToJSON for non-finite and whole values.
		writeJSONNumber(buf, trimFloatToJSON(reflected.Float()))
		return true
	default:
		// Non-float kinds are checked by the integer helpers.
		// Returning false keeps the numeric dispatcher moving.
		return false
	}
}
