package prreview

func bestPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet) PlaneResult {
	// bestPlaneResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	best := PlaneResult{Plane: plane, Status: StateNotAssessed, Usable: false, Reason: "required_plane_not_assessed", NextAction: "Run or import a reviewer result for this plane."}
	for _, result := range runs.Results {
		if result.Plane != plane {
			continue
		}
		candidate := planeResultWithModelCheck(roleByID[result.RoleID], result)
		if planeResultRank(candidate) > planeResultRank(best) {
			best = candidate
		}
	}
	return best
}
