package harnessobs

import (
	"path/filepath"

	"time"
)

func newObservedRun(ctx observationContext) Run {
	// newObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.profile.ProfileID,
		HarnessFamily:      ctx.profile.HarnessFamily,
		EventSchemaVersion: ctx.profile.EventSchemaVersion,
		SourcePath:         filepath.Base(ctx.sourcePath),
		SourceDigest:       ctx.sourceDigest,
		EventCount:         len(ctx.events),
		EventRefs:          eventRefs(ctx.events),
		CreatedAt:          ctx.now.Format(time.RFC3339),
	}
}
