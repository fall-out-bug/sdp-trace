package main

import "time"

type iterationRun struct {
	times     []time.Duration
	attempted int
	lastErr   string
}

func newIterationRun(iterations int) iterationRun {
	return iterationRun{times: make([]time.Duration, 0, iterations)}
}
