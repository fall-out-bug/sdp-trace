package harnessobs

import (
	"sort"
	"time"
)

// newSessionRun builds the initial run record before command execution or
// collection can turn cannot_verify defaults into pass evidence.
func newSessionRun(profile SessionProfile, now time.Time) SessionRun {
	actionIDs := sessionSetupActionIDs(profile)
	commit, commitState := currentSourceCommitState()
	return newSessionRunRecord(profile, now, actionIDs, commit, commitState)
}

// newSessionRunRecord copies profile fields and records construction defaults;
// source commit values are passed in from the separate source discovery helper.
// It does not validate profile paths, collect events, or normalize raw input.
func newSessionRunRecord(profile SessionProfile, now time.Time, actionIDs []string, commit, commitState string) SessionRun {
	return SessionRun{
		// Profile identity and input locations are copied from the validated
		// session profile without revalidating path safety in this slice.
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessProfilePath: profile.HarnessProfilePath,
		EventSourcePath:    profile.EventSourcePath,

		// Raw event fields are construction-time copies; normalization happens
		// in the raw-event slice.
		RawEventSourcePath: profile.RawEventSourcePath,
		RawEventFormat:     profile.RawEventFormat,
		SetupActionIDs:     actionIDs,

		// Command and process evidence can only be verified after execution.
		CommandDigestState: StateCannotVerify,
		ProcessIDState:     StateCannotVerify,
		SourceCommit:       commit,
		SourceCommitState:  commitState,

		// Collection starts open until RunSession/CollectSession observes
		// events.
		CollectionState:  StateCannotVerify,
		CollectionReason: "not_collected",
		CreatedAt:        now.Format(time.RFC3339),
	}
}

// sessionSetupActionIDs stores setup action references deterministically so
// equivalent profiles produce stable session run payloads.
func sessionSetupActionIDs(profile SessionProfile) []string {
	actionIDs := make([]string, 0, len(profile.SetupActions))
	for _, action := range profile.SetupActions {
		actionIDs = append(actionIDs, action.ID)
	}

	sort.Strings(actionIDs)
	return actionIDs
}
