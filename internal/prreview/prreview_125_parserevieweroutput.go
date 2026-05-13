package prreview

func parseReviewerOutput(base ReviewerResult, role ReviewRole, packet Packet, output []byte) (ReviewerResult, error) {
	// Parsed output must echo the packet, plane, and role. A structurally valid
	// answer for another packet is off-task, not reusable evidence.
	var parsed ReviewerResult
	if err := decodeReviewerOutput(output, &parsed); err != nil {
		return base, err
	}
	if reviewerOutputMismatched(parsed, role, packet) {
		base.Status = StatusOffTask
		base.Reason = "reviewer_output_wrong_packet_plane_or_role"
		return base, nil
	}
	return normalizeParsedReviewerOutput(parsed, base, role), nil
}
