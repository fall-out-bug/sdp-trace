package prreview

import (
	"time"
)

func runRole(packet Packet, role ReviewRole, opts RunOptions, rawDir string) (ReviewerResult, error) {
	// runRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	result := newReviewerResult(packet, role, opts.Now)
	result.CommandDigest = commandDigest(role.Command)
	baseline, ready, err := prepareRoleRunner(&result, role, opts)
	if err != nil || !ready {
		return result, err
	}
	output, timedOut, err := runRoleCommand(role, opts)
	result.EndedAt = time.Now().UTC().Format(time.RFC3339)
	result = completeRoleResult(result, role, packet, opts.WorkDir, baseline, output, timedOut, err)
	return writeRawResult(result, rawDir, output)
}
