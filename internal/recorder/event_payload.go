package recorder

import (
	"encoding/json"
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Payload normalization is the boundary between typed trace structs and the
// generic event map stored in event JSON. The conversion keeps persisted data
// aligned with the bytes used for event hash computation.

func eventFilename(sequence int, eventType trace.EventType) string {
	// Filenames duplicate sequence and type for operator inspection; the event
	// body remains the authoritative hashed subject.
	return fmt.Sprintf("%06d-%s.json", sequence, eventType)
}

func toEventPayload(payload any) (map[string]any, error) {
	// Payload structs are normalized through JSON so event hashing sees the same
	// object shape that is persisted to disk.
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	// Trace events require object payloads so canonical hashing can use stable
	// field names instead of positional array semantics.
	payloadMap, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("event payload must be a JSON object")
	}
	return payloadMap, nil
}
