package harnessobs

import (
	"path/filepath"
	"time"
)

// markSessionSourceUnavailable writes an explicit cannot_verify session instead
// of treating absent source evidence as a successful collection.
func markSessionSourceUnavailable(ctx sessionCollectionContext) (SessionRun, Run, error) {
	session := unavailableSession(ctx)

	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}

	return session, unavailableObservedRun(ctx), nil
}

// unavailableSession preserves command/session metadata while marking the
// collection result as source_unavailable.
func unavailableSession(ctx sessionCollectionContext) SessionRun {
	session := ctx.session

	session.CollectionState = StateCannotVerify
	session.CollectionReason = "source_unavailable"
	session.EndTime = ctx.now.Format(time.RFC3339)
	return session
}

// unavailableObservedRun returns the paired zero-event run summary used when
// source evidence cannot be replayed.
func unavailableObservedRun(ctx sessionCollectionContext) Run {
	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.harnessProfile.ProfileID,
		HarnessFamily:      ctx.harnessProfile.HarnessFamily,
		EventSchemaVersion: ctx.harnessProfile.EventSchemaVersion,

		SourcePath: filepath.Base(ctx.profile.EventSourcePath),
		EventCount: 0,
		CreatedAt:  ctx.now.Format(time.RFC3339),
	}
}
