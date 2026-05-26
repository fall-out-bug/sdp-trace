package main

import "io"

// printResults writes probe results as text or JSON.
func printResults(w io.Writer, results []probeResult, asJSON bool) error {
	if asJSON {
		return printResultsJSON(w, results)
	}
	return printResultsText(w, results)
}
