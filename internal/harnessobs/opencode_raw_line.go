package harnessobs

import (
	"fmt"
	"time"
)

func normalizeOpenCodeRawLine(raw map[string]any, lineNo int, now time.Time) []Event {
	signals := rawSignals(raw)
	families := openCodeFamilies(raw, signals)
	if len(families) == 0 {
		return nil
	}
	sourceRef := fmt.Sprintf("raw-%06d", lineNo)
	observedAt := openCodeObservedAt(raw, now)
	return normalizedOpenCodeEvents(sortedFamilies(families), observedAt, sourceRef, openCodeActor(raw))
}
