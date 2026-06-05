package harnessobs

import "path/filepath"

func newObservedRun(ctx observationContext) Run {
	return Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          ctx.profile.ProfileID,
		HarnessFamily:      ctx.profile.HarnessFamily,
		EventSchemaVersion: ctx.profile.EventSchemaVersion,
		SourcePath:         filepath.Base(ctx.sourcePath),
		SourceDigest:       ctx.sourceDigest,
		EventCount:         len(ctx.events),
		EventRefs:          eventRefs(ctx.events),
		CreatedAt:          ctx.now.Format("2006-01-02T15:04:05Z07:00"),
	}
}
