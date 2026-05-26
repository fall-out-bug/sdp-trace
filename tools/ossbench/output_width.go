package main

func maxNameWidth(results []benchmarkResult) int {
	maxW := defaultNameWidth
	for _, r := range results {
		if len(r.Name) > maxW {
			maxW = len(r.Name)
		}
	}
	return maxW
}
