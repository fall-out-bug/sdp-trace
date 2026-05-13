package main

import "math"

// crapScore follows the CRAP formula shared by the report and gate paths.
func crapScore(complexity int, coverage float64) float64 {
	uncovered := 1 - (coverage / 100)
	return math.Pow(float64(complexity), 2)*math.Pow(uncovered, 3) + float64(complexity)
}

// exceedsThreshold keeps strict and non-strict gate wording in one place.
func exceedsThreshold(score float64, threshold float64, strictLess bool) bool {
	if strictLess {
		// Strict mode encodes gates phrased as "CRAP < N"; equality is still a
		// failed verdict, not a boundary pass.
		return score >= threshold
	}
	return score > threshold
}
