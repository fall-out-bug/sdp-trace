package prreview

import "path/filepath"

// runReview executes reviewer roles after the public entrypoint has validated
// the profile and normalized options. It writes exactly one successful run-set
// artifact after all role results have been collected.
func runReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	rawDir, err := prepareRunDirectories(opts.OutDir)
	if err != nil {
		return RunSet{}, nil, err
	}
	results, err := runReviewRoles(packet, profile.Roles, opts, rawDir)
	if err != nil {
		return RunSet{}, nil, err
	}
	runSet := RunSet{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Results: results}
	if err := WriteJSON(filepath.Join(opts.OutDir, "results.json"), runSet); err != nil {
		return RunSet{}, nil, err
	}
	return runSet, nil, nil
}

// runReviewRoles preserves profile role order and stops on infrastructure
// errors instead of converting them into a successful partial run-set.
func runReviewRoles(packet Packet, roles []ReviewRole, opts RunOptions, rawDir string) ([]ReviewerResult, error) {
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
