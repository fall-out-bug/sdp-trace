package main

// exitCode returns 1 if any probe failed; 0 otherwise.
// not_assessed and cannot_verify do not cause a non-zero exit.
func exitCode(results []probeResult) int {
	for _, r := range results {
		if r.State == stateFail {
			return 1
		}
	}
	return 0
}
