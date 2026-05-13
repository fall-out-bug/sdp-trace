package prreview

func runReviewRoles(packet Packet, roles []ReviewRole, opts RunOptions, rawDir string) ([]ReviewerResult, error) {
	// runReviewRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	results := make([]ReviewerResult, 0, len(roles))
	for _, role := range roles {
		result, err := runRole(packet, role, opts, rawDir)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
