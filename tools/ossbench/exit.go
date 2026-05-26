package main

func exitCode(results []benchmarkResult) int {
	for _, r := range results {
		if r.Error != "" {
			return 1
		}
	}
	return 0
}
