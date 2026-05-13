package trace

import (
	"encoding/json"
	"fmt"
)

func (event Event) syncPayloadRepresentation() (Event, error) {
	// Raw Payload wins when present because it is the serialized evidence shape
	// stored in historical run artifacts.
	switch {
	case len(event.Payload) > 0:
		return event.withDecodedEventPayload()
	case event.EventPayload != nil:
		return event.withEncodedPayload()
	default:
		return event, nil
	}
}

func (event Event) withDecodedEventPayload() (Event, error) {
	// Decoding fills EventPayload for callers that need structured access while
	// preserving the original raw payload bytes.
	if event.EventPayload != nil {
		return event, nil
	}
	var decoded map[string]any
	if err := json.Unmarshal(event.Payload, &decoded); err != nil {
		return Event{}, fmt.Errorf("invalid payload: %w", err)
	}
	event.EventPayload = decoded
	return event, nil
}

func (event Event) withEncodedPayload() (Event, error) {
	// Encoding supports newer callers that construct structured payloads first
	// and then need the raw digest input.
	payload, err := json.Marshal(event.EventPayload)
	if err != nil {
		return Event{}, err
	}
	event.Payload = payload
	return event, nil
}
