package policy

import "github.com/fall_out_bug/sdp-trace/internal/trace"

// CanAdapterEmit checks whether an adapter id can emit a given event type.
func (v AuthorityPolicyValidator) CanAdapterEmit(adapterID string, eventType trace.EventType) bool {
	// Adapter authority requires both a declared identity and an explicit event
	// family grant. Presence in the policy is not enough by itself.
	for _, adapter := range v.policy.AllowedAdapters {
		if adapterCanEmit(adapter, adapterID, string(eventType)) {
			return true
		}
	}
	return false
}

func adapterCanEmit(adapter AdapterAuthorityEntry, adapterID, eventType string) bool {
	return adapter.AdapterID == adapterID && adapter.AllowedByPolicy && stringInList(adapter.AllowedEventTypes, eventType)
}
