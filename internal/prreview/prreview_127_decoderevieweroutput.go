package prreview

import (
	"encoding/json"

	"strings"
)

func decodeReviewerOutput(output []byte, parsed *ReviewerResult) error {
	// decodeReviewerOutput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	decoder := json.NewDecoder(strings.NewReader(string(output)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(parsed)
}
