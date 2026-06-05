package prreview

import (
	"encoding/json"
	"strings"
)

// parseReviewerOutput accepts only structured output bound to the reviewed
// packet, plane, and role.
func parseReviewerOutput(base ReviewerResult, role ReviewRole, packet Packet, output []byte) (ReviewerResult, error) {
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

func reviewerOutputMismatched(parsed ReviewerResult, role ReviewRole, packet Packet) bool {
	return parsed.PacketDigest != packet.PacketDigest || parsed.Plane != role.Plane || parsed.RoleID != role.RoleID
}

func decodeReviewerOutput(output []byte, parsed *ReviewerResult) error {
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(parsed)
}
