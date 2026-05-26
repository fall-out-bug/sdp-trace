package main

import "time"

func (r *iterationRun) record(dur time.Duration, errMsg string, stop bool) bool {
	r.attempted++
	if errMsg != "" {
		r.lastErr = errMsg
		return stop
	}
	r.times = append(r.times, dur)
	return false
}
