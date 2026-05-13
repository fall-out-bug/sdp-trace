package main

import "strconv"

// roundMetric matches the one-decimal precision stored in MI baseline files.
func roundMetric(value float64) float64 {
	// FormatFloat/ParseFloat avoids platform-specific printf text leaking into
	// baseline comparisons while keeping the stored value numeric.
	rounded, _ := strconv.ParseFloat(strconv.FormatFloat(value, 'f', 1, 64), 64)
	return rounded
}
