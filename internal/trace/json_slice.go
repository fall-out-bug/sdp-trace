package trace

import "encoding/json"

func mustMarshalJSONStringSlice(values []string) string {
	encoded, err := json.Marshal(values)
	if err != nil {
		// A string slice must always be JSON-marshalable; panic preserves the
		// invariant instead of silently weakening argv digest evidence.
		panic(err)
	}
	// The JSON array string is the canonical digest input for argv evidence.
	return string(encoded)
}
