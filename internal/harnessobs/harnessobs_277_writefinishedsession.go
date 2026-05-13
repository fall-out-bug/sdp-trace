package harnessobs

import (
	"path/filepath"

	"time"
)

func writeFinishedSession(outDir string, session *SessionRun, end time.Time) error {
	session.EndTime = end.Format(time.RFC3339)
	return writeJSON(filepath.Join(outDir, "session.json"), *session)
}
