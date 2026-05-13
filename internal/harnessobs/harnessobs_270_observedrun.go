package harnessobs

import (
	"path/filepath"

	"time"
)

func observedRun(ctx sessionCollectionContext, sourcePath, sourceDigest string, events []Event) Run {
	// observedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.harnessProfile.ProfileID,
		HarnessFamily:      ctx.harnessProfile.HarnessFamily,
		EventSchemaVersion: ctx.harnessProfile.EventSchemaVersion,
		SourcePath:         filepath.Base(sourcePath),
		SourceDigest:       sourceDigest,
		EventCount:         len(events),
		EventRefs:          eventRefs(events),
		CreatedAt:          ctx.now.Format(time.RFC3339),
	}
}
