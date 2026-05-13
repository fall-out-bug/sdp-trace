package trace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Canonicalization constants used by Slice A event hashes.
const (
	HashAlgSHA256        = "sha256"
	CanonicalSchemaAlgo  = "json-canonicalization-scheme"
	CanonicalAlgoVersion = "1.0.0"
)

// CanonicalJSON encodes a value in deterministic key order.
func CanonicalJSON(value any) ([]byte, error) {
	// Canonical JSON is the hash boundary for event and payload digests, so
	// callers get normalized map order and number formatting before hashing.
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
	// Marshal/unmarshal through json.Number to remove Go type-specific map and
	// numeric representations before canonical rendering.
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
	// Canonicalization dispatch is deliberately type-small: objects, lists, and
	// scalar leaves cover the replayed JSON value space.
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
	// Object keys are sorted before rendering so equivalent payloads hash the
	// same across Go map iteration orders.
	// Values are rendered through the canonical dispatcher, preserving nested
	// maps and arrays under the same rules.
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
	// Return a new key slice; callers must not depend on map iteration order for
	// digest material.
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func writeCanonicalList(buf *bytes.Buffer, value []any) error {
	// Array order is evidence content and is preserved exactly while each item is
	// canonicalized recursively.
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

func writeCanonicalScalar(buf *bytes.Buffer, value any) error {
	return writeCanonicalScalarValue(buf, value)
}

func writeCanonicalScalarValue(buf *bytes.Buffer, value any) error {
	if value == nil {
		// Nil payload fields are canonical JSON null, not omitted values.
		writeJSONNull(buf)
		return nil
	}
	return writeNonNilCanonicalScalar(buf, value)
}

func writeNonNilCanonicalScalar(buf *bytes.Buffer, value any) error {
	// Strings and booleans are rendered directly; all numeric and uncommon
	// scalar forms go through the fallback path.
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
	if writeNumericScalar(buf, value) {
		// Numeric scalars use normalized JSON forms before the generic encoder
		// can introduce implementation-dependent formatting.
		return nil
	}
	return writeJSONFallback(buf, value)
}

func invalidFloat(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}

func writeNumericScalar(buf *bytes.Buffer, value any) bool {
	// Numeric detection accepts json.Number first because decoded JSON values use
	// it to preserve source spelling until normalization.
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
	// Reflection covers typed numeric payloads that did not originate from JSON
	// decoding but still need canonical digest rendering.
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
	// Float32 values are widened before rendering so both float widths share the
	// same trim and invalid-value policy.
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
	// All signed integer widths collapse to the same base-10 canonical form.
	switch reflected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		writeJSONInteger(buf, reflected.Int())
		return true
	default:
		return false
	}
}

func writeReflectedUnsignedInteger(buf *bytes.Buffer, reflected reflect.Value) bool {
	// Unsigned integer widths also collapse to base-10 without type decoration.
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
	if value {
		// Canonical booleans are emitted directly to avoid encoder whitespace or
		// map-order side effects around scalar hashing.
		buf.WriteString("true")
		return
	}
	buf.WriteString("false")
}

func writeJSONFallback(buf *bytes.Buffer, value any) error {
	// The fallback is for scalar forms already outside the trust-critical numeric
	// paths; marshal errors still propagate to the digest caller.
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
	// Empty or exponent/decimal numbers are normalized so equivalent JSON number
	// values do not hash differently because of spelling.
	// Invalid decimal parsing falls back to the original token so decode-time
	// JSON validation remains the source of truth.
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
	// Non-finite and near-zero values are collapsed to zero to keep digest output
	// JSON-compatible and deterministic.
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

func eventHashHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// eventForCanonicalizing rebuilds events into a map excluding event_hash.
func eventForCanonicalizing(event Event) (map[string]any, error) {
	// The stored event_hash is excluded from its own digest input; all other
	// serialized event fields remain part of the canonical hash.
	// Decoding with UseNumber avoids losing integer spelling before the canonical
	// renderer normalizes numbers.
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

// computeEventHash returns a canonical hash over events excluding event_hash.
func computeEventHash(event Event) (string, error) {
	// Event hashing reuses the same canonical JSON renderer as payload hashing
	// so replay does not depend on struct or map ordering.
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

// ComputeEventHash returns the hash for trace.Event.
func ComputeEventHash(event Event) (string, error) {
	// The algorithm branch is intentionally boring today; the public API keeps a
	// stable extension point while SHA-256 remains the only accepted algorithm.
	hashHex, err := computeEventHash(event)
	if err != nil {
		return "", err
	}
	if event.HashAlgorithm == HashAlgSHA256 {
		return hashHex, nil
	}
	return hashHex, nil
}

// CanonicalEventPayloadDigest computes a deterministic digest over the payload section.
func CanonicalEventPayloadDigest(payload json.RawMessage) (string, error) {
	// Empty payload bytes are rejected because an event payload digest without
	// payload material would be unverifiable.
	// Payload bytes are decoded and re-rendered canonically before hashing, so
	// whitespace and map order do not affect the digest.
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

// ComputeEventPayloadDigest returns a deterministic digest for payload bytes.
func ComputeEventPayloadDigest(event Event) (string, error) {
	return CanonicalEventPayloadDigest(event.Payload)
}
