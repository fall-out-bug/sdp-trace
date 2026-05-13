package trace

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// This file owns final JSON scalar spelling for canonical trace bytes.
// Helpers here write directly to buffers so encoder defaults cannot change
// digest material.
// Numeric helpers call into this file after type detection has already chosen
// the canonical scalar family.

func writeJSONUnsignedInteger(buf *bytes.Buffer, value uint64) {
	// Unsigned integers use base-10 without type decoration.
	// This mirrors signed integer rendering for canonical digest bytes.
	// The writer never emits radix or width metadata.
	buf.WriteString(strconv.FormatUint(value, 10))
}

func writeJSONBool(buf *bytes.Buffer, value bool) {
	// Booleans are emitted as fixed JSON literals so payload and event digest
	// material never depends on encoder configuration.
	// The branch is intentionally explicit rather than using fmt or json.Marshal.
	if value {
		// True has exactly one JSON spelling in canonical output.
		buf.WriteString("true")
		return
	}
	// False is the only remaining boolean spelling.
	// Both branches avoid encoder allocation and whitespace.
	buf.WriteString("false")
}

func writeJSONFallback(buf *bytes.Buffer, value any) error {
	// writeJSONFallback preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	raw, err := json.Marshal(value)
	if err != nil {
		// Fallback marshal errors prevent non-numeric scalar hashing.
		return err
	}
	// json.Marshal is only used after canonical numeric/string/bool handling.
	// That fallback is for rare scalar-compatible values only.
	_, err = buf.Write(raw)
	// Surface write errors even though bytes.Buffer normally cannot fail.
	return err
}

func writeJSONString(buffer *bytes.Buffer, value string) {
	// strconv.Quote gives JSON-compatible string escaping for digest input.
	// This avoids a full encoder for a scalar leaf.
	// The quoted bytes are written directly into the canonical stream.
	buffer.WriteString(strconv.Quote(value))
}

func writeJSONInteger(buffer *bytes.Buffer, value int64) {
	// Signed integers use base-10 without type decoration.
	// No plus signs, widths, or radix prefixes are retained.
	// This mirrors JSON integer spelling expected by verifier replay.
	buffer.WriteString(strconv.FormatInt(value, 10))
}

func writeJSONNumber(buffer *bytes.Buffer, value json.Number) {
	// writeJSONNumber preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	// Decimal and exponent spellings normalize through float parsing when
	// possible; integer-like spellings remain unchanged.
	num := strings.TrimSpace(value.String())
	if num == "" {
		// Empty json.Number cannot be replayed as a source token; use zero.
		// This keeps malformed-but-decoded number values deterministic.
		buffer.WriteString("0")
		return
	}
	if strings.ContainsAny(num, ".eE") {
		// Only decimal/exponent tokens need normalization through float parsing.
		if normalized, err := strconv.ParseFloat(num, 64); err == nil {
			// Parsed decimal/exponent values share float trimming rules.
			buffer.WriteString(trimFloatToString(normalized))
			return
		}
	}
	// Integer-like or unparsable tokens keep their decoder-provided spelling.
	// The JSON decoder already accepted the token before this point.
	buffer.WriteString(num)
}
