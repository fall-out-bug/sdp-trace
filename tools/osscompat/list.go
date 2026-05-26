package main

import (
	"fmt"
	"io"
)

// listProbes prints all registered probes.
func listProbes(stdout, stderr io.Writer, reg []probe) int {
	for _, p := range reg {
		if _, err := fmt.Fprintf(stdout, "%s\t%s\n", p.Name, p.Description); err != nil {
			fmt.Fprintf(stderr, "write error: %v\n", err)
			return 2
		}
	}
	return 0
}
