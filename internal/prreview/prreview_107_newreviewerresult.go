package prreview

import (
	"time"
)

func newReviewerResult(packet Packet, role ReviewRole, now time.Time) ReviewerResult {
	// newReviewerResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	started := now.Format(time.RFC3339)

	return ReviewerResult{
		ReviewRunID:    safeID("run-" + role.RoleID),
		PacketDigest:   packet.PacketDigest,
		Plane:          role.Plane,
		RoleID:         role.RoleID,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		ObservedModel:  StateNotAssessed,
		ModelFamily:    StateNotAssessed,
		ModelVersion:   StateNotAssessed,
		Status:         StatusNotAssessed,
		StartedAt:      started,
		EndedAt:        started,
	}
}
