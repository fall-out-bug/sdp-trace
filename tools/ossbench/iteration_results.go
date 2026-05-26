package main

import "time"

func (r iterationRun) results() ([]time.Duration, int, string) {
	return r.times, r.attempted, r.lastErr
}
