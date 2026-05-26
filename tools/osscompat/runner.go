package main

// runAllProbes executes every probe in reg and returns the results.
func runAllProbes(reg []probe) []probeResult {
	results := make([]probeResult, 0, len(reg))
	for _, p := range reg {
		results = append(results, runProbe(p))
	}
	return results
}
