package main

import (
	"fmt"
	"strings"
)

func formatResultLine(r probeResult, width int) string {
	line := fmt.Sprintf("%*s %s", -width, r.Name, r.State)
	if r.Reason != "" {
		line += "  - " + strings.ReplaceAll(r.Reason, "\n", " ")
	}
	return line
}
