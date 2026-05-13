package prreview

import (
	"path/filepath"
)

func runReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	// runReview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

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
