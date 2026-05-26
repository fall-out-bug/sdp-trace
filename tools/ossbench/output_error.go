package main

import "fmt"

func formatErrorResultLine(r benchmarkResult, width int) string {
	return fmt.Sprintf("%*s  error: %s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
		-width, r.Name, r.Error, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
}
