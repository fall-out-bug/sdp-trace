package main

import "math"

// clampMetric keeps computed quality metrics in the public 0..100 range.
func clampMetric(value float64) float64 {
	return math.Max(0, math.Min(100, value))
}
