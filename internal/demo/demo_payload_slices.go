package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func payloadStringSlice(event trace.Event, key string) []string {
	// payloadStringSlice keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	value := payloadValue(event, key)
	if value == nil {
		return nil
	}
	return payloadAnyStringSlice(value)
}

func payloadAnyStringSlice(value any) []string {
	// payloadAnyStringSlice keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		return stringItems(typed)
	default:
		return nil
	}
}

func stringItems(items []any) []string {
	// stringItems keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	values := make([]string, 0, len(items))
	for _, item := range items {
		values = appendStringItem(values, item)
	}
	return values
}

func appendStringItem(values []string, item any) []string {
	// appendStringItem keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if text, ok := item.(string); ok {

		return append(values, text)
	}
	return values
}
