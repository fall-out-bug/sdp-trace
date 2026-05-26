package main

import (
	"fmt"
	"io"
)

func printResultsText(w io.Writer, results []probeResult) error {
	width := maxNameWidth(results)
	for _, r := range results {
		line := formatResultLine(r, width)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return printSummary(w, results)
}
