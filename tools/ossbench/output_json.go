package main

import (
	"encoding/json"
	"io"
)

func printJSONResults(w io.Writer, results []benchmarkResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
