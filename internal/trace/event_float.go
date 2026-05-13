package trace

import (
	"encoding/json"
	"math"
	"strconv"
)

// This file owns float spelling for canonical trace digest bytes.

func trimFloatToJSON(value float64) json.Number {
	// json.Number lets numeric rendering reuse the same writer path.
	// That keeps float and decoded-number handling under one final spelling rule.
	return json.Number(trimFloatToString(value))
}

func trimFloatToString(value float64) string {
	// trimFloatToString preserves trace canonicalization evidence and replay semantics.
	// Keep digest input, ordering, and validation behavior stable when editing.
	if invalidFloat(value) {
		// Non-finite values cannot appear in canonical JSON.
		// Zero is the deterministic replacement used for replay.
		return "0"
	}
	if value == math.Trunc(value) {
		// Whole floats are rendered as integers to avoid spelling drift.
		// This keeps 1.0 and 1 equivalent for digest replay.
		return strconv.FormatInt(int64(value), 10)
	}
	if math.Abs(value) < 1e-15 {
		// Tiny values collapse to zero to avoid platform-specific exponent form.
		// This keeps near-zero noise out of event hash material.
		return "0"
	}
	// Non-whole finite values use the shortest decimal representation.
	// The 'f' format avoids exponent output after the near-zero guard.
	return strconv.FormatFloat(value, 'f', -1, 64)
}
