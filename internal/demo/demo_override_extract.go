package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func overrideRequestsFromEvents(events []trace.Event, contract trace.Contract) []OverrideRequest {
	// overrideRequestsFromEvents keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	requests := make([]OverrideRequest, 0)
	for _, event := range events {
		if event.EventType != trace.EventPolicyOverrideRequested {
			continue
		}
		requests = append(requests, overrideRequestFromEvent(event, contract))
	}
	sortOverrideRequests(requests)
	return requests
}
func overrideRequestFromEvent(event trace.Event, contract trace.Contract) OverrideRequest {
	// overrideRequestFromEvent keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	request := OverrideRequest{
		OverrideID: payloadString(event, "override_id"),
		State:      GatePass,
		CreatedAt:  payloadString(event, "created_at"),
	}
	request.State, request.Reason = overrideRequestFieldState(event)

	request.State, request.Reason = overrideRequestReferenceState(event, contract, request.State, request.Reason)
	return request
}
