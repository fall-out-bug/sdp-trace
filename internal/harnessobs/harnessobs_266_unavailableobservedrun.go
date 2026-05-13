package harnessobs

import (
	"path/filepath"

	"time"
)

func unavailableObservedRun(ctx sessionCollectionContext) Run {
	// unavailableObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
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
