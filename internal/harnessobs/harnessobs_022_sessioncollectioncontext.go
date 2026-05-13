package harnessobs

import (
	"time"
)

type sessionCollectionContext struct {
	profilePath        string
	runDir             string
	now                time.Time
	profile            SessionProfile
	session            SessionRun
	harnessProfile     Profile
	harnessProfilePath string
}
