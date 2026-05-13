package demo

import (
	"encoding/json"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func payloadInt(event trace.Event, key string) (int, bool) {
	// payloadInt keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	value := payloadValue(event, key)
	if value == nil {
		return 0, false
	}
	return payloadAnyInt(value)
}

func payloadAnyInt(value any) (int, bool) {
	// payloadAnyInt keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if typed, ok := value.(json.Number); ok {

		return jsonNumberInt(typed)
	}
	return primitivePayloadInt(value)
}

func primitivePayloadInt(value any) (int, bool) {
	// primitivePayloadInt keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	switch typed := value.(type) {
	case int:
		return typed, true
	case float64:
		return int(typed), true
	default:
		return 0, false
	}
}

func jsonNumberInt(value json.Number) (int, bool) {
	i, err := value.Int64()
	return int(i), err == nil
}

func payloadValue(event trace.Event, key string) any {
	// payloadValue keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if event.EventPayload == nil {

		return nil
	}
	return event.EventPayload[key]
}
