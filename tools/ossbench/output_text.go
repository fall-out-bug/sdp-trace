package main

import (
	"fmt"
	"io"
)

func printTextResults(w io.Writer, results []benchmarkResult) error {
	width := maxNameWidth(results)
	for _, r := range results {
		if _, err := fmt.Fprint(w, formatResultLine(r, width)); err != nil {
			return err
		}
	}
	return nil
}
