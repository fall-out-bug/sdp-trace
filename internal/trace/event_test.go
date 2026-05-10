package trace

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestCanonicalJSONSortsKeysAndNormalizesNumbers(t *testing.T) {
	canonical, err := CanonicalJSON(map[string]any{
		"z": json.Number("1.000"),
		"a": []any{
			map[string]any{"b": true, "a": nil},
			float64(2.50),
		},
	})
	if err != nil {
		t.Fatalf("CanonicalJSON() error = %v", err)
	}
	want := `{"a":[{"a":null,"b":true},2.5],"z":1}`
	if string(canonical) != want {
		t.Fatalf("CanonicalJSON() = %s, want %s", canonical, want)
	}
}

func TestCanonicalEventPayloadDigestRejectsEmptyPayload(t *testing.T) {
	if _, err := CanonicalEventPayloadDigest(nil); err == nil {
		t.Fatalf("expected empty payload to be rejected")
	}
}

func TestWriteNumericScalarReflectDispatch(t *testing.T) {
	for name, value := range map[string]any{
		"float32": float32(3.25),
		"float64": float64(3.5),
		"int":     int(-1),
		"int8":    int8(-2),
		"int16":   int16(-3),
		"int32":   int32(-4),
		"int64":   int64(-5),
		"uint":    uint(6),
		"uint8":   uint8(7),
		"uint16":  uint16(8),
		"uint32":  uint32(9),
		"uint64":  uint64(10),
		"number":  json.Number("4.00"),
	} {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			if !writeNumericScalar(&buf, value) {
				t.Fatalf("writeNumericScalar() returned false")
			}
			if strings.TrimSpace(buf.String()) == "" {
				t.Fatalf("empty numeric output")
			}
		})
	}
}

func TestWriteJSONFallbackReportsMarshalError(t *testing.T) {
	var buf bytes.Buffer
	if err := writeJSONFallback(&buf, func() {}); err == nil {
		t.Fatalf("expected fallback marshal error")
	}
}
