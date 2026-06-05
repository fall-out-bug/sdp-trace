package harnessobs

import "encoding/json"

// addNormalizedSourceDigests hashes each normalized event after construction so
// downstream validation can replay exactly the portable JSONL bytes emitted.
func addNormalizedSourceDigests(events []Event) ([]Event, error) {
	for i := range events {
		data, err := json.Marshal(events[i])
		if err != nil {
			return nil, err
		}

		events[i].SourceDigest = digestLine(data)
	}
	return events, nil
}
