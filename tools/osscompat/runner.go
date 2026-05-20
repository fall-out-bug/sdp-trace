package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// runAllProbes executes every registered probe and returns the results.
func runAllProbes() []probeResult {
	results := make([]probeResult, 0, len(registry))
	for _, p := range registry {
		results = append(results, runProbe(p))
	}
	return results
}

// runProbe executes a single probe, checking prerequisites first.
func runProbe(p probe) probeResult {
	if p.NeedsTool != "" && !hasTool(p.NeedsTool) {
		return probeResult{
			Name:   p.Name,
			State:  stateNotAssessed,
			Reason: fmt.Sprintf("required tool %q not found in PATH", p.NeedsTool),
		}
	}
	state, reason := p.Run()
	return probeResult{
		Name:   p.Name,
		State:  state,
		Reason: reason,
	}
}

// printResults writes probe results as text or JSON.
func printResults(w io.Writer, results []probeResult, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	for _, r := range results {
		line := fmt.Sprintf("%-24s %s", r.Name, r.State)
		if r.Reason != "" {
			line += "  — " + strings.ReplaceAll(r.Reason, "\n", " ")
		}
		fmt.Fprintln(w, line)
	}
	return nil
}
