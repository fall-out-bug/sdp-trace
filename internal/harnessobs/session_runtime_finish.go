package harnessobs

import (
	"path/filepath"
	"time"
)

// collectFinishedSession writes the command end time, collects generated
// evidence, and then re-raises the command wait error if collection succeeded.
func collectFinishedSession(opts SessionOptions, session SessionRun, waitErr error, end time.Time) (SessionRun, Run, error) {
	if err := writeFinishedSession(opts.OutDir, &session, end); err != nil {
		return SessionRun{}, Run{}, err
	}

	collected, observed, err := CollectSession(SessionCollectOptions{ProfilePath: opts.ProfilePath, RunDir: opts.OutDir, Now: end})
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	if waitErr != nil {
		return collected, observed, waitErr
	}
	return collected, observed, nil
}

// writeFinishedSession stores the command completion time before collection
// reads the session record back from disk.
func writeFinishedSession(outDir string, session *SessionRun, end time.Time) error {
	session.EndTime = end.Format(time.RFC3339)
	return writeJSON(filepath.Join(outDir, "session.json"), *session)
}
