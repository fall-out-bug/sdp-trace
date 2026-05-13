package trace

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func writeCanonicalScalar(buf *bytes.Buffer, value any) error {
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
func invalidFloat(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}
func writeNumericScalar(buf *bytes.Buffer, value any) bool {
	// writeNumericScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if number, ok := value.(json.Number); ok {
		writeJSONNumber(buf, number)
		return true
	}
	reflected := reflect.ValueOf(value)
	if !reflected.IsValid() {
		return false
	}
	return writeReflectedNumericScalar(buf, reflected)
}
func writeReflectedNumericScalar(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedNumericScalar preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if writeReflectedFloat(buf, reflected) {
		return true
	}
	if writeReflectedSignedInteger(buf, reflected) {
		return true
	}
	if writeReflectedUnsignedInteger(buf, reflected) {
		return true
	}
	return false
}
func writeReflectedFloat(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedFloat preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch reflected.Kind() {
	case reflect.Float32:
		writeJSONNumber(buf, trimFloatToJSON(reflected.Convert(reflect.TypeOf(float64(0))).Float()))
		return true
	case reflect.Float64:
		writeJSONNumber(buf, trimFloatToJSON(reflected.Float()))
		return true
	default:
		return false
	}
}
func writeReflectedSignedInteger(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedSignedInteger preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch reflected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		writeJSONInteger(buf, reflected.Int())
		return true
	default:
		return false
	}
}
func writeReflectedUnsignedInteger(buf *bytes.Buffer, reflected reflect.Value) bool {
	// writeReflectedUnsignedInteger preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	switch reflected.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		writeJSONUnsignedInteger(buf, reflected.Uint())
		return true
	default:
		return false
	}
}
func writeJSONUnsignedInteger(buf *bytes.Buffer, value uint64) {
	buf.WriteString(strconv.FormatUint(value, 10))
}
func writeJSONBool(buf *bytes.Buffer, value bool) {
	// Booleans are emitted as fixed JSON literals so payload and event digest
	// material never depends on encoder configuration.
	if value {
		buf.WriteString("true")
		return
	}
	buf.WriteString("false")
}
func writeJSONFallback(buf *bytes.Buffer, value any) error {
	// writeJSONFallback preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = buf.Write(raw)
	return err
}
func writeJSONString(buffer *bytes.Buffer, value string) {
	buffer.WriteString(strconv.Quote(value))
}
func writeJSONInteger(buffer *bytes.Buffer, value int64) {
	buffer.WriteString(strconv.FormatInt(value, 10))
}
func writeJSONNumber(buffer *bytes.Buffer, value json.Number) {
	// writeJSONNumber preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// Decimal and exponent spellings normalize through float parsing when
	// possible; integer-like spellings remain unchanged.

	num := strings.TrimSpace(value.String())
	if num == "" {
		buffer.WriteString("0")
		return
	}
	if strings.ContainsAny(num, ".eE") {
		if normalized, err := strconv.ParseFloat(num, 64); err == nil {
			buffer.WriteString(trimFloatToString(normalized))
			return
		}
	}
	buffer.WriteString(num)
}
func trimFloatToJSON(value float64) json.Number {
	return json.Number(trimFloatToString(value))
}
func trimFloatToString(value float64) string {
	// trimFloatToString preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.

	if invalidFloat(value) {
		return "0"
	}
	if value == math.Trunc(value) {
		return strconv.FormatInt(int64(value), 10)
	}
	if math.Abs(value) < 1e-15 {
		return "0"
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}
