package capturedepth

import "github.com/fall_out_bug/sdp-trace/internal/adaptercapture"

type adapterEventSet map[string]bool

func missingAdapterEvents(run adaptercapture.RunEvidence) []string {
	// Missing event families are reported in required-event order so the query
	// mirrors the contract the adapter was expected to satisfy.
	seen := adapterEventSet{}
	for _, event := range run.AdapterEvents {
		seen[event.EventType] = true
	}
	missing := []string{}
	for _, required := range run.RequiredEventTypes {
		if !seen[required] {
			missing = append(missing, required)
		}
	}
	return missing
}
