package main

import "time"

func buildResult(def benchmarkDef, cmdDisplay string, argv []string, attempted int, times []time.Duration, lastErr string) benchmarkResult {
	res := baseResult(def, cmdDisplay, argv, attempted, lastErr)
	if len(times) == 0 {
		return res
	}
	return resultWithStats(res, times)
}
