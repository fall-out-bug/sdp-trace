package harnessobs

import (
	"fmt"

	"time"
)

func normalizeOpenCodeRawLine(raw map[string]any, lineNo int, now time.Time) []Event {
	// normalizeOpenCodeRawLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	signals := rawSignals(raw)
	families := openCodeFamilies(raw, signals)
	if len(families) == 0 {

		return nil
	}

	ordered := sortedFamilies(families)
	observedAt := openCodeObservedAt(raw, now)
	actor := openCodeActor(raw)

	sourceRef := fmt.Sprintf("raw-%06d", lineNo)
	return normalizedOpenCodeEvents(ordered, observedAt, sourceRef, actor)
}
