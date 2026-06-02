package harnessobs

import (
	"path/filepath"
	"time"
)

// writeObservedRun writes event artifacts before the run summary so a run file
// never points at events this helper failed to materialize.
func writeObservedRun(observedDir string, events []Event, observed Run) error {
	if err := writeObservedEvents(observedDir, events); err != nil {
		return err
	}
	return writeJSON(filepath.Join(observedDir, "run.json"), observed)
}

// writeObservedEvents keeps one event artifact per validated event ID under the
// observed events directory.
func writeObservedEvents(observedDir string, events []Event) error {
	for _, event := range events {
		if err := writeJSON(filepath.Join(observedDir, "events", event.EventID+".json"), event); err != nil {
			return err
		}
	}
	return nil
}

// observedRun builds the portable observed-run summary from the replayed event
// source digest and the validated event references.
func observedRun(ctx sessionCollectionContext, sourcePath, sourceDigest string, events []Event) Run {
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
