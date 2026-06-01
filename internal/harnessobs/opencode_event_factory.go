package harnessobs

import "fmt"

func normalizedOpenCodeEvent(family, observedAt, sourceRef, actor string) Event {
	return normalizedEvent(
		fmt.Sprintf("%s-%s", sourceRef, family),
		family,
		family+"_observed",
		observedAt,
		sourceRef,
		actor,
	)
}
