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

// summarize counts results by state.
func summarize(results []probeResult) (pass, fail, cant, na int) {
	for _, r := range results {
		switch r.State {
		case statePass:
			pass++
		case stateFail:
			fail++
		case stateCannotVerify:
			cant++
		case stateNotAssessed:
			na++
		}
	}
	return
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
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	pass, fail, cant, na := summarize(results)
	_, err := fmt.Fprintf(w, "\n%d pass, %d fail, %d cannot_verify, %d not_assessed\n",
		pass, fail, cant, na)
	return err
}
