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
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
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

func writeCanonicalList(buf *bytes.Buffer, value []any) error {
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
	switch typed := value.(type) {
	case string:
		writeJSONString(buf, typed)
	case bool:
		writeJSONBool(buf, typed)
	case nil:
		buf.WriteString("null")
	default:
		if writeNumericScalar(buf, typed) {
			return nil
		}
		return writeJSONFallback(buf, typed)
	}
	return nil
}

func writeNumericScalar(buf *bytes.Buffer, value any) bool {
	if number, ok := value.(json.Number); ok {
		writeJSONNumber(buf, number)
		return true
	}
	reflected := reflect.ValueOf(value)
	if !reflected.IsValid() {
		return false
	}
	switch reflected.Kind() {
	case reflect.Float32, reflect.Float64:
		writeJSONNumber(buf, trimFloatToJSON(reflected.Convert(reflect.TypeOf(float64(0))).Float()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		writeJSONInteger(buf, reflected.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(reflected.Uint(), 10))
	default:
		return false
	}
	return true
}

func writeJSONBool(buf *bytes.Buffer, value bool) {
	if value {
		buf.WriteString("true")
		return
	}
	buf.WriteString("false")
}

func writeJSONFallback(buf *bytes.Buffer, value any) error {
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
	if math.IsNaN(value) || math.IsInf(value, 0) {
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
