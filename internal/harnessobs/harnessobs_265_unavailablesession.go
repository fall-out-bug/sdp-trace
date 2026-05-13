package harnessobs

import (
	"time"
)

func unavailableSession(ctx sessionCollectionContext) SessionRun {
	// unavailableSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := ctx.session

	session.CollectionState = StateCannotVerify
	session.CollectionReason = "source_unavailable"
	session.EndTime = ctx.now.Format(time.RFC3339)
	return session
}
