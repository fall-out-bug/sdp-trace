package main

import "sort"

// stats returns min, max, and median of a non-empty slice.
func stats(values []float64) (min, max, median float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	min = sorted[0]
	max = sorted[len(sorted)-1]
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		median = sorted[mid]
	} else {
		median = (sorted[mid-1] + sorted[mid]) / 2
	}
	return min, max, median
}
