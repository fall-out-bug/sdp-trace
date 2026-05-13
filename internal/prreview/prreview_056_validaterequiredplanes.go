package prreview

import (
	"sort"
)

func validateRequiredPlanes(required map[string]bool, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) ([]PlaneResult, int, bool) {
	// validateRequiredPlanes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	planeResults := make([]PlaneResult, 0, len(required))
	usableCount := 0
	cannotVerify := false
	for plane := range required {
		result := requiredPlaneResult(plane, roleByID, runs, reasons, nextActions)
		usableCount += boolCount(result.Usable)
		cannotVerify = cannotVerify || planeCannotVerify(result.Status)
		planeResults = append(planeResults, result)
	}
	sort.Slice(planeResults, func(i, j int) bool { return planeResults[i].Plane < planeResults[j].Plane })
	return planeResults, usableCount, cannotVerify
}
