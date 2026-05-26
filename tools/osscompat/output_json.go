package main

import (
	"encoding/json"
	"io"
)

func printResultsJSON(w io.Writer, results []probeResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
