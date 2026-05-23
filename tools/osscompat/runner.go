package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// runAllProbes executes every probe in reg and returns the results.
func runAllProbes(reg []probe) []probeResult {
	results := make([]probeResult, 0, len(reg))
	for _, p := range reg {
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
		pass += boolToInt(r.State == statePass)
		fail += boolToInt(r.State == stateFail)
		cant += boolToInt(r.State == stateCannotVerify)
		na += boolToInt(r.State == stateNotAssessed)
	}
	return
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func maxNameWidth(results []probeResult) int {
	maxW := 24
	for _, r := range results {
		if len(r.Name) > maxW {
			maxW = len(r.Name)
		}
	}
	return maxW
}

// printResults writes probe results as text or JSON.
func printResults(w io.Writer, results []probeResult, asJSON bool) error {
	if asJSON {
		return printResultsJSON(w, results)
	}
	return printResultsText(w, results)
}

func printResultsJSON(w io.Writer, results []probeResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func printResultsText(w io.Writer, results []probeResult) error {
	width := maxNameWidth(results)
	for _, r := range results {
		line := formatResultLine(r, width)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	pass, fail, cant, na := summarize(results)
	_, err := fmt.Fprintf(w, "\n%d pass, %d fail, %d cannot_verify, %d not_assessed\n",
		pass, fail, cant, na)
	return err
}

func formatResultLine(r probeResult, width int) string {
	line := fmt.Sprintf("%*s %s", -width, r.Name, r.State)
	if r.Reason != "" {
		line += "  - " + strings.ReplaceAll(r.Reason, "\n", " ")
	}
	return line
}
