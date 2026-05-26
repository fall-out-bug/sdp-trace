package main

import "time"

func runIterations(def benchmarkDef, iterations int) ([]time.Duration, int, string) {
	run := newIterationRun(iterations)
	for i := 0; i < iterations; i++ {
		// Stop on the first failed iteration but keep durations already observed.
		if run.record(runIteration(def, i)) {
			break
		}
	}
	return run.results()
}
