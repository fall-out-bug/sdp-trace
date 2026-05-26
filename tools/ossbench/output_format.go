package main

import "fmt"

const defaultNameWidth = 24

func formatResultLine(r benchmarkResult, width int) string {
	if r.Error != "" {
		return formatErrorResultLine(r, width)
	}
	return fmt.Sprintf("%*s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
		-width, r.Name, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
}
