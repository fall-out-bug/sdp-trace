package prreview

import (
	"time"
)

// newReviewerResult initializes reviewer evidence with explicit not_assessed
// model fields before any runner output is trusted.
func newReviewerResult(packet Packet, role ReviewRole, now time.Time) ReviewerResult {
	started := now.Format(time.RFC3339)

	result := ReviewerResult{
		ReviewRunID:    safeID("run-" + role.RoleID),
		PacketDigest:   packet.PacketDigest,
		Plane:          role.Plane,
		RoleID:         role.RoleID,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		Status:         StatusNotAssessed,
		StartedAt:      started,
		EndedAt:        started,
	}
	applyReviewerModelDefaults(&result)
	return result
}

func applyReviewerModelDefaults(result *ReviewerResult) {
	result.ObservedModel = StateNotAssessed
	result.ModelFamily = StateNotAssessed
	result.ModelVersion = StateNotAssessed
}
