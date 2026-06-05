package prreview

func bestPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet) PlaneResult {
	best := PlaneResult{Plane: plane, Status: StateNotAssessed, Usable: false, Reason: "required_plane_not_assessed", NextAction: "Run or import a reviewer result for this plane."}
	// Scan all attempts for the plane: retries can improve evidence, but a
	// usable findings result must still outrank a later clean result.
	for _, result := range runs.Results {
		if result.Plane != plane {
			continue
		}
		// Model provenance checks live in the next slice, but selection must use
		// their projected PlaneResult contract.
		candidate := planeResultWithModelCheck(roleByID[result.RoleID], result)
		if betterPlaneResult(candidate, best) {
			best = candidate
		}
	}
	return best
}

func betterPlaneResult(candidate, best PlaneResult) bool {
	return planeResultRank(candidate) > planeResultRank(best)
}
