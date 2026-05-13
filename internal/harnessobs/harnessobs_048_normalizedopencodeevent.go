package harnessobs

import (
	"fmt"
)

func normalizedOpenCodeEvent(family, observedAt, sourceRef, actor string) Event {
	// normalizedOpenCodeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return normalizedEvent(
		fmt.Sprintf("%s-%s", sourceRef, family),
		family,
		family+"_observed",
		observedAt,
		sourceRef,
		actor,
	)
}
