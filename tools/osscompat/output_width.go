package main

func maxNameWidth(results []probeResult) int {
	maxW := 24
	for _, r := range results {
		if len(r.Name) > maxW {
			maxW = len(r.Name)
		}
	}
	return maxW
}
