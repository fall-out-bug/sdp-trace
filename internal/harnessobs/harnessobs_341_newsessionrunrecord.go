package harnessobs

import (
	"time"
)

func newSessionRunRecord(profile SessionProfile, now time.Time, actionIDs []string, commit, commitState string) SessionRun {
	// newSessionRunRecord keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessProfilePath: profile.HarnessProfilePath,
		EventSourcePath:    profile.EventSourcePath,

		RawEventSourcePath: profile.RawEventSourcePath,
		RawEventFormat:     profile.RawEventFormat,
		SetupActionIDs:     actionIDs,

		CommandDigestState: StateCannotVerify,
		ProcessIDState:     StateCannotVerify,
		SourceCommit:       commit,
		SourceCommitState:  commitState,

		CollectionState:  StateCannotVerify,
		CollectionReason: "not_collected",
		CreatedAt:        now.Format(time.RFC3339),
	}
}
