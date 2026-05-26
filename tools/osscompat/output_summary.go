package main

import (
	"fmt"
	"io"
)

func printSummary(w io.Writer, results []probeResult) error {
	pass, fail, cant, na := summarize(results)
	_, err := fmt.Fprintf(w, "\n%d pass, %d fail, %d cannot_verify, %d not_assessed\n",
		pass, fail, cant, na)
	return err
}
