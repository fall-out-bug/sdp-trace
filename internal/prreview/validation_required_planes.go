package prreview

import "sort"

func requiredPlaneSet(planes []string) map[string]bool {
	// Required planes are a set: duplicate declarations should not inflate
	// coverage expectations or change validation ordering.
	required := map[string]bool{}
	for _, plane := range planes {
		if plane != "" {
			required[plane] = true
		}
	}
	return required
}

func validateRequiredPlanes(required map[string]bool, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) ([]PlaneResult, int, bool) {
	planeResults := make([]PlaneResult, 0, len(required))
	usableCount := 0
	cannotVerify := false
	// Map iteration is intentionally normalized below so validation artifacts
	// stay reproducible across Go runtimes and reviewer replays.
	for plane := range required {
		result := requiredPlaneResult(plane, roleByID, runs, reasons, nextActions)
		usableCount += boolCount(result.Usable)
		cannotVerify = cannotVerify || planeCannotVerify(result.Status)
		planeResults = append(planeResults, result)
	}
	sort.Slice(planeResults, func(i, j int) bool { return planeResults[i].Plane < planeResults[j].Plane })
	return planeResults, usableCount, cannotVerify
}

func requiredPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) PlaneResult {
	result := bestPlaneResult(plane, roleByID, runs)
	appendPlaneValidationNotes(result, reasons, nextActions)
	return result
}
