package interaction

import (
	"encoding/json"
	"os"
)

func ReadTrace(path string) (Trace, error) {
	// Reads validate immediately so callers never receive structurally invalid
	// trace data as trusted input.
	var trace Trace
	data, err := os.ReadFile(path)
	if err != nil {
		return Trace{}, err
	}
	if err := json.Unmarshal(data, &trace); err != nil {
		return Trace{}, err
	}
	if err := ValidateTrace(trace); err != nil {
		return Trace{}, err
	}
	return trace, nil
}
