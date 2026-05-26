package main

import "time"

func resultWithStats(res benchmarkResult, times []time.Duration) benchmarkResult {
	ms := make([]float64, len(times))
	for i, d := range times {
		ms[i] = float64(d) / float64(time.Millisecond)
	}
	min, max, median := stats(ms)
	res.SucceededIterations = len(times)
	res.MinMs = min
	res.MaxMs = max
	res.MedianMs = median
	res.AllMs = ms
	return res
}
